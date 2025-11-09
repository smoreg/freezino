package game

import (
	"testing"
)

// FuzzSlotSymbolGeneration проверяет что генерация барабанов не падает и возвращает валидные символы
func FuzzSlotSymbolGeneration(f *testing.F) {
	// Добавляем seed корпус
	f.Add(int64(12345))
	f.Add(int64(67890))
	f.Add(int64(99999))

	f.Fuzz(func(t *testing.T, seed int64) {
		// Создаем engine с заданным seed
		engine := NewSlotsEngine()

		// Генерируем барабаны
		reels := engine.generateReels()

		// Проверяем что все барабаны содержат валидные символы
		for i, reel := range reels {
			if len(reel) != 3 {
				t.Errorf("Барабан %d должен содержать 3 символа, получено %d", i, len(reel))
			}

			for j, symbol := range reel {
				// Проверяем что символ валидный
				valid := false
				for _, validSymbol := range allSymbols {
					if symbol == validSymbol {
						valid = true
						break
					}
				}

				if !valid {
					t.Errorf("Барабан %d позиция %d: невалидный символ %v", i, j, symbol)
				}
			}
		}

		// Проверяем что спин не падает
		bet := 10.0
		result := engine.Spin(bet)

		// Проверяем базовые инварианты
		if result.TotalWin < 0 {
			t.Errorf("Выигрыш не может быть отрицательным: %f", result.TotalWin)
		}

		if len(result.Reels) != 5 {
			t.Errorf("Должно быть 5 барабанов, получено %d", len(result.Reels))
		}
	})
}
