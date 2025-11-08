# FREEZINO - 32 команды для параллельных Клодов

> Каждый Клод работает в СВОЕЙ ветке. Все параллельно. В конце мерджим → готово.

---

## КЛОД 1 → `claude/backend-setup`
```
Проект Freezino (казино-симулятор). Прочитай PLAN.md.

Создай backend на Go + Fiber:
- go mod init github.com/smoreg/freezino/backend
- Установи gofiber/fiber/v2
- Структура: cmd/server/main.go, internal/{config,middleware,router,handler}
- CORS, Logger, Recovery middleware
- GET /api/health endpoint
- Makefile (run/build/dev)

Коммит в ветку: claude/backend-setup
```

## КЛОД 2 → `claude/frontend-setup`
```
Проект Freezino. Прочитай PLAN.md.

Создай frontend React + Vite:
- npm create vite frontend -- --template react-ts
- Установи: tailwindcss, react-router-dom, axios, zustand, framer-motion
- Настрой Tailwind (casino colors)
- Структура папок + Layout + роуты

Коммит в ветку: claude/frontend-setup
```

## КЛОД 3 → `claude/database-models`
```
Проект Freezino. Прочитай PLAN.md.

Создай БД (SQLite + GORM):
- Модели: User, Transaction, Item, UserItem, WorkSession, GameSession
- Миграции + seed (тестовый юзер + 50 предметов магазина)
- backend/internal/{model,database}

Коммит в ветку: claude/database-models
```

## КЛОД 4 → `claude/docker-devops`
```
Проект Freezino. Прочитай PLAN.md.

Создай Docker:
- backend/Dockerfile (multi-stage Go)
- frontend/Dockerfile (build + nginx)
- docker-compose.yml + docker-compose.prod.yml
- nginx config (proxy /api)

Коммит в ветку: claude/docker-devops
```

## КЛОД 5 → `claude/google-auth`
```
Проект Freezino. Прочитай PLAN.md.

Google OAuth backend:
- GET /api/auth/google, /api/auth/google/callback
- JWT токены
- Middleware auth
- internal/auth/

Коммит в ветку: claude/google-auth
```

## КЛОД 6 → `claude/user-api`
```
Проект Freezino. Прочитай PLAN.md.

User API backend:
- GET /api/user/{profile,balance,stats,transactions,items}
- PATCH /api/user/profile
- internal/handler/user.go

Коммит в ветку: claude/user-api
```

## КЛОД 7 → `claude/auth-ui`
```
Проект Freezino. Прочитай PLAN.md.

Auth UI:
- /login страница (Google OAuth button)
- Auth store (Zustand)
- Protected routes
- src/pages/LoginPage.tsx, src/store/authStore.ts

Коммит в ветку: claude/auth-ui
```

## КЛОД 8 → `claude/dashboard-ui`
```
Проект Freezino. Прочитай PLAN.md.

Dashboard UI:
- /dashboard страница
- Header (баланс, аватар)
- Sidebar навигация
- Карточки игр (заглушки)

Коммит в ветку: claude/dashboard-ui
```

## КЛОД 9 → `claude/work-api`
```
Проект Freezino. Прочитай PLAN.md.

Work API:
- POST /api/work/{start,complete}
- GET /api/work/{status,history}
- Начисление 500$ за 3 мин
- internal/handler/work.go

Коммит в ветку: claude/work-api
```

## КЛОД 10 → `claude/country-stats`
```
Проект Freezino. Прочитай PLAN.md.

Country stats:
- JSON с 50+ странами (зарплаты/час)
- GET /api/stats/countries
- Расчет времени работы для 500$
- internal/data/countries.json

Коммит в ветку: claude/country-stats
```

## КЛОД 11 → `claude/work-timer-ui`
```
Проект Freezino. Прочитай PLAN.md.

Work Timer UI:
- Кнопка "Работать"
- Таймер 3 минуты
- Прогресс бар
- src/components/WorkTimer.tsx

Коммит в ветку: claude/work-timer-ui
```

## КЛОД 12 → `claude/stats-modal`
```
Проект Freezino. Прочитай PLAN.md.

Stats modal:
- Модалка после работы
- Сравнение с странами
- "В США: 16.7 мин, В России: 1.7 ч"
- src/components/StatsModal.tsx

Коммит в ветку: claude/stats-modal
```

## КЛОД 13 → `claude/game-engine`
```
Проект Freezino. Прочитай PLAN.md.

Game Engine:
- Интерфейс Game
- PlaceBet, Play, CalculateWin
- Проверка баланса, транзакции
- internal/game/engine.go

Коммит в ветку: claude/game-engine
```

