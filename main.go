package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/ArtemGontar/betting/internal/app/reader"
	"github.com/ArtemGontar/betting/internal/app/store"
	_ "github.com/lib/pq"
)

func main() {
	db, err := store.NewDB("host=localhost dbname=betting sslmode=disable user=postgres password=cce16cc03cfb49c9b247fc6faff58fa7")
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

	matchesResults := reader.ReadMatchResultsFromDataset("dataset/E0_2021_2.csv")
	store.InsertMatchResults(db, matchesResults)

	homeTeam := "Everton"
	awayTeam := "Chelsea"
	avgHomeScoredGoals, avgHomeConcededGoals, err := store.SelectHomeTeamAvgGoals(db, homeTeam)
	if err != nil {
		fmt.Println(err)
	}

	homeFullTimeResults, err := store.SelectLastFiveGamesByTeam(db, homeTeam)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Home team stats")
	fmt.Println(homeTeam, "Avg home scored goals =", avgHomeScoredGoals)
	fmt.Println(homeTeam, "Avg home conceded goals =", avgHomeConcededGoals)
	fmt.Print(homeTeam, "Last 5 games results = ")
	ProcessResults(homeFullTimeResults, homeTeam)

	avgAwayScoredGoals, avgAwayConcededGoals, _ := store.SelectAwayTeamAvgGoals(db, awayTeam)
	if err != nil {
		fmt.Println(err)
	}

	awayFullTimeResults, _ := store.SelectLastFiveGamesByTeam(db, awayTeam)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Away team stats")
	fmt.Println(awayTeam, "Avg home scored goals =", avgAwayScoredGoals)
	fmt.Println(awayTeam, "Avg home conceded goals =", avgAwayConcededGoals)
	fmt.Print(homeTeam, "Last 5 games results = ")
	ProcessResults(awayFullTimeResults, awayTeam)

	eachOtherGames, err := store.SelectAgainstEachOtherResults(db, homeTeam, awayTeam)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(eachOtherGames)
}

func ProcessResults(results []store.Results, team string) {
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
