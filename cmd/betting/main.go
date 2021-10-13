package main

import (
	"fmt"
	"io"
	"math"
	"net/http"

	"github.com/ArtemGontar/betting/internal/app/store"
	"github.com/ArtemGontar/betting/internal/app/store/sqlstore"
	_ "github.com/lib/pq"
)

var e = 2.71828182845904523536028747135266249

func main() {
	store, err := sqlstore.New("host=localhost dbname=betting sslmode=disable user=postgres password=cce16cc03cfb49c9b247fc6faff58fa7")
	if err != nil {
		fmt.Println(err)
		return
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
	league := 1
	avgLeagueHomeScoredGoals, avgLeagueAwayScoredGoals, err := store.MatchResult().SelectLeagueAvgScoredGoals(league)
	if err != nil {
		fmt.Println(err)
	}
	avgLeagueHomeConcededGoals := avgLeagueAwayScoredGoals
	avgLeagueAwayConcededGoals := avgLeagueHomeScoredGoals

	fmt.Println(league, "avg home scored goals =", avgLeagueHomeScoredGoals)
	fmt.Println(league, "avg away scored goals =", avgLeagueAwayScoredGoals)
	fmt.Println(league, "avg home conceded goals =", avgLeagueHomeConcededGoals)
	fmt.Println(league, "avg away conceded goals =", avgLeagueAwayConcededGoals)

	homeTeam := "Man United"
	awayTeam := "Everton"
	avgHomeScoredGoals, avgHomeConcededGoals, err := store.MatchResult().SelectHomeTeamAvgGoals(homeTeam)
	if err != nil {
		fmt.Println(err)
	}

	homeFullTimeResults, err := store.MatchResult().SelectLastFiveGamesByTeam(homeTeam)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Home team stats")
	fmt.Println(homeTeam, "Avg home scored goals =", avgHomeScoredGoals)
	fmt.Println(homeTeam, "Avg home conceded goals =", avgHomeConcededGoals)
	fmt.Print(homeTeam, "Last 5 games results = ")
	ProcessResults(homeFullTimeResults, homeTeam)
	//Poison theory
	homeAttackPower := avgHomeScoredGoals / avgLeagueHomeScoredGoals
	homeDefencePower := avgHomeConcededGoals / avgLeagueHomeConcededGoals
	fmt.Println(homeTeam, "attack power =", homeAttackPower)
	fmt.Println(homeTeam, "defence power =", homeDefencePower)

	avgAwayScoredGoals, avgAwayConcededGoals, _ := store.MatchResult().SelectAwayTeamAvgGoals(awayTeam)
	if err != nil {
		fmt.Println(err)
	}

	awayFullTimeResults, _ := store.MatchResult().SelectLastFiveGamesByTeam(awayTeam)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Away team stats")
	fmt.Println(awayTeam, "Avg home scored goals =", avgAwayScoredGoals)
	fmt.Println(awayTeam, "Avg home conceded goals =", avgAwayConcededGoals)
	fmt.Print(homeTeam, "Last 5 games results = ")
	ProcessResults(awayFullTimeResults, awayTeam)
	//Poison theory
	awayAttackPower := avgAwayScoredGoals / avgLeagueAwayScoredGoals
	awayDefencePower := avgAwayConcededGoals / avgLeagueAwayConcededGoals
	fmt.Println(awayTeam, "attack power =", awayAttackPower)
	fmt.Println(awayTeam, "defence power =", awayDefencePower)

	//вероятное количество голов атакующей
	homePredictScore := homeAttackPower * awayDefencePower * avgHomeScoredGoals
	fmt.Println(homeTeam, "predict score =", homePredictScore)
	awayPredictScore := awayAttackPower * homeDefencePower * avgAwayScoredGoals
	fmt.Println(awayTeam, "predict score =", awayPredictScore)
	//1.213^(5)*e^(-1.213)/(5!)
	homePoisson := []float64{}
	for i := 0; i < 5; i++ {
		homePoisson = append(homePoisson, math.Pow(homePredictScore, (float64(i)))*math.Pow(e, -homePredictScore)/float64(factorial(i)))
		fmt.Println("Goals", i, "possibility", homePoisson[i])
	}

	// matches against each other (last 5)
	eachOtherGames, err := store.MatchResult().SelectAgainstEachOtherResults(homeTeam, awayTeam)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(eachOtherGames)

	// Для расчета мат. ожидания есть даже формула:
	//(Вероятность выигрыша) х (сумму потенциального выигрыша по текущему пари) – (вероятность проигрыша) х (сумму потенциального проигрыша по текущему пари).
}

func factorial(n int) int {

	if n == 0 {
		return 1
	}
	return n * factorial(n-1)
}

func ProcessResults(results []store.Result, team string) {
	defer fmt.Println()
	for _, result := range results {
		if result.Result == "H" {
			fmt.Print("W ")
		} else if result.AwayTeam == team && result.Result == "A" {
			fmt.Print("W ")
		} else if result.HomeTeam == team && result.Result == "A" {
			fmt.Print("L ")
		} else if result.AwayTeam == team && result.Result == "H" {
			fmt.Print("L ")
		} else {
			fmt.Print("D ")
		}
	}
}

func authHttpGetRequest(url string, authToken string) []byte {
	client := http.Client{}
	req, err := http.NewRequest(
		"GET", url, nil,
	)
	req.Header.Add("Authorization", authToken)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return body
}
