package repository

import (
	"context"

	"github.com/deadshvt/kvstore/internal/entity"
)

type UserRepository interface {
	GetUser(ctx context.Context, username string) (*entity.EncryptedUser, error)
}

type PairRepository interface {
	SetPairs(ctx context.Context, pairs []*entity.EncryptedPair) (*SetPairsResult, error)
	GetPairs(ctx context.Context, keys []string) (*GetPairsResult, error)
}

type SetPairsResult struct {
	Errors []*entity.Error
}

type GetPairsResult struct {
	Pairs  []*entity.EncryptedPair
	Errors []*entity.Error
}
