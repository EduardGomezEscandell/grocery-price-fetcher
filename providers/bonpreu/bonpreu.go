package bonpreu

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/ubuntu/decorate"
)

var regex = regexp.MustCompile(`<span class="[^"]*price__StyledText[^"]*">([0-9]+,[0-9]{2}).€</span>`)

func Get(ctx context.Context, name string, args ...string) (price float32, err error) {
	defer decorate.OnError(&err, "could not get price for %s", name)

	if len(args) != 1 {
		return 0, fmt.Errorf("expected 1 argument, got %d", len(args))
	}

	url := fmt.Sprintf("https://www.compraonline.bonpreuesclat.cat/products/%s/details", args[0])
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

	return float32(euro) + float32(cent)/100, nil
}
