package service

import (
	"fmt"
	"math"

	"github.com/ArtemGontar/betting/internal/app/model"
	"github.com/ArtemGontar/betting/internal/app/store"
)

func LeagueAvgStatistics(store store.Store, league int) model.LeagueAvgStatistic {
	avgHomeScoredGoals, avgAwayScoredGoals, err := store.MatchResult().SelectLeagueAvgScoredGoals(league)
	if err != nil {
		fmt.Println(err)
	}
	avgHomeConcededGoals := avgAwayScoredGoals
	avgAwayConcededGoals := avgHomeScoredGoals
	fmt.Println(league, "avg home scored goals =", avgHomeScoredGoals)
	fmt.Println(league, "avg away scored goals =", avgAwayScoredGoals)
	fmt.Println(league, "avg home conceded goals =", avgHomeConcededGoals)
	fmt.Println(league, "avg away conceded goals =", avgAwayConcededGoals)

	return model.LeagueAvgStatistic{
		AvgHomeScoredGoals:   avgHomeScoredGoals,
		AvgHomeConcededGoals: avgHomeConcededGoals,
		AvgAwayScoredGoals:   avgAwayScoredGoals,
		AvgAwayConcededGoals: avgAwayConcededGoals,
	}
}

func TeamStatistics(store store.Store, teamStat *model.TeamStatistic, avgLeagueScoredGoals float64, avgLeagueConcededGoals float64) {
	avgScoredGoals, avgConcededGoals, err := store.MatchResult().SelectTeamAvgGoals(teamStat.TeamName, teamStat.IsHome)
	if err != nil {
		fmt.Println(err)
	}

	fullTimeResults, err := store.MatchResult().SelectLastFiveGamesByTeam(teamStat.TeamName)
	if err != nil {
		fmt.Println(err)
	}

	results := matchResults(fullTimeResults, teamStat.TeamName)
	attackPower := avgScoredGoals / avgLeagueScoredGoals
	defencePower := avgConcededGoals / avgLeagueConcededGoals

	teamStat.FullTimeResults = results
	teamStat.AttackPower = attackPower
	teamStat.DefencePower = defencePower
	teamStat.AvgScoredGoals = avgScoredGoals
	teamStat.AvgConcededGoals = avgConcededGoals
}

func PoissonDistribution(teamForStat *model.TeamStatistic, teamAgainstStat model.TeamStatistic, avgScoredGoals float64) {
	predictScore := teamForStat.AttackPower * teamAgainstStat.DefencePower * avgScoredGoals

	//1.213^(5)*e^(-1.213)/(5!)
	poissonDistribution := []float64{}
	for i := 0; i < 5; i++ {
		poissonDistribution = append(poissonDistribution, math.Pow(predictScore, (float64(i)))*math.Pow(math.E, -predictScore)/float64(factorial(i)))
	}

	teamForStat.PredictScore = predictScore
	teamForStat.PoissonDistribution = poissonDistribution
}

func matchResults(results []store.Result, team string) []string {
	var resultsArr []string
	for _, result := range results {
		if result.Result == "H" {
			resultsArr = append(resultsArr, "W")
		} else if result.AwayTeam == team && result.Result == "A" {
			resultsArr = append(resultsArr, "W")
		} else if result.HomeTeam == team && result.Result == "A" {
			resultsArr = append(resultsArr, "L")
		} else if result.AwayTeam == team && result.Result == "H" {
			resultsArr = append(resultsArr, "L")
		} else {
			resultsArr = append(resultsArr, "D")
		}
	}
	return resultsArr
}

func factorial(n int) int {

	if n == 0 {
		return 1
	}
	return n * factorial(n-1)
}
