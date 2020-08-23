package status

import "time"

//　0 : 空いてます（0-2人）
//　1 : 普通です（3-6人）
//　2 : 少し混んでます（7-8人）
//　3 : 満員です（9人）

type Status struct {
	MaleSauna   int `json:"male_sauna"`
	FemaleSauna int `json:"female_sauna"`
	Timestamp   time.Time
	CreatedAt   time.Time
}

func ToRow(s string) (string, error) {
	return "", nil
}
