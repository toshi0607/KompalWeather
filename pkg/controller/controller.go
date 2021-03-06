package controller

import (
	"context"
	"fmt"

	"github.com/toshi0607/kompal-weather/pkg/analyzer"
	"github.com/toshi0607/kompal-weather/pkg/gcs"
	"github.com/toshi0607/kompal-weather/pkg/kompal"
	"github.com/toshi0607/kompal-weather/pkg/logger"
	"github.com/toshi0607/kompal-weather/pkg/monitoring"
	"github.com/toshi0607/kompal-weather/pkg/notifier"
	"github.com/toshi0607/kompal-weather/pkg/path"
	"github.com/toshi0607/kompal-weather/pkg/report"
	"github.com/toshi0607/kompal-weather/pkg/storage"
	"golang.org/x/sync/errgroup"
)

type controller struct {
	kompal    kompal.Fetcher
	storage   storage.Storage
	notifiers []notifier.Notifier
	analyzer  analyzer.Analyzer
	monitor   *monitoring.Monitor
	gcs       *gcs.GCS
	log       logger.Logger
}

// New builds a new Controller
func New(f kompal.Fetcher, s storage.Storage, ns []notifier.Notifier, a analyzer.Analyzer, m *monitoring.Monitor, g *gcs.GCS, l logger.Logger) Controller {
	return &controller{
		kompal:    f,
		storage:   s,
		notifiers: ns,
		analyzer:  a,
		monitor:   m,
		gcs:       g,
		log:       l,
	}
}

func (c controller) Watch(ctx context.Context) (*analyzer.Result, error) {
	f, err := c.kompal.Fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch kompal status: %v", err)
	}
	c.log.Info("status fetched: %+v", *f)

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		_, err := c.storage.Save(ctx, f)
		if err != nil {
			return fmt.Errorf("failed to save status: %v", err)
		}
		c.log.Info("status saved: %+v", *f)
		return nil
	})
	eg.Go(func() error {
		if err := c.monitor.CreatePoint(ctx, f); err != nil {
			c.log.Error("failed to create point", err)
			// Keep processing
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("failed to save status: %v", err)
	}

	result, err := c.analyzer.Analyze(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze: %v", err)
	}
	c.log.Info("result: %+v", *result)

	eg2, ctx := errgroup.WithContext(context.Background())
	for _, n := range c.notifiers {
		n := n
		eg2.Go(func() error {
			return n.Notify(ctx, result)
		})
	}
	if err := eg2.Wait(); err != nil {
		return nil, fmt.Errorf("failed to notify: %v", err)
	}

	return result, nil
}

func (c controller) Trend(ctx context.Context, k report.Kind) error {
	var malePath, femalePath, msg string
	switch k {
	case report.WeekAgoReport:
		malePath = path.MaleWeekAgoReportObject()
		femalePath = path.FemaleWeekAgoReportObject()
		msg = "先週の今日の混雑傾向です！\n"
	case report.WeeklyReport:
		malePath = path.MaleWeeklyReportObject()
		femalePath = path.FemaleWeeklyReportObject()
		msg = "直近1週間の混雑傾向です！\n"
	default:
		err := fmt.Errorf("unexpected report kind: %s", k)
		c.log.Error("unexpected report kind", err)
		return err
	}
	msg += `縦軸の数字の意味はつぎのとおりです。
0: 営業時間外です
1: 空いてます
2: 普通です
3: 少し混んでます
4: 満員です
`

	var images [][]byte
	mb, err := c.gcs.Get(ctx, path.MaleWeekAgoReportObject())
	if err != nil {
		return fmt.Errorf("failed to get image, path: %s, error: %v", malePath, err)
	}
	images = append(images, mb)
	fmb, err := c.gcs.Get(ctx, femalePath)
	if err != nil {
		return fmt.Errorf("failed to get image, path: %s, error: %v", femalePath, err)
	}
	images = append(images, fmb)

	for _, n := range c.notifiers {
		n := n
		if err := n.NotifyWithMedium(ctx, msg, images); err != nil {
			return fmt.Errorf("failed to notify image with message: %v", err)
		}
	}

	return nil
}
