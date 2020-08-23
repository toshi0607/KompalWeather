package status

import "time"

//男性サウナ：普通です。 女性サウナ：空いてます。
//(8月23日 13:41現在)

type Status struct {
	Male      string
	Female    string
	Timestamp time.Time
}

func ToRow(s string) (string, error) {
	return "", nil
}
