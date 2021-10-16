package statistic_service

import (
	"fmt"
	"math"

	"github.com/ArtemGontar/betting/internal/app/model"
	"github.com/ArtemGontar/betting/internal/app/store"
)

type StatisticsService struct {
	matchRepository store.MatchResultRepository
}

func New(matchRepository store.MatchResultRepository) *StatisticsService {
	return &StatisticsService{
		matchRepository: matchRepository,
	}
}

func (ss *StatisticsService) LeagueStatistics(league int) model.LeagueStatistic {
	avgHomeScoredGoals, avgAwayScoredGoals, err := ss.matchRepository.SelectLeagueAvgScoredGoals(league)
	if err != nil {
		fmt.Println(err)
	}
	avgHomeConcededGoals := avgAwayScoredGoals
	avgAwayConcededGoals := avgHomeScoredGoals
	return model.LeagueStatistic{
		AvgHomeScoredGoals:   avgHomeScoredGoals,
		AvgHomeConcededGoals: avgHomeConcededGoals,
		AvgAwayScoredGoals:   avgAwayScoredGoals,
		AvgAwayConcededGoals: avgAwayConcededGoals,
	}
}

func (ss *StatisticsService) TeamStatistics(teamStat *model.TeamStatistic, avgLeagueScoredGoals float64, avgLeagueConcededGoals float64) {
	avgScoredGoals, avgConcededGoals, err := ss.matchRepository.SelectTeamAvgGoals(teamStat.TeamName, teamStat.IsHome)
	if err != nil {
		fmt.Println(err)
	}

	fullTimeResults, err := ss.matchRepository.SelectLastFiveGamesByTeam(teamStat.TeamName)
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

func (ss *StatisticsService) AgainstEachOtherResults(team1Stat *model.TeamStatistic, team2Stat *model.TeamStatistic) []model.Result {
	eachOtherGames, err := ss.matchRepository.SelectAgainstEachOtherResults(team1Stat.TeamName, team2Stat.TeamName)
	if err != nil {
		fmt.Println(err)
	}
	return eachOtherGames
}

func (ss *StatisticsService) PoissonDistribution(teamForStat *model.TeamStatistic, teamAgainstStat model.TeamStatistic, avgForScoredGoals float64) {
	predictScore := teamForStat.AttackPower * teamAgainstStat.DefencePower * avgForScoredGoals

	//1.213^(5)*e^(-1.213)/(5!)
	poissonDistribution := []float64{}
	for i := 0; i < 5; i++ {
		poissonDistribution = append(poissonDistribution, math.Pow(predictScore, (float64(i)))*math.Pow(math.E, -predictScore)/float64(factorial(i)))
	}

	teamForStat.PredictScore = predictScore
	teamForStat.PoissonDistribution = poissonDistribution
}

func matchResults(results []model.Result, team string) []string {
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
