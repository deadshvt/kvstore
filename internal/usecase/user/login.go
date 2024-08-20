package user

import (
	"context"

	"github.com/deadshvt/kvstore/internal/entity"
	"github.com/deadshvt/kvstore/internal/security"
	"github.com/deadshvt/kvstore/internal/usecase"
)

func (u *Usecase) Login(ctx context.Context, user *entity.User) (*usecase.LoginResult, error) {
	u.Logger.Info().Msg("Logging in...")

	repoResponse, err := u.UserRepository.GetUser(ctx, user.Username)
	if err != nil {
		return nil, err
	}

	decryptedPassword, err := security.Decrypt(repoResponse.Password, u.EncryptionKey)
	if err != nil {
		return nil, err
	}

	if user.Password != decryptedPassword {
		return &usecase.LoginResult{
			Success: false,
		}, nil
	}

	token, err := u.JWTService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &usecase.LoginResult{
		Success: true,
		Token:   token,
	}, nil
}
