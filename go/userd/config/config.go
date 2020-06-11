package config

import (
	"os"

	"github.com/GDVFox/tenjin/utils/database"
	"github.com/GDVFox/tenjin/utils/logging"
	"github.com/GDVFox/tenjin/utils/server"
	"gopkg.in/yaml.v2"
)

// Config contains daemon configuration
type Config struct {
	LoggingConfig  *logging.Config  `yaml:"logging"`
	DatabaseConfig *database.Config `yaml:"database"`
	HTTPConfig     *server.Config   `yaml:"http"`
}

// LoadConfig creates new config from file
func LoadConfig(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := &Config{}
	dec := yaml.NewDecoder(f)
	err = dec.Decode(cfg)
	return cfg, err
}
