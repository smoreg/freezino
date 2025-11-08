# Freezino API Documentation

Quick reference guide for Freezino API endpoints.

> For detailed API specification, see [OpenAPI documentation](./openapi.yaml).

## Base URL

**Development**: `http://localhost:3000/api`

## Authentication

Most endpoints require JWT authentication via the `Authorization` header:

```
Authorization: Bearer <your_jwt_token>
```

Get tokens via Google OAuth flow:
1. Redirect to `/api/auth/google`
2. Handle callback at `/api/auth/google/callback`
3. Store returned access token

## API Endpoints

### 🏥 Health Check

#### GET `/health`
Check API status.

**Response**:
```json
{
  "status": "ok",
  "timestamp": "2025-11-08T10:00:00Z"
}
```

---

### 🔐 Authentication

#### GET `/auth/google`
Initiate Google OAuth login flow.

**Response**: Redirects to Google OAuth consent screen

#### GET `/auth/google/callback`
Google OAuth callback handler.

**Query Params**:
- `code` - OAuth authorization code

**Response**: Redirects to frontend with tokens

#### POST `/auth/refresh`
Refresh access token using refresh token.

**Request**:
```json
{
  "refreshToken": "your_refresh_token"
}
```

**Response**:
```json
{
  "accessToken": "new_access_token"
}
```

#### GET `/auth/me` 🔒
Get current authenticated user.

**Response**:
```json
{
  "id": 1,
  "email": "user@example.com",
  "name": "John Doe",
  "balance": 1500.50,
  "avatar": "https://...",
  "createdAt": "2025-01-01T00:00:00Z"
}
```

#### POST `/auth/logout` 🔒
Logout and invalidate session.

**Response**: `200 OK`

---

### 👤 User

#### GET `/user/profile` 🔒
Get user profile.

**Response**: User object (same as `/auth/me`)

#### PATCH `/user/profile` 🔒
Update user profile.

**Request**:
```json
{
  "displayName": "New Name",
  "avatar": "https://new-avatar-url.com/image.jpg"
}
```

#### GET `/user/balance` 🔒
Get current balance.

**Response**:
```json
{
  "balance": 2500.75
}
```

#### GET `/user/stats` 🔒
Get user statistics.

**Response**:
```json
{
  "totalWorkTime": 180,
  "totalEarned": 5000,
  "totalLost": 3500,
  "gamesPlayed": 42,
  "favoriteGame": "roulette"
}
```

#### GET `/user/transactions` 🔒
Get transaction history.

**Query Params**:
- `limit` (default: 20)
- `offset` (default: 0)

**Response**:
```json
[
  {
    "id": 1,
    "type": "work",
    "amount": 500,
    "description": "Work session completed",
    "createdAt": "2025-11-08T10:00:00Z"
  }
]
```

#### GET `/user/items` 🔒
Get user's purchased items.

**Response**:
```json
[
  {
    "id": 1,
    "itemId": 10,
    "equipped": true,
    "purchasePrice": 1000,
    "item": {
      "id": 10,
      "name": "Blue Jeans",
      "type": "clothing",
      "rarity": "common"
    }
  }
]
```

---

### 💼 Work System

#### POST `/work/start` 🔒
Start a work session (3 minutes to earn $500).

**Response**:
```json
{
  "id": 1,
  "startTime": "2025-11-08T10:00:00Z",
  "duration": 180,
  "earned": 500,
  "completed": false
}
```

#### GET `/work/status` 🔒
Get current work session status.

**Response**:
```json
{
  "isWorking": true,
  "timeRemaining": 120,
  "session": {
    "id": 1,
    "startTime": "2025-11-08T10:00:00Z",
    "duration": 180
  }
}
```

#### POST `/work/complete` 🔒
Complete work session and earn money.

**Response**:
```json
{
  "earned": 500,
  "newBalance": 2500.50
}
```

#### GET `/work/history` 🔒
Get work session history.

**Response**: Array of work sessions

---

### 📊 Statistics

#### GET `/stats/countries`
Get list of countries with wage data.

**Response**:
```json
[
  {
    "code": "US",
    "name": "United States",
    "averageWagePerHour": 30,
    "timeToEarn500": 16.67
  }
]
```

#### GET `/stats/countries/{code}`
Get specific country data.

**Path Params**:
- `code` - Country code (e.g., "US", "RU")

**Response**: Single country object

---

### 🛍️ Shop

#### GET `/shop/items`
Get shop items.

**Query Params**:
- `type` - Filter by type (clothing, car, house, accessory)
- `rarity` - Filter by rarity (common, rare, epic, legendary)

**Response**:
```json
[
  {
    "id": 1,
    "name": "Blue Jeans",
    "description": "Comfortable denim jeans",
    "price": 500,
    "type": "clothing",
    "rarity": "common",
    "imageUrl": "/images/items/jeans.png"
  }
]
```

#### POST `/shop/buy/{itemId}` 🔒
Purchase an item.

**Path Params**:
- `itemId` - Item ID to purchase

