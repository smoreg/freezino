package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"
)

// Symbol представляет символ на барабане
type Symbol struct {
	Name       string
	Emoji      string
	Payout3    float64 // Выплата за 3 подряд
	Payout4    float64 // Выплата за 4 подряд
	Payout5    float64 // Выплата за 5 подряд
}

// Базовые символы
var (
	SymCherry  = Symbol{"Cherry", "🍒", 2.0, 10.0, 40.0}
	SymLemon   = Symbol{"Lemon", "🍋", 3.0, 15.0, 60.0}
	SymOrange  = Symbol{"Orange", "🍊", 4.0, 20.0, 80.0}
	SymGrape   = Symbol{"Grape", "🍇", 5.0, 25.0, 100.0}
	SymDiamond = Symbol{"Diamond", "💎", 8.0, 40.0, 150.0}
	SymStar    = Symbol{"Star", "⭐", 10.0, 50.0, 200.0}
	SymSeven   = Symbol{"Seven", "7️⃣", 20.0, 100.0, 500.0}

	// Низко-выигрышные символы для частых побед
	SymBar     = Symbol{"Bar", "━", 1.5, 5.0, 20.0}
	SymBell    = Symbol{"Bell", "🔔", 1.2, 4.0, 15.0}
	SymClover  = Symbol{"Clover", "🍀", 1.0, 3.0, 12.0}
)

// Reel представляет один барабан
type Reel []Symbol

// ReelSet представляет набор из 5 барабанов
type ReelSet [5]Reel

// Individual представляет одну особь в популяции
type Individual struct {
	ReelSet ReelSet
	Fitness float64
	Stats   *SimStats
}

// SimStats статистика симуляции
type SimStats struct {
	Players        int
	Spins          int
	TotalBet       float64
	TotalWon       float64
	RTP            float64
	WinRate        float64
	CompanyProfit  float64
	ProfitPercent  float64 // Процент прибыли от всех денег игроков
	Wins           int

	// Статистика для режима "until the end"
	PlayersBroke   int     // Обанкротились
	PlayersWon     int     // Ушли победителями
	AvgSpinsPerPlayer float64 // Среднее количество спинов на игрока

	// Распределение размеров выигрышей (реальные казино)
	SmallWins    int     // 1-5x ставки (должно быть ~20-25% спинов)
	MediumWins   int     // 10-50x ставки (должно быть ~2-5% спинов)
	BigWins      int     // 50-100x ставки (должно быть ~0.1-1% спинов)
	JackpotWins  int     // 100x+ ставки (должно быть ~0.001-0.01% спинов)

	SmallWinRate    float64 // % от всех спинов
	MediumWinRate   float64 // % от всех спинов
	BigWinRate      float64 // % от всех спинов
	JackpotWinRate  float64 // % от всех спинов
}

// Payline линии выплат
var Paylines = [][5]int{
	{1, 1, 1, 1, 1}, // Middle
	{0, 0, 0, 0, 0}, // Top
	{2, 2, 2, 2, 2}, // Bottom
	{0, 1, 2, 1, 0}, // V
	{2, 1, 0, 1, 2}, // ^
	{1, 0, 1, 0, 1}, // Zigzag
	{1, 2, 1, 2, 1}, // Reverse zigzag
	{0, 1, 0, 1, 0}, // W
	{2, 1, 2, 1, 2}, // M
	{0, 0, 1, 2, 2}, // Diagonal
}

// GAConfig конфигурация генетического алгоритма
type GAConfig struct {
	PopSize          int
	Generations      int
	NormalMutations  int
	StrongMutations  int
	EliteCount       int
	PlayersPerSim    int
	SpinsPerPlayer   int
	TargetRTP        float64
	TargetWinRate    float64
	UntilTheEnd      bool    // Игроки играют до конца
	StartBalance     float64 // Стартовый баланс игрока
	WinThreshold     float64 // Порог победы (множитель от стартового баланса)
	SameReels        bool    // Все 5 барабанов одинаковые
}

