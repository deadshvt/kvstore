package security

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/deadshvt/kvstore/internal/entity"
	"github.com/deadshvt/kvstore/internal/errs"
)

type JWTService struct {
	Method    jwt.SigningMethod
	SecretKey string
}

type Claims struct {
	User *entity.User
	jwt.StandardClaims
}

func NewJWTService(method jwt.SigningMethod, secretKey string) *JWTService {
	return &JWTService{
		Method:    method,
		SecretKey: secretKey,
	}
}

func NewClaims(user *entity.User, t time.Time) *Claims {
	return &Claims{
		User: &entity.User{
			Username: user.Username,
			Password: user.Password,
		},
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  t.Unix(),
			ExpiresAt: t.Add(time.Hour * 24).Unix(),
		},
	}
}

func (s *JWTService) GenerateToken(user *entity.User) (string, error) {
	token := jwt.NewWithClaims(s.Method, NewClaims(user, time.Now()))
	signedToken, err := token.SignedString([]byte(s.SecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *JWTService) VerifyToken(signedToken string) error {
	token, err := jwt.ParseWithClaims(signedToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errs.ErrInvalidSigningMethod
		}
		return []byte(s.SecretKey), nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errs.ErrInvalidToken
	}

	_, ok := token.Claims.(*Claims)
	if !ok {
		return errs.ErrInvalidToken
	}

	return nil
}
