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

type PointType string

const (
	typeMale   = PointType("custom.googleapis.com/male_status")
	typeFemale = PointType("custom.googleapis.com/female_status")
)

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

// Close closes the client connection
func (m Monitor) Close() error {
	return m.client.Close()
}

// CreatePoint creates points for each status in Cloud Monitoring
func (m Monitor) CreatePoint(ctx context.Context, s *status.Status) error {
	male := m.createPoint(s.MaleSauna, s.Timestamp.Unix())
	female := m.createPoint(s.FemaleSauna, s.Timestamp.Unix()+1)

	if err := m.createTimeSeries(ctx, typeMale, male); err != nil {
		return err
	}
	if err := m.createTimeSeries(ctx, typeFemale, female); err != nil {
		return err
	}

	return nil
}

func (m Monitor) createTimeSeries(ctx context.Context, t PointType, p *monitoringpb.Point) error {
	if err := m.client.CreateTimeSeries(ctx, &monitoringpb.CreateTimeSeriesRequest{
		Name: fmt.Sprintf("projects/%s", m.gcpProjectID),
		TimeSeries: []*monitoringpb.TimeSeries{
			{
				Metric: &metricpb.Metric{
					Type:   string(t),
					Labels: nil,
				},
				Resource: &monitoredrespb.MonitoredResource{
					Type: "global",
					Labels: map[string]string{
						"project_id": m.gcpProjectID,
					},
				},
				Points: []*monitoringpb.Point{
					p,
				},
			},
		},
	}); err != nil {
		return fmt.Errorf("failed to createTimeSeries: %v", err)
	}
	return nil
}

func (m Monitor) createPoint(s status.Sauna, timestampSec int64) *monitoringpb.Point {
	return &monitoringpb.Point{
		Interval: &monitoringpb.TimeInterval{
			EndTime: &googlepb.Timestamp{
				Seconds: timestampSec,
			},
		},
		Value: &monitoringpb.TypedValue{
			Value: &monitoringpb.TypedValue_Int64Value{
				Int64Value: int64(s),
			},
		},
	}
}
