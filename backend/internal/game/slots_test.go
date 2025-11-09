package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSlotsEngine(t *testing.T) {
	engine := NewSlotsEngine()
	assert.NotNil(t, engine)
	assert.NotNil(t, engine.rng)
}

func TestSlotsEngineGenerateReel(t *testing.T) {
	engine := NewSlotsEngine()

	reel := engine.generateReel()
	assert.Len(t, reel, 3, "reel should have 3 symbols")

	// Verify all symbols are valid
	for i, symbol := range reel {
		found := false
		for _, validSymbol := range allSymbols {
			if symbol == validSymbol {
				found = true
				break
			}
		}
		assert.True(t, found, "symbol at position %d should be valid: %s", i, symbol)
	}
}

func TestSlotsEngineGenerateReels(t *testing.T) {
	engine := NewSlotsEngine()

	reels := engine.generateReels()
	assert.Len(t, reels, 5, "should have 5 reels")

	for i, reel := range reels {
		assert.Len(t, reel, 3, "reel %d should have 3 symbols", i)
	}
}

func TestSlotsSpin(t *testing.T) {
	engine := NewSlotsEngine()

	// Test multiple spins
	for i := 0; i < 10; i++ {
		result := engine.Spin(10.0)

		require.NotNil(t, result)
		assert.Len(t, result.Reels, 5, "should have 5 reels")
		assert.NotNil(t, result.WinningLine, "winning lines should not be nil")
		assert.GreaterOrEqual(t, result.TotalWin, 0.0, "total win should be non-negative")
		assert.GreaterOrEqual(t, result.Multiplier, 0.0, "multiplier should be non-negative")

		// Verify all reels have valid symbols
		for reelIdx, reel := range result.Reels {
			for symbolIdx, symbol := range reel {
				found := false
				for _, validSymbol := range allSymbols {
					if symbol == validSymbol {
						found = true
						break
					}
				}
				assert.True(t, found, "invalid symbol at reel %d, position %d: %s", reelIdx, symbolIdx, symbol)
			}
		}
	}
}

func TestSlotsCheckPaylineNoWin(t *testing.T) {
	engine := NewSlotsEngine()

	// Create reels with no matching symbols
	reels := [5]SlotReel{
		{SymbolCherry, SymbolLemon, SymbolOrange},
		{SymbolGrape, SymbolDiamond, SymbolStar},
		{SymbolSeven, SymbolCherry, SymbolLemon},
		{SymbolOrange, SymbolGrape, SymbolDiamond},
		{SymbolStar, SymbolSeven, SymbolCherry},
	}

	payline := Payline{1, 1, 1, 1, 1} // Middle line
	bet := 10.0

	winLine := engine.checkPayline(reels, payline, 1, bet)
	assert.Nil(t, winLine, "should have no win with non-matching symbols")
}

func TestSlotsCheckPaylineThreeInRow(t *testing.T) {
	engine := NewSlotsEngine()

	// Create reels with 3 matching symbols
	reels := [5]SlotReel{
		{SymbolCherry, SymbolCherry, SymbolOrange},
		{SymbolGrape, SymbolCherry, SymbolStar},
		{SymbolSeven, SymbolCherry, SymbolLemon},
		{SymbolOrange, SymbolGrape, SymbolDiamond},
		{SymbolStar, SymbolSeven, SymbolCherry},
	}

	payline := Payline{1, 1, 1, 1, 1} // Middle line
	bet := 10.0

	winLine := engine.checkPayline(reels, payline, 1, bet)
	require.NotNil(t, winLine, "should have win with 3 matching symbols")
	assert.Equal(t, SymbolCherry, winLine.Symbol)
	assert.Equal(t, 3, winLine.Count)
	assert.Equal(t, 2.0, winLine.Multiplier)
	assert.Equal(t, 20.0, winLine.Win)
}

func TestSlotsCheckPaylineFiveInRow(t *testing.T) {
	engine := NewSlotsEngine()

	// Create reels with 5 matching symbols (jackpot!)
	reels := [5]SlotReel{
		{SymbolCherry, SymbolSeven, SymbolOrange},
		{SymbolGrape, SymbolSeven, SymbolStar},
		{SymbolSeven, SymbolSeven, SymbolLemon},
		{SymbolOrange, SymbolSeven, SymbolDiamond},
		{SymbolStar, SymbolSeven, SymbolCherry},
	}

	payline := Payline{1, 1, 1, 1, 1} // Middle line
	bet := 10.0

	winLine := engine.checkPayline(reels, payline, 1, bet)
	require.NotNil(t, winLine, "should have win with 5 sevens")
	assert.Equal(t, SymbolSeven, winLine.Symbol)
	assert.Equal(t, 5, winLine.Count)
	assert.Equal(t, 500.0, winLine.Multiplier, "five sevens should pay 500x")
	assert.Equal(t, 5000.0, winLine.Win)
}

