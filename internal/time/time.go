package time

import "time"

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

// https://golang.org/pkg/time/#Time.String
// 保存時に Asia/TokyoになるがParseできない
//var layout = "2006-01-02 15:04:05.999999999 -0700 MST"

var layout = time.RFC3339

func ToJST(t time.Time) time.Time {
	return t.In(jst)
}

func ToJSTString(t time.Time) string {
	return t.In(jst).Format(layout)
}

func ToJSTTime(s string) (time.Time, error) {
	t, err := time.ParseInLocation(layout, s, jst)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}