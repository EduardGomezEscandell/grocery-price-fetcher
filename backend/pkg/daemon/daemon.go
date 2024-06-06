package daemon

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/httputils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
)

type Settings struct {
	Address  string
	CertFile string
	KeyFile  string
}

func (s Settings) Defaults() Settings {
	return Settings{
		Address:  "localhost:443",
		CertFile: "/run/secrets/cert.pem",
		KeyFile:  "/run/secrets/key.pem",
	}
}

type Daemon struct {
	log      logger.Logger
	settings Settings

	static  map[string]string
	dynamic map[string]httputils.Handler
}

func New(logger logger.Logger, s Settings) Daemon {
	return Daemon{
		log:      logger,
		settings: s,
		static:   make(map[string]string),
		dynamic:  make(map[string]httputils.Handler),
	}
}

func (d *Daemon) RegisterStaticEndpoint(path string, contentPath string) {
	d.static[path] = contentPath
}

func (d *Daemon) RegisterDynamicEndpoint(path string, handler httputils.Handler) {
	d.log.Infof("Registering dynamic endpoint: %s", path)
	d.dynamic[path] = handler
}

func (d *Daemon) Serve(ctx context.Context) (err error) {
	tlsConfig, err := d.tlsConfig()
	if err != nil {
		return fmt.Errorf("could not load TLS config: %v", err)
	}

	sv := http.Server{
		Addr:         d.settings.Address,
		Handler:      d.multiplexer(),
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
		TLSConfig:    tlsConfig,
	}

	ln, err := (&net.ListenConfig{}).Listen(ctx, "tcp", sv.Addr)
	if err != nil {
		return fmt.Errorf("could not listen: %v", err)
	}

	d.log.Infof("Listening on %s", ln.Addr())

	if err := sv.ServeTLS(ln, "", ""); err != nil {
		return err
	}

	d.log.Infof("Server: stopped serving")
	return nil
}

func (d *Daemon) multiplexer() *http.ServeMux {
	mux := http.NewServeMux()

	for path, fsPath := range d.static {
		fs := http.FileServer(http.Dir(fsPath))
		mux.Handle(path, fs)
	}

	for path, handler := range d.dynamic {
		mux.HandleFunc(path, httputils.HandleRequest(d.log, handler))
	}

	return mux
}

func (d *Daemon) tlsConfig() (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(d.settings.CertFile, d.settings.KeyFile)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}, nil
}
