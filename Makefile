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

.PHONY: trend
trend:
	curl -sSi http://localhost:8080/trend -X POST \
	  --data '{"reportType": "weekAgo"}'

.PHONY: v
v:
	curl -sSi http://localhost:8080/visualize -X POST \
	  --data '{"reportType": "daily"}'

.PHONY: dr
dr:
	docker run -p 127.0.0.1:8080:8080/tcp --rm --mount type=bind,source=$(shell pwd),target=/app --env-file .env --name vis5 visualizer:0.1

.PHONY: db
db:
	docker build -f ./Dockerfile.visualizer -t visualizer:0.1 .

.PHONY: de
de:
	docker exec -it vis5 bash

.PHONY: dc
dc:
	docker cp vis5:/tmp/. ./tmp/
