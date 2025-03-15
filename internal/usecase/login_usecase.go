package usecase

import (
	"time"

	"github.com/timuraipov/diploma/internal/domain"
)

type loginUsecase struct {
	userRepository domain.MockRepository
	contextTimeout time.Duration
}

func NewLoginUsecase(userRepository domain.MockRepository, timeout time.Duration) domain.MockUsecase {
	return &loginUsecase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}
