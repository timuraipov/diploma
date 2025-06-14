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

type RegisterController struct {
	l               *logging.ZapLogger
	registerUsecase domain.RegisterUsecase
}

func NewRegisterController(logger *logging.ZapLogger, ru domain.RegisterUsecase, cfg *bootstrap.Config) *RegisterController {
	return &RegisterController{
		l:               logger,
		registerUsecase: ru,
	}
}

func (rc *RegisterController) Register(w http.ResponseWriter, r *http.Request) {
	var request domain.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		handleErrorResponse(rc.l, w, r, err, http.StatusBadRequest)
		return
	}

	signupResponse, err := rc.registerUsecase.Create(r.Context(), request)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyRegistered) {
			handleErrorResponse(rc.l, w, r, err, http.StatusConflict)
			return
		}
		handleErrorResponse(rc.l, w, r, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", "x-user-id "+signupResponse.AccessToken)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, signupResponse)
}
