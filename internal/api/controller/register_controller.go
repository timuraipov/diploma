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
	_, err = rc.registerUsecase.GetByLogin(r.Context(), request.Login)
	if err != nil {
		rc.l.ErrorCtx(r.Context(), "User already exists with the given login"+request.Login)
		http.Error(w, jsonError("User already exists with the given login"), http.StatusConflict)
		return
	}
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		rc.l.ErrorCtx(r.Context(), err.Error())
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}
	request.Password = string(encryptedPassword)
	user := domain.User{
		Login:    request.Login,
		Password: request.Password,
	}
	err = rc.registerUsecase.Create(r.Context(), &user)
	if err != nil {
		rc.l.ErrorCtx(r.Context(), err.Error())
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	fmt.Print("login:", request.Login, " password:", request.Password)
	w.Write([]byte(`success`))
}
