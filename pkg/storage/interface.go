package storage

import (
	"context"

	"github.com/toshi0607/kompal-weather/pkg/status"
)

type Storage interface {
	Statuses(ctx context.Context) ([]status.Status, error)
	Save(ctx context.Context, st *status.Status) error
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
