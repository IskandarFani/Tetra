# Tetra

Tetra is a small Go backend service for collecting trial access requests.
It exposes a Fiber HTTP API, stores requests in PostgreSQL through GORM, and keeps the code split into handler, service, repository, database, and model layers.

## Stack

- Go
- Fiber
- GORM
- PostgreSQL
- godotenv

## Project Structure

```text
.
|-- cmd/migrate      # database migration entrypoint
|-- database         # PostgreSQL connection setup
|-- handler          # HTTP handlers
|-- models           # GORM models
|-- repository       # database operations
|-- services         # business logic
|-- main.go          # application entrypoint
`-- routers.go       # Fiber routes and middleware
```

## Configuration

Create a `.env` file in the project root:

```env
DB_HOST=<database-host>
DB_PORT=<database-port>
DB_USER=<database-user>
DB_PASSWORD=<database-password>
DB_NAME=<database-name>
DB_SSLMODE=disable
```

The `.env` file is ignored by Git. Keep real credentials only in `.env` or system environment variables.

## Install Dependencies

```bash
go mod download
```

## Run Database Migration

```bash
go run ./cmd/migrate
```

The migration creates the `trial_requests` table from `models.TrialRequest`.

## Run The API

```bash
go run .
```

The API listens on `http://localhost:8080`.

## API

### Health

```http
GET /health
```

Liveness endpoint provided by Fiber healthcheck middleware.

### Database Check

```http
GET /checkdb
```

Checks whether the application can ping PostgreSQL.

Successful response:

```json
{
  "status": "success",
  "message": "database connection is healthy"
}
```

### Trial Access Request

```http
POST /trial-access
Content-Type: application/json
```

Request body:

```json
{
  "email": "user@example.com"
}
```

Successful response:

```json
{
  "status": "success",
  "message": "A request has been sent to user@example.com to confirm your email address"
}
```

The email is stored in PostgreSQL with a unique index.

### Fiber Monitor

```http
GET /monitor
```

Shows the built-in Fiber monitor page.

## Notes

- CORS allows `http://localhost:5173` and `http://127.0.0.1:5173` for a local frontend.
- The backend uses port `8080`; port `3000` is reserved for frontend usage.
