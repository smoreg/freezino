# FREEZINO - Все команды для параллельных Клодов

> **Как использовать**: Скопируй команду → вставь Клоду → он сделает → коммит → пуш

---

# ФАЗА 0: Setup (4 Клода)

## КЛОД 0.1: Backend Setup
```
Проект Freezino (казино-симулятор). Прочитай PLAN.md и PHASES.md.

Создай backend на Go:
- Инициализируй Go модуль: github.com/smoreg/freezino/backend
- Установи Fiber framework
- Создай структуру: cmd/server/main.go, internal/{config,middleware,router,handler}
- Middleware: CORS, Logger, Recovery
- Health endpoint: GET /api/health
- Makefile с командами run/build/dev
- .env.example с PORT=3000

Коммит: "feat(backend): initialize Go backend with Fiber"
Пуш в: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

## КЛОД 0.2: Frontend Setup
```
Проект Freezino (казино-симулятор). Прочитай PLAN.md и PHASES.md.

Создай frontend на React:
- npm create vite@latest frontend -- --template react-ts
- Установи: tailwindcss, react-router-dom, axios, zustand, framer-motion
- Настрой Tailwind (цвета: primary #DC2626, secondary #FBBF24, dark #1F2937)
- Структура: src/{components,pages,layouts,hooks,store,services,types,utils}
- Layout: Header, Sidebar, Footer
- Роуты: /, /login, /404
- Axios service с interceptors

Коммит: "feat(frontend): initialize React app with Vite"
Пуш в: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

## КЛОД 0.3: Database
```
Проект Freezino (казино-симулятор). Прочитай PLAN.md и PHASES.md.

Создай database на SQLite + GORM:
- Модели: User, Transaction, Item, UserItem, WorkSession, GameSession
- backend/internal/model/ - все модели
- backend/internal/database/ - database.go, migrate.go, seed.go
- Seed данные: тестовый юзер (test@freezino.com) + 30+ предметов магазина
- Предметы: одежда, машины, дома (цены от $500 до $1,000,000)

Коммит: "feat(backend): add database models and migrations"
Пуш в: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

## КЛОД 0.4: Docker & DevOps
```
Проект Freezino (казино-симулятор). Прочитай PLAN.md и PHASES.md.

Создай Docker setup:
- backend/Dockerfile (multi-stage: Go build → Alpine)
- frontend/Dockerfile (npm build → nginx)
- docker-compose.yml (dev)
- docker-compose.prod.yml (production)
- docker/nginx/default.conf (proxy /api → backend:3000)
- deploy.sh скрипт
- .dockerignore для обоих

Коммит: "feat(devops): add Docker and Nginx configuration"
Пуш в: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

---

# ФАЗА 1: Auth & Core API (4 Клода)

## КЛОД 1.1: Google OAuth Backend
```
Проект Freezino. Прочитай PLAN.md.

Реализуй Google OAuth в backend:
- Установи golang.org/x/oauth2
- Endpoints: GET /api/auth/google, GET /api/auth/google/callback
- JWT токены (access + refresh)
- Middleware для проверки токенов
- GET /api/auth/me - текущий юзер
- POST /api/auth/logout

Файлы: internal/auth/, internal/middleware/auth.go

Коммит: "feat(auth): implement Google OAuth authentication"
Пуш в: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

## КЛОД 1.2: User API Backend
```
Проект Freezino. Прочитай PLAN.md.

Реализуй User API:
- GET /api/user/profile - профиль юзера
- PATCH /api/user/profile - обновить профиль
- GET /api/user/balance - баланс
- GET /api/user/stats - статистика (время работы, игры)
- GET /api/user/transactions - история транзакций
- GET /api/user/items - купленные предметы

Файлы: internal/handler/user.go, internal/service/user.go

Коммит: "feat(user): add user profile and statistics API"
Пуш в: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

## КЛОД 1.3: Auth UI Frontend
```
Проект Freezino. Прочитай PLAN.md.

Реализуй Auth UI:
- Страница /login с Google OAuth кнопкой
- Auth context/store (Zustand)
- Protected routes (redirect → /login)
- Token management (localStorage)
- Automatic token refresh
- Logout функционал

Файлы: src/pages/LoginPage.tsx, src/store/authStore.ts, src/components/ProtectedRoute.tsx

Коммит: "feat(auth): add login page and auth state management"
Пуш в: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

## КЛОД 1.4: Dashboard UI Frontend
```
Проект Freezino. Прочитай PLAN.md.

Реализуй Dashboard:
- Страница /dashboard
- Header с балансом и аватаром юзера
- Sidebar с навигацией (Игры, Магазин, Профиль, Работа)
- Карточки игр (пока заглушки)
- Responsive дизайн
- Loading states

Файлы: src/pages/DashboardPage.tsx, src/components/layout/{Header,Sidebar}.tsx, src/components/GameCard.tsx

Коммит: "feat(dashboard): add dashboard layout and navigation"
Пуш в: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

---

# ФАЗА 2: Work System (4 Клода)

## КЛОД 2.1: Work API Backend
```
Проект Freezino. Прочитай PLAN.md.

Реализуй Work API:
- POST /api/work/start - начать работу (создать WorkSession)
- GET /api/work/status - статус (осталось времени)
- POST /api/work/complete - завершить (начислить 500$, создать Transaction)
- GET /api/work/history - история работы
- Валидация: нельзя работать параллельно

Файлы: internal/handler/work.go, internal/service/work.go

Коммит: "feat(work): add work system API"
Пуш в: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

## КЛОД 2.2: Country Stats Backend
```
Проект Freezino. Прочитай PLAN.md.

Реализуй статистику стран:
- JSON файл с 50+ странами (название, средняя зарплата/час)
- GET /api/stats/countries - список стран
- Функция расчета: сколько времени работать для 500$
- Сравнение с реальными зарплатами

Файлы: internal/data/countries.json, internal/handler/stats.go, internal/service/stats.go

Коммит: "feat(stats): add country wage statistics"
Пуш в: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

## КЛОД 2.3: Work Timer UI Frontend
```
Проект Freezino. Прочитай PLAN.md.

Реализуй Work Timer UI:
- Кнопка "Работать" (показывать при балансе = 0)
- Таймер 3 минуты (countdown)
- Прогресс бар с анимацией
- Нельзя закрыть пока идет таймер
- После завершения → показать модалку со статистикой

Файлы: src/components/WorkTimer.tsx, src/store/workStore.ts

Коммит: "feat(work): add work timer UI component"
Пуш в: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

## КЛОД 2.4: Stats Modal Frontend
```
Проект Freezino. Прочитай PLAN.md.

Реализуй модалку статистики:
- Показывать после завершения работы
- Заработано: 500$
- Сравнение с 5-10 странами (таблица/список)
- "В США: 16.7 мин, В России: 1.7 часа"
- Всего отработано времени
- Кнопка "Закрыть"

Файлы: src/components/StatsModal.tsx, src/pages/StatsPage.tsx

Коммит: "feat(stats): add work completion statistics modal"
Пуш в: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

---

# ФАЗА 3: Games (6 Клодов)

## КЛОД 3.1: Game Engine Core Backend
```
Проект Freezino. Прочитай PLAN.md.

Реализуй Game Engine:
- Интерфейс Game (PlaceBet, Play, CalculateWin)
- Базовые функции: проверка баланса, создание GameSession, обновление баланса
- Crypto/rand для честных случайных чисел
- Transaction для ставок и выигрышей

Файлы: internal/game/engine.go, internal/game/game.go

Коммит: "feat(game): add game engine core"
Пуш в: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

## КЛОД 3.2: Roulette
```
Проект Freezino. Прочитай PLAN.md.

Реализуй Рулетку (Backend + Frontend):
Backend:
- POST /api/games/roulette/bet
- Европейская рулетка (0-36)
- Ставки: число, цвет (red/black), четность (odd/even), дюжины
- Расчет выигрышей

Frontend:
- Анимация вращения колеса
- Betting board с всеми ставками
- История выпавших чисел

Файлы: internal/game/roulette.go, src/components/games/Roulette.tsx

Коммит: "feat(game): add roulette game"
Пуш в: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

## КЛОД 3.3: Slots
```
Проект Freezino. Прочитай PLAN.md.

Реализуй Слоты (Backend + Frontend):
Backend:
- POST /api/games/slots/spin
- 5 барабанов, символы: 🍒🍋🍊🍇💎⭐7️⃣
- Комбинации и выплаты (3 в ряд, 4 в ряд, 5 в ряд)
- Линии выплат

Frontend:
- Анимация вращения барабанов
- Кнопка SPIN
- Таблица выплат
- Выбор ставки

Файлы: internal/game/slots.go, src/components/games/Slots.tsx

Коммит: "feat(game): add slots game"
Пуш в: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

## КЛОД 3.4: Blackjack
```
Проект Freezino. Прочитай PLAN.md.

Реализуй Блэкджек (Backend + Frontend):
Backend:
- WebSocket /ws/blackjack
- Логика блэкджека (дилер, игрок)
- Действия: Hit, Stand, Double, Split
- Расчет очков (туз = 1 или 11)

Frontend:
- Стол с картами
- Карты игрока и дилера
- Кнопки Hit, Stand, Double, Split
- Счет очков

Файлы: internal/game/blackjack.go, src/components/games/Blackjack.tsx

Коммит: "feat(game): add blackjack game"
Пуш в: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

## КЛОД 3.5: Mini Games (Crash, HiLo, Wheel)
```
Проект Freezino. Прочитай PLAN.md.

Реализуй 3 простые игры (Backend + Frontend):

1. Crash: график с множителем (1.00x → crash)
2. Hi-Lo: угадай выше/ниже
3. Wheel: колесо фортуны (сектора с множителями)

Backend:
- POST /api/games/crash/bet
- POST /api/games/hilo/bet
- POST /api/games/wheel/spin

Frontend: по компоненту на каждую игру

Файлы: internal/game/{crash,hilo,wheel}.go, src/components/games/{Crash,HiLo,Wheel}.tsx

Коммит: "feat(game): add crash, hi-lo and wheel games"
Пуш в: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

## КЛОД 3.6: Game History & Stats
```
Проект Freezino. Прочитай PLAN.md.

Реализуй историю игр (Backend + Frontend):

Backend:
- GET /api/games/history?game=&limit=&offset=
- GET /api/games/stats (всего игр, выигрышей, проигрышей, любимая игра)

Frontend:
- Страница /history
- Таблица с фильтрами (по игре, дате)
- Графики выигрышей/проигрышей (recharts)
- Статистика

Файлы: internal/handler/game_history.go, src/pages/GameHistoryPage.tsx

Коммит: "feat(game): add game history and statistics"
Пуш в: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

---

# ФАЗА 4: Shop & Profile (5 Клодов)

## КЛОД 4.1: Shop API Backend
```
Проект Freezino. Прочитай PLAN.md.

Реализуй Shop API:
- GET /api/shop/items?type=&rarity= - список предметов
- POST /api/shop/buy/:itemId - купить (проверка баланса, создать UserItem, Transaction)
- POST /api/shop/sell/:itemId - продать (50% от цены)
- GET /api/shop/my-items - мои предметы
- POST /api/shop/equip/:itemId - экипировать (только 1 на категорию)

Файлы: internal/handler/shop.go, internal/service/shop.go

Коммит: "feat(shop): add shop API"
Пуш в: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

## КЛОД 4.2: Item Seeding
```
Проект Freezino. Прочитай PLAN.md.

Создай предметы для магазина (если еще нет):
- 50+ предметов в categories: clothing, car, house, accessory
- Цены от $500 до $1,000,000
- Rarity: common, rare, epic, legendary
- Добавь в seed.go или отдельный items_seed.go

Категории:
- 15 одежды ($500-$50k)
- 10 машин ($1k-$500k)
- 10 домов ($2k-$1M)
- 15 аксессуаров ($500-$20k)

Файлы: internal/database/items_seed.go

Коммит: "feat(shop): add shop items seed data"
Пуш в: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

## КЛОД 4.3: Shop UI Frontend
```
Проект Freezino. Прочитай PLAN.md.

Реализуй Shop UI:
- Страница /shop
- Сетка предметов (grid)
- Фильтры: по категории, по цене, по rarity
- Карточка предмета: фото, название, цена, кнопка "Купить"
- Модалка подтверждения покупки
- Анимация при покупке

Файлы: src/pages/ShopPage.tsx, src/components/shop/{ItemCard,ShopFilters}.tsx

Коммит: "feat(shop): add shop UI and item purchasing"
Пуш в: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

## КЛОД 4.4: Profile & Avatar Frontend
```
Проект Freezino. Прочитай PLAN.md.

Реализуй профиль с визуализацией:
- Страница /profile
- Аватар юзера (композиция из предметов)
- Слои: фон (дом), персонаж (одежда), машина
- Canvas или div композиция
- Показ экипированных предметов
- Статистика юзера

Файлы: src/pages/ProfilePage.tsx, src/components/profile/Avatar.tsx

Коммит: "feat(profile): add profile page with item visualization"
Пуш в: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

## КЛОД 4.5: Sell Mechanism Frontend
```
Проект Freezino. Прочитай PLAN.md.

Реализуй продажу предметов:
- Кнопка "Продать" на каждом предмете в профиле
- Модалка при балансе = 0: "Продайте предметы чтобы играть"
- Показ цены продажи (50% от покупки)
- Подтверждение продажи
- Обновление баланса

Файлы: src/components/shop/SellModal.tsx, src/components/profile/MyItems.tsx

Коммит: "feat(shop): add item selling mechanism"
Пуш в: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

---

# ФАЗА 5: Polish & UX (4 Клода)

## КЛОД 5.1: Animations & Transitions
```
Проект Freezino. Прочитай PLAN.md.

Добавь анимации:
- Framer Motion для всех страниц (fade in)
- Анимации кнопок (hover, active)
- Particle effects при выигрыше (конфетти)
- Loading skeletons
- Smooth transitions между роутами

Файлы: src/components/animations/, обновление всех компонентов

Коммит: "feat(ui): add animations and transitions"
Пуш в: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

## КЛОД 5.2: Sounds & Music
```
Проект Freezino. Прочитай PLAN.md.

Добавь звуки:
- Фоновая музыка (с кнопкой выключения)
- Звуки кнопок (click)
- Звуки игр (рулетка, слоты)
- Звук монет при выигрыше
- Звук работы (timer)
- Web Audio API или Howler.js

Файлы: src/utils/sounds.ts, public/sounds/

Коммит: "feat(ui): add sound effects and background music"
Пуш в: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

## КЛОД 5.3: Responsive Design
```
Проект Freezino. Прочитай PLAN.md.

Адаптация под mobile:
- Все страницы должны работать на телефонах
- Мобильное меню (hamburger)
- Touch-friendly controls для игр
- Responsive grid для магазина
- Тестирование breakpoints: sm, md, lg, xl

Обновить все компоненты с Tailwind responsive классами

Коммит: "feat(ui): add responsive mobile design"
Пуш в: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

## КЛОД 5.4: Error Handling & UX
```
Проект Freezino. Прочитай PLAN.md.

Улучши UX:
- Toast notifications (react-hot-toast)
- Error boundaries (React)
- Валидация форм
- Graceful degradation (offline mode)
- Loading states везде
- Error pages (404, 500)

Файлы: src/components/ErrorBoundary.tsx, src/components/Toast.tsx

Коммит: "feat(ui): improve error handling and user experience"
Пуш в: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

---

# ФАЗА 6: Testing & Deploy (5 Клодов)

## КЛОД 6.1: Backend Tests
```
Проект Freezino. Прочитай PLAN.md.

Напиши тесты для backend:
- Unit тесты для game logic
- Integration тесты для API endpoints
- Тестирование auth
- Тестирование транзакций
- Coverage > 70%
- Используй testify

Файлы: backend/internal/game/*_test.go, backend/internal/handler/*_test.go

Коммит: "test(backend): add unit and integration tests"
Пуш в: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

## КЛОД 6.2: Frontend Tests
```
Проект Freezino. Прочитай PLAN.md.

Напиши тесты для frontend:
- Unit тесты компонентов (Vitest + Testing Library)
- Integration тесты
- E2E тесты (Playwright) - login, игра, покупка
- Snapshot тесты

Файлы: frontend/src/**/*.test.tsx, frontend/e2e/

Коммит: "test(frontend): add unit and e2e tests"
Пуш в: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

## КЛОД 6.3: Performance Optimization
```
Проект Freezino. Прочитай PLAN.md.

Оптимизируй перформанс:
- Lazy loading для роутов (React.lazy)
- Code splitting
- Image optimization
- Bundle size анализ (vite-bundle-analyzer)
- Backend query optimization (индексы)
- Redis caching (опционально)

Обновить: vite.config.ts, оптимизировать компоненты

Коммит: "perf: optimize frontend and backend performance"
Пуш в: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

## КЛОД 6.4: Deployment
```
Проект Freezino. Прочитай PLAN.md.

Подготовь к деплою:
- Обнови docker-compose.prod.yml
- SSL конфигурация для Nginx
- PM2 конфигурация (ecosystem.config.js)
- Скрипты деплоя (deploy.sh)
- Health checks
- Логирование (winston для backend)
- Мониторинг

Файлы: deploy.sh, ecosystem.config.js, docker/nginx/ssl.conf

Коммит: "feat(deploy): add production deployment configuration"
Пуш в: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

## КЛОД 6.5: Documentation
```
Проект Freezino. Прочитай PLAN.md.

Напиши документацию:
- README.md (установка, запуск, деплой)
- API документация (Swagger/OpenAPI)
- Инструкция для пользователей
- Contributing guide
- Архитектура проекта (diagrams)

Файлы: README.md, docs/API.md, docs/ARCHITECTURE.md, openapi.yaml

Коммит: "docs: add comprehensive project documentation"
Пуш в: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

---

# Итого: 27 Клодов для всего проекта

| Фаза | Клодов | Время (последовательно) | Время (параллельно) |
|------|--------|-------------------------|---------------------|
| 0. Setup | 4 | 3-5 дней | 1 день |
| 1. Auth | 4 | 5-7 дней | 1-2 дня |
| 2. Work | 4 | 4-5 дней | 1 день |
| 3. Games | 6 | 7-10 дней | 2-3 дня |
| 4. Shop | 5 | 5-6 дней | 1-2 дня |
| 5. Polish | 4 | 3-4 дня | 1 день |
| 6. Deploy | 5 | 4-5 дней | 1 день |
| **ИТОГО** | **32** | **31-42 дня** | **8-10 дней** |

---

# Как использовать:

1. **Последовательно (1 человек)**:
   - Запускай Клодов по одному
   - Каждый Клод завершает → коммитит → пушит
   - Переходишь к следующему

2. **Параллельно (N человек)**:
   - Каждому человеку дай по команде
   - Все работают одновременно
   - В конце мержишь ветки

3. **Мердж после каждой фазы**:
   ```bash
   git merge claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
   # Тестируешь
   # Переходишь к следующей фазе
   ```

**Конфликтов не будет** - каждый Клод работает с разными файлами!
