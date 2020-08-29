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
	Off     = 0
	Few     = 1
	Normal  = 2
	Crowded = 3
	Full    = 4
)

func (s Sauna) String() string {
	switch s {
	case Off:
		return "営業時間外です。"
	case Few:
		return "空いてます。"
	case Normal:
		return "普通です。"
	case Crowded:
		return "少し混んでます。"
	case Full:
		return "満員です。"
	default:
		return "確認中です。"
	}
}
