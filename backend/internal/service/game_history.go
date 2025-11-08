package service

import (
	"errors"
	"fmt"

	"github.com/smoreg/freezino/backend/internal/database"
	"github.com/smoreg/freezino/backend/internal/model"
	"gorm.io/gorm"
)

// GameHistoryService provides business logic for game history operations
type GameHistoryService struct {
	db *gorm.DB
}

// NewGameHistoryService creates a new game history service instance
func NewGameHistoryService() *GameHistoryService {
	return &GameHistoryService{
		db: database.GetDB(),
	}
}

// GameHistoryItem represents a single game history entry
type GameHistoryItem struct {
	ID        uint    `json:"id"`
	GameType  string  `json:"game_type"`
	Bet       float64 `json:"bet"`
	Win       float64 `json:"win"`
	Profit    float64 `json:"profit"` // win - bet
	CreatedAt string  `json:"created_at"`
}

// GameHistoryResponse represents paginated game history
type GameHistoryResponse struct {
	Games  []GameHistoryItem `json:"games"`
	Total  int64             `json:"total"`
	Limit  int               `json:"limit"`
	Offset int               `json:"offset"`
}

// GameStatsResponse represents overall game statistics
type GameStatsResponse struct {
	TotalGames    int             `json:"total_games"`
	TotalWins     int             `json:"total_wins"`
	TotalLosses   int             `json:"total_losses"`
	TotalBet      float64         `json:"total_bet"`
	TotalWon      float64         `json:"total_won"`
	NetProfit     float64         `json:"net_profit"`
	FavoriteGame  string          `json:"favorite_game,omitempty"`
	WinRate       float64         `json:"win_rate"` // percentage
	BiggestWin    float64         `json:"biggest_win"`
	BiggestLoss   float64         `json:"biggest_loss"`
	GameBreakdown []GameTypeStats `json:"game_breakdown"`
}

// GameTypeStats represents statistics for a specific game type
type GameTypeStats struct {
	GameType    string  `json:"game_type"`
	GamesPlayed int     `json:"games_played"`
	TotalBet    float64 `json:"total_bet"`
	TotalWon    float64 `json:"total_won"`
	NetProfit   float64 `json:"net_profit"`
}

// GetHistory retrieves game history with optional filters
func (s *GameHistoryService) GetHistory(userID uint, gameType string, limit int, offset int) (*GameHistoryResponse, error) {
	// Verify user exists
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	query := s.db.Model(&model.GameSession{}).Where("user_id = ?", userID)

	// Apply game type filter if provided
	if gameType != "" {
		query = query.Where("game_type = ?", gameType)
	}

	// Count total matching records
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count game sessions: %w", err)
	}

	// Get paginated results
	var gameSessions []model.GameSession
	query = query.Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&gameSessions).Error; err != nil {
		return nil, fmt.Errorf("failed to get game sessions: %w", err)
	}

	// Transform to response format
	games := make([]GameHistoryItem, len(gameSessions))
	for i, session := range gameSessions {
		games[i] = GameHistoryItem{
			ID:        session.ID,
			GameType:  string(session.GameType),
			Bet:       session.Bet,
			Win:       session.Win,
			Profit:    session.Win - session.Bet,
			CreatedAt: session.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return &GameHistoryResponse{
		Games:  games,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

// GetStats retrieves overall game statistics for a user
func (s *GameHistoryService) GetStats(userID uint) (*GameStatsResponse, error) {
	// Verify user exists
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	stats := &GameStatsResponse{}

	// Get overall statistics
	var overallStats struct {
		TotalGames int
		TotalBet   float64
		TotalWon   float64
	}
	s.db.Model(&model.GameSession{}).
		Where("user_id = ?", userID).
		Select("COUNT(*) as total_games, COALESCE(SUM(bet), 0) as total_bet, COALESCE(SUM(win), 0) as total_won").
		Scan(&overallStats)

	stats.TotalGames = overallStats.TotalGames
	stats.TotalBet = overallStats.TotalBet
	stats.TotalWon = overallStats.TotalWon
	stats.NetProfit = overallStats.TotalWon - overallStats.TotalBet

	// Count wins and losses
	var wins int64
	s.db.Model(&model.GameSession{}).
		Where("user_id = ? AND win > bet", userID).
		Count(&wins)
	stats.TotalWins = int(wins)

	var losses int64
	s.db.Model(&model.GameSession{}).
		Where("user_id = ? AND win <= bet", userID).
		Count(&losses)
	stats.TotalLosses = int(losses)

	// Calculate win rate
	if stats.TotalGames > 0 {
		stats.WinRate = (float64(stats.TotalWins) / float64(stats.TotalGames)) * 100
	}

	// Find favorite game (most played)
	var favoriteGame struct {
		GameType string
	}
	s.db.Model(&model.GameSession{}).
		Where("user_id = ?", userID).
		Select("game_type").
		Group("game_type").
		Order("COUNT(*) DESC").
		Limit(1).
		Scan(&favoriteGame)
	stats.FavoriteGame = favoriteGame.GameType

	// Find biggest win
	var biggestWin struct {
		Win float64
	}
	s.db.Model(&model.GameSession{}).
		Where("user_id = ?", userID).
		Select("COALESCE(MAX(win - bet), 0) as win").
		Scan(&biggestWin)
	stats.BiggestWin = biggestWin.Win

	// Find biggest loss
	var biggestLoss struct {
		Loss float64
	}
	s.db.Model(&model.GameSession{}).
		Where("user_id = ? AND win < bet", userID).
		Select("COALESCE(MAX(bet - win), 0) as loss").
		Scan(&biggestLoss)
	stats.BiggestLoss = biggestLoss.Loss

	// Get breakdown by game type
	var gameBreakdown []struct {
		GameType    string
		GamesPlayed int
		TotalBet    float64
		TotalWon    float64
	}
	s.db.Model(&model.GameSession{}).
		Where("user_id = ?", userID).
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

	return stats, nil
}
