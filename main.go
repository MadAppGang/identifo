package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/madappgang/identifo/v2/config"
	"github.com/madappgang/identifo/v2/model"
)

func main() {
	configFlag := flag.String("config", "", "The location of a server configuration file (local file, s3 or etcd)")
	flag.Parse()

	done := make(chan bool)
	restart := make(chan bool)
	defer close(done)
	defer close(restart)

	srv, httpSrv, err := initServer(*configFlag, restart)
	if err != nil {
		log.Fatalf("Unable to start Identifo with error: %v ", err)
	}
	go startHTTPServer(httpSrv)
	log.Printf("Started the server on port: %s", srv.Settings().GetPort())
	log.Printf("You can open admin panel: %s/adminpanel or http://localhost:%s/adminpanel", srv.Settings().General.Host, srv.Settings().GetPort())

	watcher, err := config.NewConfigWatcher(srv.Settings().Config)
	if err != nil {
		log.Fatalf("Unable to start Identifo with error: %v ", err)
	}
	go func() {
		time.Sleep(time.Second)
		watcher.Watch()
	}()

	osch := make(chan os.Signal, 1)
	signal.Notify(osch, syscall.SIGINT, syscall.SIGTERM)

	restartServer := func() {
		ctx, _ := context.WithTimeout(context.Background(), time.Minute*3)
		httpSrv.Shutdown(ctx)
		srv.Close() // TODO: implement gracefull server shutdown
		srv, httpSrv, err = initServer(*configFlag, restart)
		if err != nil {
			log.Fatalf("Unable to start Identifo with error: %v ", err)
		}
		go startHTTPServer(httpSrv)
		log.Printf("Started the server on port: %s", srv.Settings().GetPort())
		log.Printf("You can open admin panel: %s/adminpanel or http://localhost:%s/adminpanel", srv.Settings().General.Host, srv.Settings().GetPort())
		log.Println("Server successfully restarted with new settings ...")
	}

	for {
		select {
		case <-watcher.WatchChan():
			log.Println("Config file has been changed, restarting ...")
			restartServer()

		case <-restart:
			log.Println("Restart signal have been received, restarting ...")
			restartServer()

		case err := <-watcher.ErrorChan():
			log.Printf("Getting error from config watcher: %v", err)

		case <-osch:
			log.Println("Received termination signal, shutting down the server ⤵️...")
			ctx, _ := context.WithTimeout(context.Background(), time.Minute*3)
			httpSrv.Shutdown(ctx)
			srv.Close()
			log.Println("The server is down, good bye 👋👋👋.")
			return
		}
	}
}

func initServer(flag string, restartChan chan<- bool) (model.Server, *http.Server, error) {
	srv, err := config.NewServerFromFlag(flag, restartChan)
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to start Identifo with error: %v ", err)
	}
	httpSrv := &http.Server{
		Addr:    srv.Settings().GetPort(),
		Handler: srv.Router(),
	}
	return srv, httpSrv, nil
}

func startHTTPServer(server *http.Server) {
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe() error: %v", err)
	}
}
