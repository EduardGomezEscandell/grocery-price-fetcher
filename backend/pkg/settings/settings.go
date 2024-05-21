package settings

import (
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services"
	"gopkg.in/yaml.v3"
)

type Settings struct {
	Verbosity int
	FrontEnd  string
	Address   string

	Services services.Settings
}

func Defaults() Settings {
	return Settings{
		Verbosity: 1,
		FrontEnd:  "/usr/share/grocery-price-fetcher/frontend",
		Address:   "http://localhost:3000",

		Services: services.Settings{}.Defaults(),
	}
}

func (s Settings) String() string {
	out, err := yaml.Marshal(s)
	if err != nil {
		panic(err.Error())
	}

	return string(out)
}
