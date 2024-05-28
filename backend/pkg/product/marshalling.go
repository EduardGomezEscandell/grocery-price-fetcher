package product

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers"
	"github.com/ubuntu/decorate"
)

func (p *Product) UnmarshalTSV(args []string) (err error) {
	defer decorate.OnError(&err, "could not unmarshal product %q", args)

	const (
		flieldName int = iota
		fieldProvider
		fieldBatchSize
		fieldArgv
		nMinimumArgs
	)

	if len(args) < nMinimumArgs {
		return fmt.Errorf("expected at least 2 fields, got %d", len(args))
	}

	if name := args[flieldName]; name == "" {
		p.Name = "COULD NOT PARSE"
		return fmt.Errorf("empty name")
	} else {
		p.Name = name
	}

	if prov, ok := providers.Lookup(args[fieldProvider]); ok {
		p.Provider = prov
	} else {
		return fmt.Errorf("could not find provider %q", args[fieldProvider])
	}

	if s, err := strconv.ParseFloat(args[fieldBatchSize], 32); err != nil {
		return fmt.Errorf("could not parse batch size: %w", err)
	} else {
		p.BatchSize = float32(s)
	}

	if pid := args[fieldArgv:]; p.Provider.ValidateID(pid) != nil {
		return fmt.Errorf("invalid product ID %q", pid)
	} else {
		p.ProductID = pid
	}

	return nil
}

type jsonHelper struct {
	Name      string   `json:"name"`
	BatchSize float32  `json:"batch_size"`
	Price     float32  `json:"price"`
	Provider  string   `json:"provider"`
	ProductID []string `json:"product_id"`
}

func (p *Product) UnmarshalJSON(b []byte) (err error) {
	var helper jsonHelper
	if err := json.Unmarshal(b, &helper); err != nil {
		return fmt.Errorf("could not unmarshal product: %w", err)
	}

	defer decorate.OnError(&err, "could not unmarshal product %+v", p)
	p.Name = helper.Name
	p.BatchSize = helper.BatchSize
	p.Price = helper.Price

	if prov, ok := providers.Lookup(helper.Provider); ok {
		p.Provider = prov
	} else {
		return fmt.Errorf("could not find provider %q", helper.Provider)
	}

	if p.Provider.ValidateID(helper.ProductID) != nil {
		return fmt.Errorf("invalid product ID %q", helper.ProductID)
	} else {
		p.ProductID = helper.ProductID
	}

	return nil
}

func (p *Product) MarshalJSON() (b []byte, err error) {
	providerName := "NoProvider"
	if p.Provider != nil {
		providerName = p.Provider.Name()
	}

	helper := jsonHelper{
		Name:      p.Name,
		BatchSize: p.BatchSize,
		Price:     p.Price,
		Provider:  providerName,
		ProductID: p.ProductID,
	}

	return json.Marshal(helper)
}
