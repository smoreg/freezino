# FREEZINO - Промпты для параллельных Клодов

> **Инструкция**: Скопируйте промпт для нужной ветки и дайте его отдельному Claude. Каждый Клод работает независимо над своей веткой.

---

## 📚 Контекст проекта (для всех Клодов)

**Freezino** - веб-приложение казино-симулятор для борьбы с игровой зависимостью.

### Основная идея:
- Пользователи играют на **виртуальные деньги** (псевдодоллары)
- Когда деньги кончаются → кнопка **"Работать"** (таймер 3 минуты = 500$)
- Показывает статистику: сколько времени нужно работать в разных странах для заработка 500$
- Можно тратить деньги на виртуальные покупки (дом, машина, одежда)
- **НЕЛЬЗЯ вводить реальные деньги** - только образовательная цель

### Технологический стек:
- **Backend**: Go (Golang) + Fiber/Gin + GORM + SQLite
- **Frontend**: React + TypeScript + Vite + TailwindCSS
- **Auth**: Google OAuth 2.0
- **Deploy**: Docker + Nginx

### Важные файлы для изучения:
- `PLAN.md` - полная спецификация проекта
- `PHASES.md` - разбивка на фазы и ветки

---

# ФАЗА 0: Infrastructure & Setup

## 🔧 ВЕТКА 0.1: Backend Scaffolding

```
КОНТЕКСТ:
Ты работаешь над проектом Freezino - казино-симулятор для борьбы с игровой зависимостью.

ТВОЯ ЗАДАЧА:
Создать базовую структуру backend на Go (Golang) с использованием Fiber framework.

ЧТО НУЖНО СДЕЛАТЬ:

1. Изучи файлы PLAN.md и PHASES.md в корне репозитория
2. Создай структуру backend приложения на Go
3. Настрой Fiber web framework
4. Реализуй базовые middleware (CORS, Logger, Recovery)
5. Создай health check endpoint
6. Настрой систему конфигурации через .env
7. Создай Makefile для удобной сборки

СТРУКТУРА ФАЙЛОВ (создай ВСЕ эти файлы):

backend/
├── cmd/
│   └── server/
│       └── main.go                 # Точка входа приложения
├── internal/
│   ├── config/
│   │   └── config.go              # Загрузка конфигурации из .env
│   ├── middleware/
│   │   ├── cors.go                # CORS middleware
│   │   ├── logger.go              # Логирование запросов
│   │   └── recovery.go            # Обработка паники
│   ├── router/
│   │   └── router.go              # Настройка роутов
│   └── handler/
│       └── health.go              # Health check handler
├── pkg/
│   └── response/
│       └── response.go            # Стандартные JSON ответы
├── go.mod                         # Go модуль
├── go.sum                         # Зависимости (будет автоматически)
├── Makefile                       # Команды для сборки/запуска
├── .env.example                   # Пример конфигурации
└── .gitignore                     # Игнорируемые файлы

ТЕХНИЧЕСКИЕ ТРЕБОВАНИЯ:

1. go.mod:
   - module: github.com/smoreg/freezino/backend
   - Go версия: 1.21+
   - Зависимости: github.com/gofiber/fiber/v2, github.com/joho/godotenv

2. main.go:
   - Загрузка конфигурации
   - Инициализация Fiber приложения
   - Подключение middleware
   - Настройка роутов
   - Graceful shutdown

3. Middleware:
   - CORS: разрешить http://localhost:5173 (Vite dev server)
   - Logger: логировать все HTTP запросы
   - Recovery: обрабатывать panic и возвращать 500

4. Health check:
   - GET /api/health
   - Ответ: {"status": "ok", "timestamp": "2024-01-01T00:00:00Z"}

5. Конфигурация (.env):
   - PORT=3000
   - ENV=development
   - CORS_ORIGIN=http://localhost:5173

6. Makefile команды:
   - make run         # Запуск приложения
   - make build       # Сборка бинарника
   - make dev         # Запуск с hot reload (air)
   - make test        # Запуск тестов
   - make clean       # Очистка

7. .gitignore:
   - .env
   - *.exe, *.dll, *.so, *.dylib
   - tmp/
   - bin/

ПРИМЕР СТРУКТУРЫ main.go:

package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/gofiber/fiber/v2"
    "github.com/smoreg/freezino/backend/internal/config"
    "github.com/smoreg/freezino/backend/internal/router"
)

func main() {
    // Загрузка конфигурации
    cfg := config.Load()

    // Создание Fiber приложения
    app := fiber.New(fiber.Config{
        AppName: "Freezino API",
    })

    // Настройка роутов
    router.Setup(app, cfg)

    // Graceful shutdown
    go func() {
        sigint := make(chan os.Signal, 1)
        signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
        <-sigint
        app.Shutdown()
    }()

    // Запуск сервера
    log.Printf("Starting server on port %s", cfg.Port)
    if err := app.Listen(":" + cfg.Port); err != nil {
        log.Fatal(err)
    }
}

КРИТЕРИИ ЗАВЕРШЕНИЯ:
✅ Backend запускается без ошибок
✅ GET /api/health возвращает {"status": "ok"}
✅ CORS настроен корректно
✅ Логи запросов работают
✅ Makefile команды работают

КОММИТ:
После завершения сделай коммит с сообщением:
"feat(backend): initialize Go backend with Fiber framework

- Setup project structure and Go modules
- Add Fiber web framework with basic middleware
- Implement CORS, Logger, Recovery middleware
- Add health check endpoint
- Create Makefile for build automation
- Add environment configuration"

Затем запуш в ветку: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

---

## 🎨 ВЕТКА 0.2: Frontend Scaffolding

```
КОНТЕКСТ:
Ты работаешь над проектом Freezino - казино-симулятор для борьбы с игровой зависимостью.

