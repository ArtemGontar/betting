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
