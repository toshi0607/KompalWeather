package monitoring

import (
	"context"
	"fmt"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	googlepb "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/toshi0607/kompal-weather/pkg/status"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoredrespb "google.golang.org/genproto/googleapis/api/monitoredres"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

// Secret represents a monitor
type Monitor struct {
	client       *monitoring.MetricClient
	gcpProjectID string
}

// Add close

// New builds new Monitor
func New(ID string) (*Monitor, error) {
	ctx := context.TODO()
	client, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return nil, err
	}
	return &Monitor{
		client:       client,
		gcpProjectID: ID,
	}, nil
}

func (m Monitor) CreatePoint(ctx context.Context, s *status.Status) error {
	male := &monitoringpb.Point{
		Interval: &monitoringpb.TimeInterval{
			EndTime: &googlepb.Timestamp{
				Seconds: s.Timestamp.Unix(),
			},
		},
		Value: &monitoringpb.TypedValue{
			Value: &monitoringpb.TypedValue_Int64Value{
				Int64Value: int64(s.MaleSauna),
			},
		},
	}
	female := &monitoringpb.Point{
		Interval: &monitoringpb.TimeInterval{
			EndTime: &googlepb.Timestamp{
				Seconds: s.Timestamp.Unix() + 1,
			},
		},
		Value: &monitoringpb.TypedValue{
			Value: &monitoringpb.TypedValue_Int64Value{
				Int64Value: int64(s.FemaleSauna),
			},
		},
	}

	if err := m.client.CreateTimeSeries(ctx, &monitoringpb.CreateTimeSeriesRequest{
		Name: fmt.Sprintf("projects/%s", m.gcpProjectID),
		TimeSeries: []*monitoringpb.TimeSeries{
			{
				Metric: &metricpb.Metric{
					Type: "custom.googleapis.com/male_status",
				},
				Resource: &monitoredrespb.MonitoredResource{
					Type: "global",
					Labels: map[string]string{
						"project_id": m.gcpProjectID,
					},
				},
				Points: []*monitoringpb.Point{
					male,
				},
			},
		},
	}); err != nil {
		return err
	}

	if err := m.client.CreateTimeSeries(ctx, &monitoringpb.CreateTimeSeriesRequest{
		Name: fmt.Sprintf("projects/%s", m.gcpProjectID),
		TimeSeries: []*monitoringpb.TimeSeries{
			{
				Metric: &metricpb.Metric{
					Type: "custom.googleapis.com/female_status",
				},
				Resource: &monitoredrespb.MonitoredResource{
					Type: "global",
					Labels: map[string]string{
						"project_id": m.gcpProjectID,
					},
				},
				Points: []*monitoringpb.Point{
					female,
				},
			},
		},
	}); err != nil {
		return err
	}

	return nil
}
