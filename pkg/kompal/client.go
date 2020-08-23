package kompal

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/toshi0607/kompal-weather/pkg/status"
)

const kompalUrl = "https://kom-pal.com/"

var _ fetcher = (*Kompal)(nil)

type Kompal struct {
}

func (k *Kompal) Fetch() (*status.Status, error) {
	resp, err := http.Get(kompalUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Print(html)
	// take status and timestamp from html
	return &status.Status{}, nil
}