ТВОЯ ЗАДАЧА:
Создать базовую структуру frontend на React + TypeScript + Vite.

ЧТО НУЖНО СДЕЛАТЬ:

1. Изучи файлы PLAN.md и PHASES.md в корне репозитория
2. Создай React приложение с Vite
3. Настрой TypeScript
4. Установи и настрой TailwindCSS
5. Создай базовую структуру папок
6. Настрой роутинг (React Router)
7. Создай базовые layout компоненты
8. Настрой Axios для API запросов

КОМАНДЫ ДЛЯ ИНИЦИАЛИЗАЦИИ:

cd /home/user/freezino
npm create vite@latest frontend -- --template react-ts
cd frontend
npm install
npm install -D tailwindcss postcss autoprefixer
npx tailwindcss init -p
npm install react-router-dom axios zustand framer-motion

СТРУКТУРА ФАЙЛОВ (создай/настрой ВСЕ эти файлы):

frontend/
├── src/
│   ├── components/
│   │   ├── ui/                    # Базовые UI компоненты
│   │   │   ├── Button.tsx
│   │   │   ├── Card.tsx
│   │   │   └── Input.tsx
│   │   └── layout/                # Layout компоненты
│   │       ├── Header.tsx
│   │       ├── Sidebar.tsx
│   │       └── Footer.tsx
│   ├── pages/
│   │   ├── HomePage.tsx           # Главная страница
│   │   ├── LoginPage.tsx          # Страница входа
│   │   └── NotFoundPage.tsx       # 404 страница
│   ├── layouts/
│   │   ├── MainLayout.tsx         # Основной layout
│   │   └── AuthLayout.tsx         # Layout для auth страниц
│   ├── hooks/
│   │   └── index.ts               # Кастомные hooks
│   ├── store/
│   │   └── authStore.ts           # Zustand store для auth
│   ├── services/
│   │   └── api.ts                 # Axios instance и API методы
│   ├── types/
│   │   └── index.ts               # TypeScript типы
│   ├── utils/
│   │   └── constants.ts           # Константы
│   ├── App.tsx                    # Главный компонент с роутингом
│   ├── main.tsx                   # Точка входа
│   └── index.css                  # Глобальные стили + Tailwind
├── public/
│   └── vite.svg
├── index.html
├── package.json
├── vite.config.ts                 # Vite конфигурация
├── tailwind.config.js             # Tailwind конфигурация
├── tsconfig.json                  # TypeScript конфигурация
├── .env.example                   # Пример env переменных
└── .gitignore

