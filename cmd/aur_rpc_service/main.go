package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/moson-mo/goaurrpc/internal/config"
	"github.com/moson-mo/goaurrpc/internal/rpc"
)

/*
	version can be overridden during build:
	go build -ldflags="-X 'main.version=v1.0.0'"
*/

var version = "v1.2.2"

func main() {
	veryVerboseLog := flag.Bool("v", false, "Very verbose logging")
	flag.Parse()

	var settings config.Settings
	err := config.LoadFromEnv(&settings)
	if err != nil {
		panic("Error loading configuration: " + err.Error())
	}

	// prep logging
	if settings.LogFile != "" {
		f, err := os.OpenFile(settings.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(fmt.Errorf("error opening log file: %s", err))
		}
		log.SetOutput(f)
	}
	log.SetFlags(log.LstdFlags | log.Lmsgprefix)
	log.SetPrefix("| ")
	slog.SetLogLoggerLevel(settings.GetLogLevel())

	// construct a new server and start listening for requests
	slog.Info("Starting aur rpc service", "version", version, "port", settings.Port)
	s, err := rpc.New(settings, *veryVerboseLog, version)
	if err != nil {
		panic(fmt.Errorf("error setting up rpc server: %s", err))
	}
	if err = s.Listen(); !errors.Is(err, http.ErrServerClosed) {
		slog.Error("The rpc service stopped unexpectedly", "error", err)
	}
	slog.Info("Stopped aur rpc service", "version", version)
}
