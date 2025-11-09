package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/smoreg/freezino/backend/internal/game"
)

// ReelWeights defines the probability weights for each symbol on a reel
type ReelWeights map[game.SlotSymbol]int

// Individual represents a candidate solution (reel weight configuration)
type Individual struct {
	Weights ReelWeights
	Fitness float64
	Stats   *SimulationStats
}

// GeneticAlgorithmConfig holds GA parameters
type GeneticAlgorithmConfig struct {
	PopulationSize  int
	Generations     int
	SimulationSpins int
	MutationRate    float64
	EliteSize       int
	TargetRTP       float64
	MinWinRate      float64 // Minimum % of spins that should win (any win)
	CrossoverRate   float64
}

// SimulationStats holds statistics from a Monte Carlo simulation
type SimulationStats struct {
	// Player statistics
	TotalPlayers      int
	PlayersWon        int     // Players who ended with more money
	PlayersLost       int     // Players who ended with less money
	PlayersBroke      int     // Players who went broke
	PlayersJackpot    int     // Players who hit 100k and cashed out

	// Spin statistics
	TotalSpins        int
	TotalBet          float64
	TotalWon          float64
	RTP               float64 // Return to Player %
	HouseEdge         float64 // House edge %
	WinCount          int
	WinRate           float64
	BigWinCount       int     // Wins >= 10x bet
	MegaWinCount      int     // Wins >= 50x bet
	JackpotCount      int     // Wins >= 100x bet

	// Financial statistics
	TotalInvested     float64 // Total money players brought
	TotalCashedOut    float64 // Total money players took
	CompanyProfit     float64 // Company's net profit
	CompanyRisk       float64 // Max potential loss (if all jackpots hit)

	LargestWin        float64
	AverageWin        float64
	WinDistribution   map[string]int
}

// PlayerSession represents a single player's gaming session
type PlayerSession struct {
	StartBalance  float64
	CurrentBalance float64
	SpinsPlayed   int
	MaxSpins      int
	BetSize       float64
	CashedOut     bool
	WentBroke     bool
	HitJackpot    bool
}

// WeightedSlotsEngine extends game.SlotsEngine with weighted symbol selection
type WeightedSlotsEngine struct {
	rng         *rand.Rand
	reelWeights ReelWeights
}

// NewWeightedSlotsEngine creates a slots engine with weighted symbols
func NewWeightedSlotsEngine(weights ReelWeights) *WeightedSlotsEngine {
	return &WeightedSlotsEngine{
		rng:         rand.New(rand.NewSource(time.Now().UnixNano())),
		reelWeights: weights,
	}
}

// generateReel generates a single reel with weighted random symbols
func (wse *WeightedSlotsEngine) generateReel() [3]game.SlotSymbol {
	var reel [3]game.SlotSymbol

	allSymbols := game.GetAllSymbols()
	totalWeight := 0
	for _, weight := range wse.reelWeights {
		totalWeight += weight
	}

	for i := 0; i < 3; i++ {
		roll := wse.rng.Intn(totalWeight)
		currentWeight := 0
		for _, symbol := range allSymbols {
			currentWeight += wse.reelWeights[symbol]
			if roll < currentWeight {
				reel[i] = symbol
				break
			}
		}
	}

	return reel
}

// SimpleSlotResult holds spin results without using game package types
type SimpleSlotResult struct {
	Reels       [5][3]game.SlotSymbol
	WinningLine []SimpleWinningLine
	TotalWin    float64
	Multiplier  float64
}

// SimpleWinningLine represents a winning line
type SimpleWinningLine struct {
	LineNumber int
	Symbol     game.SlotSymbol
	Count      int
	Multiplier float64
	Win        float64
}