ТЕХНИЧЕСКИЕ ТРЕБОВАНИЯ:

1. tailwind.config.js:
   - Добавь casino-тематические цвета:
     * primary: '#DC2626' (casino red)
     * secondary: '#FBBF24' (gold)
     * dark: '#1F2937'
   - Темная тема по умолчанию

2. vite.config.ts:
   - Настрой proxy для API: /api -> http://localhost:3000

3. src/services/api.ts:
   - Axios instance с baseURL из env
   - Interceptors для токенов
   - Обработка ошибок

4. src/App.tsx:
   - React Router с роутами:
     * / - HomePage
     * /login - LoginPage
     * * - NotFoundPage

5. .env.example:
   VITE_API_URL=http://localhost:3000

6. index.css:
   - @tailwind base/components/utilities
   - Темная тема body
   - Casino-стиль шрифты

7. MainLayout.tsx:
   - Header (logo, баланс, профиль)
   - Sidebar (навигация)
   - Main content area
   - Footer

ПРИМЕР src/services/api.ts:

import axios from 'axios';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:3000',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor для добавления токена
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor для обработки ошибок
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Redirect to login
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default api;

ПРИМЕР tailwind.config.js:

/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: '#DC2626',
        secondary: '#FBBF24',
        dark: {
          DEFAULT: '#1F2937',
          lighter: '#374151',
          darker: '#111827',
        },
      },
    },
  },
  plugins: [],
}

КРИТЕРИИ ЗАВЕРШЕНИЯ:
✅ npm run dev запускается на http://localhost:5173
✅ Роутинг работает (переходы между страницами)
✅ TailwindCSS стили применяются
✅ MainLayout отображается корректно
✅ API service готов к использованию

КОММИТ:
После завершения сделай коммит с сообщением:
"feat(frontend): initialize React app with Vite and TypeScript

- Setup Vite + React + TypeScript project
- Add TailwindCSS with casino theme colors
- Configure React Router for navigation
- Create basic layout components (Header, Sidebar, Footer)
- Setup Axios API service with interceptors
- Add Zustand store structure
- Create project folder structure"

Затем запуш в ветку: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

---

## 🗄️ ВЕТКА 0.3: Database Schema & Models

