package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/toshi0607/kompal-weather/internal/config"
	"github.com/toshi0607/kompal-weather/pkg/gcs"
	"github.com/toshi0607/kompal-weather/pkg/http"
	"github.com/toshi0607/kompal-weather/pkg/logger"
	"github.com/toshi0607/kompal-weather/pkg/secret"
	"github.com/toshi0607/kompal-weather/pkg/visualizer"
	"golang.org/x/sync/errgroup"
)

const (
	exitOK    = 0
	exitError = 1
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("recovered, error: %v", err)
			os.Exit(1)
		}
	}()
	os.Exit(realMain(os.Args))
}

func realMain(_ []string) int {
	ctx := context.Background()

	// init secret manager client
	s, err := secret.New()
	if err != nil {
		log.Printf("failed to create new secret: %v", err)
		return exitError
	}
	defer func() {
		if err := s.Close(); err != nil {
			fmt.Print(err)
		}
	}()

	// Init config
	c := config.NewVisualizer(s)
	if err := c.Init(); err != nil {
		log.Printf("failed to create new config: %v", err)
		return exitError
	}

	// Init logger
	var l logger.Logger
	if c.IsLocal() {
		l = logger.NewLog()
	} else {
		var err error
		l, err = logger.NewCloudLogging(ctx, c.GCPProjectID, c.ServiceName, c.Version, c.Environment)
		if err != nil {
			log.Printf("failed to create new logging: %v", err)
			return exitError
		}
		defer func() {
			if err := l.Close(); err != nil {
				l.Error("failed to close logging: %s", err)
			}
		}()
	}

	// Init GCS
	g, err := gcs.New(c.BucketName)
	if err != nil {
		log.Printf("failed to create new gcs: %v", err)
		return exitError
	}

	// Init visualizer
	v, err := visualizer.New(c, g, l)
	if err != nil {
		l.Error("failed to create new visualizer", err)
		return exitError
	}

	// Server start
	server := http.NewVisualizer(v, l)

	httpLn, err := net.Listen("tcp", fmt.Sprintf(":%d", c.ServerPort))
	if err != nil {
		l.Error("failed to listen port", err)
		return exitError
	}
	l.Info("http server listening, port: %d", c.ServerPort)

	wg, ctx := errgroup.WithContext(ctx)
	wg.Go(func() error { return server.Serve(httpLn) })

	// Signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)
	select {
	case <-sigCh:
		l.Info("received SIGTERM, exiting server gracefully")
	case <-ctx.Done():
		l.Info("context done")
	}

	// Graceful shutdown
	if err := server.GracefulStop(ctx); err != nil {
		l.Error("failed to stop server gracefully", err)
	} else {
		l.Info("succeeded to stop server gracefully")
	}

	if err := wg.Wait(); err != nil {
		l.Error("server failed", err)
		return exitError
	}

	return exitOK
}
