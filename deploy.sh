#!/bin/bash

# ================================================
# Freezino Deployment Script
# ================================================

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if .env file exists
check_env() {
    if [ ! -f .env ]; then
        log_error ".env file not found!"
        log_info "Please create a .env file based on .env.example"
        exit 1
    fi
    log_success ".env file found"
}

# Load environment variables
load_env() {
    if [ -f .env ]; then
        export $(cat .env | grep -v '^#' | xargs)
        log_success "Environment variables loaded"
    fi
}

# Check if Docker is running
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        log_error "Docker is not running!"
        log_info "Please start Docker and try again."
        exit 1
    fi
    log_success "Docker is running"
}

# Check if Docker Compose is installed
check_docker_compose() {
    if ! command -v docker-compose &> /dev/null; then
        log_error "docker-compose is not installed!"
        log_info "Please install docker-compose and try again."
        exit 1
    fi
    log_success "docker-compose is installed"
}

# Stop existing containers
stop_containers() {
    log_info "Stopping existing containers..."
    docker-compose -f docker-compose.prod.yml down
    log_success "Containers stopped"
}

# Build Docker images
build_images() {
    log_info "Building Docker images..."
    docker-compose -f docker-compose.prod.yml build --no-cache
    log_success "Docker images built"
}

# Start containers
start_containers() {
    log_info "Starting containers..."
    docker-compose -f docker-compose.prod.yml up -d
    log_success "Containers started"
}

# Check container health
check_health() {
    log_info "Waiting for services to be healthy..."
    sleep 10

    # Check backend health
    if docker-compose -f docker-compose.prod.yml ps | grep -q "backend.*healthy"; then
        log_success "Backend is healthy"
    else
        log_warning "Backend health check pending..."
    fi

    # Check frontend health
    if docker-compose -f docker-compose.prod.yml ps | grep -q "frontend.*healthy"; then
        log_success "Frontend is healthy"
    else
        log_warning "Frontend health check pending..."
    fi
}

# Show logs
show_logs() {
    log_info "Recent logs:"
    docker-compose -f docker-compose.prod.yml logs --tail=50
}

# Show running containers
show_status() {
    log_info "Container status:"
    docker-compose -f docker-compose.prod.yml ps
}

# Backup database
backup_db() {
    log_info "Creating database backup..."
    BACKUP_DIR="./backups"
    mkdir -p $BACKUP_DIR
    TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
    BACKUP_FILE="$BACKUP_DIR/freezino_backup_$TIMESTAMP.db"

    # Copy database from container
    docker cp freezino-backend-prod:/root/data/freezino.db "$BACKUP_FILE" 2>/dev/null || true

    if [ -f "$BACKUP_FILE" ]; then
        log_success "Database backed up to $BACKUP_FILE"
    else
        log_warning "No database found to backup (this is normal for first deployment)"
    fi
}

# Clean up old images
cleanup() {
    log_info "Cleaning up old Docker images..."
    docker image prune -f
    log_success "Cleanup complete"
}

# Main deployment function
deploy() {
    echo ""
    log_info "============================================"
    log_info "   FREEZINO DEPLOYMENT SCRIPT"
    log_info "============================================"
    echo ""

    # Pre-flight checks
    check_docker
    check_docker_compose
    check_env
    load_env

    # Backup before deployment
    backup_db

    # Deploy
    stop_containers
    build_images
    start_containers

    # Post-deployment
    check_health
    show_status
    cleanup

    echo ""
    log_success "============================================"
    log_success "   DEPLOYMENT COMPLETE!"
    log_success "============================================"
    echo ""
    log_info "Application should be available at:"
    log_info "  - Frontend: http://localhost (or your domain)"
    log_info "  - Backend API: http://localhost/api"
    echo ""
    log_info "To view logs, run:"
    log_info "  docker-compose -f docker-compose.prod.yml logs -f"
    echo ""
}

# Development deployment
deploy_dev() {
    log_info "Deploying in DEVELOPMENT mode..."
    docker-compose up -d
    docker-compose ps
    log_success "Development environment is running"
    log_info "Frontend: http://localhost:5173"
    log_info "Backend: http://localhost:3000"
    log_info "Nginx: http://localhost:8080"
}

# Parse command line arguments
case "${1:-prod}" in
    dev)
        deploy_dev
        ;;
    prod)
        deploy
        ;;
    stop)
        log_info "Stopping all containers..."
        docker-compose -f docker-compose.prod.yml down
        log_success "All containers stopped"
        ;;
    restart)
        log_info "Restarting containers..."
        docker-compose -f docker-compose.prod.yml restart
        log_success "Containers restarted"
        ;;
    logs)
        docker-compose -f docker-compose.prod.yml logs -f
        ;;
    backup)
        backup_db
        ;;
    clean)
        log_warning "This will remove all containers, volumes, and images!"
        read -p "Are you sure? (y/N) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            docker-compose -f docker-compose.prod.yml down -v
            docker system prune -af --volumes
            log_success "Cleanup complete"
        fi
        ;;
    *)
        echo "Usage: $0 {dev|prod|stop|restart|logs|backup|clean}"
        echo ""
        echo "Commands:"
        echo "  dev      - Deploy in development mode"
        echo "  prod     - Deploy in production mode (default)"
        echo "  stop     - Stop all containers"
        echo "  restart  - Restart all containers"
        echo "  logs     - Show and follow logs"
        echo "  backup   - Backup database"
        echo "  clean    - Remove all containers, volumes, and images"
        exit 1
        ;;
esac
