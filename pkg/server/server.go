package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/formatter"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/menu"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/server/httputils"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	db     *database.DB
	static map[string]string
}

type staticEntry struct {
	url     string
	content string
}

func WithStatic(urlPath, contentPath string) staticEntry {
	return staticEntry{url: urlPath, content: contentPath}
}

func New(db *database.DB, e ...staticEntry) Server {
	s := Server{db: db}
	s.static = make(map[string]string)
	for _, e := range e {
		s.static[e.url] = e.content
	}
	return s
}

func (s *Server) Serve(lis net.Listener) (err error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/hello-world", httputils.HandleRequest(s.helloWorld))
	mux.HandleFunc("/api/recipes", httputils.HandleRequest(s.handleRecipes))
	mux.HandleFunc("/api/menu", httputils.HandleRequest(s.handleMenu))

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

func (s *Server) handleRecipes(log *log.Entry, w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return httputils.Errorf(http.StatusMethodNotAllowed, "method %s not allowed", r.Method)
	}

	var names []string
	for _, r := range s.db.Recipes {
		names = append(names, r.Name)
	}

	if err := json.NewEncoder(w).Encode(names); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "failed to encode recipes: %v", err)
	}

	return nil
}

type MenuRequestData struct {
	Menu   menu.Menu          `json:",omitempty"`
	Pantry []menu.ProductData `json:",omitempty"`
	Format string             `json:",omitempty"`
}

func (s *Server) handleMenu(log *log.Entry, w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return httputils.Errorf(http.StatusMethodNotAllowed, "method %s not allowed", r.Method)
	}

	out, err := io.ReadAll(r.Body)
	if err != nil {
		return httputils.Error(http.StatusBadRequest, "failed to read request")
	}
	r.Body.Close()

	var data MenuRequestData
	if err := json.Unmarshal(out, &data); err != nil {
		return httputils.Errorf(http.StatusBadRequest, "failed to unmarshal request: %v", err)
	}

	log.Debugf("Received request with %d days and %d items in the pantry: ", len(data.Menu.Days), len(data.Pantry))

	if data.Format == "" {
		data.Format = "table"
	}
	f, err := formatter.New(data.Format)
	if err != nil {
		return httputils.Errorf(http.StatusBadRequest, "failed to create formatter: %v", err)
	}

	log.Debug("Selected formatter: ", data.Format)

	shoppingList, err := data.Menu.Compute(s.db, data.Pantry)
	if err != nil {
		return httputils.Errorf(http.StatusBadRequest, "failed to compute shopping list: %v", err)
	}

	log.Debug("Computed shopping list")

	if err := f.PrintHead(w, "Product", "Amount", "Cost"); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not write header to output: %w", err)
	}

	i := 0
	for _, p := range shoppingList {
		if p.Amount == 0 {
			continue
		}

		if err := f.PrintRow(w, map[string]any{
			"Product": p.Name,
			"Amount":  p.Amount,
			"Cost":    formatter.Euro(p.Cost),
		}); err != nil {
			return httputils.Errorf(http.StatusInternalServerError, "could not write results to output: %w", err)
		}
		i++
	}

	log.Debugf("Responded with %d items", i)

	if err := f.PrintTail(w); err != nil {
		return httputils.Errorf(http.StatusInternalServerError, "could not write footer to output: %w", err)
	}

	return nil
}

func (s *Server) helloWorld(_ *log.Entry, w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return httputils.Errorf(http.StatusMethodNotAllowed, "method %s not allowed", r.Method)
	}

	fmt.Fprintln(w, "Hello, world!")
	return nil
}
