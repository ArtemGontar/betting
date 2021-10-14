package reader

import (
	"strconv"
	"time"

	"github.com/ArtemGontar/betting/internal/app/model"
)

func arrToMatchResultsMapping(record []string) model.MatchResult {
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
	matches := model.MatchResult{
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
