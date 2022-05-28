package helpers

import (
	"ParsissCrm/internal/config"
	"fmt"
	"net/http"
	"runtime/debug"
)

var app *config.AppConfig

func NewHelpers(a *config.AppConfig) {
	app = a
}

func ClientError(rw http.ResponseWriter, status int) {
	app.InfoLog.Println("Client erroe with status of", status)
	http.Error(rw, http.StatusText(status), status)
}

func ServerError(rw http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
