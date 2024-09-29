# Golang Assessment Project

This project contains a producer and consumer service built using Go, which processes tasks and tracks metrics using Prometheus and Grafana.

## Project Overview

- **Producer Service**: Generates random tasks and sends them to the consumer.
- **Consumer Service**: Receives tasks, processes them, and stores the results in a SQLite database. It also tracks task metrics using Prometheus.

## Prerequisites

- Docker
- Docker Compose
- Go (optional for local development)

## Setup Instructions

1. **Clone the repository**:
    ```bash
    git clone https://github.com/your-repo/golang-assessment.git
    cd golang-assessment
    ```

2. **Build and run the services**:
    ```bash
    docker-compose up --build
    ```

3. **Access Grafana**:
    - Grafana is available at `http://localhost:3000`.
    - Default credentials: `admin/admin`.
    - Add Prometheus as a data source (`http://prometheus:9090`) and import the provided dashboards.

## Prometheus Metrics

The consumer exposes the following Prometheus metrics:

- `tasks_state_count`: Number of tasks in each state (`received`, `processing`, `done`).
- `tasks_processed_total`: Total number of tasks processed by type.

## Profiling

The application supports CPU and memory profiling using `pprof`. To enable profiling:

1. Start the consumer service.
2. Access the profiling interface at `http://localhost:6062/debug/pprof/`.
3. Use `go tool pprof` to analyze the profile data.

```bash
go tool pprof http://localhost:6062/debug/pprof/profile
