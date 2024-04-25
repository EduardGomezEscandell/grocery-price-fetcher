package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/provider"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/server"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/providers/bonpreu"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/providers/mercadona"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func main() {
	s := loadSettings()
	setVerbosity(s)

	log.Debugf("Settings loaded: %#v", s)

	provider.Register("Bonpreu", bonpreu.New)
	provider.Register("Mercadona", mercadona.New)

	db, err := loadDatabase(s)
	if err != nil {
		log.Fatalf("Failed to load database: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db.UpdatePrices(ctx)

	lis, err := net.Listen("tcp", s.Address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer lis.Close()

	log.Info("Listening on ", lis.Addr().String())

	sv := server.New(db, server.WithStatic("/", s.FrontEnd))
	if err := sv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	log.Infof("Exiting")
}

type Settings struct {
	Verbosity int
	Database  string
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

func setVerbosity(s Settings) {
	switch s.Verbosity {
	case 0:
		log.SetLevel(log.InfoLevel)
	case 1:
		log.SetLevel(log.DebugLevel)
	case 2:
		log.SetLevel(log.TraceLevel)
	}

	log.Debug("DEBUG mode enabled")
}

func loadDatabase(s Settings) (*database.DB, error) {
	if s.Database == "" {
		return nil, errors.New("database path is empty")
	}

	out, err := os.ReadFile(s.Database)
	if err != nil {
		return nil, err
	}

	var db database.DB
	if err := json.Unmarshal(out, &db); err != nil {
		return nil, err
	}

	return &db, nil
}
