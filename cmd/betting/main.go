package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"sort"

	"github.com/ArtemGontar/betting/internal/app/apiserver"
	"github.com/ArtemGontar/betting/internal/app/model"
	"github.com/ArtemGontar/betting/internal/app/reader/premierleague"
	"github.com/ArtemGontar/betting/internal/app/service"
	league_service "github.com/ArtemGontar/betting/internal/app/service/league"
	statistic_service "github.com/ArtemGontar/betting/internal/app/service/statistic"
	match_results_repository "github.com/ArtemGontar/betting/internal/app/store/match_results"
	"github.com/BurntSushi/toml"
	_ "github.com/lib/pq"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

func main() {
	teams := premierleague.Teams()
	//premierleague_reader.Players()

	flag.Parse()
	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()
	matchResultsRepository := match_results_repository.New(db)
	statisticService := statistic_service.New(matchResultsRepository)
	leagueService := league_service.New(matchResultsRepository)
	matchResults := simulationSeason(leagueService, statisticService, teams)
	table := leagueService.LeagueTable(matchResults)

	sort.Slice(table, func(i, j int) bool {
		return table[i].Points > table[j].Points
	})
	for _, row := range table {
		fmt.Println(row.TeamName, row.MatchCount, row.ScoredGoals, row.ConcededGoals, row.Points)
	}
}

func simulationSeason(leagueService service.LeagueService, statisticService service.StatisticService, teams []model.Team) []model.MatchResult {
	matchResults := []model.MatchResult{}
	leagueStat := statisticService.LeagueStatistics(1)
	for i := 0; i < len(teams); i++ {
		team1 := model.TeamStatistic{
			TeamName: teams[i].TeamName,
			IsHome:   true,
		}
		statisticService.TeamStatistics(&team1, leagueStat.AvgHomeScoredGoals, leagueStat.AvgHomeConcededGoals)
		for j := 0; j < len(teams); j++ {
			if teams[i].TeamName == teams[j].TeamName {
				continue
			}
			team2 := model.TeamStatistic{
				TeamName: teams[j].TeamName,
				IsHome:   false,
			}
			statisticService.TeamStatistics(&team2, leagueStat.AvgAwayScoredGoals, leagueStat.AvgAwayConcededGoals)

			statisticService.PoissonDistribution(&team1, team2, leagueStat.AvgHomeScoredGoals)
			statisticService.PoissonDistribution(&team2, team1, leagueStat.AvgAwayScoredGoals)
			team1Goals := maxArrValue(team1.PoissonDistribution)
			team2Goals := maxArrValue(team2.PoissonDistribution)
			var mr string
			if team1Goals > team2Goals {
				mr = "H"
			} else if team2Goals > team1Goals {
				mr = "A"
			} else {
				mr = "D"
			}
			matchResult := model.MatchResult{
				HomeTeam:              team1.TeamName,
				FullTimeHomeTeamGoals: team1Goals,
				AwayTeam:              team2.TeamName,
				FullTimeAwayTeamGoals: team2Goals,
				FullTimeResult:        mr,
			}
			matchResults = append(matchResults, matchResult)

			fmt.Println(team1.TeamName, team1Goals, "-", team2Goals, team2.TeamName)
		}
	}
	return matchResults
}

func maxArrValue(arr []float64) int {
	max := 0.0
	maxIndex := 0
	for i := range arr {
		if arr[i] > max {
			max = arr[i]
			maxIndex = i
		}
	}
	return maxIndex
}

func apiServer() {
	flag.Parse()
	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}

func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
