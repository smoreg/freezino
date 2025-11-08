# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Freezino is a gamified productivity platform with a Go/Fiber backend and React/TypeScript frontend. Users earn virtual currency by completing work sessions, which they can spend on casino-style games or shop items. The platform includes Google OAuth authentication, multilingual support (i18n), and a virtual economy system.

## Development Commands

### Quick Start
```bash
make init          # Install all dependencies (backend + frontend)
make dev           # Run both servers in parallel (backend:3000, frontend:5173)
```

### Backend (Go/Fiber)
```bash
cd backend
make run           # Run server (localhost:3000)
make build         # Build binary to ./bin/
make test          # Run tests
make fmt           # Format code
make lint          # Run golangci-lint (if installed)
make dev           # Run with air hot-reload (installs air if needed)
```

### Frontend (React/Vite)
```bash
cd frontend
npm run dev        # Dev server (localhost:5173)
npm run build      # TypeScript check + production build
npm run lint       # ESLint check
npm run preview    # Preview production build
```

### Testing
- **Backend tests**: `cd backend && make test`
- **Frontend has no test command** - tests not yet implemented

## Architecture

### Backend Structure (Go)

**Entry Point**: `cmd/server/main.go`
- Initializes SQLite database (`./data/freezino.db`)
- Runs migrations via GORM
- Configures Fiber app with middleware (CORS, logger, recover)
- Sets up routes via `internal/router`

**Layer Pattern**:
```
handler/ (HTTP layer)
  ↓
service/ (Business logic)
  ↓
database/ (Data access via GORM models)
```

**Key Directories**:
- `internal/model/` - GORM models (User, Item, Transaction, GameSession, etc.)
- `internal/handler/` - HTTP handlers and API endpoints
- `internal/service/` - Business logic (shop service, user service, etc.)
- `internal/middleware/` - Auth middleware (JWT validation), CORS, logger
- `internal/router/` - Route definitions (router.go contains ALL API routes)
- `internal/auth/` - Google OAuth handler
- `internal/database/` - Database initialization, migrations, seeding
- `internal/game/` - Game logic (roulette, slots, crash, hi-lo, wheel, blackjack)
- `internal/config/` - Environment configuration

**Database**:
- SQLite with GORM ORM
- Auto-migrations in `database/migrate.go`
- Seed data in `database/seed.go`
- Models use soft deletes (`gorm.DeletedAt`)

**Authentication**:
- Google OAuth 2.0 flow (`/api/auth/google` → `/api/auth/google/callback`)
- JWT tokens (access + refresh)
- Auth middleware validates JWT for protected routes
- `user_id` query param used for development (production uses JWT)

### Frontend Structure (React/TypeScript)

**Tech Stack**:
- React 19 + TypeScript
- Vite build tool
- React Router for routing
- Zustand for state management
- i18next for internationalization (EN/ES/RU)
- Tailwind CSS for styling
- Framer Motion for animations
- Recharts for data visualization

**State Management (Zustand stores)**:
- `authStore.ts` - Authentication state, user profile
- `workStore.ts` - Work timer state
- `shopStore.ts` - Shop items, filters

**Routing** (`App.tsx`):
- Public: `/login`, `/contact`, `/about`, `/terms`, `/privacy`, `/cookies`
- Protected (requires auth): `/`, `/dashboard`, `/history`, `/shop`, `/profile`
- Games are embedded in home page or dashboard

**Key Directories**:
- `pages/` - Page components (Home, Dashboard, ShopPage, ProfilePage, etc.)
- `components/` - Reusable components organized by feature:
  - `games/` - Game components (Roulette, Slots, Blackjack, etc.)
  - `shop/` - Shop UI (ItemCard, BuyModal, SellModal, ShopFilters)
  - `profile/` - Profile components (MyItemsList)
  - `layout/` - Layout components (Navbar, Footer)
- `store/` - Zustand state management
- `services/` - API client (`api.ts` with axios)
- `i18n/` - Internationalization setup
- `types/` - TypeScript type definitions
- `layouts/` - Layout wrappers (MainLayout with navigation)

**i18n Setup**:
- Translation files: `i18n/locales/{en,es,ru}.json`
- Language switcher component
- Browser language detection enabled
- All user-facing text should use `t()` function from `react-i18next`

## API Architecture

All API routes defined in `backend/internal/router/router.go`:

**Auth**: `/api/auth/*`
- Google OAuth login/callback
- Token refresh
- `/me`, `/logout` (protected)

**User**: `/api/user/*`
- Profile, balance, stats, transactions, items

**Work**: `/api/work/*`
- Start/complete work sessions, get history

**Shop**: `/api/shop/*` (see `backend/SHOP_API.md` for full docs)
- Get items (with type/rarity filters)
- Buy/sell items
- Equip items
- Get user's items

**Games**: `/api/games/*`
- Roulette, slots, blackjack, crash, hi-lo, wheel of fortune
- Game history and statistics

**Stats**: `/api/stats/countries/*`
- Country statistics (mock data)

**Contact**: `/api/contact`
- Submit contact form

## Database Models

Key GORM models in `internal/model/`:
- `User` - User account with Google ID, balance, avatar
- `Item` - Shop items (clothing, car, house, accessories) with rarity
- `UserItem` - Junction table for user's purchased items (with equipped status)
- `Transaction` - Balance changes (work, game_win, game_loss, shop_purchase, shop_sale)
- `WorkSession` - Completed work sessions with earnings
- `GameSession` - Game plays with bet/win amounts
- `ContactMessage` - Contact form submissions

## Important Patterns

### Adding a New API Endpoint

1. Define model in `internal/model/` (if needed)
2. Add migration in `internal/database/migrate.go`
3. Create handler in `internal/handler/`
4. Create service in `internal/service/` (for business logic)
5. Register route in `internal/router/router.go`

### Adding a New Frontend Page

1. Create page component in `pages/`
2. Add route in `App.tsx`
3. Add to navigation in `layouts/MainLayout.tsx` (if needed)
4. Create Zustand store in `store/` (if state needed)
5. Add i18n keys to all locale files

### Game Implementation Pattern

Games follow this structure:
- Frontend component in `components/games/`
- Backend handler in `internal/handler/` or `internal/handler/games/`
- Game logic in `internal/game/`
- WebSocket support for real-time games (like roulette)

### Transaction Pattern

All balance changes create a `Transaction` record with:
- Type: `work`, `game_win`, `game_loss`, `shop_purchase`, `shop_sale`
- Amount (positive for credits, negative for debits)
- Description
- `BalanceAfter` snapshot

## Environment Variables

Backend uses `.env` (see `.env.example`):
- `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET` - OAuth credentials
- `GOOGLE_REDIRECT_URL` - OAuth callback URL
- `JWT_SECRET` - JWT signing key
- `FRONTEND_URL` - CORS configuration
- `PORT` - Server port (default: 3000)

## Common Gotchas

- **Auth in development**: Routes accept `?user_id=1` query param for testing without OAuth
- **Database resets**: Database file is `backend/data/freezino.db` - delete to reset
- **i18n missing keys**: Always add keys to ALL locale files (en, es, ru) to avoid fallback warnings
- **CORS issues**: Backend CORS allows `http://localhost:5173` by default
- **Shop filters**: Use query params `?type=clothing&rarity=rare` - both optional
- **Soft deletes**: Models use `gorm.DeletedAt` - don't forget `.Unscoped()` for hard deletes

## Documentation

- Shop API: `backend/SHOP_API.md`
- Shop components: `frontend/src/components/shop/README.md`
