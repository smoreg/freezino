.PHONY: init dev clean

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
