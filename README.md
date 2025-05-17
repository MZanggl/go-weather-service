# Weather API

## Prerequisites

- **Node.js** (for data ingestion)
- **Golang** (for API)
- **Docker/Docker Desktop** (for database dependency)

## Getting Started

### 1. Start Dependencies

Ensure Docker/Docker Desktop is running, then start the database:

```bash
docker compose up --build  # Use the -d flag to run in detached mode
```

### 2. Start the Application

```bash
cd api
go mod tidy # Only needed the first time
APP_ENV=development go run main.go # Or use go build
```

---

## Database Management

- The database is seeded automatically on the first container startup using `db/seed.sql`.

### Reseed the Database

1. Remove the volume:
    ```bash
    docker compose down -v
    ```
2. Restart the application.

### Access the Database

```bash
docker exec -it weather_postgres psql -h localhost -U weather_api weather
```

---

## Data Ingestion

To ingest weather data from the provided file, run:

```bash
API_HOST=http://localhost:8090 API_TOKEN=abcdef node ingestion/index.mjs
```

> **Note:** The `"date"` field is unique. You can add more weather data and run the ingestion script multiple times. Existing dates will be skipped without causing failures. To start with a fresh dataset, reseed the database (see above).

---

## API Usage

### Create Weather Records manually

This is what the data ingestion script also uses.

```bash
curl -H "X-Api-Token: abcdef" -X POST -H "Content-Type: application/json" \
-d '{"humidity":23.234234423, "temperature":57.234234423, "date": "2025-01-01"}' \
http://127.0.0.1:8090/weather
```

### Retrieve Weather Records for a Given Day

```bash
curl -X GET -H "Content-Type: application/json" \
http://127.0.0.1:8090/weather/2025-01-01
```

### Retrieve Weather Records for a Range of Days

```bash
curl -X GET -H "Content-Type: application/json" \
http://127.0.0.1:8090/weather/2025-01-01/2025-01-02
```

---

## WebSocket Usage

- Connect to:
  ```
  ws://127.0.0.1:8090/ws/<some user id>
  ```
- As records are added, you will be notified on this channel.
- You can connect from multiple clients.

The simplest websocket client is this CLI that you can run in your terminal:

```bash
npx wscat -c  ws://127.0.0.1:8090/ws/<some user id>
```

When using Postman, create a new WebSocket request (not Socket.IO).

---

## Running Tests

Run all tests using:

```bash
cd api
APP_ENV=test go test ./...
```

---

## Tools & Resources

- [Fiber](https://docs.gofiber.io/)
- [GORM](https://gorm.io/)
- [Postman](https://www.postman.com/downloads/)
- [Docker](https://www.docker.com/)
- [Postgres](https://www.postgresql.org/)

## Assignment Notes

- The `"date"` field is unique. In a real-world scenario, you might also want to include location as part of the unique key. That's also why the single day retrieval API returns an array of records.
