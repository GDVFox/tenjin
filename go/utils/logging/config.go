package logging

// Config contains logger configuration
type Config struct {
	File  string `yaml:"file"`
	Level string `yaml:"level"`
}
