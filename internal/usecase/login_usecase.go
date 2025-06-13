package usecase

import (
	"context"

	"github.com/timuraipov/diploma/internal/domain"
	"github.com/timuraipov/diploma/internal/tokenutil"
	"github.com/timuraipov/diploma/pkg/logging"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type loginUsecase struct {
	l              *logging.ZapLogger
	userRepository domain.UserRepository
	authSecret     domain.AuthSecret
}

func NewLoginUsecase(logger *logging.ZapLogger, userRepository domain.UserRepository, authSecret domain.AuthSecret) domain.LoginUsecase {
	return &loginUsecase{
		l:              logger,
		userRepository: userRepository,
		authSecret:     authSecret,
	}
}

func (lu *loginUsecase) Login(ctx context.Context, credentials domain.LoginRequest) (domain.AuthResponse, error) {
	user, err := lu.userRepository.GetByLogin(ctx, credentials.Login)
	if err != nil {
		return domain.AuthResponse{}, domain.ErrUserNotFound
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)) != nil {
		return domain.AuthResponse{}, domain.ErrIncorrectLoginOrPassword
	}
	authResp, err := tokenutil.GetTokensByUser(&user, lu.authSecret.AccessTokenSecret, lu.authSecret.RefreshTokenSecret, lu.authSecret.AccessTokenExpiryHour, lu.authSecret.RefreshTokenExpiryHour)
	if err != nil {
		lu.l.ErrorCtx(ctx, "Something wrong with generating access/refresh tokens", zap.Error(err))
		return domain.AuthResponse{}, err
	}
	response := domain.AuthResponse{AccessToken: authResp.AccessToken, RefreshToken: authResp.RefreshToken}
	return response, nil
}
