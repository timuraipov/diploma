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

type loginUsecase struct {
	l              *logging.ZapLogger
	userRepository domain.UserRepository
	cfg            *bootstrap.Config
}

func NewLoginUsecase(logger *logging.ZapLogger, userRepository domain.UserRepository, cfg *bootstrap.Config) domain.LoginUsecase {
	return &loginUsecase{
		l:              logger,
		userRepository: userRepository,
		cfg:            cfg,
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
	authResp, err := tokenutil.GetTokensByUser(&user, lu.cfg.AccessTokenSecret, lu.cfg.RefreshTokenSecret, lu.cfg.AccessTokenExpiryHour, lu.cfg.RefreshTokenExpiryHour)
	if err != nil {
		lu.l.ErrorCtx(ctx, "Something wrong with generating access/refresh tokens", zap.Error(err))
		return domain.AuthResponse{}, err
	}
	response := domain.AuthResponse{AccessToken: authResp.AccessToken, RefreshToken: authResp.RefreshToken}
	return response, nil
}
