package httputils

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func HandleRequest(f func(*logrus.Entry, http.ResponseWriter, *http.Request) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New()
		log := logrus.WithField("request_id", id.String())

		log.Infof("Server: handling request %s %s", r.Method, r.URL.Path)
		err := f(log, w, r)
		if e := (RequestError{}); errors.As(err, &e) {
			log.Infof("Server: request error: %v", err)
			http.Error(w, e.Error(), e.Code)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Errorf("Server: internal error: %v", err)
			return
		}

		log.Infof("Server: request handled successfully")
	}
}

type RequestError struct {
	Code int
	Err  error
}

func Error(code int, msg string) RequestError {
	return RequestError{Code: code, Err: errors.New(msg)}
}

func Errorf(code int, format string, args ...any) RequestError {
	return RequestError{Code: code, Err: fmt.Errorf(format, args...)}
}

func (e RequestError) Error() string {
	return fmt.Sprintf("request error %d: %v", e.Code, e.Err)
}
