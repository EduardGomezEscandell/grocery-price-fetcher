package daemon

import (
	"net"
	"net/http"
	"time"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/httputils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
)

type Daemon struct {
	log logger.Logger

	static  map[string]string
	dynamic map[string]httputils.Handler
}

func New(logger logger.Logger) Daemon {
	s := Daemon{
		log: logger,
	}

	s.static = make(map[string]string)
	s.dynamic = make(map[string]httputils.Handler)

	return s
}

func (d *Daemon) RegisterStaticEndpoint(path string, contentPath string) {
	d.static[path] = contentPath
}

func (d *Daemon) RegisterDynamicEndpoint(path string, handler httputils.Handler) {
	d.log.Infof("Registering dynamic endpoint: %s", path)
	d.dynamic[path] = handler
}

func (d *Daemon) Serve(lis net.Listener) (err error) {
	mux := http.NewServeMux()

	for path, fsPath := range d.static {
		fs := http.FileServer(http.Dir(fsPath))
		mux.Handle(path, fs)
	}

	for path, handler := range d.dynamic {
		mux.HandleFunc(path, httputils.HandleRequest(d.log, handler))
	}

	d.log.Infof("Server: serving on %s", lis.Addr())

	sv := http.Server{
		Handler:      mux,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}

	if err := sv.Serve(lis); err != nil {
		return err
	}

	d.log.Infof("Server: stopped serving")
	return nil
}
