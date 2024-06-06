package main

import (
	"context"
	"fmt"
	"os"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/daemon"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/settings"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func main() {
	sett := loadSettings()
	log := newLogger(sett)

	log.Debugf("Settings loaded: %s", sett)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s, err := services.New(ctx, log, sett.Services)
	if err != nil {
		log.Fatalf("Could not initialize service: %v", err)
	}

	s.Run()
	defer s.Stop()

	daemon := daemon.New(log, sett.Daemon)
	daemon.RegisterStaticEndpoint("/", sett.FrontEnd)
	s.Register(daemon.RegisterDynamicEndpoint)

	if err := daemon.Serve(ctx); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	log.Infof("Exiting")
}

func loadSettings() settings.Settings {
	switch len(os.Args) {
	case 1:
		fmt.Println("No manifest provided, using defaults.")
		return settings.Defaults()
	case 2:
	default:
		fmt.Printf("Usage: %s [MANIFEST]\n", os.Args[0])
		os.Exit(1)
	}

	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("Usage: %s [MANIFEST]\n", os.Args[0])
		os.Exit(0)
	}

	out, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("Failed to read manifest: %v", err)
	}

	s := settings.Defaults()
	if err := yaml.Unmarshal(out, &s); err != nil {
		log.Fatalf("Failed to unmarshal manifest: %v", err)
	}

	return s
}

func newLogger(s settings.Settings) logger.Logger {
	logger := logger.New()

	switch s.Verbosity {
	case 0:
		logger.SetLevel(int(log.InfoLevel))
	case 1:
		logger.SetLevel(int(log.DebugLevel))
	case 2:
		logger.SetLevel(int(log.TraceLevel))
	}

	logger.Debug("DEBUG mode enabled")
	return logger
}
