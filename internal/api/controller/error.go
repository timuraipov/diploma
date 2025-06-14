package controller

import (
	"net/http"

	"github.com/timuraipov/diploma/pkg/logging"
)

func jsonError(message string) string {
	return `{"message": "` + message + `"}`
}

func handleErrorResponse(l *logging.ZapLogger, w http.ResponseWriter, r *http.Request, err error, status int) {
	l.ErrorCtx(r.Context(), err.Error())
	http.Error(w, jsonError(err.Error()), status)
}
