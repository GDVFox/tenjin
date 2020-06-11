package database

import "github.com/go-sql-driver/mysql"

// Config contains database connection settings
type Config struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

// DSN create dsn string from config
func (c *Config) DSN() string {
	mycfg := mysql.NewConfig()
	mycfg.User = c.User
	mycfg.Passwd = c.Password
	mycfg.Net = "tcp"
	mycfg.Addr = c.Host
	mycfg.DBName = c.Database
	mycfg.ParseTime = true
	return mycfg.FormatDSN()
}
