package usecase

import (
	"context"
	"time"

	"github.com/timuraipov/diploma/internal/domain"
	"github.com/timuraipov/diploma/internal/tokenutil"
	"github.com/timuraipov/diploma/pkg/logging"
)

type loginUsecase struct {
	l              *logging.ZapLogger
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

func NewLoginUsecase(logger *logging.ZapLogger, userRepository domain.UserRepository, timeout time.Duration) domain.LoginUsecase {
	return &loginUsecase{
		l:              logger,
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (lu *loginUsecase) GetUserByLogin(ctx context.Context, login string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, lu.contextTimeout)
	defer cancel()
	return lu.userRepository.GetByLogin(ctx, login)
}
func (lu *loginUsecase) CreateAccessToken(user *domain.User, secret string, expiry int) (accessToken string, err error) {
	return tokenutil.CreateAccessToken(user, secret, expiry)
}
func (lu *loginUsecase) CreateRefreshToken(user *domain.User, secret string, expiry int) (refreshToken string, err error) {
	return tokenutil.CreateRefreshToken(user, secret, expiry)
}
