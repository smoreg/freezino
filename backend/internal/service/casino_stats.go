package service

import (
	"github.com/smoreg/freezino/backend/internal/database"
	"github.com/smoreg/freezino/backend/internal/model"
	"gorm.io/gorm"
)

// CasinoStatsService provides business logic for casino-wide statistics
type CasinoStatsService struct {
	db *gorm.DB
}

// NewCasinoStatsService creates a new casino stats service instance
func NewCasinoStatsService() *CasinoStatsService {
	return &CasinoStatsService{
		db: database.GetDB(),
	}
}

// PlayerStats represents statistics for individual players
type PlayerStats struct {
	UserID      uint    `json:"user_id"`
	Username    string  `json:"username"`
	TotalBet    float64 `json:"total_bet"`
	TotalWon    float64 `json:"total_won"`
	NetProfit   float64 `json:"net_profit"`
	GamesPlayed int     `json:"games_played"`
}

// CasinoStatsResponse represents overall casino statistics
type CasinoStatsResponse struct {
	// Overall metrics
	TotalPlayers     int     `json:"total_players"`
	TotalGamesPlayed int     `json:"total_games_played"`
	TotalBet         float64 `json:"total_bet"`
	TotalWon         float64 `json:"total_won"`
	HouseProfit      float64 `json:"house_profit"`       // total_bet - total_won
	HouseEdgePercent float64 `json:"house_edge_percent"` // (house_profit / total_bet) * 100

	// Player profitability
	PlayersInProfit   int     `json:"players_in_profit"`  // players with net_profit > 0
	PlayersInLoss     int     `json:"players_in_loss"`    // players with net_profit < 0
	PlayersBreakEven  int     `json:"players_break_even"` // players with net_profit == 0
	ProfitablePercent float64 `json:"profitable_percent"` // (players_in_profit / total_players) * 100

	// Average metrics
	AverageBetPerGame float64 `json:"average_bet_per_game"`
	AverageWinPerGame float64 `json:"average_win_per_game"`

	// Game breakdown
	GameBreakdown []GameTypeStats `json:"game_breakdown"`

	// Top players
	TopWinners []PlayerStats `json:"top_winners"` // top 10 by net profit
	TopLosers  []PlayerStats `json:"top_losers"`  // top 10 by net loss
}

