# Weather API 🌞

<img src="https://mermaid.ink/svg/pako:eNptU1tv2jAU_ivWeeokaLisAfxQaUDXddJW1lBVmngxziH1IHZ27KjtEP99DiZFhPohir_LOedz4i1IkyJwsPi3RC1xqkRGIl9o5lchyCmpCqEde7RI5-idztA6ZXQiSRWOCcumwokjzgJx7vwyu6vUTyjcM9J-e3Froq9qifTpXD4d17WXwiK7mBnrMkL7gfQJl4mRa3STjULtQpcDxAIWTOHZkLevr_0snE2M1ii9Mqiq-J5q5OXsodRM1WCkdIqvl_kfG0wNdV16dp_MWfQSkgelxz07HXPv8Z0cO7CMUBpKg2g6btclklJKtDa6ITInJRppOBuTEakU1h1VH-X4Np_PTqtGyVoVVfyTA9h3v705zh_xFZmcGWogEXemme1XifTWiGbPsj2gK0mf0oexqyE4-57c__S0LYw_q2o-aEFGKgXuqMQW5Ei5qLawrdwL8O1yXAD3r6mg9QIWeuc9_m_5bUxe28iU2TPwldhYvyuLVLj6OryjhP4L08SU2gHvx919EeBbeAXe7lx2qtUbXcWjQdyNO8PuoPd5eNWCN-A12R8MOvGoO-qPenGvO9y14N9-Bl1uNi3AVDlDP8KdlEavVAa7_xyFMCM">

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

Connect to:
  ```
  ws://127.0.0.1:8090/ws/<some user id>
  ```
As records are added, you will be notified on this channel (you can connect from multiple clients).

The simplest websocket client is [wscat](https://github.com/websockets/wscat) that you can run from your terminal:

```bash
npx wscat -c  ws://127.0.0.1:8090/ws/<some user id>
```

When using Postman, make sure to create a WebSocket request, not a Socket.IO request.

---

## Running Tests

Run both unit and integrations tests using:

```bash
cd api
APP_ENV=test go test ./...
```

---

## Folder Structure

```
.
├── api/
│   ├── configs/         # Configuration files and environment variable loaders
│   ├── handlers/        # HTTP route handlers
│   ├── models/          # Database models
│   ├── server/          # Server setup (DB, websocket, etc.)
│   ├── services/        # Business logic and data access
│   ├── utils/           # Utility functions (date, number formatting, etc.)
│   ├── main.go          # Application entry point
│   ├── main_test.go     # API Integration tests
│   ├── go.mod           # Go module definition
│   └── go.sum           # Go module checksums
├── db/
│   ├── Dockerfile       # Database Dockerfile
│   └── seed.sql         # SQL for initial schema and seed data
├── data/
│   └── weather.dat      # Weather data for ingestion
├── ingestion/
│   └── index.mjs        # Node.js script for data ingestion
├── scripts/
│   └── wait-for-port.sh # Helper script for Docker healthcheck
├── docker-compose.yml   # Docker Compose setup for services
```

---

## Tools & Resources

- [Golang](https://go.dev/)
- [Node.js](https://nodejs.org/en)
- [Fiber Framework](https://docs.gofiber.io/)
- [gorm ORM](https://gorm.io/)
- [Postman](https://www.postman.com/downloads/)
- [Docker](https://www.docker.com/)
- [Postgres](https://www.postgresql.org/)

## Assignment Notes

- The `"date"` field is unique. In a real-world scenario, you might also want to include location as part of the unique key. That's also why the single day retrieval API returns an array of records.
