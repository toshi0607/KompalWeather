FROM golang:1.14.6

WORKDIR /project

COPY ./go.* ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go install -v ./cmd/server/kompal_weather


FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=0 /go/bin/kompal_weather /bin/kompal_weather

RUN addgroup -g 1001 app && adduser -D -G app -u 1001 app

USER 1001

CMD ["/bin/kompal_weather"]
