package bonpreu

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers"
)

var regex = regexp.MustCompile(`<span class="[^"]*price__StyledText[^"]*">([0-9]+,[0-9]{2}).€</span>`)

type Provider struct {
	log logger.Logger
}

func New(log logger.Logger) providers.Provider {
	return Provider{log: log}
}

const (
	pidProductCode = iota
	pidNFields
)

func (p Provider) Name() string {
	return "Bonpreu"
}

func (p Provider) FetchPrice(ctx context.Context, pid providers.ProductID) (float32, error) {
	if len(pid) != pidNFields {
		return 0, fmt.Errorf("expected 1 field in product ID, got %d", len(pid))
	}

	url := fmt.Sprintf(
		"https://www.compraonline.bonpreuesclat.cat/products/%s/details",
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

	content, err := io.ReadAll(r.Body)
	if err != nil {
		return 0, fmt.Errorf("could not read response: %w", err)
	}

	matches := regex.FindAllSubmatch(content, -1)

	switch len(matches) {
	case 0:
		return 0, fmt.Errorf("could not find price in response")
	case 1:
		break
	default:
		return 0, fmt.Errorf("found multiple prices in response")
	}

	var euro uint
	var cent uint

	m := string(matches[0][1])
	_, err = fmt.Sscanf(m, "%d,%d", &euro, &cent)
	if err != nil {
		return 0, fmt.Errorf("could not parse price %q: %w", m, err)
	}

	if cent > 99 {
		return 0, fmt.Errorf("invalid price: %s", m)
	}

	batchPrice := float32(euro) + float32(cent)/100

	p.log.Tracef("Got price from %s", url)

	return batchPrice, nil
}

func (p Provider) ValidateID(pid providers.ProductID) error {
	if len(pid) != pidNFields {
		return fmt.Errorf("expected %d fields, got %d", pidNFields, len(pid))
	}

	if pid[pidProductCode] == "" {
		return fmt.Errorf("product code is empty")
	}

	return nil
}
