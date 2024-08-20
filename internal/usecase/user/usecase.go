package user

import (
	"github.com/rs/zerolog"

	"github.com/deadshvt/kvstore/internal/repository"
	"github.com/deadshvt/kvstore/internal/security"
)

type Usecase struct {
	UserRepository repository.UserRepository
	JWTService     *security.JWTService
	EncryptionKey  string
	Logger         zerolog.Logger
}

func NewUsecase(r repository.UserRepository, jwts *security.JWTService, key string, logger zerolog.Logger) *Usecase {
	return &Usecase{
		UserRepository: r,
		JWTService:     jwts,
		EncryptionKey:  key,
		Logger:         logger,
	}
}
