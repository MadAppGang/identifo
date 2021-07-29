package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/madappgang/identifo/config"
	"github.com/madappgang/identifo/model"
)

func main() {
	configFlag := flag.String("config", "", "The location of a server configuration file (local file, s3 or etcd)")
	flag.Parse()

	// ignore error to fall back to default if needed
	settings, fileErr := model.ConfigStorageSettingsFromString(*configFlag)
	configStorage, err := config.InitConfigurationStorage(settings)
	if err != nil || fileErr != nil || *configFlag == "" {
		log.Printf("Unable to init config using\n\tconfig string: %s\n\twith error: %v\n",
			*configFlag,
			err,
		)
		// Trying to fall back to default settings:
		log.Printf("Trying to load default settings from env variable 'SERVER_CONFIG_PATH' or default pathes.\n")
		configStorage, err = config.DefaultStorage()
		if err != nil {
			log.Fatalf("Unable to load default config with error: %v", err)
		}
	}

	srv, err := config.NewServer(configStorage)
	if err != nil {
		log.Fatalf("error creating server: %v", err)
	}

	httpSrv := &http.Server{
		Addr:    srv.Settings().GetPort(),
		Handler: srv.Router(),
	}

	go startHTTPServer(httpSrv)

	log.Printf("Started the server on host: %s", srv.Settings().General.Host)
	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
	log.Println("shutting down the service ⤵️")

	// Stop the service gracefully.
	// TODO: implement gracefull server shutdown
	// srv.Shutdown()
	httpSrv.Shutdown(context.Background())
	log.Println("the server is gracefully stopped, bye 👋")
}

func startHTTPServer(httpSrv *http.Server) {
	if err := httpSrv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %s", err)
	}
}
