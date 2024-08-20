package database

import (
	"context"

	"github.com/deadshvt/kvstore/internal/entity"
	"github.com/deadshvt/kvstore/internal/errs"
	"github.com/deadshvt/kvstore/internal/repository"
	"github.com/deadshvt/kvstore/internal/repository/pair/database/tarantool"
)

const (
	Tarantool = "tarantool"
)

type PairDB interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error

	SetPairs(ctx context.Context, pairs []*entity.EncryptedPair) (*repository.SetPairsResult, error)
	GetPairs(ctx context.Context, keys []string) (*repository.GetPairsResult, error)
}

func NewPairDB(ctx context.Context, dbType string) (PairDB, error) {
	var db PairDB

	switch dbType {
	case Tarantool:
		db = &tarantool.Tarantool{}
	default:
		return nil, errs.ErrUnsupportedDBType
	}

	err := db.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
