package pair

import (
	"context"

	"github.com/deadshvt/kvstore/internal/entity"
	"github.com/deadshvt/kvstore/internal/security"
	"github.com/deadshvt/kvstore/internal/usecase"
)

func (u *Usecase) SetPairs(ctx context.Context, pairs []*entity.Pair) (*usecase.SetPairsResult, error) {
	u.Logger.Info().Msg("Setting pairs...")

	encryptedPairs := make([]*entity.EncryptedPair, 0, len(pairs))
	responseErrors := make([]*entity.Error, 0, len(pairs))

	for _, pair := range pairs {
		encryptedValue, err := security.Encrypt(pair.Value, u.EncryptionKey)
		if err != nil {
			responseErrors = append(responseErrors, &entity.Error{
				Key:     pair.Key,
				Message: err.Error(),
			})

			continue
		}

		encryptedPairs = append(encryptedPairs, &entity.EncryptedPair{
			Key:   pair.Key,
			Value: encryptedValue,
		})
	}

	repoResponse, err := u.PairRepository.SetPairs(ctx, encryptedPairs)
	if err != nil {
		return nil, err
	}

	if len(repoResponse.Errors) > 0 || len(responseErrors) > 0 {
		responseErrors = append(responseErrors, repoResponse.Errors...)

		return &usecase.SetPairsResult{
			Success: false,
			Errors:  responseErrors,
		}, nil
	}

	return &usecase.SetPairsResult{
		Success: true,
	}, nil
}
