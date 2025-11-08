# 🎰 Freezino

**Freezino** is an educational casino simulator web application designed to combat gambling addiction. Users play with virtual money to understand how quickly they can lose it, while learning about responsible gaming and financial literacy.

> **⚠️ Educational Purpose Only**: This application uses NO real money. All currency is virtual and for educational purposes only.

## 🎯 Project Goals

- 🎓 **Educational**: Help users understand the risks of gambling
- 📊 **Statistical Awareness**: Show real-world wage comparisons
- 🛡️ **Prevention**: Combat gambling addiction through awareness
- 🎮 **Gamification**: Engage users with realistic casino mechanics

## ✨ Features

### 🎮 Casino Games
- **Roulette** - European roulette with multiple bet types
- **Slots** - 5-reel slot machine with various symbols
- **Blackjack** - Real-time card game via WebSocket
- **Crash** - Multiplier-based betting game
- **Hi-Lo** - Higher or lower prediction game
- **Wheel** - Fortune wheel with multiplier sectors

### 👤 User Features
- **Google OAuth Authentication** - Secure login with Google
- **Virtual Currency** - Play with pseudo-dollars
- **Work System** - Earn $500 by "working" for 3 minutes
- **Shop System** - Buy virtual items (clothes, cars, houses)
- **Profile Visualization** - See your avatar with purchased items
- **Transaction History** - Track all earnings and losses
- **Country Statistics** - Compare earnings with real-world wages

### 🌍 Internationalization (i18n)
- 🇬🇧 English
- 🇷🇺 Russian
- 🇪🇸 Spanish (optional)

### 📜 Legal Compliance
- GDPR-compliant cookie consent
- Terms of Service
- Privacy Policy
- Cookie Policy

## 🏗️ Technology Stack

### Frontend
- **Framework**: React 18 + TypeScript
- **Build Tool**: Vite
- **Styling**: TailwindCSS
- **State Management**: Zustand
- **Routing**: React Router v6
- **HTTP Client**: Axios
- **Animations**: Framer Motion
- **i18n**: react-i18next
- **Forms**: react-hook-form + zod
- **Notifications**: react-hot-toast

### Backend
- **Language**: Go 1.21+
- **Framework**: Fiber v2
- **ORM**: GORM
- **Database**: SQLite
- **Authentication**: Google OAuth 2.0
- **WebSockets**: Gorilla WebSocket
- **Validation**: go-playground/validator

### DevOps
- **Containerization**: Docker + Docker Compose
- **Reverse Proxy**: Nginx
- **Development**: Air (hot reload)

## 📁 Project Structure

```
freezino/
├── backend/              # Go backend API
│   ├── cmd/
│   │   └── server/      # Main application entry point
│   ├── internal/
│   │   ├── auth/        # Authentication logic
│   │   ├── config/      # Configuration
│   │   ├── database/    # Database setup and migrations
│   │   ├── handler/     # HTTP handlers
│   │   ├── middleware/  # Middleware functions
│   │   ├── model/       # Database models
│   │   ├── router/      # Route definitions
│   │   └── service/     # Business logic
│   ├── Dockerfile
│   ├── Makefile
│   └── go.mod
├── frontend/             # React frontend
│   ├── public/          # Static assets
│   ├── src/
│   │   ├── components/  # React components
│   │   ├── pages/       # Page components
│   │   ├── store/       # Zustand stores
│   │   ├── i18n/        # Internationalization
│   │   └── utils/       # Utility functions
│   ├── Dockerfile
│   ├── Dockerfile.dev
│   └── package.json
├── docs/                 # Documentation
├── docker/              # Docker configurations
│   └── nginx/          # Nginx configs
├── docker-compose.yml
└── README.md
```

## 🚀 Getting Started

### Prerequisites

- **Node.js** 18+ and npm
- **Go** 1.21+
- **Docker** and Docker Compose (optional, recommended)
- **Google OAuth Credentials** (for authentication)

### Environment Variables

Create `.env` files in both backend and frontend directories:

#### Backend `.env`
```bash
# Server
ENV=development
PORT=3000

# Database
DATABASE_URL=file:./data/freezino.db

# JWT
JWT_SECRET=your-secret-key-change-in-production

# Google OAuth
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:5173/auth/callback

# Frontend URL (for CORS)
FRONTEND_URL=http://localhost:5173
```

#### Frontend `.env`
```bash
VITE_API_URL=http://localhost:3000/api
VITE_WS_URL=ws://localhost:3000/ws
```

### Installation

#### Option 1: Docker (Recommended)

1. Clone the repository:
```bash
git clone https://github.com/smoreg/freezino.git
cd freezino
```

2. Create `.env` files (see above)

3. Start all services:
```bash
docker-compose up
```

4. Access the application:
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:3000
   - Nginx Proxy: http://localhost:8080

#### Option 2: Manual Setup

##### Backend Setup

