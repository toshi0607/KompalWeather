package kompal

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	t "github.com/toshi0607/kompal-weather/internal/time"
	"github.com/toshi0607/kompal-weather/pkg/status"
)

var _ Fetcher = (*Kompal)(nil)

type Kompal struct {
	config *Config
}

type Config struct {
	URL string
}

type KompalResponse struct {
	MaleSauna   Sauna `json:"male_sauna"`
	FemaleSauna Sauna `json:"female_sauna"`
	Open        bool
	Timestamp   time.Time
}

func (kr *KompalResponse) ToStatus() *status.Status {
	return &status.Status{
		MaleSauna:   kr.MaleSauna.ToStatusSauna(kr.Open),
		FemaleSauna: kr.FemaleSauna.ToStatusSauna(kr.Open),
		Timestamp:   t.ToJST(kr.Timestamp),
	}
}

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

func New(c *Config) *Kompal {
	return &Kompal{c}
}

func (k *Kompal) Fetch(ctx context.Context) (*status.Status, error) {
	// TODO: use ctx,make http client replaceable
	resp, err := http.Get(k.config.URL)
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

	var kr KompalResponse
	if err := json.Unmarshal(body, &kr); err != nil {
		return nil, err
	}

	return kr.ToStatus(), nil
}
