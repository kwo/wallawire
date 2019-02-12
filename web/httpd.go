package web

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog"

	"wallawire/ctxutil"
)

type Service struct {
	httpd    *http.Server
	logger   zerolog.Logger
	addr     string
	certFile string
	keyFile  string
	caFile   string
}

func New(router http.Handler, addr, certFile, keyFile, caFile string) *Service {
	return &Service{
		httpd: &http.Server{
			Addr:    addr,
			Handler: router,
		},
		logger:   ctxutil.NewLogger("httpd", "", nil),
		addr:     addr,
		certFile: certFile,
		keyFile:  keyFile,
		caFile:   caFile,
	}
}

func (z *Service) Start(errors chan<- error) {
	go func() {
		z.logger.Info().Str("addr", z.addr).Msg("Listening...")
		if len(z.certFile) == 0 && len(z.keyFile) == 0 {
			if err := z.httpd.ListenAndServe(); err != nil {
				errors <- err
			}
		} else {
			if err := z.httpd.ListenAndServeTLS(z.certFile, z.keyFile); err != nil {
				errors <- err
			}
		}
	}()
}

func (z *Service) Stop(timeout time.Duration) {

	z.logger.Info().Msg("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if errShutdown := z.httpd.Shutdown(ctx); errShutdown != nil {
		z.logger.Error().Err(errShutdown).Msg("Server shutdown failed")
		z.logger.Warn().Msg("Forcing server close...")
		if errClose := z.httpd.Close(); errClose != nil {
			z.logger.Error().Err(errClose).Msg("Server close failed")
		}
	}
	z.httpd = nil
	z.logger.Info().Msg("Shutdown complete")

}
