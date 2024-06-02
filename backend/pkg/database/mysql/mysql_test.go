package mysql_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/mysql"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/product"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers/blank"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/testutils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/types"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	providers.Register(blank.Provider{})

	if os.Getenv("MYSQL_SKIP_TEST_MAIN") == "" {
		fmt.Println("Starting database")
		cmd := exec.Command("make", "stand-up")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatal("could not start database")
		}
		fmt.Println("Database started")

		defer func() {
			fmt.Println("Stopping database")
			cmd := exec.Command("make", "stand-down")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Fatal("Could not shut down database")
			}
			fmt.Println("Database stopped")
		}()
	}

	m.Run()
}
