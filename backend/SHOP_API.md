# Shop API Documentation

## Endpoints

### 1. Get Shop Items
```
GET /api/shop/items?type=&rarity=
```

**Query Parameters:**
- `type` (optional): Filter by item type (`clothing`, `car`, `house`, `accessories`)
- `rarity` (optional): Filter by rarity (`common`, `rare`, `epic`, `legendary`) - Will be supported when Rarity field is added

**Response:**
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": 1,
        "name": "Plain T-Shirt",
        "type": "clothing",
        "price": 500.00,
        "image_url": "https://...",
        "description": "A simple everyday t-shirt",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "count": 1
  }
}
```

### 2. Buy Item
```
POST /api/shop/buy/:itemId?user_id=1
```

**Path Parameters:**
- `itemId`: ID of the item to purchase

**Query Parameters:**
- `user_id`: ID of the user (in production, from JWT token)

**Response:**
```json
{
  "success": true,
  "data": {
    "user_item": {
      "id": 1,
      "user_id": 1,
      "item_id": 5,
      "purchased_at": "2024-01-01T12:00:00Z",
      "is_equipped": false,
      "item": {
        "id": 5,
        "name": "Designer Shirt",
        "type": "clothing",
        "price": 2000.00,
        "image_url": "https://...",
        "description": "Stylish designer shirt",
        "created_at": "2024-01-01T00:00:00Z"
      }
    },
    "new_balance": 8000.00,
    "transaction_id": 42
  },
  "message": "item purchased successfully"
}
```

**Errors:**
- `400`: Insufficient balance
- `404`: User or item not found

### 3. Sell Item
```
POST /api/shop/sell/:userItemId?user_id=1
```

**Path Parameters:**
- `userItemId`: ID of the user's item to sell

**Query Parameters:**
- `user_id`: ID of the user (in production, from JWT token)

**Response:**
```json
{
  "success": true,
  "data": {
    "sale_price": 1000.00,
    "new_balance": 9000.00,
    "transaction_id": 43
  },
  "message": "item sold successfully"
}
```

**Note:** Items are sold for 50% of their original purchase price.

### 4. Get My Items
```
GET /api/shop/my-items?user_id=1
```

**Query Parameters:**
- `user_id`: ID of the user (in production, from JWT token)

**Response:**
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": 1,
        "user_id": 1,
        "item_id": 5,
        "purchased_at": "2024-01-01T12:00:00Z",
        "is_equipped": true,
        "item": {
          "id": 5,
          "name": "Designer Shirt",
          "type": "clothing",
          "price": 2000.00,
          "image_url": "https://...",
          "description": "Stylish designer shirt",
          "created_at": "2024-01-01T00:00:00Z"
        }
      }
    ],
    "count": 1
  }
}
```

### 5. Equip Item
```
POST /api/shop/equip/:userItemId?user_id=1
```

**Path Parameters:**
- `userItemId`: ID of the user's item to equip

**Query Parameters:**
- `user_id`: ID of the user (in production, from JWT token)

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 1,
    "item_id": 5,
    "purchased_at": "2024-01-01T12:00:00Z",
    "is_equipped": true,
    "item": {
      "id": 5,
      "name": "Designer Shirt",
      "type": "clothing",
      "price": 2000.00,
      "image_url": "https://...",
      "description": "Stylish designer shirt",
      "created_at": "2024-01-01T00:00:00Z"
    }
  },
  "message": "item equipped successfully"
}
```

**Note:** Only one item per category (type) can be equipped at a time. Equipping a new item automatically unequips the previous item of the same type.

## Business Logic

### Purchase Flow
1. Verify user and item exist
2. Check user has sufficient balance
3. Deduct price from user balance
4. Create UserItem record (is_equipped = false)
5. Create Transaction record (type = "purchase", amount = -price)
6. Return user item, new balance, and transaction ID

### Sell Flow
1. Verify user owns the item
2. Calculate sale price (50% of original price)
3. Add sale price to user balance
4. Delete UserItem record
5. Create Transaction record (type = "sale", amount = +sale_price)
6. Return sale price, new balance, and transaction ID

### Equip Flow
1. Verify user owns the item
2. Unequip all items of the same type for this user
3. Set is_equipped = true for the selected item
4. Return updated user item

## Database Transactions

All mutating operations (Buy, Sell, Equip) use database transactions with row-level locking to prevent race conditions and ensure data consistency.
