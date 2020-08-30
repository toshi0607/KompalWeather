package message

import (
	"fmt"
	"strings"

	"github.com/toshi0607/kompal-weather/pkg/analyzer"
)

func Build(r *analyzer.Result) string {
	if r.MaleTrend == analyzer.Close {
		return "本日の営業は終了しました。また来てね！"
	} else if r.MaleTrend == analyzer.Open {
		return "本日の営業を開始しました！来てね！"
	}

	var m strings.Builder
	m.WriteString(maleMessage(r))
	m.WriteString(femaleMessage(r))
	m.WriteString(fmt.Sprintf("（%s現在）\n", r.LatestStatus.Timestamp.Format("2006年01月02日 15時04分")))
	m.WriteString("最新状況はHPから！ https://kom-pal.com/")

	return m.String()
}

func maleMessage(r *analyzer.Result) string {
	if r.MaleTrend == analyzer.Unknown {
		return "男湯サウナは現在確認中です。\n"
	}
	return fmt.Sprintf("男湯サウナは%s現在%s\n", r.MaleTrend, r.LatestStatus.MaleSauna)
}

func femaleMessage(r *analyzer.Result) string {
	if r.FemaleTrend == analyzer.Unknown {
		return "女湯サウナは現在確認中です。\n"
	}
	return fmt.Sprintf("女湯サウナは%s現在%s\n", r.FemaleTrend, r.LatestStatus.FemaleSauna)
}
