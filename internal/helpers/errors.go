package helpers

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

type Resource struct {
	ErrorLog *log.Logger
	InfoLog *log.Logger
}

func (r Resource) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	r.ErrorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError)
}

func (r Resource) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (r Resource) notFound(w http.ResponseWriter) {
	r.clientError(w, http.StatusNotFound)
}