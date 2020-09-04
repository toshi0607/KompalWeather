package storage

import (
	"context"
	"fmt"
	"strconv"
	"time"

	t "github.com/toshi0607/kompal-weather/internal/time"
	"github.com/toshi0607/kompal-weather/pkg/status"
	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
)

var _ Storage = (*Sheets)(nil)

const (
	maleSaunaIndex   = 0
	femaleSaunaIndex = 1
	timestampIndex   = 2
	createdAtIndex   = 3
)

// SheetsConfig is a configuration of Sheets
type SheetsConfig struct {
	SpreadSheetID string
	SheetID       uint
}

// Sheets is representation of Sheets
type Sheets struct {
	service *spreadsheet.Service
	config  *SheetsConfig
}

// NewSheets builds new Sheets
func NewSheets(c *SheetsConfig) (*Sheets, error) {
	ctx := context.TODO()
	client, err := google.DefaultClient(ctx, spreadsheet.Scope)
	if err != nil {
		return nil, fmt.Errorf("failed to build default google client: %v", err)
	}

	service := spreadsheet.NewServiceWithClient(client)
	return &Sheets{
		service: service,
		config:  c,
	}, nil
}

// Statuses returns the last two statuses from a spreadsheet
func (s *Sheets) Statuses(ctx context.Context) ([]status.Status, error) {
	ss, err := s.service.FetchSpreadsheet(s.config.SpreadSheetID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch spreadsheet: %v", err)
	}
	sheet, err := ss.SheetByID(s.config.SheetID)
	if err != nil {
		return nil, fmt.Errorf("failed to find sheet by ID: %v", err)
	}

	// Sheet example:
	//   male female timestamp createdAt
	//   3 0 2020-08-26T19:20:18+09:00 2020-08-26T20:00:03+09:00 prev
	//   3 1 2020-08-26T20:11:39+09:00 2020-08-26T20:20:00+09:00 current
	current, err := toStatus(sheet.Rows[len(sheet.Rows)-1])
	if err != nil {
		return nil, fmt.Errorf("failed to convert current row to status: %v", err)
	}
	prev, err := toStatus(sheet.Rows[len(sheet.Rows)-2])
	if err != nil {
		return nil, fmt.Errorf("failed to convert prev row to status: %v", err)
	}

	return []status.Status{
		*current,
		*prev,
	}, nil
}

func toStatus(data []spreadsheet.Cell) (*status.Status, error) {
	male, err := strconv.Atoi(data[maleSaunaIndex].Value)
	if err != nil {
		return nil, fmt.Errorf("failed to convert male sauna string to Sauna: %v", err)
	}
	female, err := strconv.Atoi(data[femaleSaunaIndex].Value)
	if err != nil {
		return nil, fmt.Errorf("failed to convert female sauna string to Sauna: %v", err)
	}
	timestamp, err := t.ToJSTTime(data[timestampIndex].Value)
	if err != nil {
		return nil, fmt.Errorf("failed to convert timestamp to JST time: %v", err)
	}
	createdAt, err := t.ToJSTTime(data[createdAtIndex].Value)
	if err != nil {
		return nil, fmt.Errorf("failed to convert createdAt to JST time: %v", err)
	}

	return &status.Status{
		MaleSauna:   status.Sauna(male),
		FemaleSauna: status.Sauna(female),
		Timestamp:   timestamp,
		CreatedAt:   createdAt,
	}, nil
}

// Save saves status in a spreadsheet
func (s *Sheets) Save(ctx context.Context, st *status.Status) (*status.Status, error) {
	ss, err := s.service.FetchSpreadsheet(s.config.SpreadSheetID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch the spreadsheet: %v", err)
	}
	sheet, err := ss.SheetByID(s.config.SheetID)
	if err != nil {
		return nil, fmt.Errorf("failed to find the sheet: %v", err)
	}

	if sheet.Properties.GridProperties.RowCount == uint(len(sheet.Rows)) {
		if err := s.service.ExpandSheet(sheet, 0, 1000); err != nil {
			return nil, fmt.Errorf("failed to expand the sheet: %v", err)
		}
	}

	targetRowIndex := len(sheet.Rows)
	now := t.ToJST(time.Now())
	sheet.Update(targetRowIndex, maleSaunaIndex, strconv.Itoa(int(st.MaleSauna)))
	sheet.Update(targetRowIndex, femaleSaunaIndex, strconv.Itoa(int(st.FemaleSauna)))
	sheet.Update(targetRowIndex, timestampIndex, t.ToJSTString(st.Timestamp))
	sheet.Update(targetRowIndex, createdAtIndex, t.ToJSTString(now))
	if err := sheet.Synchronize(); err != nil {
		return nil, fmt.Errorf("failed to update the sheet: %v", err)
	}
	st.CreatedAt = now

	return st, nil
}
