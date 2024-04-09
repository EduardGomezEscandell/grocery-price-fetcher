package mercadona

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/provider"
	"github.com/sirupsen/logrus"
)

type Provider struct {
	batchSize float32
	zoneCode  string
	id        string
}

func New() provider.Provider {
	return &Provider{}
}

func (p *Provider) UnmarshalTSV(cols ...string) error {
	if len(cols) != 3 {
		return fmt.Errorf("expected 3 arguments (batch_size, zone_code, id), got %d", len(cols))
	}

	c, err := strconv.ParseFloat(cols[0], 32)
	if err != nil {
		return fmt.Errorf("could not parse batch_size (%s): %w", cols[0], err)
	}
	if c <= 0 {
		return fmt.Errorf("invalid batch_size: %f", c)
	}

	p.batchSize = float32(c)
	p.id = cols[1]
	p.zoneCode = cols[2]

	return nil
}

func (p *Provider) UnmarshalMap(argv map[string]string) (err error) {
	if len(argv) != 3 {
		return fmt.Errorf("expected 3 arguments (batch_size, zone_code, id), got %d", len(argv))
	}

	bs, ok := argv["batch_size"]
	if !ok {
		return fmt.Errorf("missing batch_size")
	}

	p.id, ok = argv["id"]
	if !ok {
		return fmt.Errorf("missing id")
	}

	p.zoneCode, ok = argv["zone_code"]
	if !ok {
		return fmt.Errorf("missing zone_code")
	}

	batchSize, err := strconv.ParseFloat(bs, 32)
	if err != nil {
		return fmt.Errorf("could not parse batch_size (%s): %w", argv["batch_size"], err)
	}
	if batchSize <= 0 {
		return fmt.Errorf("invalid batch_size: %f", batchSize)
	}
	p.batchSize = float32(batchSize)

	return nil
}

func (p *Provider) FetchPrice(ctx context.Context) (float32, error) {
	url := fmt.Sprintf("https://tienda.mercadona.es/api/products/%s/?lang=es&wh=%s", p.id, p.zoneCode)
	logrus.Trace("fetching price from ", url)

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

	return float32(batchPrice) / p.batchSize, nil
}
