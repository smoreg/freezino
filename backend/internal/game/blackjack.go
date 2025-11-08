package game

import (
	"fmt"
	"math/rand"
	"time"
)

// Card represents a playing card
type Card struct {
	Suit  string `json:"suit"`  // hearts, diamonds, clubs, spades
	Rank  string `json:"rank"`  // A, 2-10, J, Q, K
	Value int    `json:"value"` // Numeric value for calculation
}

// Hand represents a player's or dealer's hand
type Hand struct {
	Cards []Card `json:"cards"`
	Value int    `json:"value"`
	Soft  bool   `json:"soft"` // True if ace is counted as 11
}

// BlackjackGame represents a blackjack game session
type BlackjackGame struct {
	PlayerHand Hand    `json:"player_hand"`
	DealerHand Hand    `json:"dealer_hand"`
	Deck       []Card  `json:"-"` // Don't send deck to client
	Bet        float64 `json:"bet"`
	GameOver   bool    `json:"game_over"`
	Result     string  `json:"result"`     // "player_win", "dealer_win", "push", "blackjack"
	Multiplier float64 `json:"multiplier"` // Win multiplier (1.0, 1.5, 2.0)
	CanDouble  bool    `json:"can_double"`
	CanSplit   bool    `json:"can_split"`
}

// NewBlackjackGame creates a new blackjack game
func NewBlackjackGame(bet float64) *BlackjackGame {
	game := &BlackjackGame{
		Bet:        bet,
		GameOver:   false,
		Result:     "",
		Multiplier: 0,
		CanDouble:  true,
		CanSplit:   false,
	}

	// Create and shuffle deck
	game.Deck = createDeck()
	game.shuffleDeck()

	// Deal initial cards
	game.PlayerHand.Cards = append(game.PlayerHand.Cards, game.drawCard())
	game.DealerHand.Cards = append(game.DealerHand.Cards, game.drawCard())
	game.PlayerHand.Cards = append(game.PlayerHand.Cards, game.drawCard())
	game.DealerHand.Cards = append(game.DealerHand.Cards, game.drawCard())

	// Calculate initial values
	game.calculateHandValue(&game.PlayerHand)
	game.calculateHandValue(&game.DealerHand)

	// Check for split possibility
	if len(game.PlayerHand.Cards) == 2 && game.PlayerHand.Cards[0].Rank == game.PlayerHand.Cards[1].Rank {
		game.CanSplit = true
	}

	// Check for immediate blackjack
	if game.PlayerHand.Value == 21 {
		game.GameOver = true
		if game.DealerHand.Value == 21 {
			game.Result = "push"
			game.Multiplier = 1.0
		} else {
			game.Result = "blackjack"
			game.Multiplier = 2.5 // Blackjack pays 3:2
		}
	}

	return game
}

// createDeck creates a standard 52-card deck
func createDeck() []Card {
	suits := []string{"hearts", "diamonds", "clubs", "spades"}
	ranks := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	deck := make([]Card, 0, 52)

	for _, suit := range suits {
		for _, rank := range ranks {
			value := getCardValue(rank)
			deck = append(deck, Card{
				Suit:  suit,
				Rank:  rank,
				Value: value,
			})
		}
	}

	return deck
}

// getCardValue returns the numeric value of a card
func getCardValue(rank string) int {
	switch rank {
	case "A":
		return 11 // Ace is 11 by default, adjusted later if needed
	case "J", "Q", "K":
		return 10
	default:
		// Convert string to int (2-10)
		var value int
		if _, err := fmt.Sscanf(rank, "%d", &value); err != nil {
			return 0 // Invalid rank
		}
		return value
	}
}

// shuffleDeck shuffles the deck using Fisher-Yates algorithm
func (g *BlackjackGame) shuffleDeck() {
	rand.Seed(time.Now().UnixNano())
	for i := len(g.Deck) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		g.Deck[i], g.Deck[j] = g.Deck[j], g.Deck[i]
	}
}

// drawCard draws a card from the deck
func (g *BlackjackGame) drawCard() Card {
	if len(g.Deck) == 0 {
		// Reshuffle if deck is empty
		g.Deck = createDeck()
		g.shuffleDeck()
	}
	card := g.Deck[0]
	g.Deck = g.Deck[1:]
	return card
}

// calculateHandValue calculates the value of a hand
func (g *BlackjackGame) calculateHandValue(hand *Hand) {
	value := 0
	aces := 0

	// Count all cards and number of aces
	for _, card := range hand.Cards {
		value += card.Value
		if card.Rank == "A" {
			aces++
		}
	}

	// Adjust for aces (convert from 11 to 1 if needed)
	hand.Soft = false
	for aces > 0 && value > 21 {
		value -= 10
		aces--
	}

	// If we have an ace counted as 11, it's a soft hand
	if aces > 0 {
		hand.Soft = true
	}

	hand.Value = value
}

