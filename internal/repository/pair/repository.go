package pair

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/deadshvt/kvstore/internal/entity"
	"github.com/deadshvt/kvstore/internal/repository"
	"github.com/deadshvt/kvstore/internal/repository/pair/database"
)

type Repository struct {
	DB     database.PairDB
	Logger zerolog.Logger
}

func NewRepository(db database.PairDB, logger zerolog.Logger) repository.PairRepository {
	return &Repository{
		DB:     db,
		Logger: logger,
	}
}

func (r *Repository) SetPairs(ctx context.Context, pairs []*entity.EncryptedPair) (*repository.SetPairsResult, error) {
	r.Logger.Info().Msg("Setting pairs...")

	return r.DB.SetPairs(ctx, pairs)
}

func (r *Repository) GetPairs(ctx context.Context, keys []string) (*repository.GetPairsResult, error) {
	r.Logger.Info().Msg("Getting pairs...")

	return r.DB.GetPairs(ctx, keys)
}
