.PHONY: all clean

SIMULATOR_BINARY=bin/simulator
ANALYTICS_SERVER_BINARY=bin/analytics-server
VERSION=$(shell git tag | tail -1)
BUILD=$(shell git rev-parse HEAD)

all: server simulator

.PHONY: server
server:
	go build -o $(ANALYTICS_SERVER_BINARY) ./cmd/analytics-server/main.go

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
	-rm -f $(SIMULATOR_BINARY) $(ANALYTICS_SERVER_BINARY)

# .PHONY: build/docker
build/docker:
	docker-compose build

.PHONY: build/docker/service
build/docker/service:
	docker build -t analytics-server -f Dockerfile.server .
