# Metrics Server

This project is a metrics server composed by the `metrics-server-api` (backend) and the `metrics-server-ui` (frontend).
The `metrics-server-api` exposes the API to ingest metrics and to consume aggregated metrics per granularity.
The `metrics-server-ui` shows a chart of the aggregated data.
Bonus: there is a `simulator` (CLI tool) to send multiple metrics concurrently (20 metrics every 300 ms).

## Specification #3 (better for senior)

> We need a Frontend + Backend application that allows you to post and visualize metrics.
> Each metric will have: Timestamp, name, and value. The metrics will be shown in a timeline and must show averages per minute/hour/day The metrics will be persisted in the database.

## Architecture

![architecture](docs/metrics-server-architecture.drawio.svg)

## Project Structure

.
├── README.md
├── backend
├── docker-compose.yml
├── docs
└── metrics-server-ui

## Documentation

You can find the documentation of each project (backend and frontend) on its own README.md files:
- [`metrics-server-ui` (frontend)](metrics-server-ui/README.md)
- [`metrics-server-api` (backend)](backend/README.md)

## Getting Started

### Build and run with docker

You can run the frontend, backend and the database using Docker Compose

  ```bash
  docker-compose up -d
  ```

  Check if services are up and running
  ```bash
  docker-compose ps
  ```

  You should be able see something like:

  ```
  NAME                             IMAGE                                 COMMAND                  SERVICE              CREATED         STATUS         PORTS
  metrics-server-api               simple_analytics-metrics-server-api   "/bin/sh -c '/wait &…"   metrics-server-api   2 minutes ago   Up 4 seconds   3000/tcp, 0.0.0.0:8080->8080/tcp
  metrics-server-ui                simple_analytics-metrics-server-ui    "docker-entrypoint.s…"   metrics-server-ui    2 minutes ago   Up 2 minutes   0.0.0.0:3000->3000/tcp
  simple_analytics-timescaledb-1   timescale/timescaledb:latest-pg15     "docker-entrypoint.s…"   timescaledb          2 minutes ago   Up 2 minutes   0.0.0.0:5432->5432/tcp
  ```

  You can post a metric:
  ```
  curl -i -X POST -H "Content-Type: application/json" --data "{\"timestamp\":\"2024-02-02T11:43:02.099Z\", \"name\":\"temperature\",\"value\":30}" http://localhost:8080/v1/metrics
  ```

  Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

## Architecture decisions

- As it's a POC, I avoided to introduce a Pub/Sub or Queue system, as Golang could handle efficiently 67 requests per second or 4020 RPMs (using the simulator: 20 Request * 300 ms).
- Golang: I dediced to use Golang as I'm familiar with the programming language and I believe it's a good choice for this kind of challenge in the backend.
- NextJS/Typescript: I decided to use NextJS as it's simple and effective for building simple frontend sites.
- TimescaleDB: I decided to use timescaleDB as it is built on top of PostgreSQL and supports timeseries queries.

### Next steps if requires more scalability
- It can be scaled horizontally and add an Pub/Sub or Queue to offload the database writting and process metrics asynchronously
- Add a readonly replica of the database. Split the current database in a writting and readonly databases to remove the competition between the reading and writing operations.

## Disclaimer
- I commited the `.env` file to facilitate to run this project. In real life, I would create a .env.example file instead.


## TODO
Here you can find the tasks that I planned for this project.

Unfortunately, I couldn't do all of them, but you can verify the ones that is missing.

Backend:
- [X] API to receive the metric (POST)
- [X] Database to store the metric
- [x] API to retrieve the aggregations per minute/hour/day (GET)
- [ ] API Swagger documentation
- [ ] CI/CD
- [ ] Deploy to a cloud provider
- [ ] Renomear pacote para metrics-server go.mod
- [ ] Renomear backend para metrics-server-api

Frontend:
- [x] Sparkline chart component
- [x] Filter by:
  - granularity: use for `minute`, `hour` and `day`
  - time frame: use for two different dates including times (hour/min/sec). E.g: "02-02-2024 16:50:00" and "02-02-2024 18:50:00"
- [ ] Component Tests
- [ ] CI/CD
- [ ] Deploy to a cloud provider


Producer Simulator (Bonus):
- [x] post a metrics every second or less


```
curl -i -X POST -H "Content-Type: application/json" --data "{\"timestamp\":\"2024-02-02T11:43:02.099Z\", \"name\":\"temperature\",\"value\":30}" localhost:8080/v1/metrics
```