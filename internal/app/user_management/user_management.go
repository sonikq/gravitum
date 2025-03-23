package user_management

import (
	"context"
	"errors"
	"github.com/sonikq/gravitum_test_task/internal/config"
	"github.com/sonikq/gravitum_test_task/internal/handler"
	"github.com/sonikq/gravitum_test_task/internal/repository"
	httpserv "github.com/sonikq/gravitum_test_task/internal/server/http"
	"github.com/sonikq/gravitum_test_task/internal/service"
	"github.com/sonikq/gravitum_test_task/pkg/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	conf, err := config.Load()
	if err != nil {
		log.Fatal("failed to initialize config")
	}

	lg := logger.NewLogger(conf.ServiceName)

	ctx, cancel := context.WithTimeout(context.Background(), conf.CtxTimeOut)
	defer cancel()

	repo, err := repository.New(ctx, conf)
	if err != nil {
		lg.Fatal().Err(err).Msg("failed to initialize repository")
	}
	defer repo.Close()

	serviceManager := service.New(repo)
	router := handler.NewRouter(handler.Option{
		Conf:    conf,
		Logger:  lg,
		Service: serviceManager,
	})

	server := httpserv.NewServer(conf.RunAddress, router)

	go func() {
		err = server.Run()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			lg.Info().Err(err).Msg("failed to run http server")
		}
	}()

	lg.Info().Msg("Server listening on " + conf.RunAddress)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	ctx, cancel = context.WithTimeout(context.Background(), conf.CtxTimeOut)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		lg.Error().Err(err).Msg("error in shutting down server")
	}

	lg.Info().Msg("server stopped successfully")
}