// Spin performs a slot machine spin with weighted reels
func (wse *WeightedSlotsEngine) Spin(bet float64) *SimpleSlotResult {
	// Generate weighted reels
	var reels [5][3]game.SlotSymbol
	for i := 0; i < 5; i++ {
		reels[i] = wse.generateReel()
	}

	result := &SimpleSlotResult{
		Reels:       reels,
		WinningLine: []SimpleWinningLine{},
		TotalWin:    0,
		Multiplier:  0,
	}

	// Check all paylines
	paylines := getPaylines()
	payoutTable := game.GetPayoutTable()

	for lineNum, payline := range paylines {
		// Get symbols along this payline
		var symbols [5]game.SlotSymbol
		for i := 0; i < 5; i++ {
			symbols[i] = reels[i][payline[i]]
		}

		// Count consecutive matching symbols
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
		if count >= 3 {
			multiplier := payoutTable[firstSymbol][count]
			win := bet * multiplier

			winLine := SimpleWinningLine{
				LineNumber: lineNum + 1,
				Symbol:     firstSymbol,
				Count:      count,
				Multiplier: multiplier,
				Win:        win,
			}

			result.WinningLine = append(result.WinningLine, winLine)
			result.TotalWin += win
			result.Multiplier += multiplier
		}
	}

	return result
}

