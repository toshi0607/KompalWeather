package kompal

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/toshi0607/kompal-weather/pkg/status"
)

var _ Fetcher = (*Kompal)(nil)
var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

type Kompal struct {
	url string
}

func New(url string) *Kompal {
	return &Kompal{url: url}
}

func (k *Kompal) Fetch(ctx context.Context) (*status.Status, error) {
	// TODO: use ctx,make http client replaceable
	resp, err := http.Get(k.url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var st status.Status
	if err := json.Unmarshal(body, &st); err != nil {
		return nil, err
	}

	return &status.Status{
		MaleSauna:   st.MaleSauna,
		FemaleSauna: st.MaleSauna,
		Timestamp:   st.Timestamp.In(jst),
		CreatedAt:   nil,
	}, nil
}
