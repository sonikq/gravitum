package user_management

import (
	"github.com/sonikq/gravitum_test_task/internal/config"
	"github.com/sonikq/gravitum_test_task/internal/service"
	"github.com/sonikq/gravitum_test_task/pkg/logger"
)

const (
	contentTypeHeaderKey = "Content-Type"
	contentTypeJSON      = "application/json"
	contentTypeTextPlain = "text/plain"
)

type Handler struct {
	config  config.Config
	logger  *logger.Logger
	service *service.Service
}

type HandlerConfig struct {
	Config  config.Config
	Logger  *logger.Logger
	Service *service.Service
}

func New(cfg *HandlerConfig) *Handler {
	return &Handler{
		config:  cfg.Config,
		logger:  cfg.Logger,
		service: cfg.Service,
	}
}