```
КОНТЕКСТ:
Ты работаешь над проектом Freezino - казино-симулятор для борьбы с игровой зависимостью.

ТВОЯ ЗАДАЧА:
Создать схему базы данных (SQLite) и модели для GORM.

ЧТО НУЖНО СДЕЛАТЬ:

1. Изучи файлы PLAN.md и PHASES.md в корне репозитория
2. Создай GORM модели для всех сущностей
3. Реализуй миграции
4. Создай seed данные (тестовый аккаунт + предметы для магазина)
5. Реализуй database utility функции

ПРЕДВАРИТЕЛЬНОЕ УСЛОВИЕ:
- Должна быть создана базовая структура backend (ветка 0.1)
- Если её нет, создай минимальную структуру backend/internal/

СТРУКТУРА ФАЙЛОВ:

backend/
├── internal/
│   ├── model/
│   │   ├── user.go                # User модель
│   │   ├── transaction.go         # Transaction модель
│   │   ├── item.go                # Item модель (предметы в магазине)
│   │   ├── user_item.go           # UserItem (купленные предметы)
│   │   ├── work_session.go        # WorkSession (сессии работы)
│   │   └── game_session.go        # GameSession (игровые сессии)
│   ├── database/
│   │   ├── database.go            # Инициализация БД
│   │   ├── migrate.go             # Миграции
│   │   └── seed.go                # Seed данные
│   └── repository/
│       ├── user_repository.go     # CRUD для User
│       └── repository.go          # Общий интерфейс
├── scripts/
│   └── seed.sh                    # Скрипт для seed данных
└── freezino.db                    # SQLite файл (создастся автоматически)

МОДЕЛИ (все поля в snake_case для БД):

1. User:
   - ID (uint, primary key)
   - GoogleID (string, unique, not null)
   - Email (string, unique, not null)
   - Name (string)
   - Avatar (string, URL)
   - Balance (float64, default: 1000)
   - CreatedAt (time.Time)
   - UpdatedAt (time.Time)

   Relations:
   - Transactions (hasMany)
   - UserItems (hasMany)
   - WorkSessions (hasMany)
   - GameSessions (hasMany)

2. Transaction:
   - ID (uint, primary key)
   - UserID (uint, foreign key)
   - Type (string: "work", "game_win", "game_loss", "purchase", "sale")
   - Amount (float64)
   - Description (string)
   - Metadata (string, JSON)
   - CreatedAt (time.Time)

   Relations:
   - User (belongsTo)

3. Item:
   - ID (uint, primary key)
   - Name (string)
   - Type (string: "clothing", "car", "house", "accessory")
   - Price (float64)
   - SellPrice (float64, 50% от Price)
   - ImageURL (string)
   - Description (string)
   - Rarity (string: "common", "rare", "epic", "legendary")
   - CreatedAt (time.Time)

4. UserItem:
   - ID (uint, primary key)
   - UserID (uint, foreign key)
   - ItemID (uint, foreign key)
   - IsEquipped (bool, default: false)
   - PurchasedAt (time.Time)

   Relations:
   - User (belongsTo)
   - Item (belongsTo)

5. WorkSession:
   - ID (uint, primary key)
   - UserID (uint, foreign key)
   - DurationSeconds (int, default: 180)
   - Earned (float64, default: 500)
   - CompletedAt (time.Time)

   Relations:
   - User (belongsTo)

6. GameSession:
   - ID (uint, primary key)
   - UserID (uint, foreign key)
   - GameType (string: "roulette", "slots", "blackjack", etc.)
   - Bet (float64)
   - Win (float64)
   - Profit (float64, Win - Bet)
   - Metadata (string, JSON для деталей игры)
   - CreatedAt (time.Time)

   Relations:
   - User (belongsTo)

SEED ДАННЫЕ:

1. Тестовый пользователь:
   - Email: test@freezino.com
   - Name: Test User
   - GoogleID: test_google_id_12345
   - Balance: 5000

2. Предметы магазина (минимум 30 штук):

Одежда:
- Простая футболка ($500)
- Джинсы ($800)
- Кроссовки ($1200)
- Деловой костюм ($5000)
- Дизайнерская одежда ($15000)
- Люксовый костюм ($50000)

Машины:
- Старый седан ($1000)
- Компактный хэтчбек ($5000)
- Седан бизнес-класса ($15000)
- Спортивный автомобиль ($50000)
- Премиальный SUV ($100000)
- Суперкар ($500000)

Дома:
- Комната в общежитии ($2000)
- Квартира-студия ($5000)
- Двухкомнатная квартира ($25000)
- Пентхаус ($100000)
- Загородный дом ($250000)
- Особняк ($1000000)

Аксессуары:
- Часы ($2000)
- Солнечные очки ($1500)
- Ювелирные изделия ($10000)
- Смартфон премиум ($3000)

ПРИМЕР model/user.go:

package model

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    GoogleID     string    `gorm:"uniqueIndex;not null" json:"google_id"`
    Email        string    `gorm:"uniqueIndex;not null" json:"email"`
    Name         string    `json:"name"`
    Avatar       string    `json:"avatar"`
    Balance      float64   `gorm:"default:1000" json:"balance"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`

    // Relations
    Transactions  []Transaction  `gorm:"foreignKey:UserID" json:"transactions,omitempty"`
    UserItems     []UserItem     `gorm:"foreignKey:UserID" json:"user_items,omitempty"`
    WorkSessions  []WorkSession  `gorm:"foreignKey:UserID" json:"work_sessions,omitempty"`
    GameSessions  []GameSession  `gorm:"foreignKey:UserID" json:"game_sessions,omitempty"`
}

ПРИМЕР database/database.go:

package database

import (
    "log"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "github.com/smoreg/freezino/backend/internal/model"
)

var DB *gorm.DB

func Connect() error {
    var err error
    DB, err = gorm.Open(sqlite.Open("freezino.db"), &gorm.Config{})
    if err != nil {
        return err
    }

    log.Println("Database connected successfully")
    return nil
}

func Migrate() error {
    return DB.AutoMigrate(
        &model.User{},
        &model.Transaction{},
        &model.Item{},
        &model.UserItem{},
        &model.WorkSession{},
        &model.GameSession{},
    )
}

ПРИМЕР database/seed.go:

package database

import (
    "github.com/smoreg/freezino/backend/internal/model"
)

func Seed() error {
    // Тестовый пользователь
    testUser := model.User{
        GoogleID: "test_google_id_12345",
        Email:    "test@freezino.com",
        Name:     "Test User",
        Balance:  5000,
    }
    DB.Create(&testUser)

    // Предметы магазина
    items := []model.Item{
        {Name: "Простая футболка", Type: "clothing", Price: 500, SellPrice: 250, Rarity: "common"},
        {Name: "Деловой костюм", Type: "clothing", Price: 5000, SellPrice: 2500, Rarity: "rare"},
        {Name: "Старый седан", Type: "car", Price: 1000, SellPrice: 500, Rarity: "common"},
        {Name: "Tesla Model 3", Type: "car", Price: 50000, SellPrice: 25000, Rarity: "epic"},
        {Name: "Квартира-студия", Type: "house", Price: 5000, SellPrice: 2500, Rarity: "common"},
        {Name: "Особняк", Type: "house", Price: 100000, SellPrice: 50000, Rarity: "legendary"},
        // ... добавь еще 24+ предметов
    }

    for _, item := range items {
        DB.Create(&item)
    }

    return nil
}

ТЕХНИЧЕСКИЕ ТРЕБОВАНИЯ:

1. Все модели должны использовать GORM теги
2. JSON теги для API ответов
3. Индексы для часто используемых полей (email, google_id)
4. Правильные foreign keys и relations
5. Default значения где нужно
6. Timestamps (CreatedAt, UpdatedAt) через gorm.Model или отдельно

КРИТЕРИИ ЗАВЕРШЕНИЯ:
✅ Все 6 моделей созданы
✅ Миграции работают без ошибок
✅ Seed данные создаются корректно
✅ Тестовый пользователь в БД
✅ Минимум 30 предметов в магазине
✅ freezino.db файл создается

КОММИТ:
После завершения сделай коммит с сообщением:
"feat(backend): add database models and migrations

- Create GORM models: User, Transaction, Item, UserItem, WorkSession, GameSession
- Setup SQLite database connection
- Implement auto-migrations
- Add seed data with test user and 30+ shop items
- Create repository pattern structure
- Add database utility functions"

Затем запуш в ветку: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

---

## 🐳 ВЕТКА 0.4: DevOps Setup (Docker & Nginx)

```
КОНТЕКСТ:
Ты работаешь над проектом Freezino - казино-симулятор для борьбы с игровой зависимостью.

