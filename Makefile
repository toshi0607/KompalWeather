.PHONY: run
run:
	go run cmd/server/kompal_weather/main.go

.PHONY: tidy
tidy:
	go mod tidy
