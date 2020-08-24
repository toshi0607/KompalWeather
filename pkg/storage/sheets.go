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
	//cDt, err := time.ParseInLocation("", currentRow[2].Value, jst)
	//if err != nil {
	//	return nil, err
	//}

	previousRow := sheet.Rows[len(sheet.Rows)-2]

	pMale, err := strconv.Atoi(previousRow[0].Value)
	if err != nil {
		return nil, err
	}
	pFemale, err := strconv.Atoi(previousRow[1].Value)
	if err != nil {
		return nil, err
	}
	//pDt := previousRow[2]
	return []status.Status{
		{
			MaleSauna:   status.Sauna(cMale),
			FemaleSauna: status.Sauna(cFemale),
		},
		{
			MaleSauna:   status.Sauna(pMale),
			FemaleSauna: status.Sauna(pFemale),
		},
	}, nil
}

func (s *Sheets) Save(ctx context.Context, st *status.Status) error {
	panic("implement me")
}
