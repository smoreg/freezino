package main

import (
	"math/rand"
	"testing"
	"time"
)

// TestMutationPreservesAllSymbols проверяет что мутации не удаляют символы
func TestMutationPreservesAllSymbols(t *testing.T) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	allSymbols := []Symbol{
		SymClover, SymBell, SymBar, SymCherry,
		SymLemon, SymOrange, SymGrape,
		SymDiamond, SymStar, SymSeven,
	}

	// Создаем начальный барабан
	baseReel := CreateOptimalStartReelSet(95.0, 10.0, true)[0]

	// Проверяем начальный барабан
	if !baseReel.validateReel(allSymbols) {
		t.Error("Начальный барабан не содержит всех символов")
	}

	// Тестируем 1000 мутаций
	for i := 0; i < 1000; i++ {
		// Обычная мутация
		mutated := baseReel.Mutate(rng, false)
		if !mutated.validateReel(allSymbols) {
			t.Errorf("Обычная мутация %d удалила символы", i)
			logMissingSymbols(t, mutated, allSymbols)
		}

		// Сильная мутация
		mutatedStrong := baseReel.Mutate(rng, true)
		if !mutatedStrong.validateReel(allSymbols) {
			t.Errorf("Сильная мутация %d удалила символы", i)
			logMissingSymbols(t, mutatedStrong, allSymbols)
		}

		// Используем мутированный барабан для следующих тестов
		if i%10 == 0 {
			baseReel = mutated
		}
	}
}

// TestReelSetMutationPreservesSymbols проверяет мутацию набора барабанов
func TestReelSetMutationPreservesSymbols(t *testing.T) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	allSymbols := []Symbol{
		SymClover, SymBell, SymBar, SymCherry,
		SymLemon, SymOrange, SymGrape,
		SymDiamond, SymStar, SymSeven,
	}

	testCases := []struct {
		name      string
		sameReels bool
	}{
		{"Одинаковые барабаны", true},
		{"Разные барабаны", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			baseSet := CreateOptimalStartReelSet(95.0, 10.0, tc.sameReels)

			// Проверяем начальный набор
			for i, reel := range baseSet {
				if !reel.validateReel(allSymbols) {
					t.Errorf("Начальный барабан %d не содержит всех символов", i)
				}
			}

			// Тестируем 100 мутаций набора
			for i := 0; i < 100; i++ {
				mutated := baseSet.Mutate(rng, false, tc.sameReels)
				for j, reel := range mutated {
					if !reel.validateReel(allSymbols) {
						t.Errorf("Мутация %d барабана %d удалила символы", i, j)
						logMissingSymbols(t, reel, allSymbols)
					}
				}

				// Сильная мутация
				mutatedStrong := baseSet.Mutate(rng, true, tc.sameReels)
				for j, reel := range mutatedStrong {
					if !reel.validateReel(allSymbols) {
						t.Errorf("Сильная мутация %d барабана %d удалила символы", i, j)
						logMissingSymbols(t, reel, allSymbols)
					}
				}
			}
		})
	}
}

// TestCountSymbol проверяет подсчет символов
func TestCountSymbol(t *testing.T) {
	reel := Reel{
		SymCherry, SymCherry, SymLemon,
		SymCherry, SymSeven, SymDiamond,
	}

	tests := []struct {
		symbol Symbol
		want   int
	}{
		{SymCherry, 3},
		{SymLemon, 1},
		{SymSeven, 1},
		{SymDiamond, 1},
		{SymStar, 0},
	}

	for _, tt := range tests {
		got := reel.countSymbol(tt.symbol)
		if got != tt.want {
			t.Errorf("countSymbol(%s) = %d, want %d", tt.symbol.Emoji, got, tt.want)
		}
	}
}

// logMissingSymbols выводит какие символы отсутствуют
func logMissingSymbols(t *testing.T, reel Reel, allSymbols []Symbol) {
	t.Log("Отсутствующие символы:")
	for _, symbol := range allSymbols {
		count := reel.countSymbol(symbol)
		if count == 0 {
			t.Logf("  %s (%s) - ОТСУТСТВУЕТ", symbol.Emoji, symbol.Name)
		}
	}

	t.Log("Содержимое барабана:")
	counts := make(map[string]int)
	for _, sym := range reel {
		counts[sym.Emoji]++
	}
	for emoji, count := range counts {
		t.Logf("  %s: %d", emoji, count)
	}
}
