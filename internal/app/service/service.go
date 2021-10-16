package service

import "github.com/ArtemGontar/betting/internal/app/model"

type LeagueService interface {
	LeagueTable(matchResults []model.MatchResult) []model.TableItem
}

type StatisticService interface {
	LeagueStatistics(league int) model.LeagueStatistic
	TeamStatistics(teamStat *model.TeamStatistic, avgLeagueScoredGoals float64, avgLeagueConcededGoals float64)
	AgainstEachOtherResults(team1Stat *model.TeamStatistic, team2Stat *model.TeamStatistic) []model.Result
	PoissonDistribution(teamForStat *model.TeamStatistic, teamAgainstStat model.TeamStatistic, avgForScoredGoals float64)
}
