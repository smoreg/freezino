# СТАРТ ИНСТРУКЦИЯ - Freezino

> Пошаговый план запуска Клодов. Следуй строго по порядку.

---

## ✅ ШАГ 1: ПЕРВЫЕ 4 КЛОДА (УЖЕ ЗАПУЩЕНЫ)

Эти Клоды уже работают над задачами из `PARALLEL_PROMPTS.md`:

| Клод | Задача | Ветка |
|------|--------|-------|
| Клод 1 | Backend Setup (Go + Fiber) | `claude/main-feature-011CUvjAWDDHWrR7yb7AmixU` |
| Клод 2 | Frontend Setup (React + Vite) | `claude/main-feature-011CUvjAWDDHWrR7yb7AmixU` |
| Клод 3 | Database (GORM + SQLite) | `claude/main-feature-011CUvjAWDDHWrR7yb7AmixU` |
| Клод 4 | Docker & DevOps | `claude/main-feature-011CUvjAWDDHWrR7yb7AmixU` |

**Действие**: Дождись пока все 4 закончат и запушат в одну ветку.

---

## 🔄 ШАГ 2: МЕРДЖ ФАЗЫ 0

Когда все 4 Клода закончат:

```bash
cd /home/user/freezino
git fetch origin
git merge claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

**Проверь что работает**:
```bash
# Backend
cd backend && make run  # Должен запуститься на :3000

# Frontend
cd frontend && npm run dev  # Должен запуститься на :5173

