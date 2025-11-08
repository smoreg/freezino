package service

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
)

// Country represents a country with wage information
type Country struct {
	Code          string  `json:"code"`
	Name          string  `json:"name"`
	AvgHourlyWage float64 `json:"avgHourlyWage"`
	Currency      string  `json:"currency"`
}

// CountryStatsResponse represents country statistics with work time calculation
type CountryStatsResponse struct {
	Code             string  `json:"code"`
	Name             string  `json:"name"`
	AvgHourlyWage    float64 `json:"avgHourlyWage"`
	Currency         string  `json:"currency"`
	HoursToEarn500   float64 `json:"hoursToEarn500"`
	DaysToEarn500    float64 `json:"daysToEarn500"`
	ComparisonToGame float64 `json:"comparisonToGame"` // How many times longer than game work timer (3 min = 0.05 hours)
}

// StatsService provides business logic for statistics operations
type StatsService struct {
	countries []Country
}

// NewStatsService creates a new stats service instance
func NewStatsService() (*StatsService, error) {
	service := &StatsService{}
	if err := service.loadCountries(); err != nil {
		return nil, err
	}
	return service, nil
}

// loadCountries loads country data from JSON file
func (s *StatsService) loadCountries() error {
	// Get the path to countries.json file
	dataPath := filepath.Join("backend", "internal", "data", "countries.json")

	// Try alternative paths if the first one doesn't work
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		dataPath = filepath.Join("internal", "data", "countries.json")
	}

	data, err := os.ReadFile(dataPath)
	if err != nil {
		return fmt.Errorf("failed to read countries data: %w", err)
	}

	if err := json.Unmarshal(data, &s.countries); err != nil {
		return fmt.Errorf("failed to parse countries data: %w", err)
	}

	return nil
}

// GetCountries returns all countries with calculated statistics
func (s *StatsService) GetCountries() []CountryStatsResponse {
	const targetAmount = 500.0
	const gameWorkHours = 0.05 // 3 minutes = 0.05 hours

	results := make([]CountryStatsResponse, len(s.countries))

	for i, country := range s.countries {
		hoursToEarn := targetAmount / country.AvgHourlyWage
		daysToEarn := hoursToEarn / 8.0 // Assuming 8-hour work day
		comparisonToGame := hoursToEarn / gameWorkHours

		results[i] = CountryStatsResponse{
			Code:             country.Code,
			Name:             country.Name,
			AvgHourlyWage:    math.Round(country.AvgHourlyWage*100) / 100,
			Currency:         country.Currency,
			HoursToEarn500:   math.Round(hoursToEarn*100) / 100,
			DaysToEarn500:    math.Round(daysToEarn*100) / 100,
			ComparisonToGame: math.Round(comparisonToGame*100) / 100,
		}
	}

	return results
}

// GetCountryByCode returns a specific country by its code
func (s *StatsService) GetCountryByCode(code string) (*CountryStatsResponse, error) {
	const targetAmount = 500.0
	const gameWorkHours = 0.05

	for _, country := range s.countries {
		if country.Code == code {
			hoursToEarn := targetAmount / country.AvgHourlyWage
			daysToEarn := hoursToEarn / 8.0
			comparisonToGame := hoursToEarn / gameWorkHours

			result := &CountryStatsResponse{
				Code:             country.Code,
				Name:             country.Name,
				AvgHourlyWage:    math.Round(country.AvgHourlyWage*100) / 100,
				Currency:         country.Currency,
				HoursToEarn500:   math.Round(hoursToEarn*100) / 100,
				DaysToEarn500:    math.Round(daysToEarn*100) / 100,
				ComparisonToGame: math.Round(comparisonToGame*100) / 100,
			}

			return result, nil
		}
	}

	return nil, fmt.Errorf("country not found")
}
