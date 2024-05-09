package mercadona

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/providers"
)

type Provider struct {
	log logger.Logger
}

const (
	pidProductCode = iota
	pidZoneCode
	pidNFields
)

func New(log logger.Logger) providers.Provider {
	return Provider{log: log}
}

func (p Provider) Name() string {
	return "Mercadona"
}

func (p Provider) FetchPrice(ctx context.Context, pid providers.ProductID) (float32, error) {
	if len(pid) != pidNFields {
		return 0, fmt.Errorf("expected 1 field in product ID, got %d", len(pid))
	}

	url := fmt.Sprintf(
		"https://tienda.mercadona.es/api/products/%s/?lang=es&wh=%s",
		pid[pidProductCode],
		pid[pidZoneCode],
	)

	p.log.Trace("Fetching price from ", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", r.StatusCode)
	}

	content, err := io.ReadAll(r.Body)
	if err != nil {
		return 0, fmt.Errorf("could not read response: %w", err)
	}

	var data struct {
		PriceInstructions struct {
			UnitPrice string `json:"unit_price"`
		} `json:"price_instructions"`
	}

	if err := json.Unmarshal(content, &data); err != nil {
		return 0, fmt.Errorf("could not unmarshal response: %w", err)
	}

	batchPrice, err := strconv.ParseFloat(data.PriceInstructions.UnitPrice, 32)
	if err != nil {
		return 0, fmt.Errorf("could not parse price: %w", err)
	}

	p.log.Tracef("Got price from %s", url)

	return float32(batchPrice), nil
}

func (p Provider) ValidateID(pid providers.ProductID) error {
	if len(pid) != pidNFields {
		return fmt.Errorf("expected %d fields, got %d", pidNFields, len(pid))
	}

	if pid[pidProductCode] == "" {
		return errors.New("product code is empty")
	}

	if pid[pidZoneCode] == "" {
		return errors.New("zone code is empty")
	}

	return nil
}