// CreateOptimalStartReelSet создает оптимальный стартовый набор барабанов на основе теории вероятностей
func CreateOptimalStartReelSet(targetRTP float64, targetWinRate float64, sameReels bool) ReelSet {
	// Аналитический расчет:
	// Для RTP 95% и 10% побед нужно:
	// - Много низко-выигрышных символов (частые маленькие победы)
	// - Мало высоко-выигрышных символов (редкие большие победы)

	// Расчет: с 10 paylines вероятность совпадения 3+ символов должна быть ~10%
	// Если на каждом барабане символ встречается с вероятностью p,
	// то вероятность 3+ одинаковых = p^3 (упрощенно)
	// Для 10% побед: p^3 ≈ 0.01 => p ≈ 0.215 (21.5% барабана - один символ)

	// Создаем один оптимальный барабан
	reelSize := 50 // Средний размер
	baseReel := make(Reel, 0, reelSize)

	// Распределение символов (в процентах от барабана):
	// Низкие выплаты (частые) - 60%
	// Средние выплаты - 30%
	// Высокие выплаты - 10%

	// Clover (1.0x) - 30% - ОЧЕНЬ частые маленькие победы
	for j := 0; j < 15; j++ {
		baseReel = append(baseReel, SymClover)
	}

	// Bell (1.2x) - 20%
	for j := 0; j < 10; j++ {
		baseReel = append(baseReel, SymBell)
	}

	// Bar (1.5x) - 10%
	for j := 0; j < 5; j++ {
		baseReel = append(baseReel, SymBar)
	}

	// Cherry (2.0x) - 15%
	for j := 0; j < 8; j++ {
		baseReel = append(baseReel, SymCherry)
	}

	// Lemon (3.0x) - 10%
	for j := 0; j < 5; j++ {
		baseReel = append(baseReel, SymLemon)
	}

	// Orange (4.0x) - 7%
	for j := 0; j < 3; j++ {
		baseReel = append(baseReel, SymOrange)
	}

	// Grape (5.0x) - 4%
	for j := 0; j < 2; j++ {
		baseReel = append(baseReel, SymGrape)
	}

	// Diamond (8.0x) - 2%
	baseReel = append(baseReel, SymDiamond)

	// Star (10.0x) - 1%
	baseReel = append(baseReel, SymStar)

	// Seven (20.0x) - 0.5%
	if len(baseReel)%2 == 0 {
		baseReel = append(baseReel, SymSeven)
	}

	// ВАЖНО: Гарантируем что все символы присутствуют
	allSymbols := []Symbol{SymClover, SymBell, SymBar, SymCherry, SymLemon, SymOrange, SymGrape, SymDiamond, SymStar, SymSeven}
	baseReel = baseReel.ensureAllSymbols(allSymbols)

	var reels ReelSet

	if sameReels {
		// Все барабаны одинаковые
		for i := 0; i < 5; i++ {
			reels[i] = make(Reel, len(baseReel))
			copy(reels[i], baseReel)
		}
	} else {
		// Барабаны могут отличаться (старая логика)
		for i := 0; i < 5; i++ {
			reel := make(Reel, len(baseReel))
			copy(reel, baseReel)

			// Добавляем вариативность для разных барабанов
			if i%2 == 0 {
				reel = append(reel, SymStar)
			}
			if i == 2 {
				reel = append(reel, SymSeven)
			}

			reels[i] = reel
		}
	}

	return reels
}

// Spin выполняет один спин
func (rs ReelSet) Spin(bet float64, rng *rand.Rand) float64 {
	// Получаем 3 видимых символа на каждом барабане
	var visible [5][3]Symbol

	for i := 0; i < 5; i++ {
		if len(rs[i]) < 3 {
			return 0 // Невалидный барабан
		}

		// Случайная начальная позиция
		start := rng.Intn(len(rs[i]))

		for j := 0; j < 3; j++ {
			visible[i][j] = rs[i][(start+j)%len(rs[i])]
		}
	}

	// Проверяем все линии выплат
	totalWin := 0.0

	for _, payline := range Paylines {
		// Получаем символы вдоль линии
		var lineSymbols [5]Symbol
		for i := 0; i < 5; i++ {
			lineSymbols[i] = visible[i][payline[i]]
		}

		// Считаем совпадения слева направо
		firstSym := lineSymbols[0]
		count := 1

		for i := 1; i < 5; i++ {
			if lineSymbols[i].Name == firstSym.Name {
				count++
			} else {
				break
			}
		}

		// Выплата за совпадения
		if count >= 3 {
			var payout float64
			switch count {
			case 3:
				payout = firstSym.Payout3
			case 4:
				payout = firstSym.Payout4
			case 5:
				payout = firstSym.Payout5
			}
			totalWin += bet * payout
		}
	}

	return totalWin
}

