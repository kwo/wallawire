package main

//go:generate go run tools/webui/generate.go

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli"

	"wallawire/idgen"
	"wallawire/logging"
	"wallawire/model"
	"wallawire/repository"
	"wallawire/schema"
	"wallawire/services"
	"wallawire/services/push"
	"wallawire/web"
	"wallawire/web/auth"
	"wallawire/web/router"
	"wallawire/web/sse"
	"wallawire/web/static"
	"wallawire/web/status"
	"wallawire/web/user"
)

const (
	ServiceName         = "wallawire"
	dbConnectionTimeout = time.Minute * 2
)

var (
	Version   = ""
	BuildTime = ""
	StartTime = time.Now().UTC().Truncate(time.Second)
)

func main() {

	app := cli.NewApp()
	app.Name = ServiceName
	app.Usage = "Wallawire Server"
	app.Version = Version
	app.HideHelp = true
	app.HideVersion = true
	cli.VersionPrinter = printVersion
	app.Before = initApp
	app.Action = start

	app.Commands = []cli.Command{
		{
			Name:   "help",
			Usage:  "show help",
			Action: cli.ShowAppHelp,
		},
		{
			Name:   "version",
			Usage:  "show version information",
			Action: printVersion,
		},
		{
			Name:   "migrate",
			Usage:  "upgrade database schema",
			Action: migrate,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "revert-last",
					Usage: "revert last migration",
				},
				cli.StringFlag{
					Name:   "postgres-url",
					EnvVar: "WALLAWIRE_POSTGRES_ROOT_URL",
					Usage:  "URL with which to connect to postgres as root",
				},
			},
		},
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "help, h",
			Usage: "show help",
		},
		cli.BoolTFlag{
			Name:   "production",
			EnvVar: "WALLAWIRE_PRODUCTION_MODE",
			Usage:  "production mode, when true, prevents using local ui files and log pretty",
		},
		cli.StringFlag{
			Name:   "server-addr",
			Value:  "localhost:8888",
			EnvVar: "WALLAWIRE_SERVER_ADDR",
			Usage:  "HTTP bind address as host:port",
		},
		cli.StringFlag{
			Name:   "server-cert",
			EnvVar: "WALLAWIRE_SERVER_CERT",
			Usage:  "HTTPS TLS server certificate file",
		},
		cli.StringFlag{
			Name:   "server-key",
			EnvVar: "WALLAWIRE_SERVER_KEY",
			Usage:  "HTTPS TLS server key file",
		},
		cli.StringFlag{
			Name:   "server-ca",
			EnvVar: "WALLAWIRE_SERVER_CA",
			Usage:  "HTTPS TLS server ca file",
		},
		cli.StringFlag{
			Name:   "token-password",
			EnvVar: "WALLAWIRE_TOKEN_PASSWORD",
			Usage:  "password to sign JWT tokens",
		},
		cli.StringFlag{
			Name:   "postgres-url",
			EnvVar: "WALLAWIRE_POSTGRES_URL",
			Usage:  "URL with which to connect to postgres",
		},
		cli.BoolFlag{
			Name:   "log-debug",
			EnvVar: "WALLAWIRE_LOG_DEBUG",
			Usage:  "show debug-level log messages",
		},
		cli.BoolFlag{
			Name:   "log-pretty",
			EnvVar: "WALLAWIRE_LOG_PRETTY",
			Usage:  "pretty print log messages to console",
		},
		cli.StringFlag{
			Name:   "ui-local-path",
			EnvVar: "WALLAWIRE_UI_LOCAL_PATH",
			Usage:  "if not empty serve UI files from given path, only if production=false",
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}

