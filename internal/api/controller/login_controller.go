package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/timuraipov/diploma/bootstrap"
	"github.com/timuraipov/diploma/internal/domain"
	"github.com/timuraipov/diploma/pkg/logging"
	"golang.org/x/crypto/bcrypt"
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
		http.Error(w, jsonError(err.Error()), http.StatusBadRequest)
		return
	}
	fmt.Println(request)

	user, err := l.loginUsecase.GetUserByLogin(r.Context(), request.Login)
	if err != nil {
		http.Error(w, jsonError("User not found with the given login"), http.StatusNotFound)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)) != nil {
		http.Error(w, jsonError("Invalid credentials"), http.StatusUnauthorized)
		return
	}

	accessToken, err := l.loginUsecase.CreateAccessToken(&user, l.cfg.AccessTokenSecret, l.cfg.AccessTokenExpiryHour)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	refreshToken, err := l.loginUsecase.CreateRefreshToken(&user, l.cfg.RefreshTokenSecret, l.cfg.RefreshTokenExpiryHour)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	loginResponse := domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loginResponse)
}
