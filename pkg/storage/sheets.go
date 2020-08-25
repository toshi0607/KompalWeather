package storage

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/toshi0607/kompal-weather/pkg/status"
	"gopkg.in/Iwark/spreadsheet.v2"

	"golang.org/x/oauth2/google"
)

var _ Storage = (*Sheets)(nil)
var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

// https://golang.org/pkg/time/#Time.String
var layout = "2006-01-02 15:04:05.999999999 -0700 MST"

type SheetsConfig struct {
	SpreadSheetID string
	SheetId       uint
}

type Sheets struct {
	service *spreadsheet.Service
	config  *SheetsConfig
}

func NewSheets(c *SheetsConfig) (*Sheets, error) {
	ctx := context.TODO()
	client, err := google.DefaultClient(ctx, spreadsheet.Scope)
	if err != nil {
		return nil, err
	}

	service := spreadsheet.NewServiceWithClient(client)
	return &Sheets{
		service: service,
		config:  c,
	}, nil
}

func (s *Sheets) Statuses(ctx context.Context) ([]status.Status, error) {
	ss, err := s.service.FetchSpreadsheet(s.config.SpreadSheetID)
	if err != nil {
		return nil, err
	}
	sheet, err := ss.SheetByID(s.config.SheetId)
	if err != nil {
		return nil, err
	}

	currentRow := sheet.Rows[len(sheet.Rows)-1]
	cMale, err := strconv.Atoi(currentRow[0].Value)
	if err != nil {
		return nil, err
	}
	cFemale, err := strconv.Atoi(currentRow[1].Value)
	if err != nil {
		return nil, err
	}
	cDt, err := time.ParseInLocation(layout, currentRow[2].Value, jst)
	if err != nil {
		return nil, err
	}

	previousRow := sheet.Rows[len(sheet.Rows)-2]
	pMale, err := strconv.Atoi(previousRow[0].Value)
	if err != nil {
		return nil, err
	}
	pFemale, err := strconv.Atoi(previousRow[1].Value)
	if err != nil {
		return nil, err
	}
	pDt, err := time.ParseInLocation(layout, previousRow[2].Value, jst)
	if err != nil {
		return nil, err
	}

	return []status.Status{
		{
			MaleSauna:   status.Sauna(cMale),
			FemaleSauna: status.Sauna(cFemale),
			Timestamp:   cDt,
		},
		{
			MaleSauna:   status.Sauna(pMale),
			FemaleSauna: status.Sauna(pFemale),
			Timestamp:   pDt,
		},
	}, nil
}

func (s *Sheets) Save(ctx context.Context, st *status.Status) (*status.Status, error) {
	ss, err := s.service.FetchSpreadsheet(s.config.SpreadSheetID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch the spreadsheet: %v", err)
	}
	sheet, err := ss.SheetByID(s.config.SheetId)
	if err != nil {
		return nil, fmt.Errorf("failed to find the sheet: %v", err)
	}

	if sheet.Properties.GridProperties.RowCount == uint(len(sheet.Rows)) {
		if err := s.service.ExpandSheet(sheet, 0, 1000); err != nil {
			return nil, fmt.Errorf("failed to expand the sheet: %v", err)
		}
	}

	targetRowIndex := len(sheet.Rows)
	now := time.Now().In(jst)
	sheet.Update(targetRowIndex, 0, strconv.Itoa(int(st.MaleSauna)))
	sheet.Update(targetRowIndex, 1, strconv.Itoa(int(st.FemaleSauna)))
	sheet.Update(targetRowIndex, 2, st.Timestamp.Format(layout))
	sheet.Update(targetRowIndex, 3, now.Format(layout))
	if err := sheet.Synchronize(); err != nil {
		return nil, fmt.Errorf("failed to update the sheet: %v", err)
	}
	st.CreatedAt = now

	return st, nil
}
