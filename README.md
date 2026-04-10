# Hookfy

Hookfy is a lightweight webhook inspector and testing tool built by [Scalefy](https://github.com/scalefy). It lets you receive webhooks on unique URLs and inspect their full contents — headers, body, query strings, and metadata — through a simple API.

Perfect for debugging webhook integrations during development without relying on external services.

## Features

- **Instant webhook capture** — send any POST request to `/webhooks/:hash` and it's stored automatically
- **Full request inspection** — captures method, headers, body, query parameters, content-type, and client IP
- **Hash-based inboxes** — group webhooks by hash to test multiple integrations simultaneously
- **Auto-expiration** — webhooks are automatically removed after 24 hours
- **Zero config** — runs with SQLite out of the box, no external database required
- **CORS enabled** — ready to be consumed by any frontend or external client

## Tech Stack

- [Go](https://go.dev/) 1.25+
- [Gin](https://github.com/gin-gonic/gin) — HTTP framework
- [GORM](https://gorm.io/) + SQLite — persistence
- [godotenv](https://github.com/joho/godotenv) — environment configuration

## Getting Started

### Prerequisites

- Go 1.25 or later

### Installation

```bash
git clone https://github.com/scalefy/hookfy.git
cd hookfy
go mod download
```

### Configuration

Create a `.env` file (or copy from the example):

```bash
cp .env.example .env
```

Available variables:

| Variable  | Description                  | Default              |
|-----------|------------------------------|----------------------|
| `DB_PATH` | Path to the SQLite database  | `./db/hookfy.sqlite` |

### Running

```bash
go run main.go
```

The server starts on `http://localhost:8081`.

## API

### Receive a Webhook

```
POST /webhooks/:hash
```

Send any payload to this endpoint. The `:hash` acts as an inbox identifier — use any string you want.

**Example:**

```bash
curl -X POST http://localhost:8081/webhooks/my-test-123 \
  -H "Content-Type: application/json" \
  -d '{"amount": 4000, "status": "paid"}'
```

**Response:**

```json
{
  "message": "webhook received",
  "id": 1
}
```

### Retrieve Webhooks

```
GET /webhooks/inbox
```

| Parameter | Type   | Description                          |
|-----------|--------|--------------------------------------|
| `hash`    | string | Filter by inbox hash (optional)      |
| `type`    | string | Response format: `json` or `html` (default: `json`) |

**Example:**

```bash
curl http://localhost:8081/webhooks/inbox?hash=my-test-123
```

**Response:**

```json
{
  "data": [
    {
      "ID": 1,
      "Hash": "my-test-123",
      "Method": "POST",
      "Headers": { "Content-Type": "application/json" },
      "Body": { "raw": "{\"amount\":4000,\"status\":\"paid\"}" },
      "QueryString": {},
      "ContentType": "application/json",
      "RemoteAddr": "127.0.0.1",
      "CreatedAt": "2026-04-10T12:00:00Z",
      "ExpiresAt": "2026-04-11T12:00:00Z"
    }
  ],
  "total": 1
}
```

## Project Structure

```
hookfy/
├── config/
│   └── database.go              # Database connection and setup
├── handlers/
│   └── webhook_handler.go       # API request handlers
├── models/
│   └── webhook.go               # Webhook data model
├── worker/
│   └── delete_expired_webhook.go # Expired webhook cleanup (WIP)
├── tests/                       # Bruno API test collection
├── main.go                      # Entry point
├── .env.example                 # Environment variables template
├── go.mod
└── go.sum
```

## Testing

The project includes API tests using [Bruno](https://www.usebruno.com/). Open the `tests/` folder as a Bruno collection to run them.

## Contributing

Contributions are welcome! Feel free to open issues and pull requests.

## License

This project is maintained by [Scalefy](https://github.com/scalefy).
