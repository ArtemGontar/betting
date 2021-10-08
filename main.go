package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"encoding/csv"

	"github.com/ArtemGontar/betting/internal/app/models"
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

	matchesResults := readMatchResultsFromDataset("dataset/E0.csv")
	store.InsertMatchResults(db, matchesResults)
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

func arrToMatchResultsMapping(record []string) models.MatchResult {
	fullTimeHomeTeamGoals, _ := strconv.Atoi(record[4])
	fullTimeAwayTeamGoals, _ := strconv.Atoi(record[5])

	halfTimeHomeTeamGoals, _ := strconv.Atoi(record[7])
	halfTimeAwayTeamGoals, _ := strconv.Atoi(record[8])

	homeTeamShots, _ := strconv.Atoi(record[11])
	awayTeamShots, _ := strconv.Atoi(record[12])

	homeTeamShotsOnTarget, _ := strconv.Atoi(record[13])
	awayTeamShotsOnTarget, _ := strconv.Atoi(record[14])

	homeTeamFoulsCommitted, _ := strconv.Atoi(record[15])
	awayTeamFoulsCommitted, _ := strconv.Atoi(record[16])

	homeTeamCorners, _ := strconv.Atoi(record[17])
	awayTeamCorners, _ := strconv.Atoi(record[18])

	homeTeamYellowCards, _ := strconv.Atoi(record[19])
	awayTeamYellowCards, _ := strconv.Atoi(record[20])

	homeTeamRedCards, _ := strconv.Atoi(record[21])
	awayTeamRedCards, _ := strconv.Atoi(record[22])

	dateStart, _ := time.Parse(time.RFC3339, record[1])
	matches := models.MatchResult{
		HomeTeam:               record[2],
		HomeTeamId:             0,
		FullTimeHomeTeamGoals:  fullTimeHomeTeamGoals,
		HalfTimeHomeTeamGoals:  halfTimeHomeTeamGoals,
		HomeTeamShots:          homeTeamShots,
		HomeTeamShotsOnTarget:  homeTeamShotsOnTarget,
		HomeTeamCorners:        homeTeamCorners,
		HomeTeamFoulsCommitted: homeTeamFoulsCommitted,
		HomeTeamYellowCards:    homeTeamYellowCards,
		HomeTeamRedCards:       homeTeamRedCards,
		AwayTeam:               record[3],
		AwayTeamId:             0,
		FullTimeAwayTeamGoals:  fullTimeAwayTeamGoals,
		HalfTimeAwayTeamGoals:  halfTimeAwayTeamGoals,
		AwayTeamShots:          awayTeamShots,
		AwayTeamShotsOnTarget:  awayTeamShotsOnTarget,
		AwayTeamCorners:        awayTeamCorners,
		AwayTeamFoulsCommitted: awayTeamFoulsCommitted,
		AwayTeamYellowCards:    awayTeamYellowCards,
		AwayTeamRedCards:       awayTeamRedCards,
		FullTimeResult:         record[6],
		HalfTimeResult:         record[9],
		DateStart:              dateStart,
		LeagueId:               1,
	}

	return matches
}

func readMatchResultsFromDataset(filePath string) []models.MatchResult {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	matchesResults := []models.MatchResult{}
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 23
	for {
		record, e := reader.Read()
		if e != nil {
			fmt.Println(e)
			break
		}
		matchesResults = append(matchesResults, arrToMatchResultsMapping(record))
	}

	return matchesResults
}
