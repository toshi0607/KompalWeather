package status

import (
	"testing"
)

func TestSauna_String(t *testing.T) {
	tests := map[string]struct {
		sauna Sauna
		want  string
	}{
		"Off": {
			sauna: Off,
			want:  "営業時間外です。",
		},
		"Few": {
			sauna: Few,
			want:  "空いてます。",
		},
		"Normal": {
			sauna: Normal,
			want:  "普通です。",
		},
		"Crowded": {
			sauna: Crowded,
			want:  "少し混んでます。",
		},
		"Full": {
			sauna: Full,
			want:  "満員です。",
		},
		"Invalid": {
			sauna: 99999,
			want:  "確認中です。",
		},
	}

	for name, te := range tests {
		te := te

		got := te.sauna.String()

		if got != te.want {
			t.Errorf("[%s] got: %s, want: %s", name, got, te.want)
		}
	}
}
