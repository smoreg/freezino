package handler

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/smoreg/freezino/backend/internal/game"
	"github.com/smoreg/freezino/backend/internal/model"
	"gorm.io/gorm"
)

// GameHandler manages game WebSocket connections
type GameHandler struct {
	games sync.Map // map[*websocket.Conn]*game.BlackjackGame
	db    *gorm.DB
}

// NewGameHandler creates a new game handler
func NewGameHandler(db *gorm.DB) *GameHandler {
	return &GameHandler{
		db: db,
	}
}

// WebSocketUpgrade upgrades HTTP connection to WebSocket
func (h *GameHandler) WebSocketUpgrade(c *fiber.Ctx) error {
	// IsWebSocketUpgrade returns true if the client requested upgrade to the WebSocket protocol
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

// Message types
const (
	MsgTypeNewGame    = "new_game"
	MsgTypeHit        = "hit"
	MsgTypeStand      = "stand"
	MsgTypeDouble     = "double"
	MsgTypeSplit      = "split"
	MsgTypeGameState  = "game_state"
	MsgTypeError      = "error"
	MsgTypeBalanceUpdate = "balance_update"
)

// WebSocketMessage represents a WebSocket message
type WebSocketMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// NewGamePayload represents the payload for starting a new game
type NewGamePayload struct {
	Bet    float64 `json:"bet"`
	UserID uint    `json:"user_id"`
}

// ErrorPayload represents an error message
type ErrorPayload struct {
	Message string `json:"message"`
}

// BalanceUpdatePayload represents balance update
type BalanceUpdatePayload struct {
	Balance float64 `json:"balance"`
}

// BlackjackWebSocket handles blackjack WebSocket connections
func (h *GameHandler) BlackjackWebSocket(c *websocket.Conn) {
	var (
		currentGame *game.BlackjackGame
		userID      uint
	)

	defer func() {
		h.games.Delete(c)
		c.Close()
	}()

	for {
		var msg WebSocketMessage
		if err := c.ReadJSON(&msg); err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		switch msg.Type {
		case MsgTypeNewGame:
			var payload NewGamePayload
			if err := json.Unmarshal(msg.Payload, &payload); err != nil {
				h.sendError(c, "Invalid payload")
				continue
			}

			userID = payload.UserID

			// Get user balance
			var user model.User
			if err := h.db.First(&user, userID).Error; err != nil {
				h.sendError(c, "User not found")
				continue
			}

			// Check if user has enough balance
			if user.Balance < payload.Bet {
				h.sendError(c, "Insufficient balance")
				continue
			}

			// Deduct bet from user balance
			user.Balance -= payload.Bet
			if err := h.db.Save(&user).Error; err != nil {
				h.sendError(c, "Failed to update balance")
				continue
			}

			// Create new game
			currentGame = game.NewBlackjackGame(payload.Bet)
			h.games.Store(c, currentGame)

			// Send game state
			h.sendGameState(c, currentGame)

			// Send balance update
			h.sendBalanceUpdate(c, user.Balance)

			// If game is already over (blackjack), handle payout
			if currentGame.GameOver {
				h.handleGameEnd(currentGame, userID)
			}

		case MsgTypeHit:
			if currentGame == nil {
				h.sendError(c, "No active game")
				continue
			}

			if err := currentGame.Hit(); err != nil {
				h.sendError(c, err.Error())
				continue
			}

			h.sendGameState(c, currentGame)

			if currentGame.GameOver {
				h.handleGameEnd(currentGame, userID)
			}

		case MsgTypeStand:
			if currentGame == nil {
				h.sendError(c, "No active game")
				continue
			}

			if err := currentGame.Stand(); err != nil {
				h.sendError(c, err.Error())
				continue
			}

			h.sendGameState(c, currentGame)

			if currentGame.GameOver {
				h.handleGameEnd(currentGame, userID)
			}

		case MsgTypeDouble:
			if currentGame == nil {
				h.sendError(c, "No active game")
				continue
			}

			// Check if user has enough balance to double
			var user model.User
			if err := h.db.First(&user, userID).Error; err != nil {
				h.sendError(c, "User not found")
				continue
			}

			if user.Balance < currentGame.Bet {
				h.sendError(c, "Insufficient balance to double")
				continue
			}

			// Deduct additional bet
			user.Balance -= currentGame.Bet
			if err := h.db.Save(&user).Error; err != nil {
				h.sendError(c, "Failed to update balance")
				continue
			}

			if err := currentGame.Double(); err != nil {
				h.sendError(c, err.Error())
				// Refund the additional bet
				user.Balance += currentGame.Bet / 2
				h.db.Save(&user)
				continue
			}

			h.sendGameState(c, currentGame)
			h.sendBalanceUpdate(c, user.Balance)

			if currentGame.GameOver {
				h.handleGameEnd(currentGame, userID)
			}

		case MsgTypeSplit:
			if currentGame == nil {
				h.sendError(c, "No active game")
				continue
			}

			if err := currentGame.Split(); err != nil {
				h.sendError(c, err.Error())
				continue
			}

			h.sendGameState(c, currentGame)

		default:
			h.sendError(c, "Unknown message type")
		}
	}
}

// handleGameEnd handles the end of a game (update balance, save session)
func (h *GameHandler) handleGameEnd(g *game.BlackjackGame, userID uint) {
	// Update user balance
	var user model.User
	if err := h.db.First(&user, userID).Error; err != nil {
		log.Printf("Error fetching user: %v", err)
		return
	}

	payout := g.GetPayout()
	user.Balance += payout

	if err := h.db.Save(&user).Error; err != nil {
		log.Printf("Error updating user balance: %v", err)
		return
	}

	// Save game session
	gameSession := model.GameSession{
		UserID:   userID,
		GameType: model.GameTypeBlackjack,
		Bet:      g.Bet,
		Win:      payout - g.Bet, // Net win/loss
	}

	if err := h.db.Create(&gameSession).Error; err != nil {
		log.Printf("Error saving game session: %v", err)
	}

	// Create transaction record
	transaction := model.Transaction{
		UserID:      userID,
		Type:        model.TransactionTypeGame,
		Amount:      payout - g.Bet,
		BalanceAfter: user.Balance,
		Description: "Blackjack - " + g.Result,
	}

	if err := h.db.Create(&transaction).Error; err != nil {
		log.Printf("Error creating transaction: %v", err)
	}
}

// sendGameState sends the current game state to the client
func (h *GameHandler) sendGameState(c *websocket.Conn, g *game.BlackjackGame) {
	state := g.GetGameState()
	msg := WebSocketMessage{
		Type:    MsgTypeGameState,
		Payload: mustMarshal(state),
	}
	if err := c.WriteJSON(msg); err != nil {
		log.Printf("Error sending game state: %v", err)
	}
}

// sendBalanceUpdate sends balance update to the client
func (h *GameHandler) sendBalanceUpdate(c *websocket.Conn, balance float64) {
	msg := WebSocketMessage{
		Type:    MsgTypeBalanceUpdate,
		Payload: mustMarshal(BalanceUpdatePayload{Balance: balance}),
	}
	if err := c.WriteJSON(msg); err != nil {
		log.Printf("Error sending balance update: %v", err)
	}
}

// sendError sends an error message to the client
func (h *GameHandler) sendError(c *websocket.Conn, message string) {
	msg := WebSocketMessage{
		Type:    MsgTypeError,
		Payload: mustMarshal(ErrorPayload{Message: message}),
	}
	if err := c.WriteJSON(msg); err != nil {
		log.Printf("Error sending error message: %v", err)
	}
}

// mustMarshal marshals data to JSON or panics
func mustMarshal(v interface{}) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}
