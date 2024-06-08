package httputils

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/google/uuid"
)

type Handler func(logger.Logger, http.ResponseWriter, *http.Request) error

func HandleRequest(log logger.Logger, handle Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New()
		log := log.WithField("request_id", id.String())

		log.Infof("Server: handling request %s %s", r.Method, r.URL.Path)
		err := withRecover(func() error {
			return handle(log, w, r)
		})
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

func withRecover(f func() error) (err error) {
	defer func(err *error) {
		if r := recover(); r != nil {
			*err = fmt.Errorf("panic: %v", r)
		}
	}(&err)

	return f()
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