ТВОЯ ЗАДАЧА:
Настроить Docker контейнеры и Nginx для деплоя.

ЧТО НУЖНО СДЕЛАТЬ:

1. Изучи файлы PLAN.md и PHASES.md в корне репозитория
2. Создай Dockerfile для backend (multi-stage build)
3. Создай Dockerfile для frontend (nginx)
4. Создай docker-compose.yml для development
5. Создай docker-compose.prod.yml для production
6. Настрой Nginx конфигурацию
7. Создай скрипты для деплоя

СТРУКТУРА ФАЙЛОВ:

/
├── backend/
│   ├── Dockerfile                 # Multi-stage build для Go
│   └── .dockerignore
├── frontend/
│   ├── Dockerfile                 # Build + Nginx
│   └── .dockerignore
├── docker/
│   └── nginx/
│       ├── nginx.conf             # Основная конфигурация
│       └── default.conf           # Site конфигурация
├── docker-compose.yml             # Development
├── docker-compose.prod.yml        # Production
├── deploy.sh                      # Скрипт деплоя
└── .dockerignore                  # Общий dockerignore

BACKEND DOCKERFILE (multi-stage):

# backend/Dockerfile
# Stage 1: Build
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/bin/server ./cmd/server

# Stage 2: Run
FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/bin/server .