# Database
ls backend/*.db  # Должен появиться freezino.db

# Docker
docker-compose up  # Должно всё запуститься
```

Если всё ОК → переходи к Шагу 3.

---

## 📋 ШАГ 3: ФАЗА 1 - Auth & Core (4 КЛОДА)

Теперь запускай следующие 4 Клода **параллельно**:

### Клод 5:
```
Проект Freezino. Прочитай PLAN.md.

Реализуй Google OAuth в backend:
- Установи golang.org/x/oauth2
- Endpoints: GET /api/auth/google, GET /api/auth/google/callback
- JWT токены (access + refresh)
- Middleware для проверки токенов
- GET /api/auth/me - текущий юзер
- POST /api/auth/logout

Файлы: backend/internal/auth/, backend/internal/middleware/auth.go

Работай в ветке: claude/phase1-google-auth
Коммит: "feat(auth): implement Google OAuth authentication"
```

### Клод 6:
```
Проект Freezino. Прочитай PLAN.md.

Реализуй User API:
- GET /api/user/profile - профиль юзера
- PATCH /api/user/profile - обновить профиль
- GET /api/user/balance - баланс
- GET /api/user/stats - статистика (время работы, игры)
- GET /api/user/transactions - история транзакций
- GET /api/user/items - купленные предметы

Файлы: backend/internal/handler/user.go, backend/internal/service/user.go

Работай в ветке: claude/phase1-user-api
Коммит: "feat(user): add user profile and statistics API"
```

### Клод 7:
```
Проект Freezino. Прочитай PLAN.md.

Реализуй Auth UI:
- Страница /login с Google OAuth кнопкой
- Auth context/store (Zustand)
- Protected routes (redirect → /login)
- Token management (localStorage)
- Automatic token refresh
- Logout функционал

Файлы: frontend/src/pages/LoginPage.tsx, frontend/src/store/authStore.ts, frontend/src/components/ProtectedRoute.tsx

Работай в ветке: claude/phase1-auth-ui
Коммит: "feat(auth): add login page and auth state management"
```

### Клод 8:
```
Проект Freezino. Прочитай PLAN.md.

Реализуй Dashboard:
- Страница /dashboard
- Header с балансом и аватаром юзера
- Sidebar с навигацией (Игры, Магазин, Профиль, Работа)
- Карточки игр (пока заглушки)
- Responsive дизайн
- Loading states

Файлы: frontend/src/pages/DashboardPage.tsx, frontend/src/components/layout/{Header,Sidebar}.tsx, frontend/src/components/GameCard.tsx

Работай в ветке: claude/phase1-dashboard-ui
Коммит: "feat(dashboard): add dashboard layout and navigation"
```

**После завершения Фазы 1**:
```bash
git merge claude/phase1-google-auth
git merge claude/phase1-user-api
git merge claude/phase1-auth-ui
git merge claude/phase1-dashboard-ui
```

---

## 📋 ШАГ 4: ФАЗА 2 - Work System (4 КЛОДА)

### Клод 9:
```
Проект Freezino. Прочитай PLAN.md.

Реализуй Work API:
- POST /api/work/start - начать работу (создать WorkSession)
- GET /api/work/status - статус (осталось времени)
- POST /api/work/complete - завершить (начислить 500$, создать Transaction)
- GET /api/work/history - история работы
- Валидация: нельзя работать параллельно

Файлы: backend/internal/handler/work.go, backend/internal/service/work.go

Работай в ветке: claude/phase2-work-api
Коммит: "feat(work): add work system API"
```

### Клод 10:
```
Проект Freezino. Прочитай PLAN.md.

Реализуй статистику стран:
- JSON файл с 50+ странами (название, средняя зарплата/час)
- GET /api/stats/countries - список стран
- Функция расчета: сколько времени работать для 500$
- Сравнение с реальными зарплатами

Файлы: backend/internal/data/countries.json, backend/internal/handler/stats.go, backend/internal/service/stats.go

Работай в ветке: claude/phase2-country-stats
Коммит: "feat(stats): add country wage statistics"
```

### Клод 11:
```
Проект Freezino. Прочитай PLAN.md.

Реализуй Work Timer UI:
- Кнопка "Работать" (показывать при балансе = 0)
- Таймер 3 минуты (countdown)
- Прогресс бар с анимацией
- Нельзя закрыть пока идет таймер
- После завершения → показать модалку со статистикой

Файлы: frontend/src/components/WorkTimer.tsx, frontend/src/store/workStore.ts

Работай в ветке: claude/phase2-work-timer-ui
Коммит: "feat(work): add work timer UI component"
```

### Клод 12:
```
Проект Freezino. Прочитай PLAN.md.

Реализуй модалку статистики:
- Показывать после завершения работы
- Заработано: 500$
- Сравнение с 5-10 странами (таблица/список)
- "В США вам нужно было бы работать 16.7 минут, В России - 1.7 часа"
- Всего отработано времени
- Кнопка "Закрыть"

Файлы: frontend/src/components/StatsModal.tsx, frontend/src/pages/StatsPage.tsx

Работай в ветке: claude/phase2-stats-modal
Коммит: "feat(stats): add work completion statistics modal"
```

**После завершения Фазы 2**:
```bash
git merge claude/phase2-work-api
git merge claude/phase2-country-stats
git merge claude/phase2-work-timer-ui
git merge claude/phase2-stats-modal
```

---

## 📋 ШАГ 5: ФАЗА 3 - Games (6 КЛОДОВ)

### Клод 13:
```
Проект Freezino. Прочитай PLAN.md.

Реализуй Game Engine:
- Интерфейс Game (PlaceBet, Play, CalculateWin)
- Базовые функции: проверка баланса, создание GameSession, обновление баланса
- Crypto/rand для честных случайных чисел
- Transaction для ставок и выигрышей

Файлы: backend/internal/game/engine.go, backend/internal/game/game.go

Работай в ветке: claude/phase3-game-engine
Коммит: "feat(game): add game engine core"
```

### Клод 14:
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

Файлы: backend/internal/game/roulette.go, frontend/src/components/games/Roulette.tsx

Работай в ветке: claude/phase3-game-roulette
Коммит: "feat(game): add roulette game"
```

### Клод 15:
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

Файлы: backend/internal/game/slots.go, frontend/src/components/games/Slots.tsx

Работай в ветке: claude/phase3-game-slots
Коммит: "feat(game): add slots game"
```

### Клод 16:
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

Файлы: backend/internal/game/blackjack.go, frontend/src/components/games/Blackjack.tsx

Работай в ветке: claude/phase3-game-blackjack
Коммит: "feat(game): add blackjack game"
```

### Клод 17:
```
Проект Freezino. Прочитай PLAN.md.

Реализуй 3 простые игры (Backend + Frontend):

1. Crash: график с множителем (1.00x → crash)
   - POST /api/games/crash/bet

2. Hi-Lo: угадай выше/ниже
   - POST /api/games/hilo/bet

3. Wheel: колесо фортуны (сектора с множителями)
   - POST /api/games/wheel/spin

Файлы: backend/internal/game/{crash,hilo,wheel}.go, frontend/src/components/games/{Crash,HiLo,Wheel}.tsx

Работай в ветке: claude/phase3-mini-games
Коммит: "feat(game): add crash, hi-lo and wheel games"
```

### Клод 18:
```
Проект Freezino. Прочитай PLAN.md.

Реализуй историю игр (Backend + Frontend):

Backend:
- GET /api/games/history?game=&limit=&offset=
- GET /api/games/stats (всего игр, выигрышей, проигрышей, любимая игра)

Frontend:
- Страница /history
- Таблица с фильтрами (по игре, дате)
- Графики выигрышей/проигрышей (установи recharts)
- Статистика

Файлы: backend/internal/handler/game_history.go, frontend/src/pages/GameHistoryPage.tsx

Работай в ветке: claude/phase3-game-history
Коммит: "feat(game): add game history and statistics"
```

**После завершения Фазы 3**:
```bash
git merge claude/phase3-game-engine
git merge claude/phase3-game-roulette
git merge claude/phase3-game-slots
git merge claude/phase3-game-blackjack
git merge claude/phase3-mini-games
git merge claude/phase3-game-history
```

---

## 📋 ШАГ 6: ФАЗА 4 - Shop & Profile (5 КЛОДОВ)

### Клод 19:
```
Проект Freezino. Прочитай PLAN.md.

Реализуй Shop API:
- GET /api/shop/items?type=&rarity= - список предметов
- POST /api/shop/buy/:itemId - купить (проверка баланса, создать UserItem, Transaction)
- POST /api/shop/sell/:itemId - продать (50% от цены)
- GET /api/shop/my-items - мои предметы
- POST /api/shop/equip/:itemId - экипировать (только 1 на категорию)

Файлы: backend/internal/handler/shop.go, backend/internal/service/shop.go

Работай в ветке: claude/phase4-shop-api
Коммит: "feat(shop): add shop API"
```

### Клод 20:
```
Проект Freezino. Прочитай PLAN.md.

Дополни предметы магазина (если их меньше 50):
- 50+ предметов в категориях: clothing, car, house, accessory
- Цены от $500 до $1,000,000
- Rarity: common, rare, epic, legendary

Категории:
- 15 одежды ($500-$50k)
- 10 машин ($1k-$500k)
- 10 домов ($2k-$1M)
- 15+ аксессуаров ($500-$20k)

Файлы: backend/internal/database/items_seed.go (обновить или создать)

Работай в ветке: claude/phase4-shop-items
Коммит: "feat(shop): expand shop items seed data to 50+ items"
```

### Клод 21:
```
Проект Freezino. Прочитай PLAN.md.

Реализуй Shop UI:
- Страница /shop
- Сетка предметов (grid layout)
- Фильтры: по категории, по цене, по rarity
- Карточка предмета: изображение, название, цена, rarity badge, кнопка "Купить"
- Модалка подтверждения покупки
- Анимация при покупке (конфетти если редкий предмет)

Файлы: frontend/src/pages/ShopPage.tsx, frontend/src/components/shop/{ItemCard,ShopFilters,BuyModal}.tsx

Работай в ветке: claude/phase4-shop-ui
Коммит: "feat(shop): add shop UI and item purchasing"
```

### Клод 22:
```
Проект Freezino. Прочитай PLAN.md.

Реализуй профиль с визуализацией:
- Страница /profile
- Аватар юзера (композиция из предметов)
- Слои: фон (дом), персонаж с одеждой, машина
- Canvas или div с absolute positioning
- Показ всех экипированных предметов
- Статистика юзера (баланс, время работы, игры)

Файлы: frontend/src/pages/ProfilePage.tsx, frontend/src/components/profile/Avatar.tsx

Работай в ветке: claude/phase4-profile-avatar
Коммит: "feat(profile): add profile page with item visualization"
```

### Клод 23:
```
Проект Freezino. Прочитай PLAN.md.

Реализуй продажу предметов:
- Кнопка "Продать" на каждом предмете в списке "Мои предметы"
- Модалка при балансе = 0: "У вас нет денег. Продайте предметы чтобы продолжить игру"
- Показ цены продажи (50% от цены покупки)
- Подтверждение продажи
- Обновление баланса после продажи

Файлы: frontend/src/components/shop/SellModal.tsx, frontend/src/components/profile/MyItemsList.tsx

Работай в ветке: claude/phase4-sell-mechanism
Коммит: "feat(shop): add item selling mechanism"
```

**После завершения Фазы 4**:
```bash
git merge claude/phase4-shop-api
git merge claude/phase4-shop-items
git merge claude/phase4-shop-ui
git merge claude/phase4-profile-avatar
git merge claude/phase4-sell-mechanism
```

---

## 📋 ШАГ 7: ФАЗА 5 - Polish (4 КЛОДА)

### Клод 24:
```
Проект Freezino. Прочитай PLAN.md.

Добавь анимации (Framer Motion уже установлен):
- Fade in анимации для всех страниц
- Hover/active анимации для кнопок
- Particle effects при выигрыше (установи react-confetti)
- Loading skeletons (shimmer effect)
- Smooth transitions между роутами

Обновить все основные компоненты, добавить src/components/animations/

Работай в ветке: claude/phase5-animations
Коммит: "feat(ui): add animations and transitions"
```

### Клод 25:
```
Проект Freezino. Прочитай PLAN.md.

Добавь звуки:
- Фоновая музыка казино (с кнопкой вкл/выкл в Header)
- Звуки кнопок (click)
- Звуки игр (вращение рулетки, слотов)
- Звук монет при выигрыше
- Звук таймера работы
- Используй Howler.js или Web Audio API

Найди бесплатные звуки на freesound.org или используй programматически сгенерированные

Файлы: frontend/src/utils/sounds.ts, frontend/public/sounds/

Работай в ветке: claude/phase5-sounds
Коммит: "feat(ui): add sound effects and background music"
```

### Клод 26:
```
Проект Freezino. Прочитай PLAN.md.

Адаптация под mobile (Tailwind responsive):
- Все страницы работают на телефонах (breakpoints: sm, md, lg, xl)
- Мобильное меню (hamburger burger menu)
- Touch-friendly controls для игр (увеличенные кнопки)
- Responsive grid для магазина (1 колонка на mobile, 2 на tablet, 4 на desktop)
- Проверка на всех breakpoints

Обновить все компоненты с Tailwind responsive классами (sm:, md:, lg:)

Работай в ветке: claude/phase5-responsive
Коммит: "feat(ui): add responsive mobile design"
```

### Клод 27:
```
Проект Freezino. Прочитай PLAN.md.

Улучши error handling и UX:
- Toast notifications (установи react-hot-toast)
- Error boundary (React) для отлова ошибок
- Валидация форм (установи react-hook-form + zod)
- Graceful degradation (показывать сообщение если offline)
- Loading states везде где идут API запросы
- Error pages: 404, 500

Файлы: frontend/src/components/ErrorBoundary.tsx, frontend/src/pages/{NotFound,Error}.tsx

Работай в ветке: claude/phase5-error-handling
Коммит: "feat(ui): improve error handling and user experience"
```

**После завершения Фазы 5**:
```bash
git merge claude/phase5-animations
git merge claude/phase5-sounds
git merge claude/phase5-responsive
git merge claude/phase5-error-handling
```

---

## 📋 ШАГ 8: ФАЗА 6 - Testing & Deploy (4 КЛОДА)

### Клод 28:
```
Проект Freezino. Прочитай PLAN.md.

Напиши тесты для backend:
- Unit тесты для game logic (roulette, slots, blackjack)
- Integration тесты для API endpoints
- Тесты auth flow
- Тесты транзакций и баланса
- Coverage > 70%
- Используй testify и go test

Файлы: backend/internal/**/*_test.go

