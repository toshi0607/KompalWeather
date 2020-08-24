package message

import (
	"fmt"
	"strings"

	"github.com/toshi0607/kompal-weather/pkg/analyzer"
)

func Build(r *analyzer.Result) string {
	var m strings.Builder
	m.WriteString(maleMessage(r))
	m.WriteString(femaleMessage(r))
	m.WriteString(fmt.Sprint("最新状況はHPから！ https://kom-pal.com/"))

	return m.String()
}

func maleMessage(r *analyzer.Result) string {
	if r.MaleTrend == analyzer.Unknown {
		return fmt.Sprintf("男湯サウナは現在確認中です。\n")
	}
	return fmt.Sprintf("男湯サウナは%s現在%s\n", r.MaleTrend, r.LatestStatus.MaleSauna)
}

func femaleMessage(r *analyzer.Result) string {
	if r.FemaleTrend == analyzer.Unknown {
		return fmt.Sprintf("女湯サウナは現在確認中です。\n")
	}
	return fmt.Sprintf("女湯サウナは%s現在%s\n", r.FemaleTrend, r.LatestStatus.FemaleSauna)
}