// getPaylines returns the standard paylines
func getPaylines() [][5]int {
	return [][5]int{
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
}

// RunMonteCarloSimulation runs simulation with N players each playing M spins
func RunMonteCarloSimulation(engine *WeightedSlotsEngine, numPlayers int, maxSpins int, startBalance float64, betSize float64) *SimulationStats {
	stats := &SimulationStats{
		TotalPlayers:    numPlayers,
		WinDistribution: make(map[string]int),
		TotalInvested:   float64(numPlayers) * startBalance,
	}

	const jackpotTarget = 100000.0

	// Simulate each player's session
	for p := 0; p < numPlayers; p++ {
		player := &PlayerSession{
			StartBalance:   startBalance,
			CurrentBalance: startBalance,
			MaxSpins:       maxSpins,
			BetSize:        betSize,
		}

		// Player plays until: broke, hit jackpot, or max spins reached
		for spin := 0; spin < maxSpins; spin++ {
			// Check if player can afford bet
			if player.CurrentBalance < betSize {
				player.WentBroke = true
				stats.PlayersBroke++
				break
			}

			// Place bet
			player.CurrentBalance -= betSize
			stats.TotalBet += betSize
			stats.TotalSpins++

			// Spin
			result := engine.Spin(betSize)
			player.CurrentBalance += result.TotalWin
			stats.TotalWon += result.TotalWin

			if result.TotalWin > 0 {
				stats.WinCount++

				multiplier := result.TotalWin / betSize
				if multiplier >= 100 {
					stats.JackpotCount++
				} else if multiplier >= 50 {
					stats.MegaWinCount++
				} else if multiplier >= 10 {
					stats.BigWinCount++
				}

				if result.TotalWin > stats.LargestWin {
					stats.LargestWin = result.TotalWin
				}
			}

			player.SpinsPlayed++

			// Check if player hit jackpot target
			if player.CurrentBalance >= jackpotTarget {
				player.HitJackpot = true
				player.CashedOut = true
				stats.PlayersJackpot++
				break
			}
		}

		// Player session ended - cash out
		if !player.WentBroke && !player.CashedOut {
			player.CashedOut = true
		}

		if player.CashedOut {
			stats.TotalCashedOut += player.CurrentBalance
		}

		// Track player outcomes
		if player.CurrentBalance > player.StartBalance {
			stats.PlayersWon++
		} else if player.CurrentBalance < player.StartBalance {
			stats.PlayersLost++
		}
	}

	// Calculate derived statistics
	stats.RTP = (stats.TotalWon / stats.TotalBet) * 100
	stats.HouseEdge = 100 - stats.RTP
	stats.WinRate = (float64(stats.WinCount) / float64(stats.TotalSpins)) * 100

	if stats.WinCount > 0 {
		stats.AverageWin = stats.TotalWon / float64(stats.WinCount)
	}

	// Financial outcomes
	stats.CompanyProfit = stats.TotalInvested - stats.TotalCashedOut

	// Risk calculation: what if all remaining players hit jackpot?
	potentialLoss := float64(stats.TotalPlayers-stats.PlayersBroke-stats.PlayersJackpot) * jackpotTarget
	stats.CompanyRisk = potentialLoss - stats.TotalInvested

	return stats
}

// generateRandomWeights creates a random weight configuration
// Favors common symbols heavily to achieve lower RTP
func generateRandomWeights(rng *rand.Rand) ReelWeights {
	weights := make(ReelWeights)
	allSymbols := game.GetAllSymbols()

	for _, symbol := range allSymbols {
		switch symbol {
		case game.SlotSymbol("🍒"): // Cherry - most common, lowest payout
			weights[symbol] = rng.Intn(30) + 70 // 70-100
		case game.SlotSymbol("🍋"): // Lemon
			weights[symbol] = rng.Intn(30) + 60 // 60-90
		case game.SlotSymbol("🍊"): // Orange
			weights[symbol] = rng.Intn(30) + 40 // 40-70
		case game.SlotSymbol("🍇"): // Grape
			weights[symbol] = rng.Intn(20) + 20 // 20-40
		case game.SlotSymbol("💎"): // Diamond
			weights[symbol] = rng.Intn(10) + 5 // 5-15
		case game.SlotSymbol("⭐"): // Star
			weights[symbol] = rng.Intn(5) + 1 // 1-6
		case game.SlotSymbol("7️⃣"): // Seven - rarest, highest payout
			weights[symbol] = rng.Intn(3) + 1 // 1-4
		default:
			weights[symbol] = rng.Intn(50) + 1
		}
	}

	return weights
}

// evaluateFitness calculates fitness score for an individual
func evaluateFitness(config GeneticAlgorithmConfig, weights ReelWeights) *Individual {
	engine := NewWeightedSlotsEngine(weights)

	// Simulate with realistic player sessions
	numPlayers := 1000
	maxSpins := 100
	startBalance := 1000.0
	betSize := 10.0

	stats := RunMonteCarloSimulation(engine, numPlayers, maxSpins, startBalance, betSize)

	individual := &Individual{
		Weights: weights,
		Stats:   stats,
	}

	// Calculate fitness as deviation from targets
	rtpDeviation := math.Abs(stats.RTP - config.TargetRTP)

	// Penalty if below minimum win rate
	winRatePenalty := 0.0
	if stats.WinRate < config.MinWinRate {
		winRatePenalty = (config.MinWinRate - stats.WinRate) * 20
	}

	// Additional penalty for high company risk (too many jackpots)
	riskPenalty := 0.0
	if stats.PlayersJackpot > numPlayers/100 { // More than 1% hit jackpot = bad
		riskPenalty = float64(stats.PlayersJackpot) * 5
	}

	// Penalty if company loses money
	profitPenalty := 0.0
	if stats.CompanyProfit < 0 {
		profitPenalty = math.Abs(stats.CompanyProfit) / 100
	}

	// Fitness = RTP deviation + win rate penalty + risk penalty + profit penalty
	individual.Fitness = rtpDeviation + winRatePenalty + riskPenalty + profitPenalty

	return individual
}

// crossover combines two parent weight configurations
func crossover(parent1, parent2 ReelWeights, rng *rand.Rand) ReelWeights {
	child := make(ReelWeights)
	allSymbols := game.GetAllSymbols()

	for _, symbol := range allSymbols {
		choice := rng.Float64()
		if choice < 0.4 {
			child[symbol] = parent1[symbol]
		} else if choice < 0.8 {
			child[symbol] = parent2[symbol]
		} else {
			child[symbol] = (parent1[symbol] + parent2[symbol]) / 2
			if child[symbol] < 1 {
				child[symbol] = 1
			}
		}
	}

	return child
}

// mutate randomly modifies weights
// ВАЖНО: Все символы должны иметь вес >= 1 (присутствовать на барабане)
func mutate(weights ReelWeights, mutationRate float64, rng *rand.Rand) ReelWeights {
	mutated := make(ReelWeights)
	allSymbols := game.GetAllSymbols()

	// Сначала копируем все веса
	for symbol, weight := range weights {
		mutated[symbol] = weight
	}

	// Гарантируем что все символы присутствуют
	for _, symbol := range allSymbols {
		if _, exists := mutated[symbol]; !exists {
			mutated[symbol] = 1
		}
	}

	// Мутируем
	for symbol, weight := range mutated {
		if rng.Float64() < mutationRate {
			change := rng.Float64()*0.6 - 0.3
			newWeight := int(float64(weight) * (1.0 + change))

			// КРИТИЧНО: вес не может быть меньше 1 (символ должен присутствовать)
			if newWeight < 1 {
				newWeight = 1
			}
			if newWeight > 100 {
				newWeight = 100
			}

			mutated[symbol] = newWeight
		}
	}

	return mutated
}

// selectParents performs tournament selection
func selectParents(population []*Individual, tournamentSize int, rng *rand.Rand) (*Individual, *Individual) {
	selectOne := func() *Individual {
		tournament := make([]*Individual, tournamentSize)
		for i := 0; i < tournamentSize; i++ {
			tournament[i] = population[rng.Intn(len(population))]
		}

		best := tournament[0]
		for _, ind := range tournament[1:] {
			if ind.Fitness < best.Fitness {
				best = ind
			}
		}
		return best
	}

	return selectOne(), selectOne()
}

// RunGeneticAlgorithm evolves optimal reel weights
func RunGeneticAlgorithm(config GeneticAlgorithmConfig) *Individual {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	fmt.Println(strings.Repeat("█", 80))
	fmt.Println("ГЕНЕТИЧЕСКИЙ АЛГОРИТМ - Выведение Оптимальной Конфигурации Слотов")
	fmt.Println(strings.Repeat("█", 80))
	fmt.Printf("Целевой RTP: %.2f%% | Мин. частота побед: %.2f%% | Популяция: %d | Поколения: %d\n",
		config.TargetRTP, config.MinWinRate, config.PopulationSize, config.Generations)
	fmt.Println(strings.Repeat("=", 80))

	// Initialize population
	population := make([]*Individual, config.PopulationSize)
	for i := 0; i < config.PopulationSize; i++ {
		weights := generateRandomWeights(rng)
		population[i] = evaluateFitness(config, weights)
	}

	// Sort by fitness
	sort.Slice(population, func(i, j int) bool {
		return population[i].Fitness < population[j].Fitness
	})

	best := population[0]
	fmt.Printf("\nПоколение 0: Fitness=%.4f | RTP=%.2f%% | Частота побед=%.2f%%\n",
		best.Fitness, best.Stats.RTP, best.Stats.WinRate)

	// Evolution loop
	for gen := 1; gen <= config.Generations; gen++ {
		newPopulation := make([]*Individual, 0, config.PopulationSize)

		// Elitism
		for i := 0; i < config.EliteSize && i < len(population); i++ {
			newPopulation = append(newPopulation, population[i])
		}

		// Generate offspring
		for len(newPopulation) < config.PopulationSize {
			parent1, parent2 := selectParents(population, 3, rng)

			var childWeights ReelWeights
			if rng.Float64() < config.CrossoverRate {
				childWeights = crossover(parent1.Weights, parent2.Weights, rng)
			} else {
				childWeights = parent1.Weights
			}

			childWeights = mutate(childWeights, config.MutationRate, rng)
			child := evaluateFitness(config, childWeights)
			newPopulation = append(newPopulation, child)
		}

		population = newPopulation

		// Sort by fitness
		sort.Slice(population, func(i, j int) bool {
			return population[i].Fitness < population[j].Fitness
		})

		best = population[0]

		// Log progress
		if gen%5 == 0 || gen == config.Generations {
			fmt.Printf("Поколение %d: Fitness=%.4f | RTP=%.2f%% | Край дома=%.2f%% | Частота побед=%.2f%%\n",
				gen, best.Fitness, best.Stats.RTP, best.Stats.HouseEdge, best.Stats.WinRate)
		}
	}

	// Final results
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("ЭВОЛЮЦИЯ ЗАВЕРШЕНА!")
	fmt.Println(strings.Repeat("=", 80))

	PrintStats("Лучшая Найденная Конфигурация", best.Stats)

	fmt.Println("\nОптимальные Веса Барабанов:")
	allSymbols := game.GetAllSymbols()
	for _, symbol := range allSymbols {
		fmt.Printf("  %s: %d\n", symbol, best.Weights[symbol])
	}

	fmt.Println("\nGo Код для slots.go:")
	fmt.Println("// Оптимальные веса для достижения", fmt.Sprintf("%.1f%% RTP", config.TargetRTP))
	fmt.Println("var symbolWeights = map[SlotSymbol]int{")
	for _, symbol := range allSymbols {
		fmt.Printf("    %s: %d,\n", getSymbolConstName(symbol), best.Weights[symbol])
	}
	fmt.Println("}")

	return best
}

// getSymbolConstName returns the constant name for a symbol
func getSymbolConstName(symbol game.SlotSymbol) string {
	switch symbol {
	case game.SlotSymbol("🍒"):
		return "SymbolCherry"
	case game.SlotSymbol("🍋"):
		return "SymbolLemon"
	case game.SlotSymbol("🍊"):
		return "SymbolOrange"
	case game.SlotSymbol("🍇"):
		return "SymbolGrape"
	case game.SlotSymbol("💎"):
		return "SymbolDiamond"
	case game.SlotSymbol("⭐"):
		return "SymbolStar"
	case game.SlotSymbol("7️⃣"):
		return "SymbolSeven"
	default:
		return "Unknown"
	}
}

// PrintStats prints simulation statistics
func PrintStats(name string, stats *SimulationStats) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Printf("%s - Результаты\n", name)
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n📊 СТАТИСТИКА ИГРОКОВ:")
	fmt.Printf("  Всего игроков:           %d\n", stats.TotalPlayers)
	fmt.Printf("  Выиграли:                %d (%.1f%%)\n", stats.PlayersWon, float64(stats.PlayersWon)/float64(stats.TotalPlayers)*100)
	fmt.Printf("  Проиграли:               %d (%.1f%%)\n", stats.PlayersLost, float64(stats.PlayersLost)/float64(stats.TotalPlayers)*100)
	fmt.Printf("  Обанкротились:           %d (%.1f%%)\n", stats.PlayersBroke, float64(stats.PlayersBroke)/float64(stats.TotalPlayers)*100)
	fmt.Printf("  Сорвали джекпот (100k+): %d (%.1f%%)\n", stats.PlayersJackpot, float64(stats.PlayersJackpot)/float64(stats.TotalPlayers)*100)

	fmt.Println("\n🎰 СТАТИСТИКА СПИНОВ:")
	fmt.Printf("  Всего спинов:            %d\n", stats.TotalSpins)
	fmt.Printf("  Частота побед:           %.2f%% (%d/%d)\n", stats.WinRate, stats.WinCount, stats.TotalSpins)
	fmt.Printf("  Средний выигрыш:         $%.2f\n", stats.AverageWin)
	fmt.Printf("  Макс. выигрыш:           $%.2f\n", stats.LargestWin)
	fmt.Printf("  Большие выигрыши (10x+): %d (%.2f%%)\n", stats.BigWinCount, float64(stats.BigWinCount)/float64(stats.TotalSpins)*100)
	fmt.Printf("  Мега выигрыши (50x+):    %d (%.2f%%)\n", stats.MegaWinCount, float64(stats.MegaWinCount)/float64(stats.TotalSpins)*100)
	fmt.Printf("  Джекпоты (100x+):        %d (%.2f%%)\n", stats.JackpotCount, float64(stats.JackpotCount)/float64(stats.TotalSpins)*100)

	fmt.Println("\n💰 ФИНАНСЫ:")
	fmt.Printf("  RTP:                     %.2f%%\n", stats.RTP)
	fmt.Printf("  Край дома:               %.2f%%\n", stats.HouseEdge)
	fmt.Printf("  Инвестировано игроками:  $%.2f\n", stats.TotalInvested)
	fmt.Printf("  Выплачено игрокам:       $%.2f\n", stats.TotalCashedOut)
	fmt.Printf("  Прибыль компании:        $%.2f\n", stats.CompanyProfit)

	if stats.CompanyProfit > 0 {
		fmt.Printf("  ✅ Компания в плюсе на %.2f%% от инвестиций\n", stats.CompanyProfit/stats.TotalInvested*100)
	} else {
		fmt.Printf("  ⚠️  Компания в минусе на %.2f%% от инвестиций\n", math.Abs(stats.CompanyProfit)/stats.TotalInvested*100)
	}

	fmt.Println(strings.Repeat("=", 80))
}

