package league_service

import (
	"github.com/ArtemGontar/betting/internal/app/model"
	"github.com/ArtemGontar/betting/internal/app/store"
)

type LeagueService struct {
	matchRepository store.MatchResultRepository
}

func New(matchRepository store.MatchResultRepository) *LeagueService {
	return &LeagueService{
		matchRepository: matchRepository,
	}
}

func (ls *LeagueService) LeagueTable(matchResults []model.MatchResult) []model.TableItem {
	items := make(map[string]model.TableItem)
	for _, res := range matchResults {
		//fill for home team
		fillMatchResult(items, res.HomeTeam, res.FullTimeHomeTeamGoals,
			res.FullTimeAwayTeamGoals, res.FullTimeResult)

		//fill for away team
		fillMatchResult(items, res.HomeTeam, res.FullTimeHomeTeamGoals,
			res.FullTimeAwayTeamGoals, res.FullTimeResult)
	}
	itemsAsSlice := make([]model.TableItem, 0, len(items))

	for _, item := range items {
		itemsAsSlice = append(itemsAsSlice, item)
	}
	return itemsAsSlice
}

func fillMatchResult(items map[string]model.TableItem, teamName string, scoredGoals int, concededGoals int, result string) {
	team, ok := items[teamName]
	if !ok {
		item := model.TableItem{
			TeamName:      teamName,
			MatchCount:    1,
			ScoredGoals:   scoredGoals,
			ConcededGoals: concededGoals,
			Points:        points(result),
		}
		items[teamName] = item
	} else {
		team.Points = team.Points + points(result)
		team.ScoredGoals = team.ScoredGoals + scoredGoals
		team.ConcededGoals = team.ConcededGoals + concededGoals
		team.MatchCount += 1
		items[team.TeamName] = team
	}
}

func points(res string) int {
	points := 0
	if res == "H" {
		points = 3
	} else if res == "D" {
		points = 1
	}
	return points
}
