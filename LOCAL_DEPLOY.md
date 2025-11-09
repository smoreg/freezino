# Local Deployment Guide

This guide explains how to deploy and manage the Freezino application locally using Make commands.

## Prerequisites

- Go 1.21+
- Node.js 18+
- Make

## Quick Start

```bash
# Initialize project (first time only)
make init

# Build for production
make build

# Deploy locally (build + start services)
make local-deploy
```

## Available Commands

### Development

```bash
# Start development servers (with hot reload)
make dev

# Development servers:
# - Backend: http://localhost:3000
# - Frontend: http://localhost:5173
```

### Build

```bash
# Build both backend and frontend
make build

# Build only backend
make build-backend

# Build only frontend
make build-frontend
```

### Local Deployment

```bash
# Full local deployment (production mode)
make local-deploy
# This will:
# 1. Build both services
# 2. Stop old services
# 3. Start new services in background
# 4. Check health status

# Services run on:
# - Backend: http://localhost:3000
# - Frontend: http://localhost:4173 (preview mode)
```

### Service Management

```bash
# Check service status
make status
# Shows:
# - Process status (running/stopped)
# - PID numbers
# - Health check results

# Stop all services
make stop

# Restart services (stop + deploy)
make restart
```

### Logs

```bash
# View backend logs (real-time)
make logs-backend

# View frontend logs (real-time)
make logs-frontend

# Log files location:
# - Backend: backend/logs/backend.log
# - Frontend: backend/logs/frontend.log
```

### Cleanup

```bash
# Clean build artifacts and database
make clean
```

## Service Details

### Backend

- **Binary**: `backend/bin/freezino-server`
- **PID file**: `backend/.backend.pid`
- **Log file**: `backend/logs/backend.log`
- **Port**: 3000
- **Health check**: `http://localhost:3000/api/health`

### Frontend

- **Command**: `npm run preview` (production preview)
- **PID file**: `frontend/.frontend.pid`
- **Log file**: `backend/logs/frontend.log`
- **Port**: 4173
- **URL**: `http://localhost:4173`

## Example Workflow

### Initial Setup

```bash
# Clone and initialize
git clone <repo>
cd freezino
make init
```

### Development

```bash
# Start dev servers
make dev

# Work on code with hot reload...
# Press Ctrl+C to stop
```

### Testing Local Production Build

```bash
# Build and deploy locally
make local-deploy

# Check status
make status

# View logs if needed
make logs-backend
make logs-frontend

# Stop when done
make stop
```

### Redeploying After Changes

```bash
# Rebuild and restart
make restart

# Or do it manually:
make build
make stop
make local-deploy
```

## Troubleshooting

### Services won't start

```bash
# Check if ports are already in use
lsof -i :3000  # backend
lsof -i :4173  # frontend

# Kill processes if needed
make stop
```

### Check service health

```bash
# Use status command
make status

# Manual health check
curl http://localhost:3000/api/health

# Check logs
make logs-backend
make logs-frontend
```

### Clean restart

```bash
# Stop everything
make stop

# Clean build artifacts
make clean

# Rebuild from scratch
make init
make build
make local-deploy
```

## Production Deployment

For production deployment to remote server:

```bash
# Deploy both services to freezino.online
make deploy

# Deploy only backend
make deploy-backend

# Deploy only frontend
make deploy-frontend
```

See main README.md for production deployment details.

## Notes

- PID files are stored in `backend/.backend.pid` and `frontend/.frontend.pid`
- Logs are stored in `backend/logs/`
- All temporary files are in `.gitignore`
- Development mode uses hot reload (ports 3000, 5173)
- Local deployment uses production build (ports 3000, 4173)
