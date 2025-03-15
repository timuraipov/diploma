package controller

import (
	"net/http"

	"github.com/timuraipov/diploma/bootstrap"
	"github.com/timuraipov/diploma/internal/domain"
)

type LoginController struct {
	LoginUsecase domain.MockUsecase
	Cfg          *bootstrap.Config
}

func (l *LoginController) Login(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`reponse`))
}
