.PHONY: init dev clean build build-backend build-frontend local-deploy restart stop status logs-backend logs-frontend deploy deploy-backend deploy-frontend

# Initialize project (install dependencies)
init:
	@echo "📦 Installing backend dependencies..."
	cd backend && go mod download
	@echo "📦 Installing frontend dependencies..."
	cd frontend && npm install
	@echo "✅ Project initialized successfully!"

# Run development servers
dev:
	@echo "🚀 Starting development servers..."
	@echo "Backend will run on http://localhost:3000"
	@echo "Frontend will run on http://localhost:5173"
	@make -j2 dev-backend dev-frontend

dev-backend:
	cd backend && make run

dev-frontend:
	cd frontend && npm run dev

# Build for production
build: build-backend build-frontend
	@echo "✅ Full build completed!"

build-backend:
	@echo "🔨 Building backend..."
	cd backend && make build
	@echo "✅ Backend built successfully!"

build-frontend:
	@echo "🔨 Building frontend..."
	cd frontend && npm run build
	@echo "✅ Frontend built successfully!"

# Local deployment (build and run in production mode)
local-deploy: build
	@echo "🚀 Starting local deployment..."
	@$(MAKE) stop
	@sleep 1
	@$(MAKE) -j2 local-run-backend local-run-frontend
	@sleep 3
	@$(MAKE) status

local-run-backend:
	@echo "▶️  Starting backend in production mode..."
	@cd backend && nohup ./bin/freezino-server > logs/backend.log 2>&1 </dev/null & echo $$! > .backend.pid
	@echo "Backend PID: $$(cat backend/.backend.pid)"

local-run-frontend:
	@echo "▶️  Starting frontend preview..."
	@cd frontend && nohup npm run preview > ../backend/logs/frontend.log 2>&1 </dev/null & echo $$! > .frontend.pid
	@echo "Frontend PID: $$(cat frontend/.frontend.pid)"

# Restart local services
restart: stop local-deploy
	@echo "♻️  Services restarted!"

# Stop all local services
stop:
	@echo "⏹️  Stopping services..."
	@if [ -f backend/.backend.pid ]; then \
		kill $$(cat backend/.backend.pid) 2>/dev/null || true; \
		rm backend/.backend.pid; \
		echo "Backend stopped"; \
	fi
	@if [ -f frontend/.frontend.pid ]; then \
		kill $$(cat frontend/.frontend.pid) 2>/dev/null || true; \
		rm frontend/.frontend.pid; \
		echo "Frontend stopped"; \
	fi
	@pkill -f "freezino-server" 2>/dev/null || true
	@pkill -f "vite preview" 2>/dev/null || true
	@echo "✅ All services stopped!"

# Check service status
status:
	@echo "📊 Service Status:"
	@echo ""
	@echo "Backend:"
	@if [ -f backend/.backend.pid ] && kill -0 $$(cat backend/.backend.pid) 2>/dev/null; then \
		echo "  Status: ✅ Running (PID: $$(cat backend/.backend.pid))"; \
		curl -sf http://localhost:3000/api/health | python3 -c "import sys,json; data=json.load(sys.stdin); print('  Health:', data.get('status', 'unknown'))" 2>/dev/null || echo "  Health: ❌ Unreachable"; \
	else \
		echo "  Status: ❌ Stopped"; \
	fi
	@echo ""
	@echo "Frontend:"
	@if [ -f frontend/.frontend.pid ] && kill -0 $$(cat frontend/.frontend.pid) 2>/dev/null; then \
		echo "  Status: ✅ Running (PID: $$(cat frontend/.frontend.pid))"; \
		curl -sf http://localhost:4173 >/dev/null 2>&1 && echo "  Health: ✅ Reachable" || echo "  Health: ❌ Unreachable"; \
	else \
		echo "  Status: ❌ Stopped"; \
	fi
	@echo ""

# View backend logs
logs-backend:
	@echo "📜 Backend logs (Ctrl+C to exit):"
	@tail -f backend/logs/backend.log 2>/dev/null || echo "No backend logs found"

# View frontend logs
logs-frontend:
	@echo "📜 Frontend logs (Ctrl+C to exit):"
	@tail -f backend/logs/frontend.log 2>/dev/null || echo "No frontend logs found"

# Clean build artifacts and dependencies
clean:
	@echo "🧹 Cleaning project..."
	cd backend && rm -rf data/*.db
	cd frontend && rm -rf node_modules dist
	@echo "✅ Project cleaned!"

# Deploy both backend and frontend to production
deploy: deploy-backend deploy-frontend
	@echo "✅ Full deployment completed!"

# Deploy backend to production server
deploy-backend:
	@echo "🚀 Deploying backend to freezino.online..."
	@echo "📦 Creating source archive..."
	@tar czf backend-src.tar.gz -C backend cmd internal go.mod go.sum Makefile
	@echo "📤 Uploading to server..."
	@scp backend-src.tar.gz root@freezino.online:/opt/freezino/
	@rm backend-src.tar.gz
	@echo "🔨 Building on server..."
	@ssh root@freezino.online "\
		cd /opt/freezino/backend && \
		tar xzf ../backend-src.tar.gz && \
		rm ../backend-src.tar.gz && \
		go build -o freezino-server cmd/server/main.go && \
		echo '✅ Build complete'"
	@echo "🔄 Restarting backend service..."
	@ssh root@freezino.online "\
		systemctl restart freezino-backend 2>/dev/null || \
		(cd /opt/freezino/backend && pkill freezino-server || true && \
		nohup ./freezino-server > server.log 2>&1 </dev/null &) && \
		sleep 2"
	@echo "🏥 Health check..."
	@ssh root@freezino.online "curl -sf http://localhost:3000/api/health || echo 'Warning: Health check failed'"
	@echo "✅ Backend deployed successfully!"

# Deploy frontend to production server
deploy-frontend:
	@echo "🚀 Deploying frontend to freezino.online..."
	@echo "🔨 Building frontend..."
	@cd frontend && npm run build
	@echo "📊 Verifying shop images in build..."
	@if [ -d "frontend/dist/images" ]; then \
		echo "  ✅ Images directory found in build"; \
		echo "  📸 Total images: $$(find frontend/dist/images -type f | wc -l)"; \
	else \
		echo "  ⚠️  Warning: Images directory not found in build!"; \
	fi
	@echo "📦 Creating archive (includes images)..."
	@tar czf dist.tar.gz -C frontend dist
	@echo "📤 Uploading to server..."
	@scp dist.tar.gz root@freezino.online:/opt/freezino/frontend/
	@rm dist.tar.gz
	@echo "📂 Extracting on server..."
	@ssh root@freezino.online "\
		cd /opt/freezino/frontend && \
		rm -rf dist && \
		tar xzf dist.tar.gz && \
		rm dist.tar.gz && \
		echo 'Images on server:' && \
		find dist/images -type f 2>/dev/null | wc -l"
	@echo "✅ Frontend deployed successfully!"