// Hit adds a card to player's hand
func (g *BlackjackGame) Hit() error {
	if g.GameOver {
		return fmt.Errorf("game is already over")
	}

	// Player can't double after first hit
	g.CanDouble = false
	g.CanSplit = false

	// Draw card and add to player's hand
	card := g.drawCard()
	g.PlayerHand.Cards = append(g.PlayerHand.Cards, card)
	g.calculateHandValue(&g.PlayerHand)

	// Check if player busted
	if g.PlayerHand.Value > 21 {
		g.GameOver = true
		g.Result = "dealer_win"
		g.Multiplier = 0
	}

	return nil
}

// Stand ends player's turn and dealer plays
func (g *BlackjackGame) Stand() error {
	if g.GameOver {
		return fmt.Errorf("game is already over")
	}

	g.CanDouble = false
	g.CanSplit = false

	// Dealer plays (must hit on 16 or less, stand on 17 or more)
	for g.DealerHand.Value < 17 {
		card := g.drawCard()
		g.DealerHand.Cards = append(g.DealerHand.Cards, card)
		g.calculateHandValue(&g.DealerHand)
	}

	// Determine winner
	g.GameOver = true
	g.determineWinner()

	return nil
}

// Double doubles the bet and draws one card, then stands
func (g *BlackjackGame) Double() error {
	if g.GameOver {
		return fmt.Errorf("game is already over")
	}

	if !g.CanDouble {
		return fmt.Errorf("cannot double at this point")
	}

	// Double the bet
	g.Bet *= 2

	// Draw exactly one card
	card := g.drawCard()
	g.PlayerHand.Cards = append(g.PlayerHand.Cards, card)
	g.calculateHandValue(&g.PlayerHand)

	g.CanDouble = false
	g.CanSplit = false

	// Check if player busted
	if g.PlayerHand.Value > 21 {
		g.GameOver = true
		g.Result = "dealer_win"
		g.Multiplier = 0
		return nil
	}

	// Automatically stand after double
	return g.Stand()
}

// Split splits the hand into two hands (simplified version - only returns error for now)
func (g *BlackjackGame) Split() error {
	if g.GameOver {
		return fmt.Errorf("game is already over")
	}

	if !g.CanSplit {
		return fmt.Errorf("cannot split at this point")
	}

	// Note: Full split implementation would require managing two hands
	// For now, we'll return an error indicating it's not fully implemented
	return fmt.Errorf("split is not yet fully implemented")
}

// determineWinner determines the winner of the game
func (g *BlackjackGame) determineWinner() {
	if g.DealerHand.Value > 21 {
		// Dealer busted
		g.Result = "player_win"
		g.Multiplier = 2.0
	} else if g.PlayerHand.Value > g.DealerHand.Value {
		// Player has higher value
		g.Result = "player_win"
		g.Multiplier = 2.0
	} else if g.PlayerHand.Value < g.DealerHand.Value {
		// Dealer has higher value
		g.Result = "dealer_win"
		g.Multiplier = 0
	} else {
		// Push (tie)
		g.Result = "push"
		g.Multiplier = 1.0
	}
}

// GetPayout calculates the payout based on the result
func (g *BlackjackGame) GetPayout() float64 {
	return g.Bet * g.Multiplier
}

// GetDealerVisibleCard returns the dealer's visible card (first card)
func (g *BlackjackGame) GetDealerVisibleCard() *Card {
	if len(g.DealerHand.Cards) > 0 {
		return &g.DealerHand.Cards[0]
	}
	return nil
}

// BlackjackGameState represents the game state sent to client
type BlackjackGameState struct {
	PlayerHand       Hand    `json:"player_hand"`
	DealerVisibleCard *Card  `json:"dealer_visible_card,omitempty"`
	DealerHand       *Hand   `json:"dealer_hand,omitempty"` // Only sent when game is over
	Bet              float64 `json:"bet"`
	GameOver         bool    `json:"game_over"`
	Result           string  `json:"result"`
	Payout           float64 `json:"payout"`
	CanDouble        bool    `json:"can_double"`
	CanSplit         bool    `json:"can_split"`
}

// GetGameState returns the current game state for the client
func (g *BlackjackGame) GetGameState() BlackjackGameState {
	state := BlackjackGameState{
		PlayerHand: g.PlayerHand,
		Bet:        g.Bet,
		GameOver:   g.GameOver,
		Result:     g.Result,
		Payout:     g.GetPayout(),
		CanDouble:  g.CanDouble,
		CanSplit:   g.CanSplit,
	}

	if g.GameOver {
		// Show full dealer hand when game is over
		state.DealerHand = &g.DealerHand
	} else {
		// Only show dealer's first card during game
		state.DealerVisibleCard = g.GetDealerVisibleCard()
	}

	return state
}