# Copy .env if exists
COPY .env* ./

EXPOSE 3000

CMD ["./server"]

FRONTEND DOCKERFILE:

# frontend/Dockerfile
# Stage 1: Build
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./
RUN npm ci

# Copy source code
COPY . .

# Build
RUN npm run build

# Stage 2: Serve with Nginx
FROM nginx:alpine

# Copy built files
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx config
COPY docker/nginx/default.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]

DOCKER-COMPOSE.YML (Development):

version: '3.8'

services:
  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    ports:
      - "3000:3000"
    volumes:
      - ./backend:/app
      - ./backend/freezino.db:/root/freezino.db
    environment:
      - ENV=development
      - PORT=3000
      - CORS_ORIGIN=http://localhost:5173
    restart: unless-stopped

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    ports:
      - "80:80"
    depends_on:
      - backend
    restart: unless-stopped

DOCKER-COMPOSE.PROD.YML (Production):

version: '3.8'

services:
  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    ports:
      - "3000:3000"
    volumes:
      - ./data:/root/data
      - ./backend/freezino.db:/root/freezino.db
    environment:
      - ENV=production
      - PORT=3000
      - CORS_ORIGIN=https://freezino.yourdomain.com
    restart: always
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./docker/nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl  # SSL сертификаты
    depends_on:
      - backend
    restart: always

  # Nginx reverse proxy (опционально, для SSL)
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./docker/nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./docker/nginx/default.conf:/etc/nginx/conf.d/default.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - frontend
      - backend
    restart: always

NGINX КОНФИГУРАЦИЯ:

