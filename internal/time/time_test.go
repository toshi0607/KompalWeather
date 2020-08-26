package time

import (
	"testing"
	"time"
)

func TestToJST(t *testing.T) {
	time, err := time.Parse(layout, "2006-01-02T15:04:05Z")
	if err != nil {
		t.Error(err)
	}
	got := ToJST(time)
	if got.Location() != jst {
		t.Errorf("got: %s, want: %s", got.Location(), jst)
	}
}

func TestToJSTString(t *testing.T) {
	time, err := time.Parse(layout, "2006-01-02T15:04:05Z")
	if err != nil {
		t.Error(err)
	}
	got := ToJSTString(time)
	want := "2006-01-03T00:04:05+09:00"
	if got != want {
		t.Errorf("got: %s, want: %s", got, want)
	}
}

func TestToJSTTime(t *testing.T) {
	got, err := ToJSTTime("2006-01-03T00:04:05+09:00")
	if err != nil {
		t.Error(err)
	}
	if got.Location() != jst {
		t.Errorf("got: %s, want: %s", got.Location(), jst)
	}
}
