package pair

import (
	"github.com/rs/zerolog"

	"github.com/deadshvt/kvstore/internal/repository"
)

type Usecase struct {
	PairRepository repository.PairRepository
	EncryptionKey  string
	Logger         zerolog.Logger
}

func NewUsecase(r repository.PairRepository, key string, logger zerolog.Logger) *Usecase {
	return &Usecase{
		PairRepository: r,
		EncryptionKey:  key,
		Logger:         logger,
	}
}
