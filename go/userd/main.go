package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/GDVFox/tenjin/userd/config"
	"github.com/GDVFox/tenjin/utils/database"
	"github.com/GDVFox/tenjin/utils/logging"
	"github.com/GDVFox/tenjin/utils/server"
)

var (
	configFile string
	logger     *logging.Logger
)

func init() {
	flag.StringVar(&configFile, "config", "", "Filename of config")
}

func main() {
	flag.Parse()
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("can not load config: %s", err)
	}

	logger, err = logging.NewFileLogger(cfg.LoggingConfig)
	if err != nil {
		log.Fatalf("can not create logger: %s", err)
	}

	if err := database.Open(cfg.DatabaseConfig); err != nil {
		log.Fatalf("can not open database connection: %s", err)
	}

	cancelCh := make(chan struct{})
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(signals)

	go func() {
		defer close(cancelCh)
		select {
		case s := <-signals:
			logger.Infof("got signal: %s", s.String())
		}
	}()

	s := server.NewServer(cfg.HTTPConfig, routes(), logger)
	if err = s.StartWithCancel(cancelCh); err != nil {
		log.Fatalf("can not start http server: %s", err)
	}

	logger.Info("done!")
}
