package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/timuraipov/diploma/bootstrap"
	"github.com/timuraipov/diploma/internal/domain"
	"github.com/timuraipov/diploma/pkg/logging"
)

type LoginController struct {
	l            *logging.ZapLogger
	loginUsecase domain.LoginUsecase
	cfg          *bootstrap.Config
}

func NewLoginController(logger *logging.ZapLogger, loginUsecase domain.LoginUsecase, cfg *bootstrap.Config) *LoginController {
	return &LoginController{
		l:            logger,
		loginUsecase: loginUsecase,
		cfg:          cfg,
	}
}

func (l *LoginController) Login(w http.ResponseWriter, r *http.Request) {
	var request domain.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, jsonError(err.Error()), http.StatusBadRequest)
		return
	}

	authResponse, err := l.loginUsecase.Login(r.Context(), request)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if errors.Is(err, domain.ErrIncorrectLoginOrPassword) {
			w.WriteHeader(http.StatusUnauthorized)
		}

	}

	w.Header().Set("Authorization", "x-user-id "+authResponse.AccessToken)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, authResponse)
}
