package usecase

import (
	"context"
	"time"

	"github.com/timuraipov/diploma/internal/domain"
	"github.com/timuraipov/diploma/internal/tokenutil"
	"github.com/timuraipov/diploma/pkg/logging"
)

type RegisterUsecase struct {
	l              *logging.ZapLogger
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

func NewRegisterUsecase(logger *logging.ZapLogger, userRepository domain.UserRepository, timeout time.Duration) domain.RegisterUsecase {
	return &RegisterUsecase{
		l:              logger,
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}
func (ru *RegisterUsecase) Create(ctx context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, ru.contextTimeout)
	defer cancel()
	return ru.userRepository.Create(ctx, user)
}
func (ru *RegisterUsecase) GetByLogin(ctx context.Context, login string) (domain.User, error) {
	user, err := ru.userRepository.GetByLogin(ctx, login)
	return user, err
}
func (ru *RegisterUsecase) GetByID(ctx context.Context, id int64) (domain.User, error) {
	user, err := ru.userRepository.GetByID(ctx, id)
	return user, err
}
func (ru *RegisterUsecase) CreateAccessToken(user *domain.User, secret string, expiry int) (accessToken string, err error) {
	return tokenutil.CreateAccessToken(user, secret, expiry)
}
func (ru *RegisterUsecase) CreateRefreshToken(user *domain.User, secret string, expiry int) (refreshToken string, err error) {
	return tokenutil.CreateRefreshToken(user, secret, expiry)
}
