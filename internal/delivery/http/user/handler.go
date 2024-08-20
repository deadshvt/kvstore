package user

import (
	"github.com/deadshvt/kvstore/internal/usecase"
	"github.com/rs/zerolog"
)

type Handler struct {
	Usecase usecase.UserUsecase
	Logger  zerolog.Logger
}

func NewHandler(u usecase.UserUsecase, logger zerolog.Logger) *Handler {
	return &Handler{
		Usecase: u,
		Logger:  logger,
	}
}
