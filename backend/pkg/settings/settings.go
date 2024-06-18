package settings

import (
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/daemon"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services"
	"gopkg.in/yaml.v3"
)

type Settings struct {
	Verbosity int
	Daemon    daemon.Settings
	Services  services.Settings
}

func Defaults() Settings {
	return Settings{
		Verbosity: 1,
		Daemon:    daemon.Settings{}.Defaults(),
		Services:  services.Settings{}.Defaults(),
	}
}

func (s Settings) String() string {
	out, err := yaml.Marshal(s)
	if err != nil {
		panic(err.Error())
	}

	return string(out)
}
