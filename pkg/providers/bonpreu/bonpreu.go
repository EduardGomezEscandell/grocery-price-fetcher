package bonpreu

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/providers"
)

var regex = regexp.MustCompile(`<span class="[^"]*price__StyledText[^"]*">([0-9]+,[0-9]{2}).€</span>`)

type Provider struct {
	id        string
	batchSize float32
}

func New() providers.Provider {
	return &Provider{}
}

func (p *Provider) UnmarshalTSV(cols ...string) error {
	if len(cols) != 2 {
		return fmt.Errorf("expected 2 arguments (batch_size, id), got %d", len(cols))
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
	return nil
}

func (p *Provider) UnmarshalMap(argv map[string]string) (err error) {
	if len(argv) != 2 {
		return fmt.Errorf("expected 2 arguments (batch_size, id), got %d", len(argv))
	}

	bs, ok := argv["batch_size"]
	if !ok {
		return fmt.Errorf("missing batch_size")
	}

	p.id, ok = argv["id"]
	if !ok {
		return fmt.Errorf("missing id")
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

func (p *Provider) FetchPrice(ctx context.Context, log logger.Logger) (float32, error) {
	url := fmt.Sprintf("https://www.compraonline.bonpreuesclat.cat/products/%s/details", p.id)
	log.Trace("Fetching price from ", url)

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

	return batchPrice / p.batchSize, nil
}

func (p *Provider) MarshalJSON() ([]byte, error) {
	helper := struct {
		Bonpreu map[string]string
	}{
		Bonpreu: map[string]string{
			"batch_size": fmt.Sprint(p.batchSize),
			"id":         p.id,
		},
	}

	return json.Marshal(helper)
}
