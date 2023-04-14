package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/miniclip/gonsul/app"
	"github.com/miniclip/gonsul/internal/config"
	"github.com/miniclip/gonsul/internal/exporter"
	"github.com/miniclip/gonsul/internal/importer"
	"github.com/miniclip/gonsul/internal/util"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			var recoveredError = r.(util.GonsulError)
			os.Exit(recoveredError.Code)
		}
	}()

	start()
}

func start() {
	// Build our configuration
	cfg, err := config.NewConfig()
	if err != nil {
		util.ExitError(err, util.ErrorBadParams, util.NewLogger(0))
	}

	// Build our logger
	logger := util.NewLogger(cfg.GetLogLevel())

	// Are we just printing the app version
	if cfg.IsShowVersion() {
		fmt.Println("Gonsul version: " + app.Version)
		fmt.Println("Build date: " + app.BuildDate)
		return
	}

	var certificate tls.Certificate
	var caCertPool *x509.CertPool
	if len(cfg.GetKeyFile()) != 0 && len(cfg.GetCaFile()) != 0 && len(cfg.GetCertFile()) != 0 {
		cert, err := ioutil.ReadFile(cfg.GetCaFile())
		if err != nil {
			log.Fatalf("could not open certificate file: %v", err)
		}
		caCertPool = x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(cert)

		certificate, err = tls.LoadX509KeyPair(cfg.GetCertFile(), cfg.GetKeyFile())
		if err != nil {
			log.Fatalf("could not load certificate: %v", err)
		}
	}

	// Build all dependencies for our application
	hookHttpServer := app.NewHookHttp(cfg, logger)
	httpClient := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:      caCertPool,
			Certificates: []tls.Certificate{certificate},
		},
	}, Timeout: time.Second * time.Duration(cfg.GetTimeout())}

	exp := exporter.NewExporter(cfg, logger)
	imp := importer.NewImporter(cfg, logger, httpClient)
	sigChannel := make(chan os.Signal)
	// Build our Applications
	once := app.NewOnce(cfg, logger, exp, imp)
	hook := app.NewHook(hookHttpServer, cfg, logger, once)
	poll := app.NewPoll(cfg, logger, once, 0)
	// Build our main Application container
	application := app.NewApplication(cfg, once, hook, poll, sigChannel)

	// Start our application
	application.Start()

	// We're still here, all went well, good bye
	logger.PrintInfo("Quitting... bye.")
}