// Simulate симулирует игровые сессии
func (rs ReelSet) Simulate(config GAConfig, bet float64) *SimStats {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	stats := &SimStats{
		Players: config.PlayersPerSim,
	}

	for p := 0; p < config.PlayersPerSim; p++ {
		balance := config.StartBalance
		playerSpins := 0

		if config.UntilTheEnd {
			// Режим "до конца" - играем пока не обанкротимся или не выиграем
			maxSpins := 10000 // Защита от бесконечного цикла
			winTarget := config.StartBalance * config.WinThreshold

			for s := 0; s < maxSpins; s++ {
				// Проверяем баланс
				if balance < bet {
					stats.PlayersBroke++
					break
				}

				if balance >= winTarget {
					stats.PlayersWon++
					break
				}

				// Делаем ставку
				balance -= bet
				stats.TotalBet += bet
				stats.Spins++
				playerSpins++

				// Крутим
				win := rs.Spin(bet, rng)
				balance += win
				stats.TotalWon += win

				if win > 0 {
					stats.Wins++

					// Классификация размера выигрыша
					multiplier := win / bet
					if multiplier >= 100 {
						stats.JackpotWins++ // 100x+
					} else if multiplier >= 50 {
						stats.BigWins++ // 50-100x
					} else if multiplier >= 10 {
						stats.MediumWins++ // 10-50x
					} else {
						stats.SmallWins++ // 1-10x (включая <10x)
					}
				}
			}
		} else {
			// Фиксированное количество спинов
			for s := 0; s < config.SpinsPerPlayer; s++ {
				if balance < bet {
					stats.PlayersBroke++
					break
				}

				balance -= bet
				stats.TotalBet += bet
				stats.Spins++
				playerSpins++

				win := rs.Spin(bet, rng)
				balance += win
				stats.TotalWon += win

				if win > 0 {
					stats.Wins++

					// Классификация размера выигрыша
					multiplier := win / bet
					if multiplier >= 100 {
						stats.JackpotWins++ // 100x+
					} else if multiplier >= 50 {
						stats.BigWins++ // 50-100x
					} else if multiplier >= 10 {
						stats.MediumWins++ // 10-50x
					} else {
						stats.SmallWins++ // 1-10x (включая <10x)
					}
				}
			}

			// Проверяем итог
			if balance >= config.StartBalance*config.WinThreshold {
				stats.PlayersWon++
			} else if balance < bet {
				stats.PlayersBroke++
			}
		}
	}

	if stats.TotalBet > 0 {
		stats.RTP = (stats.TotalWon / stats.TotalBet) * 100
		stats.ProfitPercent = ((stats.TotalBet - stats.TotalWon) / stats.TotalBet) * 100
	}

	if stats.Spins > 0 {
		stats.WinRate = (float64(stats.Wins) / float64(stats.Spins)) * 100
		stats.AvgSpinsPerPlayer = float64(stats.Spins) / float64(config.PlayersPerSim)

		// Рассчитываем процент каждого типа выигрыша от всех спинов
		stats.SmallWinRate = (float64(stats.SmallWins) / float64(stats.Spins)) * 100
		stats.MediumWinRate = (float64(stats.MediumWins) / float64(stats.Spins)) * 100
		stats.BigWinRate = (float64(stats.BigWins) / float64(stats.Spins)) * 100
		stats.JackpotWinRate = (float64(stats.JackpotWins) / float64(stats.Spins)) * 100
	}

	stats.CompanyProfit = stats.TotalBet - stats.TotalWon

	return stats
}

