package controller

import (
	"net/http"

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
	w.Write([]byte(`reponse`))
}
