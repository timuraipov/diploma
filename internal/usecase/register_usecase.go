package usecase

import (
	"context"

	"github.com/timuraipov/diploma/bootstrap"
	"github.com/timuraipov/diploma/internal/domain"
	"github.com/timuraipov/diploma/internal/tokenutil"
	"github.com/timuraipov/diploma/pkg/logging"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUsecase struct {
	l              *logging.ZapLogger
	userRepository domain.UserRepository
	cfg            *bootstrap.Config
}

func NewRegisterUsecase(logger *logging.ZapLogger, userRepository domain.UserRepository, cfg *bootstrap.Config) domain.RegisterUsecase {
	return &RegisterUsecase{
		l:              logger,
		userRepository: userRepository,
		cfg:            cfg,
	}
}
func (ru *RegisterUsecase) Create(ctx context.Context, userRequest domain.RegisterRequest) (domain.AuthResponse, error) {
	_, err := ru.GetByLogin(ctx, userRequest.Login)
	if err == nil {
		ru.l.ErrorCtx(ctx, "User already exists with the given login"+userRequest.Login)
		return domain.AuthResponse{}, domain.ErrUserAlreadyRegistered
	}
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		ru.l.ErrorCtx(ctx, err.Error())

		return domain.AuthResponse{}, err
	}
	userRequest.Password = string(encryptedPassword)
	user := domain.User{
		Login:    userRequest.Login,
		Password: userRequest.Password,
	}

	err = ru.userRepository.Create(ctx, &user)
	if err != nil {
		ru.l.ErrorCtx(ctx, err.Error())
		return domain.AuthResponse{}, err
	}
	authResp, err := tokenutil.GetTokensByUser(&user, ru.cfg.AccessTokenSecret, ru.cfg.RefreshTokenSecret, ru.cfg.AccessTokenExpiryHour, ru.cfg.RefreshTokenExpiryHour)
	if err != nil {
		ru.l.ErrorCtx(ctx, "Something wrong with generating access/refresh tokens", zap.Error(err))
		return domain.AuthResponse{}, err
	}
	response := domain.AuthResponse{AccessToken: authResp.AccessToken, RefreshToken: authResp.RefreshToken}
	return response, nil
}
func (ru *RegisterUsecase) GetByLogin(ctx context.Context, login string) (domain.User, error) {
	user, err := ru.userRepository.GetByLogin(ctx, login)
	return user, err
}