Работай в ветке: claude/phase6-backend-tests
Коммит: "test(backend): add unit and integration tests"
```

### Клод 29:
```
Проект Freezino. Прочитай PLAN.md.

Напиши тесты для frontend:
- Unit тесты компонентов (установи vitest, @testing-library/react)
- Integration тесты
- E2E тесты (установи @playwright/test):
  * Login flow
  * Play game flow
  * Buy item flow
  * Work flow
- Snapshot тесты для UI компонентов

Файлы: frontend/src/**/*.test.tsx, frontend/e2e/**/*.spec.ts

Работай в ветке: claude/phase6-frontend-tests
Коммит: "test(frontend): add unit and e2e tests"
```

### Клод 30:
```
Проект Freezino. Прочитай PLAN.md.

Оптимизируй performance:
Frontend:
- Lazy loading для роутов (React.lazy + Suspense)
- Code splitting
- Image optimization
- Bundle size анализ (установи vite-bundle-visualizer)

Backend:
- Database индексы для часто используемых queries
- Query optimization
- Кэширование (опционально Redis)

Обновить: frontend/vite.config.ts, frontend/src/App.tsx, backend queries

Работай в ветке: claude/phase6-performance
Коммит: "perf: optimize frontend and backend performance"
```

### Клод 31:
```
Проект Freezino. Прочитай PLAN.md.

