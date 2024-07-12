package session

type Settings struct {
	Enable bool `yaml:"enable"`
}

func (Settings) Defaults() Settings {
	return Settings{
		Enable: true,
	}
}
