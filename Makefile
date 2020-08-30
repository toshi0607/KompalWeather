.PHONY: run
run:
	go run cmd/server/kompal_weather/main.go

.PHONY: build
build:
	go build -v .

.PHONY: test.s
test.s:
	go test ./... -v -short

.PHONY: tidy
tidy:
	go mod tidy
