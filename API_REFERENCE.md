# Travel Sync API Reference

Comprehensive documentation for frontend developers. All protected endpoints require the HTTP-only cookie `jwt_token` (set by the auth flow).

Base URL examples:
- Local: `http://localhost:8080`

Auth uses Google OAuth2 and JWT cookie. See `SECURITY_README.md`.

## Conventions

- Content-Type: `application/json`
- Success envelope: `{ "success": true, "data": ... }`
- Error envelope: `{ "success": false, "error": "message" }` or `{ "error": "message" }` for some auth endpoints

---

## Health

GET `/health`
- Auth: none
- Request: none
- Response 200:
```json
{ "status": "healthy", "message": "Travel Sync API is running" }
```

---

## Authentication

### Start Google Login
GET `/auth/google/login`
- Auth: none
- Behavior: redirects to Google OAuth (302/307). No JSON body.

### OAuth Callback
GET `/auth/google/callback?code=...&state=...`
- Auth: none
- Behavior: validates state, exchanges code, sets `jwt_token` cookie (HTTP-only). Redirects to `FRONTEND_URL` (no JSON body).
- Errors 400/500: `{ "error": "invalid oauth state" }` or `{ "error": "..." }`

### Logout
POST `/auth/logout`
- Auth: none
- Behavior: clears `jwt_token` cookie.
- Response 200:
```json
{ "message": "Logged out successfully" }
```

### Current User (Protected)
GET `/auth/me`
- Auth: JWT cookie required
- Response 200:
```json
{ "user_id": 123, "email": "user@example.com" }
```
- Response 401:
```json
{ "error": "User not authenticated" }
```

---

## Users (Protected)

Base: `/api/user`

### Get All Users
GET `/api/user`
- Response 200:
```json
{ "success": true, "data": [
  { "id": 1, "name": "Alice", "email": "alice@example.com", "batch": "2025", "phone_number": "9876543210", "created_at": "2025-09-01T10:00:00Z", "updated_at": "2025-09-02T10:00:00Z" }
] }
```
- Errors 500:
```json
{ "success": false, "error": "Failed to get all users" }
```

### Get User By ID
GET `/api/user/:id`
- Params: `id` (int)
- Response 200:
```json
{ "success": true, "data": { "id": 1, "name": "Alice", "email": "alice@example.com", "batch": "2025", "phone_number": "9876543210", "created_at": "2025-09-01T10:00:00Z", "updated_at": "2025-09-02T10:00:00Z" } }
```
- Errors 400/404:
```json
{ "success": false, "error": "Id parsing failed" }
```
```json
{ "success": false, "error": "User not found" }
```

### Update User
PUT `/api/user/:id`
- Params: `id` (int)
- Body:
```json
{ "name": "Alice B", "phone_number": "9998887777" }
```
- Response 200:
```json
{ "success": true, "data": { "id": 1, "name": "Alice B", "email": "alice@example.com", "batch": "2025", "phone_number": "9998887777", "created_at": "2025-09-01T10:00:00Z", "updated_at": "2025-10-01T10:00:00Z" } }
```
- Errors 400/500:
```json
{ "success": false, "error": "Invalid request body" }
```
```json
{ "success": false, "error": "..." }
```

### Delete User
DELETE `/api/user/:id`
- Params: `id` (int)
- Response 200:
```json
{ "success": true, "data": "User deleted successfully" }
```
- Errors 400/500:
```json
{ "success": false, "error": "Id parsing failed" }
```
```json
{ "success": false, "error": "Failed to delete user" }
```

---

## Travel Tickets (Protected)

Base: `/api/travel`

### Create Ticket
POST `/api/travel`
- Body:
```json
{
  "source": "BLR",
  "destination": "GOI",
  "departure_at": "2025-10-01T14:30:00Z",
  "time_diff_mins": 30,
  "empty_seats": 2,
  "phone_number": "9876543210"
}
```
- Response 201:
```json
{ "success": true, "data": {
  "id": 10,
  "source": "BLR",
  "destination": "GOI",
  "empty_seats": 2,
  "departure_at": "2025-10-01T14:30:00Z",
  "time_diff_mins": 30,
  "user_id": 123,
  "phone_number": "9876543210",
  "status": "open",
  "created_at": "2025-10-01T10:00:00Z",
  "updated_at": "2025-10-01T10:00:00Z"
} }
```
- Errors 400/401:
```json
{ "success": false, "error": "unauthorized" }
```
```json
{ "success": false, "error": "invalid request body" }
```

