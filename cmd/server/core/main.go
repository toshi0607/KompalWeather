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
	"github.com/toshi0607/kompal-weather/pkg/analyzer"
	"github.com/toshi0607/kompal-weather/pkg/controller"
	"github.com/toshi0607/kompal-weather/pkg/http"
	"github.com/toshi0607/kompal-weather/pkg/kompal"
	"github.com/toshi0607/kompal-weather/pkg/logger"
	"github.com/toshi0607/kompal-weather/pkg/monitoring"
	"github.com/toshi0607/kompal-weather/pkg/notifier"
	"github.com/toshi0607/kompal-weather/pkg/secret"
	"github.com/toshi0607/kompal-weather/pkg/slack"
	"github.com/toshi0607/kompal-weather/pkg/storage"
	"github.com/toshi0607/kompal-weather/pkg/twitter"
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
	c := config.NewCore(s)
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

	// Init kompal
	k := kompal.New(c.Kompal)

	// Init storage
	sheets, err := storage.NewSheets(c.Sheets)
	if err != nil {
		l.Error("failed to init sheets", err)
		return exitError
	}

	// Init notifiers
	// Init twitter
	tw := twitter.New(c.Twitter, l)

	// Init slack
	sl := slack.New(c.Slack, l)

	// Init analyzer
	an := analyzer.New(sheets)

	// Init monitor
	m, err := monitoring.New(c.GCPProjectID)
	if err != nil {
		l.Error("failed to init monitoring", err)
		return exitError
	}
	defer func() {
		if err := l.Close(); err != nil {
			l.Error("failed to close monitoring: %s", err)
		}
	}()

	// CoreServer start
	server := http.NewCore(controller.New(k, sheets, []notifier.Notifier{sl, tw}, an, m, l), l)

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
