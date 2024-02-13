# Metrics Server UI

The `metrics-server-ui` shows a chart of the aggregated data.

## Built with

- [Next.js](https://nextjs.org/)
- Typescript
- [Carbon Charts](https://charts.carbondesignsystem.com/)

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
  NAME                IMAGE                                 COMMAND                  SERVICE             CREATED          STATUS         PORTS
  metrics-server-ui   metrics-server-ui-metrics-server-ui   "docker-entrypoint.s…"   metrics-server-ui   31 seconds ago   Up 5 seconds   0.0.0.0:3000->3000/tcp
  ```

### Build and run with NPM

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.
