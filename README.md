# url-encoder

Simple, pluggable URL encoding/decoding service for URL-shortening-style systems.

It exposes HTTP endpoints and uses the Strategy Pattern to generate short codes. The current implementation ships with a Base36 encoder and these ID strategies:

- `random`: generates a random uint64 and Base36-encodes it
- `sequential`: increments an in-memory counter and Base36-encodes it
- `tenant`: deterministic per `(tenant_id, original_url)` using a SHA-256 hash, encoded with a caller-provided minimum length

Note: `length` is only used by the `tenant` strategy.

## API

- `POST /encoder/api/v1/encode`
  - Request JSON:
    - `original_url` (required, URL)
    - `strategy` (required: `random | sequential | tenant`)
    - `tenant_id` (optional if `X-Tenant-ID` header is set)
    - `length` (optional, integer 4..10; only used for `tenant`)
  - Response JSON:
    - `short_url`
    - `original_url`
    - `tenant_id`

- `GET /encoder/api/v1/resolve/:short_url`
  - Response JSON:
    - `short_url`
    - `original_url`

## Configuration

Environment variables:

- `PORT` (default: `8081`)
- `DB_DRIVER` (default: `sqlite`; supported: `sqlite`, `postgres`)
- `DATABASE_URL` (default: `file:encodeurl.db?cache=shared&mode=rwc`)
- `CACHE_TTL` (default: `30m`)

## Examples

Encode:
```bash
curl -s -X POST "http://localhost:8081/encoder/api/v1/encode" \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: tenant-a" \
  -d '{"original_url":"https://example.com","strategy":"tenant","length":6}'
```

DB-backed sequential (restart-safe):
```bash
curl -s -X POST "http://localhost:8081/encoder/api/v1/encode" \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: tenant-a" \
  -d '{"original_url":"https://example.com","strategy":"sequential_db"}'
```

Resolve:
```bash
curl -s "http://localhost:8081/encoder/api/v1/resolve/{short_url}"
```
