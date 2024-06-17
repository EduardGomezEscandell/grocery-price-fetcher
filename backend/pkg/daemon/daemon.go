package daemon

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"path"
	"time"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
)

type Settings struct {
	Host     string
	CertFile string
	KeyFile  string
}

func (s Settings) Defaults() Settings {
	return Settings{
		Host:     "localhost",
		CertFile: "/run/secrets/cert.pem",
		KeyFile:  "/run/secrets/key.pem",
	}
}

type Daemon struct {
	log      logger.Logger
	settings Settings

	endpoints map[string]func(http.ResponseWriter, *http.Request)
}

func New(logger logger.Logger, s Settings) Daemon {
	return Daemon{
		log:       logger,
		settings:  s,
		endpoints: make(map[string]func(http.ResponseWriter, *http.Request)),
	}
}

func (d *Daemon) RegisterEndpoint(path string, handler func(http.ResponseWriter, *http.Request)) {
	d.log.Infof("Registering endpoint: %s", path)
	d.endpoints[path] = handler
}

func (d *Daemon) Serve(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ch := make(chan error)
	go func() { ch <- d.serveHTTPS(ctx) }()
	go func() { ch <- d.serveHTTP(ctx) }()

	var err error
	for range 2 {
		err = errors.Join(err, <-ch)
		cancel()
	}
	close(ch)

	if err != nil {
		d.log.Errorf("Server: stopped serving: %v", err)
		return err
	}

	d.log.Infof("Server: stopped serving")
	return nil
}

func (d *Daemon) serveHTTPS(ctx context.Context) error {
	tlsConfig, err := d.tlsConfig()
	if err != nil {
		return fmt.Errorf("could not load TLS config: %v", err)
	}

	sv := http.Server{
		Addr:         net.JoinHostPort(d.settings.Host, "443"),
		Handler:      d.multiplexer(),
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
		TLSConfig:    tlsConfig,
	}

	context.AfterFunc(ctx, func() {
		_ = sv.Shutdown(context.Background())
	})

	ln, err := (&net.ListenConfig{}).Listen(ctx, "tcp", sv.Addr)
	if err != nil {
		return fmt.Errorf("could not listen: %v", err)
	}

	d.log.Infof("Listening on %s", ln.Addr())

	if err := sv.ServeTLS(ln, "", ""); errors.Is(err, http.ErrServerClosed) {
		return nil
	} else if err != nil {
		return fmt.Errorf("error serving HTTPS: %v", err)
	}

	return nil
}

// serveHTTP serves a redirect from HTTP to HTTPS.
func (d *Daemon) serveHTTP(ctx context.Context) error {
	sv := http.Server{
		Addr: net.JoinHostPort(d.settings.Host, "80"),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://"+path.Join(r.Host, r.RequestURI), http.StatusMovedPermanently)
		}),
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}

	context.AfterFunc(ctx, func() {
		_ = sv.Shutdown(context.Background())
	})

	ln, err := (&net.ListenConfig{}).Listen(ctx, "tcp", sv.Addr)
	if err != nil {
		return fmt.Errorf("could not listen: %v", err)
	}

	d.log.Infof("Listening on %s", ln.Addr())

	if err := sv.Serve(ln); errors.Is(err, http.ErrServerClosed) {
		return nil
	} else if err != nil {
		return fmt.Errorf("error serving HTTPS: %v", err)
	}

	return nil
}

func (d *Daemon) multiplexer() *http.ServeMux {
	mux := http.NewServeMux()

	for path, handler := range d.endpoints {
		mux.HandleFunc(path, handler)
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
