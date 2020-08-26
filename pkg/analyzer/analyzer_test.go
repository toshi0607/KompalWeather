package analyzer

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/toshi0607/kompal-weather/pkg/status"
	"github.com/toshi0607/kompal-weather/pkg/storage"
)

var _ storage.Storage = (*memory)(nil)

type memory struct {
	statuses []status.Status
}

func (m *memory) Statuses(ctx context.Context) ([]status.Status, error) {
	return m.statuses, nil
}

func (m *memory) Save(ctx context.Context, st *status.Status) (*status.Status, error) {
	return nil, nil
}

func TestAnalyzer_Analyze(t *testing.T) {
	var jst = time.FixedZone("Asia/Tokyo", 9*60*60)
	var layout = time.RFC3339
	time1, err := time.ParseInLocation(layout, "2020-08-22T00:03:30+09:00", jst)
	if err != nil {
		t.Error(err)
	}
	time2, err := time.ParseInLocation(layout, "2020-08-23T00:03:30+09:00", jst)
	if err != nil {
		t.Error(err)
	}
	time3, err := time.ParseInLocation(layout, "2020-08-24T00:03:30+09:00", jst)
	if err != nil {
		t.Error(err)
	}
	time4, err := time.ParseInLocation(layout, "2020-08-25T00:03:30+09:00", jst)
	if err != nil {
		t.Error(err)
	}

	tests := map[string]struct {
		statuses []status.Status
		want     *Result
	}{
		"Male Increasing, Female Constant": {
			statuses: []status.Status{
				{
					MaleSauna:   status.Few,
					FemaleSauna: status.Normal,
					Timestamp:   time1,
					CreatedAt:   time2,
				},
				{
					MaleSauna:   status.Normal,
					FemaleSauna: status.Normal,
					Timestamp:   time3,
					CreatedAt:   time4,
				},
			},
			want: &Result{
				MaleTrend:   Increasing,
				FemaleTrend: Constant,
				LatestStatus: status.Status{
					MaleSauna:   status.Normal,
					FemaleSauna: status.Normal,
					Timestamp:   time3,
					CreatedAt:   time4,
				},
			},
		},
		"Male Constant, Female Decreasing": {
			statuses: []status.Status{
				{
					MaleSauna:   status.Few,
					FemaleSauna: status.Normal,
					Timestamp:   time3,
					CreatedAt:   time4,
				},
				{
					MaleSauna:   status.Few,
					FemaleSauna: status.Crowded,
					Timestamp:   time1,
					CreatedAt:   time2,
				},
			},
			want: &Result{
				MaleTrend:   Constant,
				FemaleTrend: Decreasing,
				LatestStatus: status.Status{
					MaleSauna:   status.Few,
					FemaleSauna: status.Normal,
					Timestamp:   time3,
					CreatedAt:   time4,
				},
			},
		},
	}

	for name, te := range tests {
		te := te
		ctx := context.TODO()
		a := New(&memory{te.statuses})

		r, err := a.Analyze(ctx)

		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(r, te.want); diff != "" {
			fmt.Printf("[%s] result r != te.want\n%s\n", name, diff)
		}
	}
}
