package mercadona

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers"
)

type Provider struct {
	log logger.Logger
}

const (
	pidProductCode = iota
	pidZoneCode
	pidNFields
)

func New(log logger.Logger) Provider {
	return Provider{log: log}
}

func (p Provider) Name() string {
	return "Mercadona"
}

func (p Provider) FetchPrice(ctx context.Context, pid providers.ProductCode) (float32, error) {
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

func (p Provider) ValidateCode(code providers.ProductCode) error {
	if code[pidProductCode] == "" {
		return errors.New("product code (ID 0) should not be empty")
	}

	if code[pidZoneCode] == "" {
		return errors.New("zone code (ID 1) should not be empty")
	}

	if code[2] != "" {
		return fmt.Errorf("unexpected field at index 2: %q", code[2])
	}

	return nil
}
