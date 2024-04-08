package mercadona

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/ubuntu/decorate"
)

func Get(ctx context.Context, name string, args ...string) (price float32, err error) {
	defer decorate.OnError(&err, "could not get price for %s", name)

	if len(args) != 1 {
		return 0, fmt.Errorf("expected 1 argument, got %d", len(args))
	}

	url := fmt.Sprintf("https://tienda.mercadona.es/api/products/%s/?lang=es&wh=bcn1", args[0])
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer r.Body.Close()

	content, err := io.ReadAll(r.Body)
	if err != nil {
		return 0, fmt.Errorf("could not read response: %v", err)
	}

	var data struct {
		PriceInstructions struct {
			UnitPrice string `json:"unit_price"`
		} `json:"price_instructions"`
	}

	if err := json.Unmarshal(content, &data); err != nil {
		return 0, fmt.Errorf("could not unmarshal response: %v", err)
	}

	p, err := strconv.ParseFloat(data.PriceInstructions.UnitPrice, 32)
	if err != nil {
		return 0, fmt.Errorf("could not parse price: %v", err)
	}

	return float32(p), nil
}