Напиши документацию:
- README.md: установка, запуск dev, запуск prod, деплой
- API документация (создай OpenAPI/Swagger spec или используй комментарии)
- CONTRIBUTING.md: как контрибьютить
- docs/ARCHITECTURE.md: архитектура проекта (диаграммы)
- User Guide: как пользоваться приложением

Файлы: README.md, docs/, openapi.yaml

Работай в ветке: claude/phase6-documentation
Коммит: "docs: add comprehensive project documentation"
```

**После завершения Фазы 6**:
```bash
git merge claude/phase6-backend-tests
git merge claude/phase6-frontend-tests
git merge claude/phase6-performance
git merge claude/phase6-documentation
```

---

## 🎉 ФИНАЛЬНЫЙ ШАГ: Запуск проекта

```bash
# Проверка
cd /home/user/freezino

# Запуск через Docker
docker-compose up --build

# Или запуск отдельно
cd backend && make run  # :3000
cd frontend && npm run dev  # :5173

# Production деплой
docker-compose -f docker-compose.prod.yml up -d
```

**Открой в браузере**: http://localhost (или http://localhost:5173 для dev)

---

## 📊 ПРОГРЕСС ТРЕКИНГ

### Фаза 0: Setup ✅ (ЗАВЕРШЕНА)
- [x] Backend Setup
- [x] Frontend Setup
- [x] Database Models
- [x] Docker & DevOps

### Фаза 1: Auth & Core ⏳ (ТЕКУЩАЯ)
- [ ] Google OAuth
- [ ] User API
- [ ] Auth UI
- [ ] Dashboard UI

### Фаза 2: Work System
- [ ] Work API
- [ ] Country Stats
- [ ] Work Timer UI
- [ ] Stats Modal

### Фаза 3: Games
- [ ] Game Engine
- [ ] Roulette
- [ ] Slots
- [ ] Blackjack
- [ ] Mini Games
- [ ] Game History

### Фаза 4: Shop & Profile
- [ ] Shop API
- [ ] Shop Items
- [ ] Shop UI
- [ ] Profile & Avatar
- [ ] Sell Mechanism

### Фаза 5: Polish
- [ ] Animations
- [ ] Sounds
- [ ] Responsive
- [ ] Error Handling

### Фаза 6: Testing & Deploy
- [ ] Backend Tests
- [ ] Frontend Tests
- [ ] Performance
- [ ] Documentation

---

## ⚡ БЫСТРЫЙ СПРАВОЧНИК

**Сколько времени?**
- Последовательно: 31-42 дня
- 4 Клода параллельно: 8-10 дней
- Все 31 Клод параллельно: 2-3 дня

**Сколько всего Клодов?**
- Фаза 0: 4 (✅ сделано)
- Фаза 1: 4
- Фаза 2: 4
- Фаза 3: 6
- Фаза 4: 5
- Фаза 5: 4
- Фаза 6: 4
- **Итого: 31 Клод**

**Какие файлы смотреть?**
- `START_HERE.md` ← ты здесь
- `PLAN.md` ← полная спецификация проекта
- `PHASES.md` ← детальное описание фаз

**Проблемы?**
- Конфликты? Не должно быть - каждый Клод работает с разными файлами
- Ошибки? Каждый Клод должен протестировать свою часть перед коммитом
- Вопросы? Читай PLAN.md и PHASES.md

---

**Начинай с Шага 3!** (Фаза 0 уже завершена) 🚀
