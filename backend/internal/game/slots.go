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
	// Новые символы для оптимизированной конфигурации (RTP 95%, WinRate 25%)
	SymbolClover  SlotSymbol = "🍀"
	SymbolBell    SlotSymbol = "🔔"
	SymbolBar     SlotSymbol = "━"
)

// SlotReel represents a single reel with 3 visible symbols
type SlotReel [3]SlotSymbol

// WinTier represents the tier/level of the win (for animations)
type WinTier string

const (
	WinTierNone   WinTier = "none"    // No win
	WinTierSmall  WinTier = "small"   // 1-10x (мелкие выигрыши)
	WinTierMedium WinTier = "medium"  // 10-50x (средние выигрыши)
	WinTierBig    WinTier = "big"     // 50-100x (большие выигрыши)
	WinTierJackpot WinTier = "jackpot" // 100x+ (джекпот)
)

// SlotResult represents the result of a slot spin
type SlotResult struct {
	Reels       [5]SlotReel       `json:"reels"`        // 5 reels, each with 3 symbols
	WinningLine []WinningLine     `json:"winning_line"` // Details of winning lines
	TotalWin    float64           `json:"total_win"`    // Total winnings
	Multiplier  float64           `json:"multiplier"`   // Total multiplier
	WinTier     WinTier           `json:"win_tier"`     // Tier of win (for animations)
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
	// All available symbols (включая оптимизированные для RTP 95% и WinRate 25%)
	allSymbols = []SlotSymbol{
		SymbolCherry, SymbolLemon, SymbolOrange, SymbolGrape,
		SymbolDiamond, SymbolStar, SymbolSeven,
		SymbolClover, SymbolBell, SymbolBar,
	}

	// Оптимизированные веса символов (на основе генетической оптимизации)
	// RTP: 95.21%, Win Rate: 24.96%
	// Распределение: Мелкие 22.6%, Средние 2.32%, Большие 0%, Джекпот 0.04%
	symbolWeights = map[SlotSymbol]int{
		SymbolClover:  9, // 27.3% - самый частый (мелкие выигрыши)
		SymbolBell:    7, // 21.2%
		SymbolGrape:   4, // 12.1%
		SymbolDiamond: 3, // 9.1%
		SymbolBar:     3, // 9.1%
		SymbolLemon:   2, // 6.1%
		SymbolCherry:  2, // 6.1%
		SymbolOrange:  1, // 3.0%
		SymbolSeven:   1, // 3.0% - редкий джекпот
		SymbolStar:    1, // 3.0%
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
	// Оптимизировано для RTP 95% и WinRate 25%
	payoutTable = map[SlotSymbol]map[int]float64{
		SymbolSeven: {
			5: 500.0, // 5 sevens = 500x bet (Jackpot)
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
		// Новые символы для частых мелких выигрышей
		SymbolBar: {
			5: 20.0,
			4: 5.0,
			3: 1.5,
		},
		SymbolBell: {
			5: 15.0,
			4: 4.0,
			3: 1.2,
		},
		SymbolClover: {
			5: 12.0,
			4: 3.0,
			3: 1.0,
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
		WinTier:     WinTierNone,
	}

	// Check all paylines for wins
	for lineNum, payline := range paylines {
		if winLine := se.checkPayline(result.Reels, payline, lineNum+1, bet); winLine != nil {
			result.WinningLine = append(result.WinningLine, *winLine)
			result.TotalWin += winLine.Win
			result.Multiplier += winLine.Multiplier
		}
	}

	// Определяем тир выигрыша для анимаций
	if result.TotalWin > 0 && bet > 0 {
		winMultiplier := result.TotalWin / bet

		if winMultiplier >= 100 {
			result.WinTier = WinTierJackpot // 100x+ = Джекпот 🎉🎉🎉
		} else if winMultiplier >= 50 {
			result.WinTier = WinTierBig // 50-100x = Большой выигрыш 🎉
		} else if winMultiplier >= 10 {
			result.WinTier = WinTierMedium // 10-50x = Средний выигрыш
		} else {
			result.WinTier = WinTierSmall // 1-10x = Мелкий выигрыш
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

// generateReel generates a single reel with 3 weighted random symbols
// Использует оптимизированные веса для достижения RTP 95% и WinRate 25%
func (se *SlotsEngine) generateReel() SlotReel {
	var reel SlotReel

	// Вычисляем общий вес
	totalWeight := 0
	for _, weight := range symbolWeights {
		totalWeight += weight
	}

	for i := 0; i < 3; i++ {
		// Генерируем случайное число от 0 до totalWeight
		roll := se.rng.Intn(totalWeight)

		// Выбираем символ на основе веса
		currentWeight := 0
		for _, symbol := range allSymbols {
			currentWeight += symbolWeights[symbol]
			if roll < currentWeight {
				reel[i] = symbol
				break
			}
		}
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

// PaytableEntry represents a single entry in the paytable for API
type PaytableEntry struct {
	Symbol      SlotSymbol `json:"symbol"`
	ThreeOfKind float64    `json:"three_of_kind"`
	FourOfKind  float64    `json:"four_of_kind"`
	FiveOfKind  float64    `json:"five_of_kind"`
}

// GetPaytableForAPI returns the paytable in a format suitable for API response
// Symbols are ordered by their maximum payout (descending)
func GetPaytableForAPI() []PaytableEntry {
	entries := make([]PaytableEntry, 0, len(payoutTable))

	for symbol, payouts := range payoutTable {
		entry := PaytableEntry{
			Symbol:      symbol,
			ThreeOfKind: payouts[3],
			FourOfKind:  payouts[4],
			FiveOfKind:  payouts[5],
		}
		entries = append(entries, entry)
	}

	// Sort by five_of_kind payout (descending)
	// Simple bubble sort since we have only 10 symbols
	for i := 0; i < len(entries)-1; i++ {
		for j := 0; j < len(entries)-i-1; j++ {
			if entries[j].FiveOfKind < entries[j+1].FiveOfKind {
				entries[j], entries[j+1] = entries[j+1], entries[j]
			}
		}
	}

	return entries
}
