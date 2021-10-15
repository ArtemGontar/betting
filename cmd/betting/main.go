package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/ArtemGontar/betting/internal/app/apiserver"
	"github.com/ArtemGontar/betting/internal/app/model"
	"github.com/ArtemGontar/betting/internal/app/reader/premierleague"
	"github.com/ArtemGontar/betting/internal/app/service"
	"github.com/ArtemGontar/betting/internal/app/store/sqlstore"
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
	for i := range teams {
		fmt.Println(teams[i].TeamName)
	}
	//premierleague_reader.Players()

	// flag.Parse()
	// config := apiserver.NewConfig()
	// _, err := toml.DecodeFile(configPath, config)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(config.DatabaseURL)
	// db, err := newDB(config.DatabaseURL)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// defer db.Close()
	// store := sqlstore.New(db)
	// simulationSeason(store, teams)
}

func simulationSeason(store *sqlstore.Store, teams []model.Team) {
	leagueStat := service.LeagueStatistics(store, 1)
	for i := 0; i < len(teams); i++ {
		team1 := model.TeamStatistic{
			TeamName: teams[i].TeamName,
			IsHome:   true,
		}
		service.TeamStatistics(store, &team1, leagueStat.AvgHomeScoredGoals, leagueStat.AvgHomeConcededGoals)
		for j := 0; j < len(teams); j++ {
			if teams[i].TeamName == teams[j].TeamName {
				break
			}
			team2 := model.TeamStatistic{
				TeamName: teams[j].TeamName,
				IsHome:   false,
			}
			service.TeamStatistics(store, &team2, leagueStat.AvgAwayScoredGoals, leagueStat.AvgAwayConcededGoals)

			service.PoissonDistribution(&team1, team2, leagueStat.AvgHomeScoredGoals)
			service.PoissonDistribution(&team2, team1, leagueStat.AvgAwayScoredGoals)

			team1Goals := maxArrValue(team1.PoissonDistribution)
			team2Goals := maxArrValue(team2.PoissonDistribution)
			fmt.Println(team1.TeamName, team1Goals, "-", team2Goals, team2.TeamName)
		}

	}
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