## КЛОД 14 → `claude/game-roulette`
```
Проект Freezino. Прочитай PLAN.md.

Рулетка (backend + frontend):
- POST /api/games/roulette/bet
- Европейская рулетка (0-36)
- Анимация + betting board
- internal/game/roulette.go, src/components/games/Roulette.tsx

Коммит в ветку: claude/game-roulette
```

## КЛОД 15 → `claude/game-slots`
```
Проект Freezino. Прочитай PLAN.md.

Слоты (backend + frontend):
- POST /api/games/slots/spin
- 5 барабанов, символы
- Анимация
- internal/game/slots.go, src/components/games/Slots.tsx

Коммит в ветку: claude/game-slots
```

## КЛОД 16 → `claude/game-blackjack`
```
Проект Freezino. Прочитай PLAN.md.

Блэкджек (backend + frontend):
- WebSocket /ws/blackjack
- Hit, Stand, Double, Split
- Карточный стол UI
- internal/game/blackjack.go, src/components/games/Blackjack.tsx

Коммит в ветку: claude/game-blackjack
```

## КЛОД 17 → `claude/game-crash`
```
Проект Freezino. Прочитай PLAN.md.

Crash (backend + frontend):
- POST /api/games/crash/bet
- График с множителем
- internal/game/crash.go, src/components/games/Crash.tsx

Коммит в ветку: claude/game-crash
```

## КЛОД 18 → `claude/game-hilo`
```
Проект Freezino. Прочитай PLAN.md.

Hi-Lo (backend + frontend):
- POST /api/games/hilo/bet
- Угадай выше/ниже
- internal/game/hilo.go, src/components/games/HiLo.tsx

Коммит в ветку: claude/game-hilo
```

## КЛОД 19 → `claude/game-wheel`
```
Проект Freezino. Прочитай PLAN.md.

Wheel (backend + frontend):
- POST /api/games/wheel/spin
- Колесо фортуны
- internal/game/wheel.go, src/components/games/Wheel.tsx

Коммит в ветку: claude/game-wheel
```

## КЛОД 20 → `claude/game-history`
```
Проект Freezino. Прочитай PLAN.md.

Game history:
- GET /api/games/history
- GET /api/games/stats
- Страница /history с фильтрами и графиками
- internal/handler/game_history.go, src/pages/GameHistoryPage.tsx

Коммит в ветку: claude/game-history
```

## КЛОД 21 → `claude/shop-api`
```
Проект Freezino. Прочитай PLAN.md.

Shop API:
- GET /api/shop/items
- POST /api/shop/{buy,sell,equip}/:itemId
- GET /api/shop/my-items
- internal/handler/shop.go

Коммит в ветку: claude/shop-api
```

## КЛОД 22 → `claude/shop-items`
```
Проект Freezino. Прочитай PLAN.md.

Shop items seed:
- 50+ предметов (одежда, машины, дома, аксессуары)
- Цены $500-$1M
- Rarity levels
- internal/database/items_seed.go

Коммит в ветку: claude/shop-items
```

## КЛОД 23 → `claude/shop-ui`
```
Проект Freezino. Прочитай PLAN.md.

Shop UI:
- Страница /shop
- Сетка предметов
- Фильтры
- Карточки + покупка
- src/pages/ShopPage.tsx

Коммит в ветку: claude/shop-ui
```

## КЛОД 24 → `claude/profile-avatar`
```
Проект Freezino. Прочитай PLAN.md.

Profile + Avatar:
- Страница /profile
- Визуализация экипированных предметов (дом, одежда, машина)
- Canvas/div композиция
- src/pages/ProfilePage.tsx, src/components/profile/Avatar.tsx

Коммит в ветку: claude/profile-avatar
```

## КЛОД 25 → `claude/shop-sell`
```
Проект Freezino. Прочитай PLAN.md.

Sell mechanism:
- Кнопка продать
- Модалка при балансе = 0
- Цена продажи 50%
- src/components/shop/SellModal.tsx

Коммит в ветку: claude/shop-sell
```

## КЛОД 26 → `claude/animations`
```
Проект Freezino. Прочитай PLAN.md.

Animations:
- Framer Motion на всех страницах
- Particle effects при выигрыше
- Loading skeletons
- Button animations
- src/components/animations/

Коммит в ветку: claude/animations
```

## КЛОД 27 → `claude/sounds`
```
Проект Freezino. Прочитай PLAN.md.

Sounds:
- Фоновая музыка
- Звуки кнопок, игр, монет
- Howler.js
- src/utils/sounds.ts, public/sounds/

Коммит в ветку: claude/sounds
```

