package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/hawry/gomocka/settings"

	"comail.io/go/colog"

	"gopkg.in/alecthomas/kingpin.v2"
)

var buildVersion string

var (
	version        = kingpin.Flag("version", "print version of gock").Default("false").Bool()
	config         = kingpin.Flag("config", "configuration file to create endpoints from").Short('c').Default("settings.json").String()
	verbose        = kingpin.Flag("verbose", "enabled verbose logging. if --silent is used, --verbose will be ignore").Short('v').Default("false").Bool()
	silent         = kingpin.Flag("silent", "disabled all output except for errors. overrides --verbose if set").Short('s').Default("false").Bool()
	generateConfig = kingpin.Flag("generate", "generate a sample configuration").Short('g').Default("false").Bool()
)

func main() {
	kingpin.Parse()

	if *version {
		fmt.Printf("gock version %s\n", buildVersion)
		return
	}

	if *generateConfig {
		if _, err := settings.CreateDefault(); err != nil {
			log.Fatalf("error: %v", err)
		}
		fmt.Printf("example configuration created\n")
		return
	}

	initLogging()

	f, err := os.Open(*config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	settings, err := settings.New(f)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	r := mux.NewRouter()
	setupHandlers(settings, r)

	srv := &http.Server{Addr: fmt.Sprintf(":%d", settings.Port()), Handler: r}
	srv.RegisterOnShutdown(func() {
		log.Printf("info: shutting down server")
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for range c {
			srv.Shutdown(context.Background())
		}
	}()
	log.Printf("info: starting server at %s", srv.Addr)
	err = srv.ListenAndServe()
	if err.Error() != "http: Server closed" {
		log.Fatalf("error: %v", err)
	}
}

func setupHandlers(s *settings.Settings, r *mux.Router) {
	for _, m := range s.Mocks {
		f := createHandler(m.ResponseCode, m.ResponseBody, m.Headers, m)
		r.HandleFunc(m.Path, f).Methods(m.Method)
		log.Printf("debug: creating mock for %s %s - returning %d, %s using %v", m.Method, m.Path, m.ResponseCode, m.ResponseBody, f)
	}
	log.Printf("info: registered %d paths", len(s.Mocks))
}

func initLogging() {
	colog.Register()
	colog.SetDefaultLevel(colog.LDebug)
	if *silent {
		colog.SetMinLevel(colog.LError)
		return
	}

	if *verbose {
		colog.SetMinLevel(colog.LDebug)
	} else {
		colog.SetMinLevel(colog.LInfo)
	}
	log.Printf("debug: enable verbose logging")
}
