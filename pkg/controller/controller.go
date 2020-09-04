package controller

import (
	"context"
	"fmt"

	"github.com/toshi0607/kompal-weather/pkg/analyzer"
	"github.com/toshi0607/kompal-weather/pkg/kompal"
	"github.com/toshi0607/kompal-weather/pkg/logger"
	"github.com/toshi0607/kompal-weather/pkg/monitoring"
	"github.com/toshi0607/kompal-weather/pkg/notifier"
	"github.com/toshi0607/kompal-weather/pkg/storage"
	"golang.org/x/sync/errgroup"
)

type controller struct {
	kompal    kompal.Fetcher
	storage   storage.Storage
	notifiers []notifier.Notifier
	analyzer  analyzer.Analyzer
	monitor   *monitoring.Monitor
	log       logger.Logger
}

// New builds a new Controller
func New(f kompal.Fetcher, s storage.Storage, ns []notifier.Notifier, a analyzer.Analyzer, m *monitoring.Monitor, l logger.Logger) Controller {
	return &controller{
		kompal:    f,
		storage:   s,
		notifiers: ns,
		analyzer:  a,
		monitor:   m,
		log:       l,
	}
}

func (c controller) Watch(ctx context.Context) (*analyzer.Result, error) {
	f, err := c.kompal.Fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch kompal status: %v", err)
	}
	c.log.Info("fetched: %+v", *f)

	st, err := c.storage.Save(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("failed to save status: %v", err)
	}
	c.log.Info("saved: %+v", *st)

	if err := c.monitor.CreatePoint(ctx, st); err != nil {
		c.log.Error("failed to create point", err)
		// Keep processing
	}

	result, err := c.analyzer.Analyze(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze: %v", err)
	}
	c.log.Info("result: %+v", *result)

	eg, ctx := errgroup.WithContext(ctx)
	for _, n := range c.notifiers {
		n := n
		eg.Go(func() error {
			return n.Notify(ctx, result)
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("failed to notify: %v", err)
	}

	return result, nil
}
