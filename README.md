# Travel Sync API

Find your travel buddy. This service exposes RESTful APIs for authentication, users, and travel tickets.

## Quick Start

1. Create a `.env` with required variables:
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

## Global Middleware

- CORS: allows localhost origins, credentials enabled.
- Rate limiting: general limiter applied globally.

## Authentication

Google OAuth2 for login; a JWT is set as an HTTP-only cookie named `jwt_token`. Protected routes require this cookie. See `SECURITY_README.md` for details.

## Endpoints

### Health
- GET `/health` → Service status

### Auth (public unless noted)
- GET `/auth/google/login` → Redirect to Google OAuth
- GET `/auth/google/callback` → OAuth callback; sets `jwt_token` cookie
- POST `/auth/logout` → Clears auth cookie
- GET `/auth/me` (protected) → Current user from JWT

### Users (all protected; base `/api/user`)
- GET `/api/user` → List users
- GET `/api/user/:id` → Get user by id
- PUT `/api/user/:id` → Update user
- DELETE `/api/user/:id` → Delete user

### Travel Tickets (all protected; base `/api/travel`)
- POST `/api/travel` → Create ticket
- GET `/api/travel` → List tickets
- GET `/api/travel/:id` → Get ticket by id
- PUT `/api/travel/:id` → Update ticket
- DELETE `/api/travel/:id` → Delete ticket
- GET `/api/travel/:id/recommendations` → Recommendations (stricter rate limit)
- GET `/api/travel/user-responses` → Current user’s responses

## Example Usage

Login flow (browser):
```bash
open http://localhost:8080/auth/google/login
```

Get current user (after login):
```bash
curl -i http://localhost:8080/auth/me
```

Create travel ticket:
```bash
curl -i -X POST http://localhost:8080/api/travel \
  -H 'Content-Type: application/json' \
  --cookie "jwt_token=YOUR_JWT" \
  -d '{"title":"Trip to Goa","from":"BLR","to":"GOI"}'
```

## Rate Limiting

- Global general limiter for all routes.
- Stricter limiter on `/auth/*`.
- Stricter limiter for `/api/travel/:id/recommendations`.

See `RATE_LIMITING.md` for configuration.