func printVersion(c *cli.Context) {

	st, errStatus := makeStatus()
	if errStatus != nil {
		fmt.Printf("cannot print status: %s", errStatus)
		return
	}

	st.PopulateNow()

	fmt.Printf("Service:      %s\n", st.ServiceName)
	fmt.Printf("Version:      %s\n", st.Version)
	fmt.Printf("Runtime:      %s\n", st.Runtime)
	fmt.Printf("Build Time:   %s\n", st.BuildTime)
	fmt.Printf("Start Time:   %s\n", st.StartTime)
	fmt.Printf("System Time:  %s\n", st.SystemTime)
	fmt.Printf("Uptime:       %s\n", st.Uptime)

}

func initApp(c *cli.Context) error {
	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if c.GlobalBool("log-debug") {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if c.GlobalBool("log-pretty") {
		zerolog.TimestampFunc = func() time.Time { return time.Now() }
		zerolog.TimeFieldFormat = "15:04:05.000"
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}
	return nil
}

func migrate(c *cli.Context) error {

	logger := logging.New(nil, "main", "migrate")
	logger.Info().Msg("starting...")

	revertLast := c.Bool("revert-last")

	dbx, errConnect := connectDatabase(c, logger)
	if errConnect != nil {
		return errConnect
	}
	defer func() {
		if err := dbx.Close(); err != nil {
			logger.Warn().Err(err).Msg("cannot close database")
		}
	}()

	if n, err := schema.Migrate(dbx.DB, revertLast); err != nil {
		logger.Error().Err(err).Msg("migration failed")
	} else if revertLast {
		logger.Info().Msg("revert schema successful")
	} else {
		logger.Info().Int("migrations", n).Msg("upgrade schema successful")
	}

	logger.Info().Msg("done")

	return nil

}

func start(c *cli.Context) error {

	if c.GlobalBool("help") {
		return cli.ShowAppHelp(c)
	}

	stat, errStatus := makeStatus()
	if errStatus != nil {
		return errStatus
	}
	stat.PopulateNow()

	logger := logging.New(nil, "main", "start")
	logger.Info().Interface("status", stat).Interface("params", formatParams(c)).Msg("starting...")

	productionMode := c.GlobalBool("production")

	// database
	db, errDB := connectDatabase(c, logger)
	if errDB != nil {
		logger.Error().Err(errDB).Msg("cannot connect to database")
		return errDB
	}
	log.Info().Msg("database connected")
	sqlDB := repository.NewDatabase(db)

	// repository
	repoid := idgen.NewUUIDGenerator()
	repo := repository.New(repoid)

	// push messaging
	pushMessenger := instantiatePushMessenger()
	heartbeatService := instantiateHeartbeatService(pushMessenger, stat)
	pushMessenger.AddOnClientConnectTrigger(func(userID, sessionID string) {
		heartbeatService.SendHeartbeat(time.Now().Truncate(time.Second), userID, sessionID)
	})

	// ui
	uiLocalPath := c.String("ui-local-path")
	var assetStore static.AssetStore
	if productionMode || len(uiLocalPath) == 0 {
		assetStore, _ = static.NewEmbeddedStore()
	} else {
		logger.Debug().Str("path", uiLocalPath).Msg("using local ui files")
		ls, errLocalStore := static.NewLocalStore(uiLocalPath)
		if errLocalStore != nil {
			return errLocalStore
		}
		assetStore = ls
	}

	// services
	idgenService := idgen.NewIdGenerator()
	userService := services.NewUserService(sqlDB, repo, idgenService)

	// router
	routerHandler, errRouter := instantiateRouter(c, userService, idgenService, assetStore, pushMessenger, stat)
	if errRouter != nil {
		return errRouter
	}

	errors := make(chan error, 1)

	// http server
	webService, errWebService := instantiateWebService(c, routerHandler)
	if errWebService != nil {
		return errWebService
	}
	webService.Start(errors)
	log.Info().Msg("webservice running")

	go heartbeatService.Start(time.Second * 60)

	// all started

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	select {
	case e := <-errors:
		logger.Error().Err(e).Msg("abort")
	case <-signals:
		fmt.Println()
		logger.Info().Msg("caught signal")
	}

	logger.Info().Msg("stopping...")
	webService.Stop(60 * time.Second)
	if err := db.Close(); err != nil {
		logger.Warn().Err(err).Msg("cannot close database.")
	}
	logger.Info().Msg("database disconnected")

	logger.Info().Msg("exited")

	return nil

}

func instantiatePushMessenger() *push.PushMessenger {
	return push.New()
}

func instantiateHeartbeatService(messageBus *push.PushMessenger, status *model.Status) *push.HeartbeatService {
	return push.NewHeartbeatService(messageBus, status)
}

func instantiateRouter(c *cli.Context, userService *services.UserService, idg *idgen.IdGenerator, assetStore static.AssetStore, pushMessenger *push.PushMessenger, stat *model.Status) (http.Handler, error) {

	tokenPassword := c.String("token-password")
	loginHandler := auth.Login(userService, tokenPassword)
	logoutHandler := auth.Logout()
	whoami := auth.Whoami()
	changepassword := user.ChangePassword(userService, tokenPassword)
	changeusername := user.ChangeUsername(userService, tokenPassword)
	changeprofile := user.ChangeProfile(userService, tokenPassword)

	authenticator := auth.NewAuthenticator(tokenPassword)
	authorizerUsers := auth.NewAuthorizer(model.RoleNameUser)

	staticHandler := static.Handler(assetStore)

	statusHandler, errStatusHandler := status.Handler(stat)
	if errStatusHandler != nil {
		return nil, errStatusHandler
	}

	sseHandler := sse.Handler(pushMessenger)

	return router.Router(router.Options{
		Authenticator:   authenticator,
		AuthorizerUsers: authorizerUsers,
		ChangePassword:  changepassword,
		ChangeUsername:  changeusername,
		ChangeProfile:   changeprofile,
		IdGenerator:     idg,
		Login:           loginHandler,
		Logout:          logoutHandler,
		Notifier:        sseHandler,
		Static:          staticHandler,
		Status:          statusHandler,
		Whoami:          whoami,
	})
}

func instantiateWebService(c *cli.Context, handler http.Handler) (*web.Service, error) {
	serverAddr := c.String("server-addr")
	serverCertFile := c.String("server-cert")
	serverKeyFile := c.String("server-key")
	serverCAFile := c.String("server-ca")
	return web.New(handler, serverAddr, serverCertFile, serverKeyFile, serverCAFile), nil
}

func connectDatabase(c *cli.Context, logger *zerolog.Logger) (db *sqlx.DB, err error) {

	postgresURL := c.String("postgres-url")
	ctx, _ := context.WithTimeout(context.Background(), dbConnectionTimeout)
	backoff := time.Second

Loop:
	for {

		db, err = sqlx.Open("postgres", postgresURL)
		if err == nil {
			err = db.Ping()
			if err == nil {
				return
			}
		}
		logger.Warn().Err(err).Msgf("retry in %s", backoff)

		t := time.NewTicker(backoff)
		select {
		case <-t.C:
			t.Stop()
			// logger.Debug().Str("backoff", backoff.String()).Msg("backoff triggered")
			backoff *= 2
			continue
		case <-ctx.Done():
			t.Stop()
			// logger.Debug().Msg("cancelled")
			break Loop
		}

	} // loop

	return

}

func makeStatus() (*model.Status, error) {

	toTime := func(value string) time.Time {
		t, _ := time.Parse(time.RFC3339, value)
		return t
	}

	return &model.Status{
		ServiceName: ServiceName,
		Version:     Version,
		BuildTime:   toTime(BuildTime),
		Runtime:     fmt.Sprintf("%s %s %s", runtime.Version(), runtime.GOOS, runtime.GOARCH),
		StartTime:   StartTime,
	}, nil

}

func formatParams(c *cli.Context) map[string]string {
	result := make(map[string]string)
	for _, flag := range c.App.Flags {
		result[flag.GetName()] = c.GlobalString(flag.GetName())
	}
	for _, flag := range c.Command.Flags {
		result[flag.GetName()] = c.GlobalString(flag.GetName())
	}
	return result
}
