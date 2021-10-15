package store

import "github.com/ArtemGontar/betting/internal/app/model"

type MatchResultRepository interface {
	InsertMatchResults([]model.MatchResult)
	SelectAgainstEachOtherResults(string, string) ([]model.Result, error)
	SelectLastFiveGamesByTeam(string) ([]model.Result, error)
	SelectLeagueAvgScoredGoals(int) (float64, float64, error)
	SelectTeamAvgGoals(string, bool) (float64, float64, error)
}
