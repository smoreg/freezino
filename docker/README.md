# Docker Configuration

This directory contains Docker and deployment configuration for Freezino.

## 📁 Structure

```
docker/
├── nginx/
│   └── default.conf    # Nginx reverse proxy configuration
└── README.md           # This file
```

## 🐳 Docker Setup

### Development

```bash
# Start development environment
docker-compose up -d

# View logs
docker-compose logs -f

# Stop containers
docker-compose down
```

**Services:**
- Frontend (Vite Dev Server): http://localhost:5173
- Backend (Go API): http://localhost:3000
- Nginx (Reverse Proxy): http://localhost:8080

### Production

```bash
# Deploy to production
./deploy.sh prod

# Or manually:
docker-compose -f docker-compose.prod.yml up -d

# View logs
docker-compose -f docker-compose.prod.yml logs -f

# Stop containers
docker-compose -f docker-compose.prod.yml down
```

**Services:**
- Application: http://localhost (or your domain)
- API: http://localhost/api

## 🔧 Configuration

### Environment Variables

Copy `.env.example` to `.env` and update the values:

```bash
cp .env.example .env
```

Required variables:
- `JWT_SECRET` - Secret for JWT token signing
- `GOOGLE_CLIENT_ID` - Google OAuth client ID
- `GOOGLE_CLIENT_SECRET` - Google OAuth client secret
- `FRONTEND_URL` - Your frontend domain
- `VITE_API_URL` - API endpoint for frontend

### Nginx Configuration

The `nginx/default.conf` file configures:
- Reverse proxy to backend API (`/api` → `backend:3000`)
- WebSocket support (`/ws` → `backend:3000`)
- Static file serving for frontend
- Gzip compression
- Security headers
- Caching for static assets

## 📦 Deployment Script

The `deploy.sh` script provides easy deployment management:

```bash
# Production deployment
./deploy.sh prod

# Development deployment
./deploy.sh dev

# Stop all containers
./deploy.sh stop

# Restart containers
./deploy.sh restart

# View logs
./deploy.sh logs

# Backup database
./deploy.sh backup

# Clean up everything
./deploy.sh clean
```

## 🔒 SSL/HTTPS

For production, you'll need to:

1. Generate SSL certificates (Let's Encrypt recommended):
   ```bash
   certbot certonly --standalone -d yourdomain.com
   ```

2. Update `nginx/default.conf` to include SSL configuration:
   ```nginx
   server {
       listen 443 ssl http2;
       ssl_certificate /etc/nginx/ssl/fullchain.pem;
       ssl_certificate_key /etc/nginx/ssl/privkey.pem;
       # ... rest of config
   }
   ```

3. Mount SSL certificates in `docker-compose.prod.yml`

## 📊 Monitoring

### View Container Status

```bash
docker-compose -f docker-compose.prod.yml ps
```

### View Logs

```bash
# All services
docker-compose -f docker-compose.prod.yml logs -f

# Specific service
docker-compose -f docker-compose.prod.yml logs -f backend
```

### Check Health

```bash
# Backend health
curl http://localhost/api/health

# Frontend health
curl http://localhost/
```

## 🗄️ Database

### Backup

```bash
# Using deploy script
./deploy.sh backup

# Manual backup
docker cp freezino-backend-prod:/root/data/freezino.db ./backup.db
```

### Restore

```bash
# Copy backup to container
docker cp backup.db freezino-backend-prod:/root/data/freezino.db

# Restart backend
docker-compose -f docker-compose.prod.yml restart backend
```

## 🐛 Troubleshooting

### Containers won't start

Check logs:
```bash
docker-compose -f docker-compose.prod.yml logs
```

### Port already in use

Change ports in `docker-compose.yml` or stop conflicting services:
```bash
# Check what's using port 80
sudo lsof -i :80

# Or use different ports
# Edit docker-compose.prod.yml and change port mappings
```

### Database locked error

Ensure only one backend instance is accessing the database:
```bash
docker-compose -f docker-compose.prod.yml down
docker-compose -f docker-compose.prod.yml up -d
```

### Permission issues

Ensure proper file permissions:
```bash
chmod +x deploy.sh
sudo chown -R $USER:$USER .
```

## 🔄 Updates

To update the application:

```bash
# Pull latest changes
git pull

# Rebuild and restart
./deploy.sh prod
```

## 📝 Notes

- SQLite database is stored in a Docker volume (`backend-data`)
- Logs are rotated automatically (max 10MB, 3 files)
- Health checks run every 30 seconds
- Containers restart automatically on failure
