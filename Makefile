.PHONY: init dev clean deploy deploy-backend deploy-frontend

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
	@tar czf backend-src.tar.gz -C backend cmd internal go.mod go.sum
	@echo "📤 Uploading to server..."
	@scp backend-src.tar.gz root@freezino.online:/opt/freezino/
	@rm backend-src.tar.gz
	@echo "🔨 Building on server..."
	@ssh root@freezino.online "\
		cd /opt/freezino && \
		tar xzf backend-src.tar.gz -C backend --strip-components=0 && \
		rm backend-src.tar.gz && \
		cd backend && \
		go build -o freezino-server cmd/server/main.go"
	@echo "🔄 Restarting backend..."
	@ssh root@freezino.online "\
		pkill -9 freezino-server || true && \
		cd /opt/freezino/backend && \
		nohup ./freezino-server > server.log 2>&1 </dev/null & \
		sleep 2 && \
		curl -s http://localhost:3000/api/health && echo ''"
	@echo "✅ Backend deployed successfully!"

# Deploy frontend to production server
deploy-frontend:
	@echo "🚀 Deploying frontend to freezino.online..."
	@echo "🔨 Building frontend..."
	@cd frontend && npm run build
	@echo "📦 Creating archive..."
	@tar czf dist.tar.gz -C frontend dist
	@echo "📤 Uploading to server..."
	@scp dist.tar.gz root@freezino.online:/opt/freezino/frontend/
	@rm dist.tar.gz
	@echo "📂 Extracting on server..."
	@ssh root@freezino.online "\
		cd /opt/freezino/frontend && \
		rm -rf dist && \
		tar xzf dist.tar.gz && \
		rm dist.tar.gz"
	@echo "✅ Frontend deployed successfully!"
