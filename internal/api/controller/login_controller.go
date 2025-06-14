package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/timuraipov/diploma/internal/domain"
	"github.com/timuraipov/diploma/pkg/logging"
)

type LoginController struct {
	l            *logging.ZapLogger
	loginUsecase domain.LoginUsecase
}

func NewLoginController(logger *logging.ZapLogger, loginUsecase domain.LoginUsecase) *LoginController {
	return &LoginController{
		l:            logger,
		loginUsecase: loginUsecase,
	}
}

func (l *LoginController) Login(w http.ResponseWriter, r *http.Request) {
	var request domain.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		handleErrorResponse(l.l, w, r, err, http.StatusBadRequest)
		return
	}

	authResponse, err := l.loginUsecase.Login(r.Context(), request)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			handleErrorResponse(l.l, w, r, err, http.StatusNotFound)
			return
		}
		if errors.Is(err, domain.ErrIncorrectLoginOrPassword) {
			handleErrorResponse(l.l, w, r, err, http.StatusUnauthorized)
		}

	}

	w.Header().Set("Authorization", "x-user-id "+authResponse.AccessToken)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, authResponse)
}