**Response**:
```json
{
  "success": true,
  "newBalance": 500,
  "userItem": {...}
}
```

**Error Responses**:
- `400` - Insufficient balance
- `409` - Item already owned

#### POST `/shop/sell/{userItemId}` 🔒
Sell an item for 50% of purchase price.

**Path Params**:
- `userItemId` - User item ID to sell

**Response**:
```json
{
  "success": true,
  "soldFor": 250,
  "newBalance": 750
}
```

#### GET `/shop/my-items` 🔒
Get user's owned items.

**Response**: Array of user items (same as `/user/items`)

#### POST `/shop/equip/{userItemId}` 🔒
Equip an item.

**Path Params**:
- `userItemId` - User item ID to equip

**Response**: `200 OK`

---

### 🎰 Games - Roulette

#### POST `/games/roulette/bet` 🔒
Place a roulette bet.

**Request**:
```json
{
  "amount": 100,
  "betType": "number",
  "value": 7
}
```

**Bet Types**:
- `number` - Specific number (0-36), requires `value`
- `red` - Red color
- `black` - Black color
- `odd` - Odd numbers
- `even` - Even numbers
- `dozen1` - 1-12
- `dozen2` - 13-24
- `dozen3` - 25-36

**Response**:
```json
{
  "winningNumber": 7,
  "won": true,
  "payout": 3500,
  "newBalance": 4500
}
```

#### GET `/games/roulette/history` 🔒
Get roulette bet history.

#### GET `/games/roulette/recent`
Get recent winning numbers.

**Response**:
```json
{
  "recentNumbers": [7, 14, 0, 23, 18]
}
```

---

### 🎰 Games - Slots

#### POST `/games/slots/spin` 🔒
Spin the slot machine.

**Request**:
```json
{
  "amount": 50
}
```

**Response**:
```json
{
  "reels": [
    ["🍒", "🍋", "🍊"],
    ["🍋", "🍋", "🍇"],
    ["💎", "⭐", "🍒"],
    ["🍒", "🍒", "🍋"],
    ["7️⃣", "💎", "⭐"]
  ],
  "won": true,
  "payout": 100,
  "newBalance": 1100
}
```

#### GET `/games/slots/payouts`
Get payout table.

---

### 🎰 Games - Other Games

#### POST `/games/crash/bet` 🔒
Place crash game bet.

**Request**:
```json
{
  "amount": 100,
  "cashoutMultiplier": 2.0
}
```

#### POST `/games/hilo/bet` 🔒
Place Hi-Lo bet.

**Request**:
```json
{
  "amount": 100,
  "prediction": "higher"
}
```

**Predictions**: `higher`, `lower`

#### POST `/games/wheel/spin` 🔒
Spin the fortune wheel.

**Request**:
```json
{
  "amount": 100
}
```

---

### 📈 Game History

#### GET `/games/history` 🔒
Get game session history.

**Query Params**:
- `game` - Filter by game type
- `limit` (default: 20)
- `offset` (default: 0)

**Response**:
```json
[
  {
    "id": 1,
    "gameType": "roulette",
    "betAmount": 100,
    "payout": 200,
    "won": true,
    "createdAt": "2025-11-08T10:00:00Z"
  }
]
```

#### GET `/games/stats` 🔒
Get game statistics.

**Response**:
```json
{
  "totalGames": 50,
  "totalWins": 20,
  "totalLosses": 30,
  "favoriteGame": "roulette",
  "biggestWin": 1000,
  "biggestLoss": 500
}
```

---

### 📧 Contact

#### POST `/contact`
Submit contact form message.

**Request**:
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "message": "Hello, I have a question..."
}
```

**Response**: `200 OK`

---

### 🎮 WebSocket - Blackjack

#### WS `/ws/blackjack` 🔒
WebSocket connection for live blackjack.

**Connection**:
```javascript
const ws = new WebSocket('ws://localhost:3000/ws/blackjack');
```

**Message Format**:
```json
{
  "action": "bet",
  "amount": 100
}
```

**Actions**:
- `bet` - Place bet
- `hit` - Take a card
- `stand` - Keep current hand
- `double` - Double bet and take one card
- `split` - Split pairs

**Server Messages**:
```json
{
  "type": "gameState",
  "playerHand": [{"rank": "A", "suit": "♠️"}],
  "dealerHand": [{"rank": "K", "suit": "♥️"}],
  "playerTotal": 11,
  "dealerTotal": 10
}
```

---

## Error Responses

All endpoints may return these error codes:

- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Missing or invalid authentication
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource conflict (e.g., item already owned)
- `500 Internal Server Error` - Server error

**Error Format**:
```json
{
  "error": true,
  "message": "Insufficient balance"
}
```

---

## Rate Limiting

*(Future feature)*

- 100 requests per minute per user
- 1000 requests per hour per IP

---

## Changelog

### Version 1.0.0 (2025-11-08)
- Initial API release
- All core endpoints implemented
- WebSocket support for Blackjack

---

For detailed schemas and examples, see the [OpenAPI specification](./openapi.yaml).