# docker/nginx/default.conf
server {
    listen 80;
    server_name _;

    root /usr/share/nginx/html;
    index index.html;

    # Frontend
    location / {
        try_files $uri $uri/ /index.html;
    }

    # API proxy to backend
    location /api {
        proxy_pass http://backend:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # WebSocket support
    location /ws {
        proxy_pass http://backend:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }

    # Gzip compression
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript application/x-javascript application/javascript application/xml+rss application/json;
}

.DOCKERIGNORE (backend):

# backend/.dockerignore
.git
.env
*.db
tmp/
*.log
node_modules/
.DS_Store

.DOCKERIGNORE (frontend):

# frontend/.dockerignore
.git
node_modules/
.env
.env.local
dist/
.DS_Store
*.log

DEPLOY СКРИПТ:

# deploy.sh
#!/bin/bash

set -e

echo "🚀 Starting Freezino deployment..."

# Pull latest code
echo "📥 Pulling latest changes..."
git pull origin main

# Build and start containers
echo "🏗️  Building Docker images..."
docker-compose -f docker-compose.prod.yml build

echo "🔄 Restarting containers..."
docker-compose -f docker-compose.prod.yml down
docker-compose -f docker-compose.prod.yml up -d

# Check status
echo "✅ Deployment complete!"
echo "📊 Container status:"
docker-compose -f docker-compose.prod.yml ps

# Show logs
echo "📜 Recent logs:"
docker-compose -f docker-compose.prod.yml logs --tail=50

ДОПОЛНИТЕЛЬНЫЕ ФАЙЛЫ:

1. .github/workflows/docker.yml (опционально):
   - GitHub Actions для автоматической сборки Docker образов

2. Makefile обновления (добавь в backend/Makefile):
   - make docker-build
   - make docker-run
   - make docker-stop

3. README.md секция для Docker:
   - Инструкции по запуску
   - Команды для development и production

ТЕХНИЧЕСКИЕ ТРЕБОВАНИЯ:

1. Multi-stage builds для минимизации размера образов
2. Правильные .dockerignore для исключения лишних файлов
3. Health checks в docker-compose
4. Volume mapping для persistence данных (БД)
5. Environment variables через .env файлы
6. Логирование в production
7. Graceful shutdown
8. Nginx gzip compression
9. WebSocket support в Nginx
10. SSL ready (закомментированная конфигурация)

КОМАНДЫ ДЛЯ ТЕСТИРОВАНИЯ:

# Development
docker-compose up --build

# Production
docker-compose -f docker-compose.prod.yml up --build -d

# Проверка логов
docker-compose logs -f backend
docker-compose logs -f frontend

# Остановка
docker-compose down

КРИТЕРИИ ЗАВЕРШЕНИЯ:
✅ Backend Dockerfile работает (multi-stage)
✅ Frontend Dockerfile работает (build + nginx)
✅ docker-compose.yml запускает оба сервиса
✅ Nginx проксирует /api на backend
✅ Frontend доступен на http://localhost
✅ Backend доступен через /api
✅ deploy.sh скрипт работает

КОММИТ:
После завершения сделай коммит с сообщением:
"feat(devops): add Docker and Nginx configuration

- Create multi-stage Dockerfile for Go backend
- Create Dockerfile for React frontend with Nginx
- Add docker-compose.yml for development
- Add docker-compose.prod.yml for production
- Configure Nginx reverse proxy for API
- Add WebSocket support in Nginx
- Create deployment script
- Add .dockerignore files"

Затем запуш в ветку: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU
```

---

# ВАЖНЫЕ ЗАМЕЧАНИЯ ДЛЯ ВСЕХ КЛОДОВ

## 🔄 Координация работы

1. **Независимость**: Каждая ветка работает с РАЗНЫМИ файлами - конфликтов не будет
2. **Порядок**: Лучше начать с веток по порядку (0.1 → 0.2 → 0.3 → 0.4)
3. **Зависимости**:
   - 0.1 и 0.2 полностью независимы (можно параллельно)
   - 0.3 желательно после 0.1 (но не критично)
   - 0.4 лучше после всех (но можно параллельно)

## ✅ Чеклист перед началом

Каждый Клод должен:
- [ ] Прочитать PLAN.md
- [ ] Прочитать PHASES.md
- [ ] Понять свою задачу
- [ ] Создать ВСЕ файлы из списка
- [ ] Протестировать работу
- [ ] Сделать коммит с правильным сообщением
- [ ] Запушить в ветку: claude/main-feature-011CUvjAWDDHWrR7yb7AmixU

## 🚫 Что НЕ делать

- ❌ Не редактируй чужие файлы
- ❌ Не меняй структуру, описанную в промпте
- ❌ Не пропускай файлы из списка
- ❌ Не коммить без тестирования
- ❌ Не пуш в другие ветки

## 📞 Вопросы

Если что-то непонятно:
1. Перечитай PLAN.md и PHASES.md
2. Посмотри примеры в промпте
3. Следуй структуре строго

## 🎯 Цель Фазы 0

После завершения ВСЕХ 4 веток должно получиться:
- ✅ Backend сервер на Go запускается
- ✅ Frontend на React запускается
- ✅ База данных с моделями и seed данными
- ✅ Docker контейнеры собираются и работают
- ✅ Всё готово к разработке Фазы 1 (аутентификация)

**Время выполнения Фазы 0**: 3-5 дней (последовательно) или 1-2 дня (параллельно)

---

## 📋 Следующие фазы

После завершения Фазы 0, будут созданы промпты для:
- **Фаза 1**: Auth & Core API (4 ветки)
- **Фаза 2**: Work System (4 ветки)
- **Фаза 3**: Game Engine (6 веток)
- **Фаза 4**: Shop System (5 веток)
- **Фаза 5**: Polish (4 ветки)
- **Фаза 6**: Deploy (5 веток)

Хотите промпты для следующей фазы? Дайте знать! 🚀
