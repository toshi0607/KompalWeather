package status

import "time"

// Status represents how Komparu-yu is crowded
type Status struct {
	MaleSauna   Sauna `json:"male_sauna"`
	FemaleSauna Sauna `json:"female_sauna"`
	Timestamp   time.Time
	CreatedAt   time.Time
}

// Sauna is status of sauna
type Sauna int

const (
	// Off means out of business hours
	Off = 0
	// Few means few
	Few = 1
	// Normal means normal
	Normal = 2
	// Crowded means crowded
	Crowded = 3
	// Full means full
	Full = 4
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
