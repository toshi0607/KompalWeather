.PHONY: run
run:
	go run cmd/server/kompal_weather/main.go

.PHONY: build
build:
	go build -v cmd/server/kompal_weather/main.go

.PHONY: test.s
test.s:
	go test ./... -v -short

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: lint
lint:
	golangci-lint run -v --tests ./...

.PHONY: watch
watch:
	curl -sS http://localhost:8080/watch -X POST
