package repos

import (
	"os"

	"gopkg.in/yaml.v2"
)

type ActionConfig struct {
	URL          string   `yaml:"url"`
	ContactEmail string   `yaml:"contactEmail"`
	UseCase      string   `yaml:"useCase"`
	Repos        []string `yaml:"repos,omitempty"`
}

// LoadConfig load yaml file from path
// unmarshal into Config struct
// return Config struct
func LoadConfig(path string) (config ActionConfig, err error) {
	source, err := os.ReadFile(path)

	if err != nil {
		return config, err
	}

	if err := yaml.Unmarshal(source, &config); err != nil {
		return config, err
	}

	return config, err
}
