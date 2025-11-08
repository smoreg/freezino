# Performance Optimization - Phase 6

This document describes all performance optimizations implemented in Phase 6.

## Frontend Optimizations

### 1. Lazy Loading & Code Splitting

**Implementation**: `frontend/src/App.tsx`

- All route components are now lazy-loaded using `React.lazy()`
- Wrapped routes with `<Suspense>` component with loading fallback
- Reduces initial bundle size by splitting code per route
- Users only download code for pages they visit

**Benefits**:
- Faster initial page load
- Smaller main bundle size
- Better user experience on slow connections

### 2. Build Optimization

**Implementation**: `frontend/vite.config.ts`

#### Manual Chunk Splitting
Separated vendor libraries into optimized chunks:
- `react-vendor`: React core libraries (react, react-dom, react-router-dom)
- `ui-vendor`: UI libraries (framer-motion, react-hot-toast, react-confetti)
- `i18n-vendor`: Internationalization libraries (i18next, react-i18next)
- `game-components`: All game components grouped together

**Benefits**:
- Better caching (vendor chunks rarely change)
- Parallel downloads of chunks
- Reduced redundancy

#### Build Settings
- **Minification**: Terser with aggressive compression
- **Tree-shaking**: Removes unused code
- **Console removal**: Drops console.logs in production
- **Source maps**: Disabled for smaller builds

#### Bundle Analyzer
- Installed `rollup-plugin-visualizer`
- Generates `dist/stats.html` after build
- Analyzes bundle size with gzip/brotli metrics
- Run `npm run build` to generate visualization

### 3. Dependency Optimization

**Pre-bundling** of frequently used dependencies:
- React ecosystem
- UI libraries
- i18n libraries

**Benefits**:
- Faster dev server startup
- Improved HMR (Hot Module Replacement)

---

## Backend Optimizations

### 1. Database Indexes

Composite indexes added for frequently queried columns to speed up database operations.

#### GameSession Model
**File**: `backend/internal/model/game_session.go`

- `idx_user_game_type`: Composite index on (user_id, game_type)
  - Optimizes: Filtering games by user and type
  - Used in: Game history with filters

- `idx_user_created`: Composite index on (user_id, created_at)
  - Optimizes: Sorting user's games by date
  - Used in: Game history pagination

- `idx_user_win`: Index on win column
  - Optimizes: Finding wins/losses
  - Used in: Statistics calculations

#### Transaction Model
**File**: `backend/internal/model/transaction.go`

- `idx_user_type`: Composite index on (user_id, type)
  - Optimizes: Filtering transactions by user and type
  - Used in: Transaction history with filters

- `idx_user_created`: Composite index on (user_id, created_at)
  - Optimizes: Sorting user's transactions by date
  - Used in: Transaction history pagination

#### WorkSession Model
**File**: `backend/internal/model/work_session.go`

- `idx_user_completed`: Composite index on (user_id, completed_at)
  - Optimizes: Sorting work sessions by completion date
  - Used in: Work history queries

#### UserItem Model
**File**: `backend/internal/model/user_item.go`

- `idx_user_purchased`: Composite index on (user_id, purchased_at)
  - Optimizes: Sorting user's items by purchase date
  - Used in: Item inventory display

- `idx_user_equipped`: Composite index on (user_id, is_equipped)
  - Optimizes: Finding equipped items
  - Used in: Avatar rendering, profile display

### 2. Query Optimizations

Existing service implementations already use:
- **Select()**: Fetching only required columns (e.g., user balance)
- **Preload()**: Eager loading relations where needed
- **Pagination**: Limit/Offset for large result sets
- **Aggregations**: Using DB for COUNT/SUM instead of loading all records

---

## Performance Metrics

### Expected Improvements

#### Frontend
- **Initial bundle size**: ~40-50% reduction
- **Route load time**: 2-3x faster for lazy-loaded routes
- **Cache hit rate**: Improved due to chunk splitting

#### Backend
- **Query speed**: 2-10x faster for indexed queries
- **Database load**: Reduced due to efficient indexes
- **API response time**: 20-40% improvement on paginated endpoints

### Monitoring

To analyze bundle size:
```bash
cd frontend
npm run build
# Open dist/stats.html in browser
```

To test database performance:
```bash
# Enable SQLite query logging
# Add to backend config: db.Debug() for development
```

---

## Best Practices for Future Development

### Frontend
1. Always use lazy loading for new routes
2. Group related components in manual chunks
3. Avoid importing entire libraries (use tree-shakeable imports)
4. Use React.memo() for expensive components
5. Implement virtualization for long lists

### Backend
1. Add indexes for columns used in WHERE clauses
2. Use composite indexes for multi-column filters
3. Always paginate large result sets
4. Use Select() to fetch only needed columns
5. Use Preload() carefully (avoid N+1 queries)
6. Consider caching for expensive, repeated queries

---

## Migration Required

After pulling these changes, run database migration to create new indexes:

```bash
cd backend
# Indexes will be created automatically by GORM AutoMigrate
make run
```

Or manually run migration if needed:
```bash
cd backend
go run cmd/server/main.go
```

---

## Notes

- All optimizations are backward compatible
- No breaking changes to API
- Frontend changes require rebuilding client bundle
- Backend changes auto-migrate on server start
- These optimizations work well with other parallel Phase 6 tasks
