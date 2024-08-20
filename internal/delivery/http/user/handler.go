package user

import (
	"github.com/rs/zerolog"

	"github.com/deadshvt/kvstore/internal/usecase"
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
