# Travel Sync

Find your travel buddy. Travel Sync lets students post upcoming trips and discover compatible co-travelers based on time, origin/destination, and preferences. It uses Google OAuth2 for login, JWT cookies for sessions, and a recommendation engine with rate limiting.

## What this app does

- Create a travel ticket with origin, destination, departure time, empty seats, and contact.
- Discover recommendations to form a group. Tickets “close” once a group is formed, so they’re no longer suggested.
- Manage your tickets and profile. You can only modify your own data.

## How it works (flow)

1. User clicks Login → `GET /auth/google/login` (Google OAuth).
2. Callback sets an HTTP-only cookie `jwt_token` and redirects to the frontend.
3. User creates a ticket → `POST /api/travel` (status defaults to `open`).
4. Recommendation engine finds matching tickets:
   - Excludes `closed` tickets and your own tickets
   - Asymmetric time window: before within your `time_diff_mins`, after within 60 minutes
   - Returns scored results with minimal owner info (name, batch), without revealing ticket/user IDs
5. When a group is formed, ticket owner sets `status` to `closed` to stop new requests.

## Features

- Google OAuth2 + JWT cookie auth
- Tickets CRUD with ownership checks (only owners can update/delete)
- Recommendation engine with asymmetric time window and redacted result fields
- CORS and rate limiting (global, auth-specific, recommendations-specific)

## Architecture

- Gin HTTP server, layered into routes → handlers → services → repositories
- PostgreSQL via GORM (auto-migrations)
- Middlewares: CORS, JWT auth, rate limiting

## Tech stack

- Go, Gin, GORM, Google OAuth2, golang-jwt, ulule/limiter

## Setup

1. Create a `.env` file:
```bash
PORT=8080
DATABASE_URL=postgres://user:pass@localhost:5432/travel_sync?sslmode=disable
JWT_SECRET=your_jwt_secret_key
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
FRONTEND_URL=http://localhost:3000
```
2. Run the server:
```bash
go run ./cmd
```
3. Health check:
```bash
curl http://localhost:8080/health
```

## Developer docs

- API reference: see `API_REFERENCE.md`
- Security details (OAuth, JWT, CORS): see `SECURITY_README.md`
- Rate limiting details: see `RATE_LIMITING.md`

## Notes

- New tickets are created with `status: open`; set to `closed` to stop recommendations.
- Recommendation results: redact `ticket.id` and `ticket.user_id`; include minimal user `{name, batch}`.
