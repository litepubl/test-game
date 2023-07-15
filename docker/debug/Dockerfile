# Step 1: Modules caching
FROM golang:1.20.5-alpine3.17 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download
RUN CGO_ENABLED=0 go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@latest

#migrate cmd
RUN go install -tags=postgres,clickhouse github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Step 2: Builder
FROM golang:1.20.5-buster as builder
COPY --from=modules /go/pkg /go/pkg
COPY --from=modules /go/bin/dlv /go/bin/dlv
COPY --from=modules /go/bin/migrate /go/bin/migrate
COPY . /app
WORKDIR /app
RUN mkdir -p /logs

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags migrate -gcflags "all=-N -l" -o /app/app ./cmd/app

COPY config /config
VOLUME /logs
EXPOSE 8080 4000
CMD [ "/go/bin/dlv", "--listen=:4000", "--headless=true", "--log=true", "--accept-multiclient", "--api-version=2", "exec", "/app/app" ]