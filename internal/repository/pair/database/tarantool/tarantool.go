package tarantool

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/tarantool/go-tarantool/v2"

	"github.com/deadshvt/kvstore/config"
	"github.com/deadshvt/kvstore/internal/entity"
	"github.com/deadshvt/kvstore/internal/errs"
	"github.com/deadshvt/kvstore/internal/repository"
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
			Address:  os.Getenv("TARANTOOL_ADDR"),
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

func (db *Tarantool) SetPairs(ctx context.Context, pairs []*entity.EncryptedPair) (*repository.SetPairsResult, error) {
	db.Logger.Info().Msg("Setting pairs...")

	response := &repository.SetPairsResult{
		Errors: make([]*entity.Error, 0, len(pairs)),
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var wg sync.WaitGroup
		var mu sync.Mutex

		const maxRetries = 3
		const delay = time.Second

		for i := range pairs {
			wg.Add(1)
			go func(pair *entity.EncryptedPair) {
				defer wg.Done()

				var err error
				for i := 0; i < maxRetries; i++ {
					select {
					case <-ctx.Done():
						mu.Lock()
						response.Errors = append(response.Errors, &entity.Error{
							Key:     pair.Key,
							Message: ctx.Err().Error(),
						})
						mu.Unlock()
						return
					default:
						_, err = db.Conn.Replace("kv_store", []interface{}{pair.Key, pair.Value})
						if err == nil {
							break
						}

						time.Sleep(delay)
					}
				}

				mu.Lock()
				defer mu.Unlock()

				if err != nil {
					response.Errors = append(response.Errors, &entity.Error{
						Key:     pair.Key,
						Message: errs.WrapError(errs.ErrSetPair, err).Error(),
					})
					return
				}
			}(pairs[i])
		}

		wg.Wait()
	}

	return response, nil
}

func (db *Tarantool) GetPairs(ctx context.Context, keys []string) (*repository.GetPairsResult, error) {
	db.Logger.Info().Msg("Getting pairs...")

	response := &repository.GetPairsResult{
		Pairs:  make([]*entity.EncryptedPair, 0, len(keys)),
		Errors: make([]*entity.Error, 0, len(keys)),
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var wg sync.WaitGroup
		var mu sync.Mutex

		const maxRetries = 3
		const delay = time.Second

		for i := range keys {
			wg.Add(1)
			go func(key string) {
				defer wg.Done()

				var resp tarantool.Response
				var err error

				for i := 0; i < maxRetries; i++ {
					select {
					case <-ctx.Done():
						mu.Lock()
						response.Errors = append(response.Errors, &entity.Error{
							Key:     key,
							Message: ctx.Err().Error(),
						})
						mu.Unlock()
						return
					default:
						resp, err = db.Conn.Do(tarantool.NewSelectRequest("kv_store").
							Index("primary").
							Limit(1).
							Iterator(tarantool.IterEq).
							Key([]interface{}{key})).GetResponse()
						if err == nil {
							break
						}

						time.Sleep(delay)
					}
				}

				mu.Lock()
				defer mu.Unlock()

				if err != nil {
					response.Errors = append(response.Errors, &entity.Error{
						Key:     key,
						Message: errs.WrapError(errs.ErrGetPair, err).Error(),
					})
					return
				}

				_, ok := resp.(*tarantool.SelectResponse)
				if !ok {
					response.Errors = append(response.Errors, &entity.Error{
						Key:     key,
						Message: errs.ErrCastValue.Error(),
					})
					return
				}

				var pairResponse []interface{}
				err = resp.DecodeTyped(&pairResponse)
				if err != nil {
					response.Errors = append(response.Errors, &entity.Error{
						Key:     key,
						Message: err.Error(),
					})
					return
				}

				if len(pairResponse) == 0 {
					response.Errors = append(response.Errors, &entity.Error{
						Key:     key,
						Message: errs.ErrKeyNotFound.Error(),
					})
					return
				}

				pairSlice, ok := pairResponse[0].([]interface{})
				if !ok {
					response.Errors = append(response.Errors, &entity.Error{
						Key:     key,
						Message: errs.ErrCastValue.Error(),
					})
					return
				}

				if len(pairSlice) != 2 {
					response.Errors = append(response.Errors, &entity.Error{
						Key:     key,
						Message: errs.ErrKeyNotFound.Error(),
					})
					return
				}

				value, ok := pairSlice[1].(string)
				if !ok {
					response.Errors = append(response.Errors, &entity.Error{
						Key:     key,
						Message: errs.ErrCastValue.Error(),
					})
					return
				}

				response.Pairs = append(response.Pairs, &entity.EncryptedPair{
					Key:   key,
					Value: value,
				})
			}(keys[i])
		}

		wg.Wait()
	}

	return response, nil
}