// EvaluateFitness оценивает приспособленность особи
func EvaluateFitness(ind *Individual, config GAConfig) {
	ind.Stats = ind.ReelSet.Simulate(config, 10.0)

	// Fitness = отклонение от целей (меньше = лучше)
	rtpDiff := math.Abs(ind.Stats.RTP - config.TargetRTP)
	winRateDiff := math.Abs(ind.Stats.WinRate - config.TargetWinRate)

	// Штраф если компания теряет деньги
	profitPenalty := 0.0
	if ind.Stats.CompanyProfit < 0 {
		profitPenalty = math.Abs(ind.Stats.CompanyProfit) / 100
	}

	// Валидация распределения размеров выигрышей (реальные казино)
	// Целевые значения:
	// - Мелкие (1-10x): 20-25% спинов
	// - Средние (10-50x): 2-5% спинов
	// - Большие (50-100x): 0.1-1% спинов
	// - Джекпот (100x+): 0.001-0.01% спинов
	distributionPenalty := 0.0

	// Пенальти за слишком много или слишком мало мелких выигрышей
	targetSmallWin := 22.5 // Середина диапазона 20-25%
	if ind.Stats.SmallWinRate < 15.0 || ind.Stats.SmallWinRate > 30.0 {
		distributionPenalty += math.Abs(ind.Stats.SmallWinRate-targetSmallWin) * 0.5
	}

	// Пенальти за слишком много средних выигрышей
	targetMediumWin := 3.5 // Середина диапазона 2-5%
	if ind.Stats.MediumWinRate > 8.0 {
		distributionPenalty += (ind.Stats.MediumWinRate - targetMediumWin) * 2.0
	}

	// Пенальти за слишком много больших выигрышей
	targetBigWin := 0.5 // Середина диапазона 0.1-1%
	if ind.Stats.BigWinRate > 2.0 {
		distributionPenalty += (ind.Stats.BigWinRate - targetBigWin) * 5.0
	}

	// Пенальти за слишком много джекпотов
	targetJackpot := 0.005 // Середина диапазона 0.001-0.01%
	if ind.Stats.JackpotWinRate > 0.05 {
		distributionPenalty += (ind.Stats.JackpotWinRate - targetJackpot) * 10.0
	}

	ind.Fitness = rtpDiff + winRateDiff*2 + profitPenalty + distributionPenalty
}

// countSymbol подсчитывает количество экземпляров символа на барабане
func (reel Reel) countSymbol(symbol Symbol) int {
	count := 0
	for _, s := range reel {
		if s.Name == symbol.Name {
			count++
		}
	}
	return count
}

// validateReel проверяет что все символы присутствуют хотя бы раз
func (reel Reel) validateReel(allSymbols []Symbol) bool {
	for _, symbol := range allSymbols {
		if reel.countSymbol(symbol) < 1 {
			return false
		}
	}
	return true
}

// ensureAllSymbols гарантирует что все символы присутствуют хотя бы раз
func (reel Reel) ensureAllSymbols(allSymbols []Symbol) Reel {
	newReel := make(Reel, len(reel))
	copy(newReel, reel)

	for _, symbol := range allSymbols {
		if newReel.countSymbol(symbol) < 1 {
			// Добавляем отсутствующий символ
			newReel = append(newReel, symbol)
		}
	}

	return newReel
}

// Mutate мутирует барабан
// ВАЖНО: Все символы должны присутствовать хотя бы раз на барабане
func (reel Reel) Mutate(rng *rand.Rand, strong bool) Reel {
	newReel := make(Reel, len(reel))
	copy(newReel, reel)

	mutations := 1
	if strong {
		mutations = rng.Intn(5) + 3 // 3-7 мутаций
	}

	allSymbols := []Symbol{SymClover, SymBell, SymBar, SymCherry, SymLemon, SymOrange, SymGrape, SymDiamond, SymStar, SymSeven}

	for m := 0; m < mutations; m++ {
		action := rng.Intn(3)

		switch action {
		case 0: // Заменить случайный символ
			if len(newReel) > 0 {
				idx := rng.Intn(len(newReel))
				newReel[idx] = allSymbols[rng.Intn(len(allSymbols))]
			}

		case 1: // Добавить символ
			if len(newReel) < 100 {
				newReel = append(newReel, allSymbols[rng.Intn(len(allSymbols))])
			}

		case 2: // Удалить символ
			if len(newReel) > 10 {
				idx := rng.Intn(len(newReel))
				symbolToRemove := newReel[idx]

				// КРИТИЧНО: проверяем что это не последний экземпляр символа
				if newReel.countSymbol(symbolToRemove) > 1 {
					newReel = append(newReel[:idx], newReel[idx+1:]...)
				}
				// Если это последний экземпляр - пропускаем удаление
			}
		}
	}

	// Гарантируем что все символы присутствуют
	newReel = newReel.ensureAllSymbols(allSymbols)

	return newReel
}

