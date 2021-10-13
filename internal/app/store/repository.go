package store

import (
	"github.com/ArtemGontar/betting/internal/app/models"
)

type MatchResultRepository interface {
	InsertMatchResults([]models.MatchResult)
	SelectAgainstEachOtherResults(string, string) ([]Result, error)
	SelectLastFiveGamesByTeam(string) ([]Result, error)
	SelectLeagueAvgScoredGoals(int) (float64, float64, error)
	SelectHomeTeamAvgGoals(string) (float64, float64, error)
	SelectAwayTeamAvgGoals(string) (float64, float64, error)
}
