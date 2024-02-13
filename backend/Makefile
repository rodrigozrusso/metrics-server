.PHONY: all clean

SIMULATOR_BINARY=bin/simulator
METRICS_SERVER_API_BINARY=bin/metrics-server-api
VERSION=$(shell git tag | tail -1)
BUILD=$(shell git rev-parse HEAD)

all: server simulator

.PHONY: server
server:
	go build -o $(METRICS_SERVER_API_BINARY) ./cmd/metrics-server-api/main.go

.PHONY: simulator
simulator:
	go build -o $(SIMULATOR_BINARY) ./cmd/simulator/main.go

.PHONY: test
test:
	go test -v -cover -coverprofile=coverage.out  $$(go list ./... | grep -v '/vendor/')

.PHONY: dependency
dependency:
	go mod download && go mod verify

.PHONY: clean
clean:
	-rm -f $(SIMULATOR_BINARY) $(METRICS_SERVER_API_BINARY)

# .PHONY: build/docker
build/docker:
	docker-compose build

.PHONY: build/docker/service
build/docker/service:
	docker build -t metrics-server-api -f Dockerfile.server .
