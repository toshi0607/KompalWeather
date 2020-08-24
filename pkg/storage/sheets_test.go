package storage

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
)

func TestNewSheets(t *testing.T) {
	ctx := context.TODO()
	spreadSheetID := os.Getenv("SPREAD_SHEET_ID")
	si := os.Getenv("SHEET_ID")
	sheetId, err := strconv.ParseUint(si, 10, 32)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	c, err := NewSheets(&SheetsConfig{
		SpreadSheetID: spreadSheetID,
		SheetId:       uint(sheetId),
	})
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	ss, err := c.Statuses(ctx)
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	fmt.Sprint(ss)
}
