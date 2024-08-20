package user

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/deadshvt/kvstore/internal/entity"
	"github.com/deadshvt/kvstore/internal/repository/user/database"
)

type Repository struct {
	DB     database.UserDB
	Logger zerolog.Logger
}

func NewRepository(db database.UserDB, logger zerolog.Logger) *Repository {
	return &Repository{
		DB:     db,
		Logger: logger,
	}
}

func (r *Repository) GetUser(ctx context.Context, username string) (*entity.EncryptedUser, error) {
	r.Logger.Info().Msg("Getting user...")

	return r.DB.GetUser(ctx, username)
}
