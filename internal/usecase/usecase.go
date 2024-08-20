package usecase

import (
	"context"

	"github.com/deadshvt/kvstore/internal/entity"
)

type UserUsecase interface {
	Login(ctx context.Context, user *entity.User) (*LoginResult, error)
}

type PairUsecase interface {
	SetPairs(ctx context.Context, pairs []*entity.Pair) (*SetPairsResult, error)
	GetPairs(ctx context.Context, keys []string) (*GetPairsResult, error)
}

type LoginResult struct {
	Success bool
	Token   string
}

type SetPairsResult struct {
	Success bool
	Errors  []*entity.Error
}

type GetPairsResult struct {
	Pairs  []*entity.Pair
	Errors []*entity.Error
}
