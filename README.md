# 🚀 Travel Sync

> Smart travel companion platform connecting students with compatible co-travelers through intelligent matching algorithms.

<div align="center">

[![Live Demo](https://img.shields.io/badge/Live-Demo-success?style=for-the-badge)](https://www.travelsync.space)
[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![JWT](https://img.shields.io/badge/JWT-black?style=for-the-badge&logo=JSON%20web%20tokens)](https://jwt.io/)

**[Live Application](https://www.travelsync.space)** | **[API Docs](./API_REFERENCE.md)** | **[Security Guide](./SECURITY_README.md)**

</div>

---

## 📖 Overview

Travel Sync helps students find compatible travel companions by intelligently matching their journey details, schedules, and preferences. Post your travel plans, discover matching co-travelers, and coordinate journeys seamlessly.

### ✨ Key Features

- 🔐 **Secure Google OAuth 2.0** authentication with JWT sessions
- 🎫 **Travel Ticket Management** - Create, update, and manage trip details
- 🤖 **Smart Recommendation Engine** - Implemented Greedy algorithm with scoring
- 🛡️ **Multi-tier Rate Limiting** - Protection against abuse
- 📱 **Real-time Status Updates** - Open/closed ticket management
- 🔒 **Privacy-First Design** - Redacted IDs, minimal data exposure

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────┐
│              Frontend (Client)                   │
│            HTTP + Cookies (JWT)                  │
└──────────────────┬──────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────┐
│         Middleware Layer                         │
│  ┌──────┐  ┌─────────┐  ┌──────────────────┐   │
│  │ CORS │  │   JWT   │  │  Rate Limiting   │   │
│  └──────┘  │  Auth   │  │ (3-tier system)  │   │
│            └─────────┘  └──────────────────┘   │
└──────────────────┬──────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────┐
│   Handlers → Services → Repositories            │
│                                                  │
│  ┌─────────────┐  ┌───────────────────────┐    │
│  │ Auth        │  │  Recommendation       │    │
│  │ Service     │  │  Engine (Scoring)     │    │
│  └─────────────┘  └───────────────────────┘    │
└──────────────────┬──────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────┐
│           PostgreSQL (GORM)                      │
│     Users | Travel Tickets | Sessions           │
└──────────────────────────────────────────────────┘
```

### Layered Design Benefits
- **Separation of Concerns** - Each layer handles specific responsibilities
- **Testability** - Mock dependencies for unit testing
- **Scalability** - Easy to extend and modify
- **Maintainability** - Clear code organization

---

## 💻 Tech Stack

| Category | Technology |
|----------|-----------|
| **Backend** | Go, Gin Framework |
| **Database** | PostgreSQL , GORM |
| **Auth** | Google OAuth 2.0, JWT (golang-jwt) |
| **Security** | HTTP-only Cookies, CORS, Rate Limiting |
| **Rate Limiting** | ulule/limiter (in-memory) |

---

## 🔒 Security Implementation

### OAuth 2.0 + JWT Flow

```mermaid
sequenceDiagram
    User->>+Server: GET /auth/google/login
    Server->>+Google: Redirect to OAuth
    Google->>+User: Authorization Page
    User->>+Google: Grant Permission
    Google->>+Server: Callback with code
    Server->>+Server: Verify State (CSRF)
    Server->>+Server: Generate JWT
    Server->>+User: Set HTTP-only Cookie
    User->>+Server: Authenticated Requests
```

### Security Features

✅ **HTTP-only Cookies** - Prevents XSS attacks  
✅ **CSRF Protection** - OAuth state parameter validation  
✅ **JWT Validation** - Signature verification on every request  
✅ **Rate Limiting** - Prevents brute force attacks  
✅ **Ownership Checks** - Users can only modify their own data  
✅ **Input Sanitization** - Validates all user inputs

**JWT Details:**
- Algorithm: HS256
- Expiration: 8 days
- Claims: `user_id`, `email`, standard claims
- Storage: HTTP-only cookie (`jwt_token`)

---

## ⚡ Rate Limiting Strategy

Three-tier protection system:

| Tier | Limit | Burst | Applied To |
|------|-------|-------|------------|
| **General** | 100/hour | 200 | All routes |
| **Auth** | 5/minute | 10 | Login/Auth endpoints |
| **Recommendations** | 20/minute | 30 | Recommendation engine |

**Key Strategy:**
- Authenticated: `prefix:uid:<user_id>`
- Unauthenticated: `prefix:<client_ip>`

**Response Headers:**
```http
X-RateLimit-Limit: 20
X-RateLimit-Remaining: 15
X-RateLimit-Reset: 1729684800
```

---

## 🎯 Recommendation Engine

Smart matching algorithm with multi-factor scoring:

### Matching Logic

**1. Route Matching**
- Exact source and destination match

**2. Asymmetric Time Window**
```
Target: 14:00, flexibility: ±30 mins

Acceptable Range:
├─ Before: 13:30 - 14:00 (configurable)
└─ After:  14:00 - 15:00 (fixed 60 mins)
```

**3. Status Filtering**
- Only `status: "open"` tickets
- Excludes user's own tickets

### Scoring Algorithm

```go
Score = (time_proximity × 0.5) + 
        (seat_availability × 0.3) + 
        (batch_similarity × 0.2)
```

**Result Tiers:**
- `best_match` - Highest scored recommendation
- `best_group` - Optimal group formation
- `other_alternatives` - Additional options

**Privacy Protection:** Redacts `ticket.id` and `user_id` in responses

---

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL 
- Google OAuth credentials

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/yourusername/travel-sync.git
cd travel-sync
```

2. **Configure environment**
```bash
cp .env.example .env
# Edit .env with your credentials
```

`.env` file:
```env
PORT=8080
DATABASE_URL=postgres://user:pass@localhost:5432/travel_sync?sslmode=disable
JWT_SECRET=your_secret_key_min_32_chars
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
FRONTEND_URL=http://localhost:3000
```

3. **Install dependencies**
```bash
go mod download
```

4. **Run the server**
```bash
go run ./cmd
```

5. **Health check**
```bash
curl http://localhost:8080/health
# Response: {"status":"healthy","message":"Travel Sync API is running"}
```

---

## 📚 API Endpoints

### Authentication
```http
GET  /auth/google/login          # Initiate Google OAuth
GET  /auth/google/callback       # OAuth callback
GET  /auth/me                    # Get current user (protected)
POST /auth/logout                # Clear session
```

### Travel Tickets
```http
POST   /api/travel                    # Create ticket
GET    /api/travel                    # List all tickets
GET    /api/travel/my                 # Get my tickets
GET    /api/travel/:id                # Get ticket by ID
PUT    /api/travel/:id                # Update ticket
DELETE /api/travel/:id                # Delete ticket
GET    /api/travel/:id/recommendations # Get matches (rate limited)
```

### Example: Create Ticket
```bash
curl -X POST http://localhost:8080/api/travel \
  -H "Content-Type: application/json" \
  -b "jwt_token=your_token" \
  -d '{
    "source": "BLR",
    "destination": "GOI",
    "departure_at": "2025-10-01T14:30:00Z",
    "time_diff_mins": 30,
    "empty_seats": 2,
    "phone_number": "9876543210"
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 10,
    "source": "BLR",
    "destination": "GOI",
    "status": "open",
    "departure_at": "2025-10-01T14:30:00Z"
  }
}
```

Complete API documentation: [API_REFERENCE.md](./API_REFERENCE.md)

---

## 📁 Project Structure

```
travel-sync/
├── cmd/                    # Application entry point
│   └── main.go
├── internal/               # Private application code
│   ├── middleware/         # CORS, rate limiting
│   ├── security/           # OAuth, JWT, auth handlers
│   ├── travel/             # Travel domain logic
│   │   ├── handlers/       # HTTP handlers
│   │   ├── services/       # Business logic + recommendation engine
│   │   ├── repositories/   # Data access
│   │   └── models/         # Data models
│   ├── user/               # User management
│   └── server/             # Server configuration
├── docs/                   # Documentation
├── .env                    # Environment variables
├── go.mod                  # Go dependencies
└── README.md
```

---

## 🎨 Architecture Visualization

### Technology Stack
<div align="center">

| Frontend | Backend | Database | Security |
|:--------:|:-------:|:--------:|:--------:|
| ![React](https://raw.githubusercontent.com/devicons/devicon/master/icons/react/react-original.svg) | ![Go](https://raw.githubusercontent.com/devicons/devicon/master/icons/go/go-original.svg) | ![PostgreSQL](https://raw.githubusercontent.com/devicons/devicon/master/icons/postgresql/postgresql-original.svg) | ![JWT](https://img.shields.io/badge/JWT-000000?style=flat&logo=jsonwebtokens&logoColor=white) |
| React | Gin Framework | PostgreSQL | JWT Auth |

</div>

### Key Components

<table>
<tr>
<td width="50%">

**🔐 Authentication System**
- Google OAuth 2.0 Integration
- JWT Token Management
- HTTP-only Cookie Storage
- CSRF Protection

</td>
<td width="50%">

**🤖 Recommendation Engine**
- Multi-factor Scoring Algorithm
- Asymmetric Time Windows
- Privacy-First Results
- Real-time Matching

</td>
</tr>
<tr>
<td width="50%">

**⚡ Rate Limiting**
- 3-Tier Protection System
- Per-user & Per-IP Tracking
- Burst Capacity Support
- Header-based Feedback

</td>
<td width="50%">

**🎫 Ticket Management**
- CRUD Operations
- Ownership Validation
- Status Tracking
- Duplicate Prevention

</td>
</tr>
</table>

> **Note:** Add your application screenshots here by uploading them to your repository in a `/docs/images/` folder

---

## 🔧 Configuration Options

| Variable | Description | Required |
|----------|-------------|----------|
| `PORT` | Server port | Yes |
| `DATABASE_URL` | PostgreSQL connection | Yes |
| `JWT_SECRET` | JWT signing key (min 32 chars) | Yes |
| `GOOGLE_CLIENT_ID` | OAuth client ID | Yes |
| `GOOGLE_CLIENT_SECRET` | OAuth secret | Yes |
| `FRONTEND_URL` | Frontend URL for redirects | Yes |

---

## 📈 Performance & Scalability

- **In-memory Rate Limiting** - Fast, per-instance tracking
- **GORM Query Optimization** - Efficient database access
- **Layered Architecture** - Easy horizontal scaling
- **Auto-migrations** - Seamless schema updates
  
---

## 🤝 Contributing

Contributions welcome! Please:
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---
