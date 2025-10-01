# Security Implementation Guide

This document explains the security implementation for the Travel Sync application.

## Overview

- Google OAuth2 for authentication
- JWT for session management (cookie `jwt_token`, HTTP-only)
- JWT middleware injects `user_id`, `user_email`, and `jwt_claims` into context
- CORS enabled for localhost development origins with credentials

## Auth Flow

1. User visits `/auth/google/login` → redirected to Google with a CSRF-protecting `state` cookie
2. Callback `/auth/google/callback` exchanges code, issues JWT, sets `jwt_token` (HTTP-only), and redirects to `FRONTEND_URL`
3. Client sends cookie automatically on subsequent requests

## JWT Details

- Algorithm: HS256
- Claims: `user_id`, `email`, standard registered claims
- Expiration: 8 days
- Middleware: `internal/security/config/jwtMiddleware.go`

### Middleware behavior
- Reads token from `jwt_token` cookie
- Validates signature using `JWT_SECRET`
- On success: sets `user_id`, `user_email`, `jwt_claims` in Gin context
- On failure: 401 JSON error and abort

## Environment Variables

```bash
# Google OAuth Configuration
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret

# JWT Configuration
JWT_SECRET=your_jwt_secret_key

# Frontend URL (for redirects)
FRONTEND_URL=http://localhost:3000
```

## Endpoints

Public:
- GET `/health`
- GET `/auth/google/login`
- GET `/auth/google/callback`
- POST `/auth/logout`

Protected (JWT required):
- GET `/auth/me`
- All `/api/user/*`
- All `/api/travel/*`

## CORS

`internal/middleware/cors.go` allows localhost origins (`3000`, `8080`, and `127.0.0.1` variants), methods GET/POST/PUT/DELETE/OPTIONS, common headers, and credentials.

## Production Hardening

- Serve over HTTPS and set cookie `Secure=true`
- Rotate `JWT_SECRET` using a proper secret manager
- Restrict OAuth redirect URIs to trusted origins only
- Validate inputs and enforce payload size limits
- Minimize PII in logs; never log tokens

## Examples

Protect routes:
```go
protected := router.Group("/api")
protected.Use(config.JWTMiddleware(jwtService))
{
    protected.GET("/protected-route", handler.ProtectedHandler)
}
```

Access current user in handlers:
```go
func ProtectedHandler(c *gin.Context) {
    userID, _ := c.Get("user_id")
    userEmail, _ := c.Get("user_email")
    c.JSON(200, gin.H{"user_id": userID, "email": userEmail})
}
```
