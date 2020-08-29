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
	"github.com/toshi0607/kompal-weather/pkg/http"
	"github.com/toshi0607/kompal-weather/pkg/kompal"
	"github.com/toshi0607/kompal-weather/pkg/logger"
	"github.com/toshi0607/kompal-weather/pkg/notifier"
	"github.com/toshi0607/kompal-weather/pkg/secret"
	"github.com/toshi0607/kompal-weather/pkg/slack"
	"github.com/toshi0607/kompal-weather/pkg/storage"
	"github.com/toshi0607/kompal-weather/pkg/twitter"
	"golang.org/x/sync/errgroup"
)

const (
	exitOK = iota
	exitError
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "Error:\n%s\n", err)
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
		fmt.Print(err)
		return exitError
	}

	// Init config
	c := config.New(s)
	if err := c.Init(); err != nil {
		fmt.Print(err)
		return exitError
	}

	// Init logger
	l, err := logger.New(ctx, c.GCPProjectID, c.ServiceName, c.Version)
	if err != nil {
		fmt.Print(err)
		return exitError
	}
	defer func() {
		if err := l.Close(); err != nil {
			log.Printf("failed to close: %s", err)
		}
	}()

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
	twitter := twitter.New(c.Twitter, l)

	// Init slack
	slack := slack.New(c.Slack, l)

	// Init analyzer
	analyzer := analyzer.New(sheets)

	// Server start
	server := http.New(k, sheets, []notifier.Notifier{slack, twitter}, analyzer, l)

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
		l.Info("contest done")
	}

	// Graceful shutdown
	if err := server.GracefulStop(ctx); err != nil {
		l.Info("failed to stop server gracefully")
	} else {
		l.Info("succeeded to stop server gracefully")
	}

	if err := wg.Wait(); err != nil {
		l.Error("server failed", err)
		return exitError
	}

	return exitOK
}