// GetCasinoStats retrieves overall casino statistics
func (s *CasinoStatsService) GetCasinoStats() (*CasinoStatsResponse, error) {
	stats := &CasinoStatsResponse{}

	// Get overall game statistics
	var overallStats struct {
		TotalGames int
		TotalBet   float64
		TotalWon   float64
	}
	s.db.Model(&model.GameSession{}).
		Select("COUNT(*) as total_games, COALESCE(SUM(bet), 0) as total_bet, COALESCE(SUM(win), 0) as total_won").
		Scan(&overallStats)

	stats.TotalGamesPlayed = overallStats.TotalGames
	stats.TotalBet = overallStats.TotalBet
	stats.TotalWon = overallStats.TotalWon
	stats.HouseProfit = overallStats.TotalBet - overallStats.TotalWon

	// Calculate house edge percentage
	if stats.TotalBet > 0 {
		stats.HouseEdgePercent = (stats.HouseProfit / stats.TotalBet) * 100
	}

	// Calculate average bet and win per game
	if stats.TotalGamesPlayed > 0 {
		stats.AverageBetPerGame = stats.TotalBet / float64(stats.TotalGamesPlayed)
		stats.AverageWinPerGame = stats.TotalWon / float64(stats.TotalGamesPlayed)
	}

	// Get player statistics (grouped by user)
	var playerStats []struct {
		UserID      uint
		Username    string
		TotalBet    float64
		TotalWon    float64
		NetProfit   float64
		GamesPlayed int
	}

	s.db.Model(&model.GameSession{}).
		Select(`
			game_sessions.user_id,
			users.username,
			COALESCE(SUM(game_sessions.bet), 0) as total_bet,
			COALESCE(SUM(game_sessions.win), 0) as total_won,
			COALESCE(SUM(game_sessions.win - game_sessions.bet), 0) as net_profit,
			COUNT(*) as games_played
		`).
		Joins("LEFT JOIN users ON users.id = game_sessions.user_id").
		Group("game_sessions.user_id, users.username").
		Scan(&playerStats)

	// Count total players who have played
	stats.TotalPlayers = len(playerStats)

	// Categorize players by profitability
	for _, player := range playerStats {
		if player.NetProfit > 0 {
			stats.PlayersInProfit++
		} else if player.NetProfit < 0 {
			stats.PlayersInLoss++
		} else {
			stats.PlayersBreakEven++
		}
	}

	// Calculate profitable player percentage
	if stats.TotalPlayers > 0 {
		stats.ProfitablePercent = (float64(stats.PlayersInProfit) / float64(stats.TotalPlayers)) * 100
	}

	// Get breakdown by game type
	var gameBreakdown []struct {
		GameType    string
		GamesPlayed int
		TotalBet    float64
		TotalWon    float64
	}
	s.db.Model(&model.GameSession{}).
		Select("game_type, COUNT(*) as games_played, COALESCE(SUM(bet), 0) as total_bet, COALESCE(SUM(win), 0) as total_won").
		Group("game_type").
		Order("games_played DESC").
		Scan(&gameBreakdown)

	stats.GameBreakdown = make([]GameTypeStats, len(gameBreakdown))
	for i, item := range gameBreakdown {
		stats.GameBreakdown[i] = GameTypeStats{
			GameType:    item.GameType,
			GamesPlayed: item.GamesPlayed,
			TotalBet:    item.TotalBet,
			TotalWon:    item.TotalWon,
			NetProfit:   item.TotalWon - item.TotalBet,
		}
	}

	// Get top 10 winners (highest net profit)
	var topWinnersData []struct {
		UserID      uint
		Username    string
		TotalBet    float64
		TotalWon    float64
		NetProfit   float64
		GamesPlayed int
	}

	s.db.Model(&model.GameSession{}).
		Select(`
			game_sessions.user_id,
			users.username,
			COALESCE(SUM(game_sessions.bet), 0) as total_bet,
			COALESCE(SUM(game_sessions.win), 0) as total_won,
			COALESCE(SUM(game_sessions.win - game_sessions.bet), 0) as net_profit,
			COUNT(*) as games_played
		`).
		Joins("LEFT JOIN users ON users.id = game_sessions.user_id").
		Group("game_sessions.user_id, users.username").
		Order("net_profit DESC").
		Limit(10).
		Scan(&topWinnersData)

	stats.TopWinners = make([]PlayerStats, len(topWinnersData))
	for i, winner := range topWinnersData {
		stats.TopWinners[i] = PlayerStats{
			UserID:      winner.UserID,
			Username:    winner.Username,
			TotalBet:    winner.TotalBet,
			TotalWon:    winner.TotalWon,
			NetProfit:   winner.NetProfit,
			GamesPlayed: winner.GamesPlayed,
		}
	}

	// Get top 10 losers (lowest net profit)
	var topLosersData []struct {
		UserID      uint
		Username    string
		TotalBet    float64
		TotalWon    float64
		NetProfit   float64
		GamesPlayed int
	}

	s.db.Model(&model.GameSession{}).
		Select(`
			game_sessions.user_id,
			users.username,
			COALESCE(SUM(game_sessions.bet), 0) as total_bet,
			COALESCE(SUM(game_sessions.win), 0) as total_won,
			COALESCE(SUM(game_sessions.win - game_sessions.bet), 0) as net_profit,
			COUNT(*) as games_played
		`).
		Joins("LEFT JOIN users ON users.id = game_sessions.user_id").
		Group("game_sessions.user_id, users.username").
		Order("net_profit ASC").
		Limit(10).
		Scan(&topLosersData)

	stats.TopLosers = make([]PlayerStats, len(topLosersData))
	for i, loser := range topLosersData {
		stats.TopLosers[i] = PlayerStats{
			UserID:      loser.UserID,
			Username:    loser.Username,
			TotalBet:    loser.TotalBet,
			TotalWon:    loser.TotalWon,
			NetProfit:   loser.NetProfit,
			GamesPlayed: loser.GamesPlayed,
		}
	}

	return stats, nil
}
