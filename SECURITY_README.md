# Security Implementation Guide

This document explains the security implementation for the Travel Sync application.

## Overview

The application uses Google OAuth2 for authentication with JWT tokens for session management. Only users with `sst.scaler.com` email domain are allowed to authenticate.

## Security Flow

1. **User Login**: User clicks "Login with Google" → redirected to Google OAuth
2. **Domain Validation**: Only `sst.scaler.com` emails are accepted
3. **User Creation**: First-time users are automatically created in the database
4. **JWT Generation**: JWT token is generated with user ID and email
5. **Cookie Storage**: JWT is stored in HTTP-only cookie
6. **Route Protection**: Protected routes require valid JWT token

## Environment Variables Required

```bash
# Google OAuth Configuration
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret

# JWT Configuration
JWT_SECRET=your_jwt_secret_key

# Frontend URL (for redirects)
FRONTEND_URL=http://localhost:3000
```

## API Endpoints

### Public Endpoints (No Authentication Required)

- `GET /health` - Health check
- `GET /auth/google/login` - Initiate Google OAuth login
- `GET /auth/google/callback` - Google OAuth callback
- `POST /auth/logout` - Logout (clears JWT cookie)

### Protected Endpoints (Authentication Required)

- `GET /auth/me` - Get current user information
- `GET /api/user/:id` - Get user by ID
- `PUT /api/user/:id` - Update user
- `DELETE /api/user/:id` - Delete user
- `GET /api/user/getAll/:id` - Get all users

## Usage Examples

### Protecting Routes

```go
// Apply JWT middleware to a route group
protected := router.Group("/api")
protected.Use(config.JWTMiddleware(jwtService))
{
    protected.GET("/protected-route", handler.ProtectedHandler)
}
```

### Getting Current User in Handlers

```go
func ProtectedHandler(c *gin.Context) {
    // Get user ID from context
    userID, err := config.GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }
    
    // Get user email from context
    email, err := config.GetUserEmailFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }
    
    // Use user information
    c.JSON(http.StatusOK, gin.H{
        "user_id": userID,
        "email": email,
    })
}
```

### Optional Authentication

```go
// Use OptionalJWTMiddleware for routes that work with or without authentication
optional := router.Group("/api")
optional.Use(config.OptionalJWTMiddleware(jwtService))
{
    optional.GET("/public-or-private", handler.OptionalAuthHandler)
}
```

## JWT Token Structure

```json
{
  "user_id": 123,
  "email": "user@sst.scaler.com",
  "sub": "user@sst.scaler.com",
  "exp": 1234567890,
  "iat": 1234567890
}
```

## Security Features

1. **Domain Restriction**: Only `sst.scaler.com` emails allowed
2. **HTTP-Only Cookies**: JWT tokens stored in HTTP-only cookies (not accessible via JavaScript)
3. **Token Expiration**: JWT tokens expire after 8 days
4. **State Parameter**: OAuth state parameter prevents CSRF attacks
5. **Secure Cookie Settings**: Cookies are configured for security

## Testing the Security

### 1. Test Health Check
```bash
curl http://localhost:8080/health
```

### 2. Test Login Flow
1. Visit `http://localhost:8080/auth/google/login`
2. Complete Google OAuth flow
3. Should redirect to frontend with JWT cookie set

### 3. Test Protected Endpoint
```bash
# This should fail without authentication
curl http://localhost:8080/auth/me

# This should work after login (cookie will be set)
curl -b "jwt_token=your_jwt_token" http://localhost:8080/auth/me
```

### 4. Test Logout
```bash
curl -X POST http://localhost:8080/auth/logout
```

## Troubleshooting

### Common Issues

1. **"JWT_SECRET not set"**: Set the JWT_SECRET environment variable
2. **"login with your Scaler Student Email"**: User must use `@sst.scaler.com` email
3. **"JWT token not found in cookies"**: User needs to login first
4. **"Invalid JWT token"**: Token may be expired or tampered with

### Debug Mode

To enable debug logging, set:
```bash
GIN_MODE=debug
```

## Security Best Practices

1. **Environment Variables**: Never commit secrets to version control
2. **HTTPS in Production**: Use HTTPS in production and set `secure=true` for cookies
3. **Token Rotation**: Consider implementing token refresh mechanism
4. **Rate Limiting**: Implement rate limiting for authentication endpoints
5. **Logging**: Log authentication attempts and failures
6. **CORS**: Configure CORS properly for your frontend domain

## File Structure

```
internal/security/
├── config/
│   ├── jwtMiddleware.go    # JWT authentication middleware
│   └── authHelpers.go      # Helper functions for auth
├── handler/
│   └── oauthHandler.go     # OAuth login/logout handlers
├── service/
│   ├── authService.go      # User authentication service
│   ├── jwtService.go       # JWT token service
│   └── customOAuth2Service.go # Google OAuth service
└── routes/
    └── authRoutes.go       # Authentication routes
```
