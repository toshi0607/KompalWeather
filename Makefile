.PHONY: run.c
run.c:
	go run cmd/server/core/main.go

.PHONY: run.v
run.v:
	go run cmd/server/visualizer/main.go

.PHONY: build.c
build.c:
	go build -v cmd/server/core/main.go

.PHONY: build.v
build.v:
	go build -v cmd/server/visualizer/main.go

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

.PHONY: v
v:
	curl -sS http://localhost:8080/visualize -X POST \
	  --data '{"reportType": "daily"}'
