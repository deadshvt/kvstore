package pair

import (
	"github.com/rs/zerolog"

	"github.com/deadshvt/kvstore/internal/usecase"
)

type Handler struct {
	Usecase usecase.PairUsecase
	Logger  zerolog.Logger
}

func NewHandler(u usecase.PairUsecase, logger zerolog.Logger) *Handler {
	return &Handler{
		Usecase: u,
		Logger:  logger,
	}
}
