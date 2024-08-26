package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/madappgang/identifo/v2/config"
	"github.com/madappgang/identifo/v2/localization"
	"github.com/madappgang/identifo/v2/logging"
	"github.com/madappgang/identifo/v2/model"
)

func main() {
	// load default translations
	localization.LoadDefaultCatalog()

	// load config
	configFlag := flag.String("config", "", "The location of a server configuration file (local file, s3 or etcd)")
	flag.Parse()

	done := make(chan bool)
	restart := make(chan bool)
	defer close(done)
	defer close(restart)

	defaultLogger := logging.NewDefaultLogger()

	srv, httpSrv, err := initServer(defaultLogger, *configFlag, restart)
	if err != nil {
		defaultLogger.Error("Unable to start Identifo with error", logging.FieldError, err)
		os.Exit(1)
	}

	settings := srv.Settings()

	logSettings := settings.Logger
	logger := logging.NewLogger(logSettings.Format, logSettings.Common.Level).
		With(logging.FieldComponent, logging.ComponentCommon)

	logging.DefaultLogger = logger

	go startHTTPServer(logger, httpSrv)

	logger.Info("Started the server", "port", srv.Settings().GetPort())
	logger.Info("You can open admin panel",
		"host", fmt.Sprintf("%s/adminpanel", settings.General.Host),
		"url", fmt.Sprintf("http://localhost:%s/adminpanel", settings.GetPort()))

	watcher, err := config.NewConfigWatcher(logger, srv.Settings().Config)
	if err != nil {
		logger.Error("Unable to start Identifo", logging.FieldError, err)
		os.Exit(3)
	}
	watcher.Watch()

	osch := make(chan os.Signal, 1)
	signal.Notify(osch, syscall.SIGINT, syscall.SIGTERM)

	restartServer := func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
		httpSrv.Shutdown(closeCtx)
		cancel()
		srv.Close() // TODO: implement graceful server shutdown

		srv, httpSrv, err = initServer(logger, *configFlag, restart)
		if err != nil {
			logger.Error("Unable to restart Identifo", logging.FieldError, err)
			os.Exit(4)
		}
		go startHTTPServer(logger, httpSrv)

		logger.Info("Started the server", "port", srv.Settings().GetPort())
		logger.Info("You can open admin panel",
			"host", fmt.Sprintf("%s/adminpanel", settings.General.Host),
			"url", fmt.Sprintf("http://localhost:%s/adminpanel", settings.GetPort()))
		logger.Info("Server successfully restarted with new settings ...")
	}

	for {
		select {
		case <-watcher.WatchChan():
			logger.Info("Config file has been changed, restarting ...")
			restartServer()

		case <-restart:
			logger.Info("Restart signal have been received, restarting ...")
			restartServer()

		case err := <-watcher.ErrorChan():
			logger.Error("Getting error from config watcher", logging.FieldError, err)

		case <-osch:
			logger.Info("Received termination signal, shutting down the server ⤵️...")
			closeCtx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
			httpSrv.Shutdown(closeCtx)
			cancel()
			srv.Close()
			logger.Info("The server is down, good bye 👋👋👋.")
			return
		}
	}
}

func initServer(logger *slog.Logger, flag string, restartChan chan<- bool) (model.Server, *http.Server, error) {
	srv, err := config.NewServerFromFlag(logger, flag, restartChan)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to start Identifo with error: %w", err)
	}

	httpSrv := &http.Server{
		Addr:    srv.Settings().GetPort(),
		Handler: srv.Router(),
	}

	return srv, httpSrv, nil
}

func startHTTPServer(
	logger *slog.Logger,
	server *http.Server,
) {
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Error("ListenAndServe()", logging.FieldError, err)
		os.Exit(2)
	}

	logger.Info("Server stopped")
}
