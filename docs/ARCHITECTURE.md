# Freezino Architecture Documentation

This document provides a comprehensive overview of the Freezino application architecture, design patterns, and technical decisions.

## 📋 Table of Contents

1. [System Overview](#system-overview)
2. [Architecture Diagram](#architecture-diagram)
3. [Backend Architecture](#backend-architecture)
4. [Frontend Architecture](#frontend-architecture)
5. [Database Schema](#database-schema)
6. [API Design](#api-design)
7. [Authentication Flow](#authentication-flow)
8. [Game Engine Architecture](#game-engine-architecture)
9. [State Management](#state-management)
10. [Security Architecture](#security-architecture)
11. [Deployment Architecture](#deployment-architecture)
12. [Design Patterns](#design-patterns)

## 🏗️ System Overview

Freezino is a full-stack web application built with a modern tech stack:

- **Frontend**: React + TypeScript + Vite
- **Backend**: Go + Fiber framework
- **Database**: SQLite with GORM ORM
- **Real-time**: WebSockets for live games (Blackjack)
- **Authentication**: Google OAuth 2.0 + JWT
- **Deployment**: Docker + Docker Compose + Nginx

### Core Principles

1. **Separation of Concerns**: Clear separation between frontend, backend, and data layers
2. **Stateless Backend**: JWT-based authentication, no server-side sessions
3. **API-First Design**: RESTful API with OpenAPI documentation
4. **Type Safety**: TypeScript on frontend, strong typing in Go
5. **Educational Focus**: Every feature designed with educational goals

## 🎨 Architecture Diagram

### High-Level System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         Client Browser                       │
│  ┌─────────────────────────────────────────────────────┐   │
│  │         React Application (SPA)                      │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────────────┐  │   │
│  │  │  Pages   │  │Components│  │  Zustand Store   │  │   │
│  │  └──────────┘  └──────────┘  └──────────────────┘  │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────────────┐  │   │
│  │  │  Hooks   │  │  Utils   │  │  i18n (Locales)  │  │   │
│  │  └──────────┘  └──────────┘  └──────────────────┘  │   │
│  └─────────────────────────────────────────────────────┘   │
│           │                              │                  │
│           │ HTTP/HTTPS                   │ WebSocket        │
│           ▼                              ▼                  │
└───────────┼──────────────────────────────┼──────────────────┘
            │                              │
┌───────────┼──────────────────────────────┼──────────────────┐
│           ▼                              ▼                  │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Nginx Reverse Proxy                     │   │
│  │         (SSL Termination, Load Balancing)            │   │
│  └─────────────────────────────────────────────────────┘   │
│           │                              │                  │
│           ▼                              ▼                  │
│  ┌─────────────────────────────────────────────────────┐   │
│  │           Go Backend (Fiber Framework)               │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────────────┐  │   │
│  │  │ Handlers │  │Services  │  │   Middleware     │  │   │
│  │  └──────────┘  └──────────┘  └──────────────────┘  │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────────────┐  │   │
│  │  │  Auth    │  │ Router   │  │   Game Engine    │  │   │
│  │  └──────────┘  └──────────┘  └──────────────────┘  │   │
│  └─────────────────────────────────────────────────────┘   │
│                          │                                  │
│                          ▼                                  │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              GORM ORM Layer                          │   │
│  └─────────────────────────────────────────────────────┘   │
│                          │                                  │
│                          ▼                                  │
│  ┌─────────────────────────────────────────────────────┐   │
│  │            SQLite Database                           │   │
│  │  ┌─────────┐ ┌─────────┐ ┌──────────┐ ┌─────────┐  │   │
│  │  │  Users  │ │  Items  │ │GameSess. │ │  Trans. │  │   │
│  │  └─────────┘ └─────────┘ └──────────┘ └─────────┘  │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│                    Docker Container Network                 │
└─────────────────────────────────────────────────────────────┘
            │
            ▼
┌─────────────────────────────────────────────────────────────┐
│                   External Services                          │
│  ┌──────────────────┐        ┌──────────────────────┐       │
│  │  Google OAuth    │        │   (Future: Redis,     │       │
│  │    Provider      │        │    Analytics, etc.)   │       │
│  └──────────────────┘        └──────────────────────┘       │
└─────────────────────────────────────────────────────────────┘
```

### Request Flow Diagram

```
User Action → React Component → API Service → Backend Handler → Service Layer → Database
                     ↓                                  ↓
              State Update ← JSON Response ←────────────┘
```

## 🔧 Backend Architecture

### Directory Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── auth/
│   │   └── google.go            # Google OAuth implementation
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── database/
│   │   ├── database.go          # Database initialization
│   │   └── migrations.go        # Database migrations
│   ├── handler/
│   │   ├── user.go              # User endpoints
│   │   ├── work.go              # Work system endpoints
│   │   ├── shop.go              # Shop endpoints
│   │   ├── roulette.go          # Roulette game
│   │   ├── slots.go             # Slots game
│   │   ├── game_handler.go      # Blackjack WebSocket
│   │   └── games/               # Other games
│   ├── middleware/
│   │   ├── auth.go              # JWT authentication
│   │   └── middleware.go        # CORS, logging, etc.
│   ├── model/
│   │   ├── user.go              # User model
│   │   ├── item.go              # Shop item model
│   │   ├── transaction.go       # Transaction model
│   │   └── ...                  # Other models
│   ├── router/
│   │   └── router.go            # Route definitions
│   └── service/
│       ├── user.go              # User business logic
│       ├── work.go              # Work system logic
│       ├── shop.go              # Shop logic
│       └── ...                  # Game services
├── Dockerfile
├── Makefile
└── go.mod
```

### Layered Architecture

```
┌────────────────────────────────────────┐
│        HTTP/WebSocket Layer            │
│         (Handlers)                     │
│  - Parse requests                      │
│  - Validate input                      │
│  - Call services                       │
│  - Return responses                    │
└────────────────────────────────────────┘
                  │
                  ▼
┌────────────────────────────────────────┐
│         Business Logic Layer           │
│           (Services)                   │
│  - Implement game logic                │
│  - Calculate payouts                   │
│  - Manage transactions                 │
│  - Enforce business rules              │
└────────────────────────────────────────┘
                  │
                  ▼
┌────────────────────────────────────────┐
│          Data Access Layer             │
│            (Models + GORM)             │
│  - Database operations                 │
│  - Data persistence                    │
│  - Relationships                       │
└────────────────────────────────────────┘
                  │
                  ▼
┌────────────────────────────────────────┐
│           Database Layer               │
│            (SQLite)                    │
└────────────────────────────────────────┘
```

### Key Backend Patterns

1. **Handler-Service Pattern**
   - Handlers: HTTP request/response logic
   - Services: Business logic implementation
   - Clear separation of concerns

2. **Dependency Injection**
   - Configuration injected via constructors
   - Database connections shared via singleton

3. **Middleware Chain**
   - Recovery → Logger → CORS → Auth → Handler

## 🎨 Frontend Architecture

### Directory Structure

```
frontend/
├── src/
│   ├── components/
│   │   ├── layout/
│   │   │   ├── Header.tsx        # App header
│   │   │   ├── Sidebar.tsx       # Navigation sidebar
│   │   │   └── Footer.tsx        # App footer
│   │   ├── games/
│   │   │   ├── Roulette.tsx      # Roulette game
│   │   │   ├── Slots.tsx         # Slots game
│   │   │   └── ...
│   │   ├── shop/
│   │   │   ├── ItemCard.tsx      # Shop item card
│   │   │   ├── ShopFilters.tsx   # Filter controls
│   │   │   └── ...
│   │   └── ...
│   ├── pages/
│   │   ├── LoginPage.tsx         # Login page
│   │   ├── DashboardPage.tsx     # Main dashboard
│   │   ├── ShopPage.tsx          # Shop page
│   │   ├── ProfilePage.tsx       # User profile
│   │   └── ...
│   ├── store/
│   │   ├── authStore.ts          # Auth state (Zustand)
│   │   ├── userStore.ts          # User state
│   │   └── ...
│   ├── hooks/
│   │   ├── useAuth.ts            # Auth hook
│   │   ├── useBalance.ts         # Balance hook
│   │   └── ...
│   ├── services/
│   │   ├── api.ts                # Axios instance
│   │   ├── authService.ts        # Auth API calls
│   │   ├── gameService.ts        # Game API calls
│   │   └── ...
│   ├── i18n/
│   │   ├── config.ts             # i18next config
│   │   └── locales/
│   │       ├── en.json           # English translations
│   │       ├── ru.json           # Russian translations
│   │       └── ...
│   ├── utils/
│   │   ├── sounds.ts             # Sound effects
│   │   └── ...
│   ├── App.tsx                   # Root component
│   └── main.tsx                  # Entry point
├── public/
│   └── sounds/                   # Audio files
├── Dockerfile
└── package.json
```

### Component Hierarchy

```
App
├── Router
│   ├── LoginPage
│   ├── ProtectedRoute
│   │   ├── DashboardPage
│   │   │   ├── Header (balance, user menu)
│   │   │   ├── Sidebar (navigation)
│   │   │   └── GameCard[] (game list)
│   │   ├── ShopPage
│   │   │   ├── ShopFilters
│   │   │   └── ItemCard[]
│   │   ├── ProfilePage
│   │   │   ├── Avatar (visual items)
│   │   │   └── UserStats
│   │   └── GamePages
│   │       ├── RoulettePage
│   │       ├── SlotsPage
│   │       └── BlackjackPage
│   └── LegalPages
│       ├── TermsPage
│       ├── PrivacyPage
│       └── CookiesPage
└── Global Components
    ├── CookieConsent
    ├── LanguageSwitcher
    └── ErrorBoundary
```

### State Management Flow

```
User Action
    ↓
Component Event Handler
    ↓
Store Action (Zustand)
    ↓
API Service Call
    ↓
Backend API
    ↓
Update Store State
    ↓
Re-render Components
    ↓
UI Update
```

## 🗄️ Database Schema

### Entity Relationship Diagram

```
┌─────────────────┐
│      Users      │
├─────────────────┤
│ id (PK)         │
│ google_id       │
│ email           │
│ name            │
│ avatar          │
│ balance         │
│ total_work_time │
│ created_at      │
│ updated_at      │
└─────────────────┘
        │
        │ 1:N
        ├──────────────────┐
        │                  │
        ▼                  ▼
┌─────────────────┐  ┌─────────────────┐
│  Transactions   │  │  UserItems      │
├─────────────────┤  ├─────────────────┤
│ id (PK)         │  │ id (PK)         │
│ user_id (FK)    │  │ user_id (FK)    │
│ type            │  │ item_id (FK)    │
│ amount          │  │ equipped        │
│ description     │  │ purchase_price  │
│ created_at      │  │ purchased_at    │
└─────────────────┘  └─────────────────┘
                              │
        ┌─────────────────────┘
        │
        ▼
┌─────────────────┐
│      Items      │
├─────────────────┤
│ id (PK)         │
│ name            │
│ description     │
│ price           │
│ type            │
│ rarity          │
│ image_url       │
└─────────────────┘

┌─────────────────┐       ┌─────────────────┐
│  WorkSessions   │       │  GameSessions   │
├─────────────────┤       ├─────────────────┤
│ id (PK)         │       │ id (PK)         │
│ user_id (FK)    │       │ user_id (FK)    │
│ start_time      │       │ game_type       │
│ end_time        │       │ bet_amount      │
│ duration        │       │ payout          │
│ earned          │       │ won             │
│ completed       │       │ created_at      │
└─────────────────┘       └─────────────────┘
```

### Key Tables

**Users**: Core user data and authentication
**Transactions**: Financial transaction history
**Items**: Shop items catalog
**UserItems**: User's purchased items
**WorkSessions**: Work history tracking
**GameSessions**: Game play history

## 🔌 API Design

### RESTful Principles

- **Resource-based URLs**: `/api/user/profile`, `/api/shop/items`
- **HTTP Methods**: GET (read), POST (create), PATCH (update), DELETE (remove)
- **Status Codes**: 200 OK, 201 Created, 400 Bad Request, 401 Unauthorized, 404 Not Found, 500 Error
- **JSON Format**: All requests/responses use JSON
- **Versioning**: Future versioning via `/api/v2/...`

### API Groups

```
/api
├── /health              # Health check
├── /auth                # Authentication
│   ├── /google
│   ├── /google/callback
│   ├── /refresh
│   ├── /me
│   └── /logout
├── /user                # User management
│   ├── /profile
│   ├── /balance
│   ├── /stats
│   ├── /transactions
│   └── /items
├── /work                # Work system
│   ├── /start
│   ├── /status
│   ├── /complete
│   └── /history
├── /stats               # Statistics
│   └── /countries
├── /shop                # Item shop
│   ├── /items
│   ├── /buy/:id
│   ├── /sell/:id
│   ├── /my-items
│   └── /equip/:id
└── /games               # Casino games
    ├── /roulette
    ├── /slots
    ├── /crash
    ├── /hilo
    ├── /wheel
    ├── /history
    └── /stats

/ws                      # WebSocket
└── /blackjack          # Live blackjack game
```

## 🔐 Authentication Flow

### OAuth 2.0 Flow

```
1. User clicks "Login with Google"
   ↓
2. Frontend redirects to /api/auth/google
   ↓
3. Backend redirects to Google OAuth consent
   ↓
4. User approves on Google
   ↓
5. Google redirects to /api/auth/google/callback?code=xxx
   ↓
6. Backend exchanges code for user info
   ↓
7. Backend creates/updates user in DB
   ↓
8. Backend generates JWT tokens (access + refresh)
   ↓
9. Backend redirects to frontend with tokens
   ↓
10. Frontend stores tokens in localStorage
    ↓
11. Frontend includes token in Authorization header
```

### JWT Token Structure

**Access Token** (short-lived, 15 min):
```json
{
  "sub": "user_id",
  "email": "user@example.com",
  "exp": 1234567890,
  "iat": 1234567000
}
```

**Refresh Token** (long-lived, 7 days):
```json
{
  "sub": "user_id",
  "type": "refresh",
  "exp": 1234999999,
  "iat": 1234567000
}
```

### Protected Endpoints

Middleware checks:
1. Authorization header exists
2. Token format is valid
3. Token signature is valid
4. Token is not expired
5. User exists in database

## 🎮 Game Engine Architecture

### Game Interface

All games implement a common interface:

```go
type Game interface {
    PlaceBet(userID uint, amount float64, params map[string]interface{}) (Result, error)
    ValidateBet(amount float64, params map[string]interface{}) error
    CalculatePayout(result interface{}) float64
}
```

### Game Flow

```
1. User places bet (POST /api/games/{game}/bet)
   ↓
2. Handler validates request
   ↓
3. Service checks user balance
   ↓
4. Service deducts bet amount
   ↓
5. Game logic executes (random number generation)
   ↓
6. Service calculates payout
   ↓
7. Service updates user balance (if win)
   ↓
8. Service creates transaction record
   ↓
9. Service creates game session record
   ↓
10. Return result to client
```

### Random Number Generation

- Uses `crypto/rand` for secure randomness
- House edge built into payout calculations
- Fair but slightly favors the house (realistic casino behavior)

### Example: Roulette

```go
// 1. Generate winning number (0-36)
winningNumber := generateSecureRandom(0, 36)

// 2. Check if bet wins
won := checkBetWin(betType, betValue, winningNumber)

// 3. Calculate payout
if won {
    multiplier := getMultiplier(betType)
    payout = betAmount * multiplier
}

// 4. Return result
return RouletteResult{
    WinningNumber: winningNumber,
    Won: won,
    Payout: payout,
}
```

## 📦 State Management

### Zustand Stores

**authStore**: Authentication state
```typescript
interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  login: (tokens: Tokens) => void;
  logout: () => void;
  refreshToken: () => Promise<void>;
}
```

**userStore**: User data and balance
```typescript
interface UserState {
  balance: number;
  stats: UserStats;
  fetchBalance: () => Promise<void>;
  updateBalance: (amount: number) => void;
}
```

**gameStore**: Active game state
```typescript
interface GameState {
  currentGame: Game | null;
  isPlaying: boolean;
  startGame: (gameType: string) => void;
  endGame: () => void;
}
```

## 🔒 Security Architecture

### Security Layers

1. **Transport Security**
   - HTTPS in production
   - Secure WebSocket (WSS)

2. **Authentication**
   - OAuth 2.0 with Google
   - JWT tokens with expiration
   - Refresh token rotation

3. **Authorization**
   - Middleware validates tokens
   - User-specific resource access

4. **Input Validation**
   - Backend validates all inputs
   - Type checking with TypeScript/Go
   - Zod schemas on frontend

5. **Output Sanitization**
   - React automatic XSS protection
   - JSON encoding prevents injection

6. **Database Security**
   - GORM prevents SQL injection
   - Parameterized queries only

7. **Rate Limiting** (future)
   - API endpoint rate limits
   - Per-user request limits

## 🚀 Deployment Architecture

### Docker Containers

```
┌──────────────────────────────────────────┐
│         Docker Network: freezino          │
│                                           │
│  ┌────────────┐  ┌──────────────┐        │
│  │  Frontend  │  │   Backend    │        │
│  │  (Vite)    │  │   (Go)       │        │
│  │  :5173     │  │   :3000      │        │
│  └────────────┘  └──────────────┘        │
│                                           │
│  ┌─────────────────────────────┐         │
│  │    Nginx Reverse Proxy      │         │
│  │         :80, :443           │         │
│  └─────────────────────────────┘         │
│                                           │
│  Volume: backend-data (SQLite DB)        │
└──────────────────────────────────────────┘
```

### Production Considerations

- Multi-stage Docker builds (smaller images)
- Health checks for containers
- Volume persistence for database
- Environment-based configuration
- Nginx for SSL termination
- Log aggregation (future)
- Monitoring (future)

## 🏛️ Design Patterns Used

### Backend Patterns

1. **Repository Pattern**: Data access abstraction (GORM models)
2. **Service Layer**: Business logic separation
3. **Factory Pattern**: Game creation
4. **Singleton**: Database connection
5. **Middleware Chain**: Request processing pipeline
6. **Strategy Pattern**: Different game implementations

### Frontend Patterns

1. **Component Composition**: Reusable UI components
2. **Custom Hooks**: Reusable logic (useAuth, useBalance)
3. **Provider Pattern**: Context/state distribution (i18n)
4. **Container/Presenter**: Smart/dumb components
5. **Observer Pattern**: Zustand state subscriptions

## 🔄 Future Enhancements

- Redis for session storage and caching
- PostgreSQL for production database
- GraphQL API option
- Server-side rendering (SSR)
- Progressive Web App (PWA)
- Real-time notifications (WebSocket push)
- Microservices architecture (if scaling needed)
- Event sourcing for game history

---

**Last Updated**: 2025-11-08
**Version**: 1.0.0