func main() {
	// Command-line flags
	// Параметры основаны на реальных онлайн-казино:
	// RTP 94-96% - стандарт онлайн-слотов
	// Win Rate 22-28% - каждый 4-5й спин выигрывает (игрок доволен)
	targetRTP := flag.Float64("rtp", 95.0, "Целевой RTP в процентах (реальные казино: 94-96%)")
	minWinRate := flag.Float64("win-rate", 25.0, "Минимальная частота побед в процентах (реальные казино: 22-28%)")
	generations := flag.Int("generations", 50, "Количество поколений для эволюции")
	population := flag.Int("population", 100, "Размер популяции")

	flag.Parse()

	config := GeneticAlgorithmConfig{
		PopulationSize:  *population,
		Generations:     *generations,
		SimulationSpins: 0, // Not used anymore, we use player sessions
		MutationRate:    0.15,
		EliteSize:       10,
		TargetRTP:       *targetRTP,
		MinWinRate:      *minWinRate,
		CrossoverRate:   0.7,
	}

	fmt.Println("\n🎰 ОПТИМИЗАТОР КОНФИГУРАЦИИ СЛОТОВ 🎰")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("Целевые параметры:\n")
	fmt.Printf("  • RTP: %.2f%% (игрок в среднем теряет %.2f%%)\n", config.TargetRTP, 100-config.TargetRTP)
	fmt.Printf("  • Частота побед: как минимум %.2f%% (1 из %.0f спинов)\n", config.MinWinRate, 100/config.MinWinRate)
	fmt.Printf("\nМодель симуляции:\n")
	fmt.Printf("  • 1000 параллельных игроков\n")
	fmt.Printf("  • Каждый игрок стартует с $1000\n")
	fmt.Printf("  • Каждый играет до 100 спинов по $10\n")
	fmt.Printf("  • Игрок выходит если: обанкротился или заработал $100,000+\n")
	fmt.Printf("\nПараметры алгоритма:\n")
	fmt.Printf("  • Поколения: %d\n", config.Generations)
	fmt.Printf("  • Популяция: %d\n", config.PopulationSize)
	fmt.Println(strings.Repeat("=", 80))

	best := RunGeneticAlgorithm(config)

	fmt.Println("\n" + strings.Repeat("🎉", 40))
	fmt.Println("ГОТОВО! Оптимальная конфигурация найдена.")
	fmt.Printf("Финальный RTP: %.2f%% | Частота побед: %.2f%%\n", best.Stats.RTP, best.Stats.WinRate)
	fmt.Printf("Прибыль компании: $%.2f из $%.2f инвестированных (%.1f%%)\n",
		best.Stats.CompanyProfit, best.Stats.TotalInvested,
		best.Stats.CompanyProfit/best.Stats.TotalInvested*100)
	fmt.Println(strings.Repeat("🎉", 40))
}
