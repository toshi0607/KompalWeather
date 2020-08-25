package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/toshi0607/kompal-weather/pkg/analyzer"
	"github.com/toshi0607/kompal-weather/pkg/config"
	"github.com/toshi0607/kompal-weather/pkg/http"
	"github.com/toshi0607/kompal-weather/pkg/kompal"
	"github.com/toshi0607/kompal-weather/pkg/notifier"
	"github.com/toshi0607/kompal-weather/pkg/secret"
	"github.com/toshi0607/kompal-weather/pkg/slack"
	"github.com/toshi0607/kompal-weather/pkg/storage"
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

	// Init kompal
	k := kompal.New(c.Kompal)

	// Init storage
	sheets, err := storage.NewSheets(c.Sheets)
	if err != nil {
		fmt.Print(err)
		return exitError
	}

	// Init notifiers
	// Init twitter

	// Init slack
	slack := slack.New(c.Slack)

	// Init analyzer
	analyzer := analyzer.New(sheets)

	// Server start
	server := http.New(k, sheets, []notifier.Notifier{slack}, analyzer)

	httpLn, err := net.Listen("tcp", fmt.Sprintf(":%d", c.ServerPort))
	if err != nil {
		fmt.Print(err)
		return exitError
	}
	fmt.Print("http server listening")

	wg, ctx := errgroup.WithContext(ctx)
	wg.Go(func() error { return server.Serve(httpLn) })

	// Signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)
	select {
	case <-sigCh:
		fmt.Print("received SIGTERM, exiting server gracefully")
	case <-ctx.Done():
	}

	// Graceful shutdown
	if err := server.GracefulStop(ctx); err != nil {
		fmt.Sprint()
	} else {
		fmt.Sprint()
	}

	if err := wg.Wait(); err != nil {
		return exitError
	}

	return exitOK
}
