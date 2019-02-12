package main

//go:generate go run util/webui/generate.go

// TODO: pass NodeID to ID generator

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

	"wallawire/ctxutil"
	"wallawire/model"
	"wallawire/repositories"
	"wallawire/services"
	"wallawire/web"
	"wallawire/web/auth"
	"wallawire/web/router"
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
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "help, h",
			Usage: "show help",
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
		cli.StringFlag{
			Name:   "postgres-cert",
			EnvVar: "WALLAWIRE_POSTGRES_CERT",
			Usage:  "postgres user certificate file",
		},
		cli.StringFlag{
			Name:   "postgres-key",
			EnvVar: "WALLAWIRE_POSTGRES_KEY",
			Usage:  "postgres user key file",
		},
		cli.StringFlag{
			Name:   "postgres-ca",
			EnvVar: "WALLAWIRE_POSTGRES_CA",
			Usage:  "postgres ca cert file",
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
	}

	// sort.Sort(cli.FlagsByName(app.Flags))
	// sort.Sort(cli.CommandsByName(app.Commands))

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
		zerolog.TimeFieldFormat = "15:04:03.000"
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}
	return nil
}

func start(c *cli.Context) error {

	if c.GlobalBool("help") {
		return cli.ShowAppHelp(c)
	}

	logger := log.With().Str("component", "main").Logger()
	logger.Info().Msg("Starting...")

	pURL := c.String("postgres-url")
	pCert := c.String("postgres-cert")
	pKey := c.String("postgres-key")
	pCA := c.String("postgres-ca")
	postgresURL := repositories.BuildPostgresURL(pURL, pCert, pKey, pCA)

	// attempt to connect to postgres in a loop
	dbCtx, _ := context.WithTimeout(context.Background(), dbConnectionTimeout)
	db, errDB := connectDatabase(dbCtx, postgresURL)
	if errDB != nil {
		logger.Error().Err(errDB).Msg("cannot connect to database")
		return errDB
	}
	log.Info().Msg("database connected")

	sqlDB := repositories.NewDatabase(db)
	tokenPassword := c.String("token-password")

	routerHandler, errRouter := instantiateRouter(sqlDB, tokenPassword)
	if errRouter != nil {
		return errRouter
	}

	serverAddr := c.String("server-addr")
	serverCertFile := c.String("server-cert")
	serverKeyFile := c.String("server-key")
	serverCAFile := c.String("server-ca")

	webService, errWebService := instantiateWebService(serverAddr, serverCertFile, serverKeyFile, serverCAFile, routerHandler)
	if errWebService != nil {
		return errWebService
	}

	errors := make(chan error, 1)
	webService.Start(errors)
	log.Info().Msg("webservice running")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	select {
	case e := <-errors:
		logger.Error().Err(e).Msg("Abort")
	case <-signals:
		fmt.Println()
		logger.Info().Msg("Caught signal")
	}

	logger.Info().Msg("Stopping...")
	webService.Stop(60 * time.Second)
	if err := db.Close(); err != nil {
		logger.Warn().Err(err).Msg("Cannot close database.")
	}
	logger.Info().Msg("Database disconnected")

	logger.Info().Msg("Exited")

	return nil

}

func instantiateRouter(db model.Database, tokenPassword string) (http.Handler, error) {

	idgen, errIdGen := ctxutil.NewIdGenerator()
	if errIdGen != nil {
		return nil, errIdGen
	}

	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(db, userRepo, idgen)
	loginHandler := auth.Login(userService, tokenPassword)
	logoutHandler := auth.Logout()
	whoami := auth.Whoami()
	changepassword := user.ChangePassword(userService, tokenPassword)
	changeusername := user.ChangeUsername(userService, tokenPassword)
	changeprofile := user.ChangeProfile(userService, tokenPassword)

	authenticator := auth.NewAuthenticator(tokenPassword)
	authorizerUsers := auth.NewAuthorizer(model.RoleNameUser)

	staticHandler := static.Handler()
	stat, errStatus := makeStatus()
	if errStatus != nil {
		return nil, errStatus
	}

	statusHandler, errStatusHandler := status.Handler(stat)
	if errStatusHandler != nil {
		return nil, errStatusHandler
	}

	return router.Router(router.Options{
		Authenticator:   authenticator,
		AuthorizerUsers: authorizerUsers,
		ChangePassword:  changepassword,
		ChangeUsername:  changeusername,
		ChangeProfile:   changeprofile,
		IdGenerator:     idgen,
		Login:           loginHandler,
		Logout:          logoutHandler,
		Static:          staticHandler,
		Status:          statusHandler,
		Whoami:          whoami,
	})
}

func instantiateWebService(serverAddr, serverCertFile, serverKeyFile, serverCAFile string, handler http.Handler) (*web.Service, error) {
	return web.New(handler, serverAddr, serverCertFile, serverKeyFile, serverCAFile), nil
}

func connectDatabase(ctx context.Context, postgresURL string) (db *sqlx.DB, err error) {

	logger := log.With().Str("component", "main").Str("fn", "connectDatabase").Logger()
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
			logger.Debug().Str("backoff", backoff.String()).Msg("backoff triggered")
			backoff *= 2
			continue
		case <-ctx.Done():
			t.Stop()
			logger.Debug().Msg("cancel triggered")
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
