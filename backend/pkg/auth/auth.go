package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/auth/google"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/dbtypes"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
)

type Manager struct {
	ctx    context.Context
	cancel context.CancelFunc

	done chan struct{}

	log    logger.Logger
	db     database.DB
	google google.Application

	gcRate time.Duration
}

type Settings struct {
	GoogleAuth google.Settings `yaml:"google"`
	GCRate     time.Duration   `yaml:"gc-rate"`
}

func (Settings) Defaults() Settings {
	return Settings{
		GoogleAuth: google.Settings{
			RedirectURI: "http://localhost",
		},
		GCRate: 6 * time.Hour,
	}
}

func NewManager(ctx context.Context, s Settings, log logger.Logger, db database.DB) (*Manager, error) {
	app, err := google.NewApplication(s.GoogleAuth)
	if err != nil {
		return nil, fmt.Errorf("could not create Google manager: %v", err)
	}

	ctx, cancel := context.WithCancel(ctx)
	return &Manager{
		ctx:    ctx,
		cancel: cancel,
		db:     db,
		google: app,
		log:    log,
		gcRate: s.GCRate,
	}, nil
}

func (m *Manager) Start() error {
	if err := m.checkStopped(); err != nil {
		return err
	}

	m.done = make(chan struct{})

	go func() {
		defer close(m.done)

		tk := time.NewTicker(m.gcRate)
		defer tk.Stop()

		for {
			select {
			case <-m.ctx.Done():
				return
			case <-tk.C:
			}

			m.log.Debugf("Checking for expired sessions")

			if err := m.db.PurgeSessions(); err != nil {
				m.log.Warningf("could not purge sessions: %v", err)
			}

			m.log.Debugf("Done checking for expired sessions")
		}
	}()

	return nil
}

func (m *Manager) Stop() {
	m.cancel()
	<-m.done
}

func (m *Manager) checkStopped() error {
	select {
	case <-m.ctx.Done():
		return m.ctx.Err()
	default:
		return nil
	}
}

func (m *Manager) NewSession(code string) (dbtypes.Session, error) {
	if err := m.checkStopped(); err != nil {
		return dbtypes.Session{}, err
	}

	session, err := m.google.Login(m.ctx, code)
	if err != nil {
		return dbtypes.Session{}, fmt.Errorf("could not login to Google: %v", err)
	}

	exists, err := m.db.LookupUser(session.User)
	if err != nil {
		return dbtypes.Session{}, fmt.Errorf("could not lookup user: %v", err)
	} else if !exists {
		err = m.db.SetUser(session.User)
		if err != nil {
			return dbtypes.Session{}, fmt.Errorf("could not create user: %v", err)
		}
	}

	if err := m.db.SetSession(session); err != nil {
		return dbtypes.Session{}, fmt.Errorf("could not save session: %v", err)
	}

	return session, nil
}

func (m *Manager) Remove(token string) error {
	if err := m.checkStopped(); err != nil {
		return err
	}

	if err := m.db.DeleteSession(token); err != nil {
		return fmt.Errorf("could not delete session: %v", err)
	}

	return nil
}

func (m *Manager) Get(key string) (dbtypes.Session, error) {
	if err := m.checkStopped(); err != nil {
		return dbtypes.Session{}, err
	}

	s, err := m.db.LookupSession(key)
	if err != nil {
		return dbtypes.Session{}, err
	}

	if time.Now().Before(s.NotAfter) {
		return s, nil
	}

	// Session has expired, refresh it

	s, err = m.google.Refresh(m.ctx, s)
	if err != nil {
		return dbtypes.Session{}, fmt.Errorf("could not refresh expired session: %v", err)
	}

	if err := m.db.SetSession(s); err != nil {
		return dbtypes.Session{}, fmt.Errorf("could not update refreshed session: %v", err)
	}

	return s, nil
}

type Getter interface {
	GetUserID(r *http.Request) (string, error)
}

func (m *Manager) GetUserID(r *http.Request) (string, error) {
	if err := m.checkStopped(); err != nil {
		return "", err
	}

	auth := r.Header.Get("Authorization")
	if auth == "" {
		return "", errors.New("missing token")
	}

	var key string
	if _, err := fmt.Sscanf(auth, "Bearer %s", &key); err != nil {
		return "", errors.New("could not parse token")
	}

	if key == "" {
		return "", errors.New("empty token")
	}

	s, err := m.Get(key)
	if err != nil {
		return "", fmt.Errorf("could not get session: %v", err)
	}

	return s.User, nil
}
