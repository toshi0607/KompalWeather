package storage

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/toshi0607/kompal-weather/pkg/status"
)

func TestSheets_Statuses(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping external test")
	}

	ctx := context.TODO()
	spreadSheetID := os.Getenv("SPREAD_SHEET_ID")
	si := os.Getenv("SHEET_ID")
	sheetID, err := strconv.ParseUint(si, 10, 32)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	c, err := NewSheets(&SheetsConfig{
		SpreadSheetID: spreadSheetID,
		SheetID:       uint(sheetID),
	})
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	ss, err := c.Statuses(ctx)
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if len(ss) != 2 {
		t.Errorf("[%s] statuses got: %d, want: 2", "TestNewSheets", len(ss))
	}
	if ss[0].Timestamp.IsZero() || ss[1].Timestamp.IsZero() {
		t.Fatal("Timestamp should not be zero value")
	}
}

func TestSheets_Save(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping external test")
	}

	ctx := context.TODO()
	spreadSheetID := os.Getenv("SPREAD_SHEET_ID")
	si := os.Getenv("SHEET_ID")
	sheetID, err := strconv.ParseUint(si, 10, 32)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	c, err := NewSheets(&SheetsConfig{
		SpreadSheetID: spreadSheetID,
		SheetID:       uint(sheetID),
	})
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	s, err := c.Save(ctx, &status.Status{
		MaleSauna:   status.Normal,
		FemaleSauna: status.Crowded,
		Timestamp:   time.Now(),
	})
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if s.CreatedAt.IsZero() {
		t.Fatal("Timestamp should not be zero value")
	}
}
