package e2e_test

import (
	"context"
	"log"
	"os"
	"testing"

	e2e "github.com/EduardGomezEscandell/grocery-price-fetcher/end-to-end"
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := e2e.Setup(ctx); err != nil {
		log.Fatalf("Setup: %v", err)
	}

	e := m.Run()
	defer os.Exit(e)

	if err := e2e.Cleanup(ctx); err != nil {
		log.Fatalf("Cleanup: %v", err)
	}
}
