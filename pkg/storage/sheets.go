package storage

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/toshi0607/kompal-weather/pkg/status"
	"gopkg.in/Iwark/spreadsheet.v2"

	"golang.org/x/oauth2/google"
)

var _ Storage = (*Sheets)(nil)

//https://developers.google.com/sheets/api/quickstart/go?authuser=1
//https://github.com/Iwark/spreadsheet

const (
	spreadSheetID = ""
	sheetId       = 120
)

type Sheets struct {
	service *spreadsheet.Service
}

func NewSheets() (*Sheets, error) {
	//https://github.com/golang/oauth2/blob/master/google/default.go 自分のdefaultかCloud RunのSAがいいなぁ…
	data, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		return nil, err
	}

	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	if err != nil {
		return nil, err
	}

	client := conf.Client(context.TODO())
	service := spreadsheet.NewServiceWithClient(client)
	return &Sheets{
		service: service,
	}, nil
}

func (s *Sheets) Statuses(ctx context.Context) ([]status.Status, error) {
	ss, err := s.service.FetchSpreadsheet(spreadSheetID)
	if err != nil {
		return nil, err
	}

	sheet, err := ss.SheetByID(sheetId)
	if err != nil {
		return nil, err
	}
	//l := sheet.Rows
	//l[0][0].Value

	currentRow := sheet.Rows[:len(sheet.Rows)-1][0]
	cMale := currentRow[0]
	cFemale := currentRow[1]
	cDt := currentRow[2]

	previousRow := sheet.Rows[:len(sheet.Rows)-2]
	pMale := previousRow[0]
	pFemale := previousRow[1]
	pDt := previousRow[2]
	fmt.Print(cMale, cFemale, cDt, pMale, pFemale, pDt)
	return []status.Status{
		//{
		//	Male: cMale.Value,
		//},
	}, nil
}

func (s *Sheets) Save(ctx context.Context, st *status.Status) error {
	panic("implement me")
}
