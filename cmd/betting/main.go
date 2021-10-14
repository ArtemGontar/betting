package main

import (
	"flag"
	"log"

	"github.com/ArtemGontar/betting/internal/app/apiserver"
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
	flag.Parse()
	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}

	// rawData := authHttpGetRequest("https://api.betting-api.com/1xbet/football/line/leagues",
	// 	"ef005110a34b481d9c1931075d779c71ed33058d80834b48b852b469bb5e7742")
	// leagues := byteArrToLeagueMapping(rawData)
	// insertLeagues(db, leagues)

	// rawData := authHttpGetRequest("https://api.betting-api.com/1xbet/football/line/league/2286681/matches",
	// 	"ef005110a34b481d9c1931075d779c71ed33058d80834b48b852b469bb5e7742")
	// fmt.Println(string(rawData))

	//matchesResults := reader.ReadMatchResultsFromDataset("dataset/E0_2021_2.csv")
	//store.InsertMatchResults(db, matchesResults)

	// League avg statistic
	// league := 1
	// leagueAvgStat := service.LeagueAvgStatistics(store, league)
	// fmt.Println(league, "avg home scored goals =", avgHomeScoredGoals)
	// fmt.Println(league, "avg away scored goals =", avgAwayScoredGoals)
	// fmt.Println(league, "avg home conceded goals =", avgHomeConcededGoals)
	// fmt.Println(league, "avg away conceded goals =", avgAwayConcededGoals)

	// homeTeamStat := models.TeamStatistic{
	// 	TeamName: "Man United",
	// 	IsHome:   true,
	// }
	// awayTeamStat := models.TeamStatistic{
	// 	TeamName: "Leicester",
	// 	IsHome:   false,
	// }

	// //calculate team stats
	// service.TeamStatistics(store, &homeTeamStat, leagueAvgStat.AvgHomeScoredGoals, leagueAvgStat.AvgHomeConcededGoals)
	// service.TeamStatistics(store, &awayTeamStat, leagueAvgStat.AvgAwayScoredGoals, leagueAvgStat.AvgAwayConcededGoals)

	// //calculate poisson distribution
	// service.PoissonDistribution(&homeTeamStat, awayTeamStat, leagueAvgStat.AvgHomeScoredGoals)
	// service.PoissonDistribution(&awayTeamStat, homeTeamStat, leagueAvgStat.AvgAwayScoredGoals)

	// //print stats
	// fmt.Println("Home team stats")
	// PrintTeamStats(homeTeamStat)
	// fmt.Println("Away team stats")
	// PrintTeamStats(awayTeamStat)

	// // matches against each other (last 5)
	// eachOtherGames, err := store.MatchResult().SelectAgainstEachOtherResults(homeTeamStat.TeamName, awayTeamStat.TeamName)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(eachOtherGames)

	// Для расчета мат. ожидания есть даже формула:
	//(Вероятность выигрыша) х (сумму потенциального выигрыша по текущему пари) – (вероятность проигрыша) х (сумму потенциального проигрыша по текущему пари).
}

// func PrintTeamStats(teamStat models.TeamStatistic) {
// 	fmt.Println(teamStat.TeamName, "avg scored goals =", teamStat.AvgScoredGoals)
// 	fmt.Println(teamStat.TeamName, "avg conceded goals =", teamStat.AvgConcededGoals)
// 	fmt.Println(teamStat.TeamName, "last 5 games results = ", teamStat.FullTimeResults)
// 	fmt.Println(teamStat.TeamName, "attack power =", teamStat.AttackPower)
// 	fmt.Println(teamStat.TeamName, "defence power =", teamStat.DefencePower)
// 	fmt.Println(teamStat.TeamName, "predict score =", teamStat.PredictScore)
// 	fmt.Println("Goals", "possibility", teamStat.PoissonDistribution)
// }
