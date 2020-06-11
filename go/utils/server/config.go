package server

import "strconv"

// Config represents server cofig
type Config struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// Address returns server address
func (c *Config) Address() string {
	return c.Host + ":" + strconv.Itoa(c.Port)
}
