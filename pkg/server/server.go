package server

import (
	"net"
	"net/http"
	"time"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/server/api"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	db     *database.DB
	static map[string]string
}

type StaticEntry struct {
	url     string
	content string
}

func WithStatic(urlPath, contentPath string) StaticEntry {
	return StaticEntry{url: urlPath, content: contentPath}
}

func New(db *database.DB, e ...StaticEntry) Server {
	s := Server{db: db}
	s.static = make(map[string]string)
	for _, e := range e {
		s.static[e.url] = e.content
	}
	return s
}

func (s *Server) Serve(lis net.Listener) (err error) {
	mux := api.RegisterEndpoints(s.db)

	for path, content := range s.static {
		fs := http.FileServer(http.Dir(content))
		mux.Handle(path, fs)
	}

	log.Infof("Server: serving on %s", lis.Addr())

	sv := http.Server{
		Handler:      mux,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}

	if err := sv.Serve(lis); err != nil {
		return err
	}

	log.Infof("Server: stopped serving")
	return nil
}
