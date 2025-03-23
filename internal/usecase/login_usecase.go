package usecase

import (
	"context"
	"time"

	"github.com/timuraipov/diploma/internal/domain"
	"github.com/timuraipov/diploma/pkg/logging"
)

type loginUsecase struct {
	l              *logging.ZapLogger
	userRepository domain.MockRepository
	contextTimeout time.Duration
}

func NewLoginUsecase(logger *logging.ZapLogger, userRepository domain.MockRepository, timeout time.Duration) domain.LoginUsecase {
	return &loginUsecase{
		l:              logger,
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (lu *loginUsecase) GetUserByLogin(c context.Context, login string) (domain.User, error) {
	return domain.User{}, nil
}
func (lu *loginUsecase) CreateAccessToken(user *domain.User, secret string, expiry int) (accessToken string, err error) {
	return "", nil
}
func (lu *loginUsecase) CreateRefreshToken(user *domain.User, secret string, expiry int) (refreshToken string, err error) {
	return "", nil
}
