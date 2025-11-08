package game

import (
	"math/rand"
	"time"
)

// SlotSymbol represents a symbol on the slot machine
type SlotSymbol string

const (
	SymbolCherry  SlotSymbol = "🍒"
	SymbolLemon   SlotSymbol = "🍋"
	SymbolOrange  SlotSymbol = "🍊"
	SymbolGrape   SlotSymbol = "🍇"
	SymbolDiamond SlotSymbol = "💎"
	SymbolStar    SlotSymbol = "⭐"
	SymbolSeven   SlotSymbol = "7️⃣"
)

// SlotReel represents a single reel with 3 visible symbols
type SlotReel [3]SlotSymbol

// SlotResult represents the result of a slot spin
type SlotResult struct {
	Reels       [5]SlotReel       `json:"reels"`        // 5 reels, each with 3 symbols
	WinningLine []WinningLine     `json:"winning_line"` // Details of winning lines
	TotalWin    float64           `json:"total_win"`    // Total winnings
	Multiplier  float64           `json:"multiplier"`   // Total multiplier
}

// WinningLine represents a winning payline
type WinningLine struct {
	LineNumber int        `json:"line_number"` // Which payline (1-10)
	Symbol     SlotSymbol `json:"symbol"`      // Winning symbol
	Count      int        `json:"count"`       // How many in a row (3, 4, or 5)
	Multiplier float64    `json:"multiplier"`  // Multiplier for this line
	Win        float64    `json:"win"`         // Win amount for this line
}

// Payline represents a payline pattern
// Each number is the row index (0=top, 1=middle, 2=bottom) for each of the 5 reels
type Payline [5]int

var (
	// All available symbols
	allSymbols = []SlotSymbol{
		SymbolCherry, SymbolLemon, SymbolOrange, SymbolGrape,
		SymbolDiamond, SymbolStar, SymbolSeven,
	}

	// Paylines - standard 10 paylines for 5-reel slots
	paylines = []Payline{
		{1, 1, 1, 1, 1}, // Line 1: Middle horizontal
		{0, 0, 0, 0, 0}, // Line 2: Top horizontal
		{2, 2, 2, 2, 2}, // Line 3: Bottom horizontal
		{0, 1, 2, 1, 0}, // Line 4: V shape
		{2, 1, 0, 1, 2}, // Line 5: Inverted V
		{1, 0, 1, 0, 1}, // Line 6: Zigzag
		{1, 2, 1, 2, 1}, // Line 7: Reverse zigzag
		{0, 1, 0, 1, 0}, // Line 8: W shape
		{2, 1, 2, 1, 2}, // Line 9: M shape
		{0, 0, 1, 2, 2}, // Line 10: Diagonal
	}

	// Payout table: symbol -> count -> multiplier
	// count can be 3, 4, or 5 (number of symbols in a row)
	payoutTable = map[SlotSymbol]map[int]float64{
		SymbolSeven: {
			5: 500.0, // 5 sevens = 500x bet
			4: 100.0, // 4 sevens = 100x bet
			3: 20.0,  // 3 sevens = 20x bet
		},
		SymbolStar: {
			5: 200.0,
			4: 50.0,
			3: 10.0,
		},
		SymbolDiamond: {
			5: 150.0,
			4: 40.0,
			3: 8.0,
		},
		SymbolGrape: {
			5: 100.0,
			4: 25.0,
			3: 5.0,
		},
		SymbolOrange: {
			5: 80.0,
			4: 20.0,
			3: 4.0,
		},
		SymbolLemon: {
			5: 60.0,
			4: 15.0,
			3: 3.0,
		},
		SymbolCherry: {
			5: 40.0,
			4: 10.0,
			3: 2.0,
		},
	}
)

// SlotsEngine handles slot machine game logic
type SlotsEngine struct {
	rng *rand.Rand
}

// NewSlotsEngine creates a new slots engine
func NewSlotsEngine() *SlotsEngine {
	return &SlotsEngine{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Spin performs a slot machine spin
func (se *SlotsEngine) Spin(bet float64) *SlotResult {
	result := &SlotResult{
		Reels:       se.generateReels(),
		WinningLine: []WinningLine{},
		TotalWin:    0,
		Multiplier:  0,
	}

	// Check all paylines for wins
	for lineNum, payline := range paylines {
		if winLine := se.checkPayline(result.Reels, payline, lineNum+1, bet); winLine != nil {
			result.WinningLine = append(result.WinningLine, *winLine)
			result.TotalWin += winLine.Win
			result.Multiplier += winLine.Multiplier
		}
	}

	return result
}

// generateReels generates 5 random reels
func (se *SlotsEngine) generateReels() [5]SlotReel {
	var reels [5]SlotReel

	for i := 0; i < 5; i++ {
		reels[i] = se.generateReel()
	}

	return reels
}

// generateReel generates a single reel with 3 random symbols
func (se *SlotsEngine) generateReel() SlotReel {
	var reel SlotReel

	for i := 0; i < 3; i++ {
		reel[i] = allSymbols[se.rng.Intn(len(allSymbols))]
	}

	return reel
}

// checkPayline checks if a payline is a winner
func (se *SlotsEngine) checkPayline(reels [5]SlotReel, payline Payline, lineNumber int, bet float64) *WinningLine {
	// Get the symbols along this payline
	var symbols [5]SlotSymbol
	for i := 0; i < 5; i++ {
		symbols[i] = reels[i][payline[i]]
	}

	// Count consecutive matching symbols from left to right
	firstSymbol := symbols[0]
	count := 1

	for i := 1; i < 5; i++ {
		if symbols[i] == firstSymbol {
			count++
		} else {
			break
		}
	}

	// Need at least 3 in a row to win
	if count < 3 {
		return nil
	}

	// Get multiplier from payout table
	multiplier := payoutTable[firstSymbol][count]
	win := bet * multiplier

	return &WinningLine{
		LineNumber: lineNumber,
		Symbol:     firstSymbol,
		Count:      count,
		Multiplier: multiplier,
		Win:        win,
	}
}

// GetPayoutTable returns the payout table for display
func GetPayoutTable() map[SlotSymbol]map[int]float64 {
	return payoutTable
}

// GetAllSymbols returns all available symbols
func GetAllSymbols() []SlotSymbol {
	return allSymbols
}
