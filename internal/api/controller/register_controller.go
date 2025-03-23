package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/timuraipov/diploma/bootstrap"
	"github.com/timuraipov/diploma/internal/domain"
)

type RegisterController struct {
	RegisterUsecase domain.MockUsecase
	Cfg             *bootstrap.Config
}

func (rc *RegisterController) Register(w http.ResponseWriter, r *http.Request) {
	var request domain.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Print("login:", request.Login, "password:", request.Password)
	w.Write([]byte(`success`))
}
