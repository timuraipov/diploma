package controller

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/timuraipov/diploma/bootstrap"
	"github.com/timuraipov/diploma/internal/domain"
	"github.com/timuraipov/diploma/pkg/logging"
)

type RegisterController struct {
	l               *logging.ZapLogger
	registerUsecase domain.RegisterUsecase
	cfg             *bootstrap.Config
}

func NewRegisterController(logger *logging.ZapLogger, ru domain.RegisterUsecase, cfg *bootstrap.Config) *RegisterController {
	return &RegisterController{
		l:               logger,
		registerUsecase: ru,
		cfg:             cfg,
	}
}
func (rc *RegisterController) Register(w http.ResponseWriter, r *http.Request) {
	var request domain.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusBadRequest)
		return
	}

	user, err := rc.registerUsecase.Create(r.Context(), request)
	if err != nil {
		rc.l.ErrorCtx(r.Context(), err.Error())
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	accessToken, err := rc.registerUsecase.CreateAccessToken(&user, rc.cfg.AccessTokenSecret, rc.cfg.AccessTokenExpiryHour)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	refreshToken, err := rc.registerUsecase.CreateRefreshToken(&user, rc.cfg.RefreshTokenSecret, rc.cfg.RefreshTokenExpiryHour)
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}
	signupResponse := domain.RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	w.Header().Set("Authorization", "x-user-id "+accessToken)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, signupResponse)
}
