package tarantool

import (
	"context"
	"github.com/tarantool/go-tarantool/v2"
	"os"

	"github.com/rs/zerolog"

	"github.com/deadshvt/kvstore/config"
	"github.com/deadshvt/kvstore/internal/entity"
	"github.com/deadshvt/kvstore/internal/errs"
)

type Tarantool struct {
	Conn   *tarantool.Connection
	Logger zerolog.Logger
}

func (db *Tarantool) Connect(ctx context.Context) error {
	db.Logger.Info().Msg("Connecting to tarantool...")

	config.Load(".env")

	var err error

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		dialer := tarantool.NetDialer{
			Address:  os.Getenv("TARANTOOL_HOST"),
			User:     os.Getenv("TARANTOOL_USER"),
			Password: os.Getenv("TARANTOOL_USER_PASSWORD"),
		}

		db.Conn, err = tarantool.Connect(ctx, dialer, tarantool.Opts{})
		if err != nil {
			return err
		}
	}

	_, err = db.Conn.Ping()

	return err
}

func (db *Tarantool) Disconnect(ctx context.Context) error {
	db.Logger.Info().Msg("Disconnecting from tarantool...")

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if db.Conn == nil {
			return nil
		}
	}

	return db.Conn.Close()
}

func (db *Tarantool) GetUser(ctx context.Context, username string) (*entity.EncryptedUser, error) {
	db.Logger.Info().Msg("Getting user...")

	user := &entity.EncryptedUser{
		Username: username,
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		resp, err := db.Conn.Do(tarantool.NewSelectRequest("user_store").
			Index("primary").
			Limit(1).
			Iterator(tarantool.IterEq).
			Key([]interface{}{username})).GetResponse()
		if err != nil {
			return nil, err
		}

		_, ok := resp.(*tarantool.SelectResponse)
		if !ok {
			return nil, errs.ErrCastValue
		}

		var userResponse []interface{}
		err = resp.DecodeTyped(&userResponse)
		if err != nil {
			return nil, err
		}

		if len(userResponse) == 0 {
			return nil, errs.ErrUserNotFound
		}

		userSlice, ok := userResponse[0].([]interface{})
		if !ok {
			return nil, errs.ErrCastValue
		}

		if len(userSlice) != 2 {
			return nil, errs.ErrUserNotFound
		}

		password, ok := userSlice[1].(string)
		if !ok {
			return nil, errs.ErrCastValue
		}

		user.Password = password
	}

	return user, nil
}
