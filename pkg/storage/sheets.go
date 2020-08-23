package storage

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/toshi0607/kompal-weather/pkg/status"

	"golang.org/x/oauth2/google"
)

var _ Storage = (*Sheets)(nil)

//https://developers.google.com/sheets/api/quickstart/go?authuser=1
//https://developers.google.com/sheets/api/quickstart/go?authuser=1

type Sheets struct {
	service *spreadsheet.Service
}

func NewSheetsClient() (*Sheets, error) {
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

func (s *Sheets) Statuses() ([]status.Status, error) {
	const id = ""
	ss, _ := s.service.FetchSpreadsheet(id)

	var sheetId uint = 120
	sheet, _ := ss.SheetByID(sheetId)

	currentRow := sheet.Rows[:len(sheet.Rows)-1]
	cMale := currentRow[0]
	cFemale := currentRow[1]
	cDt := currentRow[2]

	previousRow := sheet.Rows[:len(sheet.Rows)-2]
	pMale := previousRow[0]
	pFemale := previousRow[1]
	pDt := previousRow[2]
	fmt.Print(cMale, cFemale, cDt, pMale, pFemale, pDt)
	return nil, nil
}

func (s *Sheets) Save(st status.Status) error {
	panic("implement me")
}