### List Tickets
GET `/api/travel`
- Response 200:
```json
{ "success": true, "data": [
  { "id": 10, "source": "BLR", "destination": "GOI", "empty_seats": 2, "departure_at": "2025-10-01T14:30:00Z", "time_diff_mins": 30, "user_id": 123, "phone_number": "9876543210", "created_at": "2025-10-01T10:00:00Z", "updated_at": "2025-10-01T10:00:00Z" }
] }
```
- Errors 500:
```json
{ "success": false, "error": "failed to fetch tickets" }
```

### Get Ticket By ID
GET `/api/travel/:id`
- Params: `id` (int)
- Response 200:
```json
{ "success": true, "data": { "id": 10, "source": "BLR", "destination": "GOI", "empty_seats": 2, "departure_at": "2025-10-01T14:30:00Z", "time_diff_mins": 30, "user_id": 123, "phone_number": "9876543210", "created_at": "2025-10-01T10:00:00Z", "updated_at": "2025-10-01T10:00:00Z" } }
```
- Errors 400/404:
```json
{ "success": false, "error": "invalid id" }
```
```json
{ "success": false, "error": "ticket not found" }
```

### Update Ticket
PUT `/api/travel/:id`
- Params: `id` (int)
- Body (all optional fields):
```json
{
  "source": "BLR",
  "destination": "GOI",
  "departure_at": "2025-10-02T14:30:00Z",
  "time_diff_mins": 45,
  "empty_seats": 3,
  "phone_number": "9876543210",
  "status": "closed"
}
```
- Response 200:
```json
{ "success": true, "data": { "id": 10, "source": "BLR", "destination": "GOI", "empty_seats": 3, "departure_at": "2025-10-02T14:30:00Z", "time_diff_mins": 45, "user_id": 123, "phone_number": "9876543210", "status": "closed", "created_at": "2025-10-01T10:00:00Z", "updated_at": "2025-10-02T10:00:00Z" } }
```
- Errors 400/500:
```json
{ "success": false, "error": "invalid request body" }
```
```json
{ "success": false, "error": "..." }
```

### Delete Ticket
DELETE `/api/travel/:id`
- Params: `id` (int)
- Response 200:
```json
{ "success": true, "data": "ticket deleted" }
```
- Errors 500:
```json
{ "success": false, "error": "failed to delete ticket" }
```

### Get Recommendations (Rate Limited)
GET `/api/travel/:id/recommendations`
- Params: `id` (int)
- Behavior: only considers tickets with `status = "open"` as candidates
- Response 200:
```json
{ "success": true, "data": {
  "best_match": { "ticket": { "id": 11, "source": "BLR", "destination": "GOI", "empty_seats": 1, "departure_at": "2025-10-01T16:00:00Z", "time_diff_mins": 15, "user_id": 999, "phone_number": "9876543211", "created_at": "2025-10-01T10:00:00Z", "updated_at": "2025-10-01T10:00:00Z" }, "score": 0.92, "date": "2025-10-01", "time": "16:00" },
  "best_group": [ { "ticket": { "id": 12, "source": "BLR", "destination": "GOI", "empty_seats": 2, "departure_at": "2025-10-01T16:15:00Z", "time_diff_mins": 20, "user_id": 1000, "phone_number": "9876543212", "created_at": "2025-10-01T10:00:00Z", "updated_at": "2025-10-01T10:00:00Z" }, "score": 0.85, "date": "2025-10-01", "time": "16:15" } ],
  "other_alternatives": []
} }
```
- Errors 429:
```json
{ "success": false, "error": "Rate limit exceeded. Please try again later.", "retry_after": 1696166400 }
```

### Get Current User Responses
GET `/api/travel/user-responses`
- Response 200:
```json
{ "success": true, "data": [
  { "id": 10, "student_name": "Alice", "student_batch": "2025", "source": "BLR", "destination": "GOI", "date": "2025-10-01", "time": "14:30", "empty_seats": 2, "phone_number": "9876543210" }
] }
```
- Errors 401/500:
```json
{ "success": false, "error": "unauthorized" }
```
```json
{ "success": false, "error": "..." }
```

---

## Rate Limiting

Responses may include headers:
- `X-RateLimit-Limit`
- `X-RateLimit-Remaining`
- `X-RateLimit-Reset`

Errors when exceeded: HTTP 429 with `{ "success": false, "error": "Rate limit exceeded. Please try again later.", "retry_after": <unix_ts> }`.

See `RATE_LIMITING.md` for details.