## КЛОД 28 → `claude/responsive`
```
Проект Freezino. Прочитай PLAN.md.

Responsive design:
- Mobile адаптация всех страниц
- Hamburger menu
- Touch-friendly controls
- Обновить все компоненты (Tailwind breakpoints)

Коммит в ветку: claude/responsive
```

## КЛОД 29 → `claude/error-handling`
```
Проект Freezino. Прочитай PLAN.md.

Error handling:
- Toast notifications
- Error boundaries
- Валидация форм
- 404/500 pages
- src/components/{ErrorBoundary,Toast}.tsx

Коммит в ветку: claude/error-handling
```

## КЛОД 30 → `claude/backend-tests`
```
Проект Freezino. Прочитай PLAN.md.

Backend tests:
- Unit тесты (game logic)
- Integration тесты (API)
- Coverage > 70%
- backend/**/*_test.go

Коммит в ветку: claude/backend-tests
```

## КЛОД 31 → `claude/frontend-tests`
```
Проект Freezino. Прочитай PLAN.md.

Frontend tests:
- Unit тесты (Vitest)
- E2E тесты (Playwright)
- frontend/src/**/*.test.tsx, frontend/e2e/

Коммит в ветку: claude/frontend-tests
```

## КЛОД 32 → `claude/performance`
```
Проект Freezino. Прочитай PLAN.md.

Performance:
- Lazy loading
- Code splitting
- Image optimization
- Bundle analyzer
- Обновить vite.config.ts

Коммит в ветку: claude/performance
```

## КЛОД 33 → `claude/deployment`
```
Проект Freezino. Прочитай PLAN.md.

Deployment:
- SSL nginx config
- PM2 config
- deploy.sh скрипт
- Health checks
- Логирование

Коммит в ветку: claude/deployment
```

## КЛОД 34 → `claude/documentation`
```
Проект Freezino. Прочитай PLAN.md.

Documentation:
- README.md (установка, запуск)
- API docs (Swagger)
- User guide
- Architecture diagrams

Коммит в ветку: claude/documentation
```

---

# После завершения ВСЕХ:

```bash
# Мердж всех веток
git merge claude/backend-setup
git merge claude/frontend-setup
git merge claude/database-models
git merge claude/docker-devops
git merge claude/google-auth
git merge claude/user-api
git merge claude/auth-ui
git merge claude/dashboard-ui
git merge claude/work-api
git merge claude/country-stats
git merge claude/work-timer-ui
git merge claude/stats-modal
git merge claude/game-engine
git merge claude/game-roulette
git merge claude/game-slots
git merge claude/game-blackjack
git merge claude/game-crash
git merge claude/game-hilo
git merge claude/game-wheel
git merge claude/game-history
git merge claude/shop-api
git merge claude/shop-items
git merge claude/shop-ui
git merge claude/profile-avatar
git merge claude/shop-sell
git merge claude/animations
git merge claude/sounds
git merge claude/responsive
git merge claude/error-handling
git merge claude/backend-tests
git merge claude/frontend-tests
git merge claude/performance
git merge claude/deployment
git merge claude/documentation

# Запуск
docker-compose up --build

# 🎉 ГОТОВО!
```

---

# Скрипт автомержа:

```bash
#!/bin/bash
branches=(
  "claude/backend-setup"
  "claude/frontend-setup"
  "claude/database-models"
  "claude/docker-devops"
  "claude/google-auth"
  "claude/user-api"
  "claude/auth-ui"
  "claude/dashboard-ui"
  "claude/work-api"
  "claude/country-stats"
  "claude/work-timer-ui"
  "claude/stats-modal"
  "claude/game-engine"
  "claude/game-roulette"
  "claude/game-slots"
  "claude/game-blackjack"
  "claude/game-crash"
  "claude/game-hilo"
  "claude/game-wheel"
  "claude/game-history"
  "claude/shop-api"
  "claude/shop-items"
  "claude/shop-ui"
  "claude/profile-avatar"
  "claude/shop-sell"
  "claude/animations"
  "claude/sounds"
  "claude/responsive"
  "claude/error-handling"
  "claude/backend-tests"
  "claude/frontend-tests"
  "claude/performance"
  "claude/deployment"
  "claude/documentation"
)

for branch in "${branches[@]}"; do
  echo "Merging $branch..."
  git merge $branch --no-edit
done

echo "✅ All branches merged!"
```

Сохрани как `merge_all.sh`, дай права `chmod +x merge_all.sh`, запусти `./merge_all.sh`

**ГОТОВО!** 🚀
