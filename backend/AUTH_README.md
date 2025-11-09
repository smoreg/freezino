# Authentication System

## Test User Credentials

**Username:** `testuser123`
**Password:** `testuser123`

## Features

- ✅ Username/Password authentication (local auth)
- ✅ Google OAuth authentication
- ✅ JWT tokens (access + refresh)
- ✅ Automatic test user seeding
- ✅ Password hashing with bcrypt

## API Endpoints

### Registration
```bash
POST /api/auth/register
Content-Type: application/json

{
  "username": "username",
  "email": "email@example.com",
  "password": "password",
  "name": "Full Name"
}
```

### Login
```bash
POST /api/auth/login
Content-Type: application/json

{
  "username": "testuser123",
  "password": "testuser123"
}
```

Response:
```json
{
  "success": true,
  "data": {
    "user": {...},
    "access_token": "eyJ...",
    "refresh_token": "eyJ..."
  }
}
```

### Using Tokens
```bash
GET /api/auth/me
Authorization: Bearer YOUR_ACCESS_TOKEN
```

### Token Refresh
```bash
POST /api/auth/refresh
Content-Type: application/json

{
  "refresh_token": "YOUR_REFRESH_TOKEN"
}
```

## Database Seeding

Test user is automatically created on server startup:
- Username: `testuser123`
- Password: `testuser123` (hashed with bcrypt)
- Email: `testuser123@test.com`
- Initial balance: 1000

To reset the database:
```bash
rm data/freezino.db
make run
```

## Development

The test user is created in `internal/database/seed.go`:
- Password is hashed using bcrypt
- Username is unique (enforced by database index)
- Email is unique (enforced by database index)
