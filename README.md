# Distributed Rate Limiter

A Go/Gin rate limiter service that registers endpoint-specific rate-limit configuration, checks incoming requests, and forwards allowed requests to backend services through a reverse proxy.

Current implementation is local/in-memory. It is structured so DB-backed endpoint config and Redis-backed distributed counters can be added later.

## Features

- Register endpoint-specific rate-limit rules.
- Supports strategy selection by name.
- Current strategies:
  - `per_unit_time`
  - `leaky_bucket` placeholder
- Tracks per-unit-time request windows in memory.
- Reads backend service targets from `services.yaml`.
- Forwards allowed requests with `httputil.NewSingleHostReverseProxy`.
- Logs requests and application events with Zap.
- Includes Postman collection: `rateLimiter.collection.json`.
- Includes VS Code debugger config in `.vscode/launch.json`.

## Project structure

```text
.
├── config/                 # strategy config and services.yaml loading
├── handler/                # endpoint registration and rate-limit handling
├── limiter/                # RequestLimiter interface and strategy registry
├── limiter/strategy/       # concrete limiter strategies
├── proxy/                  # reverse proxy forwarding
├── zapLogger/              # Zap logger setup and middleware
├── docs/                   # ERD docs by feature
├── services.yaml           # backend service mapping
└── rateLimiter.collection.json
```

## Run locally

```bash
go run main.go
```

Server starts on:

```text
http://localhost:8080
```

If port `8080` is already in use:

```bash
lsof -i :8080
kill <PID>
```

## Configuration

Backend service mapping is loaded from `services.yaml`:

```yaml
services:
  user: "http://localhost:8081"
```

The service URL must include the scheme, for example `http://`.

## API

### Health check

```bash
curl http://localhost:8080/health
```

Response:

```json
{
  "message": "OK"
}
```

### Register endpoint

```bash
curl --location 'http://localhost:8080/register-endpoint' \
  --header 'Content-Type: application/json' \
  --data '{
    "path": "/ww",
    "method": "GET",
    "target_service": "user",
    "strategy": {
      "name": "per_unit_time",
      "limit": 10,
      "interval": 60
    }
  }'
```

Fields:

| Field | Meaning |
|---|---|
| `path` | Public path handled by the rate limiter |
| `method` | HTTP method used to build the registry key |
| `target_service` | Service name from `services.yaml` |
| `strategy.name` | Limiter strategy, e.g. `per_unit_time` |
| `strategy.limit` | Max requests allowed in the interval |
| `strategy.interval` | Window size in seconds |

Current registry key:

```text
METHOD:/path
```

Example:

```text
GET:/ww
```

### Call registered endpoint

The current per-unit-time limiter reads `client_id` from the query string:

```bash
curl 'http://localhost:8080/ww?client_id=test-user'
```

If the endpoint is registered and the request is allowed, the request is forwarded to the target service configured in `services.yaml`.

If the endpoint is not registered:

```json
{
  "error": "Endpoint not registered"
}
```

If the rate limit is exceeded:

```json
{
  "error": "Rate limit exceeded"
}
```

## Runtime flow

```text
request
  ↓
Gin route / NoRoute
  ↓
HandleLimit
  ↓
lookup METHOD:/path in in-memory registry
  ↓
call strategy Allow(...)
  ↓
if allowed, resolve target_service from services.yaml map
  ↓
forward request to backend service
```

## Current data model

Endpoint config is stored in memory:

```go
LimitRegistery map[string]limitConfig
```

Strategy instances are cached by strategy name:

```go
strategyRegistry map[string]RequestLimiter
```

Per-unit-time request windows are stored in memory:

```go
requestWindow map[string]requestWindow
```

Current request-window key:

```text
METHOD:/path:client_id
```

## ERD docs

- [Endpoint Registration](docs/erd-endpoint-registration.md)
- [Strategy Registry](docs/erd-strategy-registry.md)
- [Runtime Rate Limiting](docs/erd-rate-limit-runtime.md)
- [Service Routing and Proxying](docs/erd-service-routing.md)
- [Logging and Observability](docs/erd-logging-observability.md)

## Debugging

VS Code debugger config exists at:

```text
.vscode/launch.json
```

Use:

```text
Debug API Server
```

Then trigger a request:

```bash
curl 'http://localhost:8080/ww?client_id=test-user'
```

## Postman

Import:

```text
rateLimiter.collection.json
```

It contains requests for:

- health check
- endpoint registration
- calling a registered endpoint

## Current limitations

- Endpoint registration is in memory only.
- Rate-limit counters are in memory only.
- Restarting the app clears registrations and counters.
- Multiple service instances will not share counters.
- `leaky_bucket` is currently a placeholder.
- There is no DB persistence yet.
- There is no Redis-backed distributed counter yet.

## Recommended next steps

1. Persist endpoint registrations in a DB.
2. Load endpoint configs into an in-memory cache at startup.
3. Use Redis for distributed request counters.
4. Add cache invalidation when endpoint config changes.
5. Implement the leaky bucket strategy.
6. Add tests for registration, missing endpoint, allowed request, and rejected request.
