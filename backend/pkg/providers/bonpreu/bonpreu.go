package bonpreu

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers"
	"golang.org/x/exp/maps"
)

type Provider struct {
	log logger.Logger
}

func New(log logger.Logger) Provider {
	return Provider{log: log}
}

const (
	pidProductCode = iota
)

func (p Provider) Name() string {
	return "Bonpreu"
}

func (p Provider) FetchPrice(ctx context.Context, pid providers.ProductCode) (float32, error) {
	url := fmt.Sprintf(
		"https://www.compraonline.bonpreuesclat.cat/api/v5/products/search?&term=%s",
		pid[pidProductCode],
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

	var data struct {
		Entities struct {
			Product map[string]struct {
				RetailerProductID string `json:"retailerProductId"`
				Price             struct {
					Current struct {
						Amount string `json:"amount"`
					} `json:"current"`
				} `json:"price"`
			} `json:"product"`
		} `json:"entities"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return 0, fmt.Errorf("could not decode response: %v", err)
	}

	if len(data.Entities.Product) == 0 {
		return 0, fmt.Errorf("product not found")
	}

	v := maps.Values(data.Entities.Product)[0]
	if v.RetailerProductID != pid[pidProductCode] {
		return 0, fmt.Errorf("product not found")
	}

	batchPrice, err := strconv.ParseFloat(v.Price.Current.Amount, 32)
	if err != nil {
		return 0, fmt.Errorf("could not parse price: %v", err)
	}

	ret := float32(math.Round(batchPrice*100) / 100)
	p.log.Tracef("Got price from %s: %g", url, ret)
	return ret, nil
}

func (p Provider) ValidateCode(code providers.ProductCode) error {
	if code[pidProductCode] == "" {
		return fmt.Errorf("product code (index 0) is empty")
	}

	if code[1] != "" {
		return fmt.Errorf("unexpected field at index 1: %q", code[2])
	}

	if code[2] != "" {
		return fmt.Errorf("unexpected field at index 2: %q", code[2])
	}

	return nil
}
