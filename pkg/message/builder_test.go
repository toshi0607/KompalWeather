package message

import (
	"testing"

	"github.com/toshi0607/kompal-weather/pkg/analyzer"
	"github.com/toshi0607/kompal-weather/pkg/status"
)

func TestBuild(t *testing.T) {
	tests := map[string]struct {
		result *analyzer.Result
		want   string
	}{
		"male + crowded, female - few": {
			result: &analyzer.Result{
				MaleTrend:   analyzer.Increasing,
				FemaleTrend: analyzer.Decreasing,
				LatestStatus: status.Status{
					MaleSauna:   status.Crowded,
					FemaleSauna: status.Few,
				},
			},
			want: "男湯サウナは混んできました。現在少し混んでます。\n女湯サウナは空いてきました。現在空いてます。\n（0001年01月01日 00時00分現在）\n最新状況はHPから！ https://kom-pal.com/",
		},
		"male unknown, female + constant": {
			result: &analyzer.Result{
				MaleTrend:   analyzer.Unknown,
				FemaleTrend: analyzer.Constant,
				LatestStatus: status.Status{
					MaleSauna:   status.Normal,
					FemaleSauna: status.Full,
				},
			},
			want: "男湯サウナは現在確認中です。\n女湯サウナは変わりありません。現在満員です。\n（0001年01月01日 00時00分現在）\n最新状況はHPから！ https://kom-pal.com/",
		},
		"Close": {
			result: &analyzer.Result{
				MaleTrend:   analyzer.Close,
				FemaleTrend: analyzer.Close,
				LatestStatus: status.Status{
					MaleSauna:   status.Off,
					FemaleSauna: status.Off,
				},
			},
			want: "本日の営業は終了しました。また来てね！",
		},
		"Open": {
			result: &analyzer.Result{
				MaleTrend:   analyzer.Open,
				FemaleTrend: analyzer.Open,
				LatestStatus: status.Status{
					MaleSauna:   status.Few,
					FemaleSauna: status.Few,
				},
			},
			want: "本日の営業を開始しました！来てね！",
		},
	}

	for name, te := range tests {
		te := te
		m := Build(te.result)

		if m != te.want {
			t.Errorf("[%s] status got: %s, want: %s", name, m, te.want)
		}
	}
}
