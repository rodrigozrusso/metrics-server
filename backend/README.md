# Metrics Server API

The `metrics-server-api` exposes the API to ingest metrics and to consume aggregated metrics per granularity.

_Bonus_: there is a `simulator` (CLI tool) to send multiple metrics concurrently (20 metrics every 300 ms).

## Built with

- [Golang 1.21](https://go.dev/blog/go1.21)
- [TimescaleDB](https://www.timescale.com/)
- [API Swagger Documentation]()


## Getting Started

### Build and run with docker

  ```bash
  docker-compose up -d
  ```

  Check if services are up and running
  ```bash
  docker-compose ps
  ```

  You should be able see something like:

  ```
  NAME                    IMAGE                               COMMAND                  SERVICE              CREATED          STATUS          PORTS
  backend-timescaledb-1   timescale/timescaledb:latest-pg15   "docker-entrypoint.s…"   timescaledb          27 minutes ago   Up 27 minutes   0.0.0.0:5432->5432/tcp
  metrics-server-api      backend-metrics-server-api          "/bin/sh -c '/wait &…"   metrics-server-api   27 minutes ago   Up 4 seconds    0.0.0.0:8080->8080/tcp
  ```


### Build and run with Go

1. Compile binaries
  ```bash
  make
  ```

2. Run the `metrics-server-api`
  ```bash
  ./bin/metrics-server-api
  ```

3. You can run the simulator to inject metrics
  ```bash
  ./bin/simulator
  ```

#### Build and run tests
  ```bash
  make test
  ```

### Actions

For more details on the API spec, verify the live documentation with swagger.

#### Posting a metric

```
curl -i -X POST -H "Content-Type: application/json" --data "{\"timestamp\":\"2024-02-02T11:43:02.099Z\", \"name\":\"temperature\",\"value\":30}" localhost:8080/v1/metrics
```

#### Listing all metric names

```
curl -i -X GET localhost:8080/v1/metrics
```

#### Get metric's average per minute/hour/day

```
curl -i -X GET http://localhost:8080/v1/metrics/temperature/hour/2024-02-01/2024-02-03
```

