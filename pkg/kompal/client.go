package kompal

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	t "github.com/toshi0607/kompal-weather/internal/time"
	"github.com/toshi0607/kompal-weather/pkg/status"
)

var _ Fetcher = (*Kompal)(nil)

// Kompal represents Kompal-yu API
type Kompal struct {
	config *Config
}

// Config represents a configuration for Kompal-yu API
type Config struct {
	URL string
}

// Response represents a response from Kompal-yu API
type Response struct {
	MaleSauna   Sauna `json:"male_sauna"`
	FemaleSauna Sauna `json:"female_sauna"`
	Open        bool
	Timestamp   time.Time
}

// ToStatus converts Response to status.Status
func (kr *Response) ToStatus() *status.Status {
	return &status.Status{
		MaleSauna:   kr.MaleSauna.ToStatusSauna(kr.Open),
		FemaleSauna: kr.FemaleSauna.ToStatusSauna(kr.Open),
		Timestamp:   t.ToJST(kr.Timestamp),
	}
}

// Sauna represents status of sauna
type Sauna int

//　0 : 空いてます（0-2人）
//　1 : 普通です（3-6人）
//　2 : 少し混んでます（7-8人）
//　3 : 満員です（9人）
const (
	Few     = 0
	Normal  = 1
	Crowded = 2
	Full    = 3
)

// ToStatusSauna converts kompal.Sauna to status.Sauna
func (s Sauna) ToStatusSauna(open bool) status.Sauna {
	if !open {
		return status.Off
	}
	switch s {
	case Few:
		return status.Few
	case Normal:
		return status.Normal
	case Crowded:
		return status.Crowded
	case Full:
		return status.Full
	default:
		return status.Off
	}
}

// New builds new Kompal
func New(c *Config) *Kompal {
	return &Kompal{c}
}

// Fetch fetches the latest status of Kompal-yu
func (k *Kompal) Fetch(ctx context.Context) (*status.Status, error) {
	// TODO: use ctx,make http client replaceable
	resp, err := http.Get(k.config.URL)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response status code: %d", resp.StatusCode)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close resp body: %s", err)
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var kr Response
	if err := json.Unmarshal(body, &kr); err != nil {
		return nil, err
	}

	return kr.ToStatus(), nil
}