// Mutate мутирует набор барабанов
func (rs ReelSet) Mutate(rng *rand.Rand, strong bool, sameReels bool) ReelSet {
	newRS := ReelSet{}

	if sameReels {
		// Если барабаны одинаковые - мутируем один и копируем на все
		mutatedReel := rs[0].Mutate(rng, strong)
		for i := 0; i < 5; i++ {
			newRS[i] = make(Reel, len(mutatedReel))
			copy(newRS[i], mutatedReel)
		}
	} else {
		// Мутируем 1-3 случайных барабана
		reelsToMutate := rng.Intn(3) + 1
		if strong {
			reelsToMutate = rng.Intn(3) + 2 // 2-4 барабана
		}

		mutated := make(map[int]bool)
		for i := 0; i < reelsToMutate; i++ {
			reelIdx := rng.Intn(5)
			mutated[reelIdx] = true
		}

		for i := 0; i < 5; i++ {
			if mutated[i] {
				newRS[i] = rs[i].Mutate(rng, strong)
			} else {
				newRS[i] = make(Reel, len(rs[i]))
				copy(newRS[i], rs[i])
			}
		}
	}

	return newRS
}

// RunGA запускает генетический алгоритм
func RunGA(config GAConfig) *Individual {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	fmt.Println(strings.Repeat("█", 80))
	fmt.Println("ГЕНЕТИЧЕСКИЙ АЛГОРИТМ v2 - Оптимизация Барабанов")
	fmt.Println(strings.Repeat("█", 80))
	fmt.Printf("Цель: RTP %.1f%%, Win Rate %.1f%%\n", config.TargetRTP, config.TargetWinRate)
	fmt.Printf("Популяция: %d | Поколения: %d\n", config.PopSize, config.Generations)

	if config.SameReels {
		fmt.Printf("Режим барабанов: ВСЕ 5 ОДИНАКОВЫЕ\n")
	} else {
		fmt.Printf("Режим барабанов: Разные\n")
	}

	if config.UntilTheEnd {
		fmt.Printf("Симуляция: До конца (старт $%.0f, победа при $%.0f)\n",
			config.StartBalance, config.StartBalance*config.WinThreshold)
		fmt.Printf("Игроков: %d (играют до обанкротления или победы)\n", config.PlayersPerSim)
	} else {
		fmt.Printf("Симуляция: %d игроков × %d спинов\n", config.PlayersPerSim, config.SpinsPerPlayer)
	}

	fmt.Println(strings.Repeat("=", 80))

	// Создаем начальную популяцию
	population := make([]*Individual, config.PopSize)

	// Первые 3 особи - аналитически рассчитанные оптимальные
	for i := 0; i < 3 && i < config.PopSize; i++ {
		population[i] = &Individual{
			ReelSet: CreateOptimalStartReelSet(config.TargetRTP, config.TargetWinRate, config.SameReels),
		}
		EvaluateFitness(population[i], config)
	}

	// Остальные - случайные мутации оптимального
	optimalRS := CreateOptimalStartReelSet(config.TargetRTP, config.TargetWinRate, config.SameReels)
	for i := 3; i < config.PopSize; i++ {
		population[i] = &Individual{
			ReelSet: optimalRS.Mutate(rng, rng.Float64() < 0.3, config.SameReels),
		}
		EvaluateFitness(population[i], config)
	}

	// Сортируем по fitness
	sort.Slice(population, func(i, j int) bool {
		return population[i].Fitness < population[j].Fitness
	})

	best := population[0]
	fmt.Printf("\nПоколение 0: Fitness=%.2f | RTP=%.2f%% | WinRate=%.2f%% | Прибыль=%.2f%% ($%.0f)\n",
		best.Fitness, best.Stats.RTP, best.Stats.WinRate, best.Stats.ProfitPercent, best.Stats.CompanyProfit)

	// Эволюция
	for gen := 1; gen <= config.Generations; gen++ {
		newPop := make([]*Individual, 0, config.PopSize)

		// Топ-3 переходят без изменений
		for i := 0; i < config.EliteCount && i < len(population); i++ {
			elite := &Individual{ReelSet: ReelSet{}}
			for j := 0; j < 5; j++ {
				elite.ReelSet[j] = make(Reel, len(population[i].ReelSet[j]))
				copy(elite.ReelSet[j], population[i].ReelSet[j])
			}
			elite.Fitness = population[i].Fitness
			elite.Stats = population[i].Stats
			newPop = append(newPop, elite)
		}

		// 20 обычных мутаций лучших
		for i := 0; i < config.NormalMutations; i++ {
			parent := population[rng.Intn(len(population)/2)] // Выбираем из лучшей половины
			child := &Individual{
				ReelSet: parent.ReelSet.Mutate(rng, false, config.SameReels),
			}
			EvaluateFitness(child, config)
			newPop = append(newPop, child)
		}

		// 10 сильных мутаций
		for i := 0; i < config.StrongMutations; i++ {
			parent := population[rng.Intn(len(population)/3)] // Выбираем из лучшей трети
			child := &Individual{
				ReelSet: parent.ReelSet.Mutate(rng, true, config.SameReels),
			}
			EvaluateFitness(child, config)
			newPop = append(newPop, child)
		}

		population = newPop

		// Сортируем
		sort.Slice(population, func(i, j int) bool {
			return population[i].Fitness < population[j].Fitness
		})

		best = population[0]

		if gen%5 == 0 || gen == config.Generations {
			fmt.Printf("Поколение %d: Fitness=%.2f | RTP=%.2f%% | WinRate=%.2f%% | Прибыль=%.2f%% ($%.0f)\n",
				gen, best.Fitness, best.Stats.RTP, best.Stats.WinRate, best.Stats.ProfitPercent, best.Stats.CompanyProfit)
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("ЭВОЛЮЦИЯ ЗАВЕРШЕНА!")
	fmt.Println(strings.Repeat("=", 80))

	return best
}

// PrintBest выводит лучший результат
func PrintBest(best *Individual, config GAConfig) {
	fmt.Println("\n📊 ЛУЧШАЯ КОНФИГУРАЦИЯ:")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("RTP: %.2f%%\n", best.Stats.RTP)
	fmt.Printf("Win Rate: %.2f%%\n", best.Stats.WinRate)
	fmt.Printf("Всего ставок: $%.2f\n", best.Stats.TotalBet)
	fmt.Printf("Всего выплат: $%.2f\n", best.Stats.TotalWon)
	fmt.Printf("💰 Прибыль компании: $%.2f (%.2f%% от всех денег игроков)\n",
		best.Stats.CompanyProfit, best.Stats.ProfitPercent)
	fmt.Printf("Всего спинов: %d\n", best.Stats.Spins)
	fmt.Printf("Побед: %d\n", best.Stats.Wins)

	fmt.Println("\n💰 РАСПРЕДЕЛЕНИЕ РАЗМЕРОВ ВЫИГРЫШЕЙ:")
	fmt.Println("  (Реальные казино: Мелкие 20-25%, Средние 2-5%, Большие 0.1-1%, Джекпот 0.001-0.01%)")
	fmt.Printf("  Мелкие (1-10x):   %d (%.2f%% спинов)\n",
		best.Stats.SmallWins, best.Stats.SmallWinRate)
	fmt.Printf("  Средние (10-50x): %d (%.2f%% спинов)\n",
		best.Stats.MediumWins, best.Stats.MediumWinRate)
	fmt.Printf("  Большие (50-100x): %d (%.2f%% спинов)\n",
		best.Stats.BigWins, best.Stats.BigWinRate)
	fmt.Printf("  Джекпот (100x+):   %d (%.3f%% спинов)\n",
		best.Stats.JackpotWins, best.Stats.JackpotWinRate)

	if config.UntilTheEnd {
		playersLeft := best.Stats.Players - best.Stats.PlayersBroke - best.Stats.PlayersWon

		fmt.Printf("\n👥 СТАТИСТИКА ИГРОКОВ:\n")
		fmt.Printf("Обанкротились: %d (%.1f%%)\n",
			best.Stats.PlayersBroke,
			float64(best.Stats.PlayersBroke)/float64(best.Stats.Players)*100)
		fmt.Printf("Ушли победителями: %d (%.1f%%)\n",
			best.Stats.PlayersWon,
			float64(best.Stats.PlayersWon)/float64(best.Stats.Players)*100)
		fmt.Printf("Остались играть: %d (%.1f%%)\n",
			playersLeft,
			float64(playersLeft)/float64(best.Stats.Players)*100)
		fmt.Printf("Среднее спинов на игрока: %.1f\n", best.Stats.AvgSpinsPerPlayer)
	}

	fmt.Println("\n🎰 БАРАБАНЫ:")

	if config.SameReels {
		// Показываем только первый барабан, т.к. все одинаковые
		fmt.Printf("\n⚠️  ВСЕ 5 БАРАБАНОВ ОДИНАКОВЫЕ (длина: %d):\n", len(best.ReelSet[0]))
		reel := best.ReelSet[0]

		// Подсчет символов
		counts := make(map[string]int)
		for _, sym := range reel {
			counts[sym.Emoji]++
		}

		// Сортируем по частоте
		type kv struct {
			Symbol string
			Count  int
		}
		var sorted []kv
		for k, v := range counts {
			sorted = append(sorted, kv{k, v})
		}
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Count > sorted[j].Count
		})

		for _, item := range sorted {
			pct := float64(item.Count) / float64(len(reel)) * 100
			fmt.Printf("  %s: %d (%.1f%%)\n", item.Symbol, item.Count, pct)
		}
	} else {
		// Показываем все барабаны
		for i, reel := range best.ReelSet {
			fmt.Printf("\nБарабан %d (длина: %d):\n", i+1, len(reel))

			// Подсчет символов
			counts := make(map[string]int)
			for _, sym := range reel {
				counts[sym.Emoji]++
			}

			// Сортируем по частоте
			type kv struct {
				Symbol string
				Count  int
			}
			var sorted []kv
			for k, v := range counts {
				sorted = append(sorted, kv{k, v})
			}
			sort.Slice(sorted, func(i, j int) bool {
				return sorted[i].Count > sorted[j].Count
			})

			for _, item := range sorted {
				pct := float64(item.Count) / float64(len(reel)) * 100
				fmt.Printf("  %s: %d (%.1f%%)\n", item.Symbol, item.Count, pct)
			}
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
}

func main() {
	// Параметры основаны на реальных онлайн-казино:
	// RTP 94-96% - стандарт онлайн-слотов
	// Win Rate 22-28% - каждый 4-5й спин выигрывает (игрок доволен)
	targetRTP := flag.Float64("rtp", 95.0, "Целевой RTP % (реальные казино: 94-96%)")
	targetWinRate := flag.Float64("winrate", 25.0, "Целевая частота побед % (реальные казино: 22-28%)")
	generations := flag.Int("gen", 50, "Количество поколений")
	popSize := flag.Int("pop", 33, "Размер популяции (топ3 + 20 мутаций + 10 сильных)")
	untilEnd := flag.Bool("until-end", false, "Игроки играют до конца (обанкротятся или выиграют)")
	startBalance := flag.Float64("balance", 1000.0, "Стартовый баланс игрока")
	winMultiplier := flag.Float64("win-mult", 3.0, "Множитель для победы (игрок уходит победителем)")
	sameReels := flag.Bool("same-reels", false, "Все 5 барабанов одинаковые")

	flag.Parse()

	config := GAConfig{
		PopSize:         *popSize,
		Generations:     *generations,
		NormalMutations: 20,
		StrongMutations: 10,
		EliteCount:      3,
		PlayersPerSim:   50,
		SpinsPerPlayer:  100,
		TargetRTP:       *targetRTP,
		TargetWinRate:   *targetWinRate,
		UntilTheEnd:     *untilEnd,
		StartBalance:    *startBalance,
		WinThreshold:    *winMultiplier,
		SameReels:       *sameReels,
	}

	best := RunGA(config)
	PrintBest(best, config)
}