```bash
cd backend

# Install dependencies
make install

# Run database migrations (automatic on first run)
make run

# Or use development mode with hot reload
make dev
```

The backend will start on `http://localhost:3000`

##### Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

The frontend will start on `http://localhost:5173`

## 📖 Development

### Backend Commands

```bash
cd backend

# Run the application
make run

# Development mode with hot reload (requires air)
make dev

# Build the application
make build

# Run tests
make test

# Format code
make fmt

# Run linter
make lint

# Clean build artifacts
make clean

# Tidy dependencies
make tidy

# Show help
make help
```

### Frontend Commands

```bash
cd frontend

# Start development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Run linter
npm run lint
```

### Docker Commands

```bash
# Start all services
docker-compose up

# Start in detached mode
docker-compose up -d

# Stop all services
docker-compose down

# View logs
docker-compose logs -f

# Rebuild containers
docker-compose up --build

# Stop and remove volumes
docker-compose down -v
```

## 🧪 Testing

### Backend Tests

```bash
cd backend
make test
```

### Frontend Tests

```bash
cd frontend
npm run test        # Unit tests (when configured)
npm run test:e2e    # E2E tests with Playwright (when configured)
```

## 📦 Production Deployment

### Build for Production

#### Backend
```bash
cd backend
make build
# Binary will be in ./bin/freezino-server
```

#### Frontend
```bash
cd frontend
npm run build
# Build output will be in ./dist
```

### Docker Production

Create `docker-compose.prod.yml` with optimized settings:

```bash
# Build and start production containers
docker-compose -f docker-compose.prod.yml up -d --build

# View production logs
docker-compose -f docker-compose.prod.yml logs -f
```

### Deployment Checklist

- [ ] Set strong `JWT_SECRET` in production
- [ ] Configure Google OAuth for production domain
- [ ] Set `ENV=production` in backend
- [ ] Update `FRONTEND_URL` and `GOOGLE_REDIRECT_URL`
- [ ] Enable HTTPS with SSL certificates
- [ ] Configure Nginx for production
- [ ] Set up database backups
- [ ] Configure logging and monitoring
- [ ] Enable rate limiting
- [ ] Review security headers

## 🔐 Security

- **Authentication**: Secure Google OAuth 2.0
- **Authorization**: JWT-based access tokens
- **CORS**: Configured for specific origins
- **Rate Limiting**: API endpoint protection
- **Input Validation**: Server-side validation with go-playground/validator
- **SQL Injection**: Protected via GORM ORM
- **XSS Protection**: React's built-in protection
- **CSRF**: Token-based protection

## 📊 API Documentation

For detailed API documentation, see:
- [OpenAPI Specification](./docs/openapi.yaml)
- [API Endpoints Documentation](./docs/API.md)

Quick reference:

- **Base URL**: `http://localhost:3000/api`
- **Health Check**: `GET /api/health`
- **Auth**: `GET /api/auth/google`, `GET /api/auth/google/callback`
- **User**: `GET /api/user/profile`, `PATCH /api/user/profile`
- **Work**: `POST /api/work/start`, `POST /api/work/complete`
- **Shop**: `GET /api/shop/items`, `POST /api/shop/buy/:itemId`
- **Games**: `POST /api/games/roulette/bet`, `POST /api/games/slots/spin`
- **WebSocket**: `ws://localhost:3000/ws/blackjack`

## 📚 Documentation

- [Architecture Overview](./docs/ARCHITECTURE.md)
- [User Guide](./docs/USER_GUIDE.md)
- [Contributing Guidelines](./CONTRIBUTING.md)
- [API Documentation](./docs/openapi.yaml)

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for details on:
- Code of Conduct
- Development workflow
- Coding standards
- Pull request process
- Testing requirements

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [React](https://react.dev/), [Go](https://go.dev/), and [Fiber](https://gofiber.io/)
- UI components inspired by modern casino interfaces
- Educational mission inspired by responsible gaming initiatives

## 📧 Contact

- **Project Repository**: [github.com/smoreg/freezino](https://github.com/smoreg/freezino)
- **Issues**: [github.com/smoreg/freezino/issues](https://github.com/smoreg/freezino/issues)
- **Website**: http://localhost:5173 (development)

## 🎯 Roadmap

See [PHASES.md](./PHASES.md) for detailed development phases.

Current Status:
- ✅ Phase 0: Setup (Complete)
- ✅ Phase 1: Auth & Core (Complete)
- ✅ Phase 2: Work System (Complete)
- ✅ Phase 2.5: Legal & i18n (Complete)
- ✅ Phase 3: Games (Complete)
- ✅ Phase 4: Shop & Profile (Complete)
- ✅ Phase 5: Polish (Complete)
- 🔄 Phase 6: Testing & Deploy (In Progress)

---

**Remember**: This is an educational tool. If you or someone you know has a gambling problem, please seek help from professional organizations.
