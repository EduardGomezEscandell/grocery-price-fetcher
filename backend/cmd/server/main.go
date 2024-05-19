package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/daemon"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func main() {
	sett := loadSettings()
	log := newLogger(sett)

	log.Debugf("Settings loaded: %#v", sett)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s, err := services.New(ctx, log, sett.Database)
	if err != nil {
		log.Fatalf("Could not initialize service: %v", err)
	}

	s.Run()
	defer s.Stop()

	lis, err := net.Listen("tcp", sett.Address)
	if err != nil {
		log.Fatalf("Could not listen: %v", err)
	}
	defer lis.Close()

	log.Info("Listening on ", lis.Addr().String())

	daemon := daemon.New(log)

	daemon.RegisterStaticEndpoint("/", sett.FrontEnd)
	s.Register(daemon.RegisterDynamicEndpoint)

	if err := daemon.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	log.Infof("Exiting")
}

type Settings struct {
	Verbosity int
	Database  database.Settings
	FrontEnd  string
	Address   string
}

func loadSettings() Settings {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <manifest>\n", os.Args[0])
	}

	if os.Args[1] == "" {
		log.Fatalf("Manifest path is empty")
	}

	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("Usage: %s <manifest>\n", os.Args[0])
		os.Exit(0)
	}

	out, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("Failed to read manifest: %v", err)
	}

	var s Settings
	if err := yaml.Unmarshal(out, &s); err != nil {
		log.Fatalf("Failed to unmarshal manifest: %v", err)
	}

	return s
}

func newLogger(s Settings) logger.Logger {
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
