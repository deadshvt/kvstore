package database

import (
	"context"

	"github.com/deadshvt/kvstore/internal/entity"
	"github.com/deadshvt/kvstore/internal/errs"
	"github.com/deadshvt/kvstore/internal/repository/user/database/tarantool"
)

const (
	Tarantool = "tarantool"
)

type UserDB interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error

	GetUser(ctx context.Context, username string) (*entity.EncryptedUser, error)
}

func NewUserDB(ctx context.Context, dbType string) (UserDB, error) {
	var db UserDB

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
