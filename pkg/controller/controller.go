package controller

import (
	"context"

	"github.com/toshi0607/kompal-weather/pkg/analyzer"
	"github.com/toshi0607/kompal-weather/pkg/kompal"
	"github.com/toshi0607/kompal-weather/pkg/notifier"
	"github.com/toshi0607/kompal-weather/pkg/storage"
	"golang.org/x/sync/errgroup"
)

type controller struct {
	kompal    kompal.Fetcher
	storage   storage.Storage
	notifiers []notifier.Notifier
	analyzer  analyzer.Analyzer
}

func New(f kompal.Fetcher, s storage.Storage, ns []notifier.Notifier, a analyzer.Analyzer) *controller {
	return &controller{
		kompal:    f,
		storage:   s,
		notifiers: ns,
		analyzer:  a,
	}
}

func (c controller) Run(ctx context.Context) error {
	s, err := c.kompal.Fetch(ctx)
	if err != nil {
		return err
	}

	if err := c.storage.Save(ctx, s); err != nil {
		return err
	}

	r, err := c.analyzer.Analyze(ctx)
	if err != nil {
		return err
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, n := range c.notifiers {
		eg.Go(func() error {
			return n.Notify(ctx, r)
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}
