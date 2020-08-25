package storage

import (
	"context"
	"strconv"
	"time"

	"github.com/toshi0607/kompal-weather/pkg/status"
	"gopkg.in/Iwark/spreadsheet.v2"

	"golang.org/x/oauth2/google"
)

var _ Storage = (*Sheets)(nil)
var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

// https://golang.org/pkg/time/#Time.String
var layout = "2006-01-02 15:04:05.999999999 -0700"

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

func (s *Sheets) Save(ctx context.Context, st *status.Status) error {
	panic("implement me")
}