func TestSlotsGetPayoutTable(t *testing.T) {
	table := GetPayoutTable()

	require.NotNil(t, table)
	assert.Len(t, table, 10, "should have payouts for 10 symbols (7 original + 3 new)")

	// Verify Seven has highest payout
	sevenPayouts := table[SymbolSeven]
	require.NotNil(t, sevenPayouts)
	assert.Equal(t, 500.0, sevenPayouts[5])
	assert.Equal(t, 100.0, sevenPayouts[4])
	assert.Equal(t, 20.0, sevenPayouts[3])

	// Verify Clover has lowest payout (new symbol for small wins)
	cloverPayouts := table[SymbolClover]
	require.NotNil(t, cloverPayouts)
	assert.Equal(t, 12.0, cloverPayouts[5])
	assert.Equal(t, 3.0, cloverPayouts[4])
	assert.Equal(t, 1.0, cloverPayouts[3])
}

func TestSlotsGetAllSymbols(t *testing.T) {
	symbols := GetAllSymbols()

	assert.Len(t, symbols, 10, "should have 10 symbols")
	assert.Contains(t, symbols, SymbolCherry)
	assert.Contains(t, symbols, SymbolLemon)
	assert.Contains(t, symbols, SymbolOrange)
	assert.Contains(t, symbols, SymbolGrape)
	assert.Contains(t, symbols, SymbolDiamond)
	assert.Contains(t, symbols, SymbolStar)
	assert.Contains(t, symbols, SymbolSeven)
	assert.Contains(t, symbols, SymbolClover)
	assert.Contains(t, symbols, SymbolBell)
	assert.Contains(t, symbols, SymbolBar)
}

func TestSlotsPaylines(t *testing.T) {
	assert.Len(t, paylines, 10, "should have 10 paylines")

	// Verify each payline has 5 positions (for 5 reels)
	for i, payline := range paylines {
		assert.Len(t, payline, 5, "payline %d should have 5 positions", i+1)

		// Each position should be 0, 1, or 2 (row index)
		for j, rowIndex := range payline {
			assert.GreaterOrEqual(t, rowIndex, 0, "payline %d position %d should be >= 0", i+1, j)
			assert.LessOrEqual(t, rowIndex, 2, "payline %d position %d should be <= 2", i+1, j)
		}
	}
}

func TestSlotsMultipleWinningLines(t *testing.T) {
	engine := NewSlotsEngine()

	// Create reels with multiple winning lines
	reels := [5]SlotReel{
		{SymbolCherry, SymbolCherry, SymbolCherry},
		{SymbolCherry, SymbolCherry, SymbolCherry},
		{SymbolCherry, SymbolCherry, SymbolCherry},
		{SymbolGrape, SymbolGrape, SymbolGrape},
		{SymbolLemon, SymbolLemon, SymbolLemon},
	}

	bet := 10.0

	// Check middle horizontal line (should win)
	winLine := engine.checkPayline(reels, paylines[0], 1, bet)
	require.NotNil(t, winLine)
	assert.Equal(t, SymbolCherry, winLine.Symbol)
	assert.Equal(t, 3, winLine.Count)

	// Check top horizontal line (should also win)
	winLine = engine.checkPayline(reels, paylines[1], 2, bet)
	require.NotNil(t, winLine)
	assert.Equal(t, SymbolCherry, winLine.Symbol)
	assert.Equal(t, 3, winLine.Count)

	// Check bottom horizontal line (should also win)
	winLine = engine.checkPayline(reels, paylines[2], 3, bet)
	require.NotNil(t, winLine)
	assert.Equal(t, SymbolCherry, winLine.Symbol)
	assert.Equal(t, 3, winLine.Count)
}

func TestSlotsFairness(t *testing.T) {
	engine := NewSlotsEngine()

	totalWins := 0
	totalSpins := 100
	bet := 10.0

	for i := 0; i < totalSpins; i++ {
		result := engine.Spin(bet)
		if result.TotalWin > 0 {
			totalWins++
		}
	}

	// Should have some wins but not all (house edge)
	assert.Greater(t, totalWins, 0, "should have some wins")
	assert.Less(t, totalWins, totalSpins, "should not win every spin")

	// Win rate should be reasonable (not 0%, not 100%)
	winRate := float64(totalWins) / float64(totalSpins)
	assert.Greater(t, winRate, 0.05, "win rate should be > 5%")
	assert.Less(t, winRate, 0.95, "win rate should be < 95%")
}

func TestSlotsPayoutMultipliers(t *testing.T) {
	engine := NewSlotsEngine()

	// Test with different bet amounts
	bets := []float64{1.0, 10.0, 100.0}

	for _, bet := range bets {
		t.Run("", func(t *testing.T) {
			result := engine.Spin(bet)

			// If there's a win, verify the calculation is correct
			if len(result.WinningLine) > 0 {
				calculatedWin := 0.0
				for _, line := range result.WinningLine {
					calculatedWin += line.Win
					// Win should equal bet * multiplier
					expectedWin := bet * line.Multiplier
					assert.InDelta(t, expectedWin, line.Win, 0.01, "win calculation should be correct")
				}
				assert.InDelta(t, calculatedWin, result.TotalWin, 0.01, "total win should match sum of lines")
			}
		})
	}
}
