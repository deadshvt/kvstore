package pair

import (
	"context"

	"github.com/deadshvt/kvstore/internal/entity"
	"github.com/deadshvt/kvstore/internal/security"
	"github.com/deadshvt/kvstore/internal/usecase"
)

func (u *Usecase) GetPairs(ctx context.Context, keys []string) (*usecase.GetPairsResult, error) {
	u.Logger.Info().Msg("Getting pairs...")

	repoResponse, err := u.PairRepository.GetPairs(ctx, keys)
	if err != nil {
		return nil, err
	}

	pairs := make([]*entity.Pair, 0, len(repoResponse.Pairs))
	responseErrors := make([]*entity.Error, 0, len(keys))

	for _, pair := range repoResponse.Pairs {
		decryptedValue, err := security.Decrypt(pair.Value, u.EncryptionKey)
		if err != nil {
			responseErrors = append(responseErrors, &entity.Error{
				Key:     pair.Key,
				Message: err.Error(),
			})

			continue
		}

		pairs = append(pairs, &entity.Pair{
			Key:   pair.Key,
			Value: decryptedValue,
		})
	}

	responseErrors = append(responseErrors, repoResponse.Errors...)

	return &usecase.GetPairsResult{
		Pairs:  pairs,
		Errors: responseErrors,
	}, nil
}
