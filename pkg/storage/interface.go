package storage

import (
	"github.com/toshi0607/kompal-weather/pkg/status"
)

type Storage interface {
	Statuses() ([]status.Status, error)
	Save(st status.Status) error
}

//func test() {
//	var s Storage
//	hoge, _ := s.Statuses()
//	fmt.Print(s)
//
//	service, _ := NewSheetsClient()
//	_, _ := service.Statuses()
//
//}
