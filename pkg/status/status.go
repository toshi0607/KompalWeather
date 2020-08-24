package status

import "time"

type Status struct {
	MaleSauna   Sauna `json:"male_sauna"`
	FemaleSauna Sauna `json:"female_sauna"`
	Timestamp   time.Time
	CreatedAt   time.Time
}

type Sauna int

const (
	Few     = 0
	Normal  = 1
	Crowded = 2
	Full    = 3
)

//　0 : 空いてます（0-2人）
//　1 : 普通です（3-6人）
//　2 : 少し混んでます（7-8人）
//　3 : 満員です（9人）
func (s Sauna) String() string {
	switch s {
	case Few:
		return "空いてます。"
	case Normal:
		return "普通です。"
	case Crowded:
		return "少し混んでます"
	case Full:
		return "満員です。"
	default:
		return "Invalid"
	}
}

func ToRow(s string) (string, error) {
	return "", nil
}
