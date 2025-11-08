# ФАЗА 0: Команды для параллельных Клодов

> Скопируй команду → дай Клоду → он всё сделает → закоммитит → запушит

---

## КЛОД 1: Backend Setup

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

---

## КЛОД 2: Frontend Setup

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

---

## КЛОД 3: Database

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

---

## КЛОД 4: Docker & DevOps

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

## После завершения всех 4:

```bash
git merge claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

Готово! Переходим к Фазе 1.
