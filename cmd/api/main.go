package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/anuj0x16/email-dispatcher/internal/dispatcher"
)

type config struct {
	port         int
	jobQueueSize int
	nworkers     int
}

type application struct {
	config     config
	logger     *slog.Logger
	dispatcher *dispatcher.Dispatcher
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.IntVar(&cfg.jobQueueSize, "job-queue-size", 100, "Job queue size")
	flag.IntVar(&cfg.nworkers, "nworkers", 5, "Number of workers to start")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	dispatcher := dispatcher.New(cfg.jobQueueSize, cfg.nworkers)
	dispatcher.Start()

	app := &application{
		config:     cfg,
		logger:     logger,
		dispatcher: dispatcher,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
