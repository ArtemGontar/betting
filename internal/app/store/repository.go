package store

import (
	"database/sql"
	"fmt"

	"github.com/ArtemGontar/betting/internal/app/models"
)

func InsertLeagues(db *sql.DB, leagues []models.League) {
	for _, league := range leagues {
		fmt.Println(league)
		db.Exec(
			"INSERT INTO public.leagues (league_id, country_id, title) VALUES ($1, $2, $3)",
			league.LeagueId,
			league.CountryId,
			league.Title,
		)
	}
}

func InsertMatchResults(db *sql.DB, matchResults []models.MatchResults) {
	for _, matchResult := range matchResults {
		fmt.Println(matchResult)
		_, err := db.Exec(`INSERT INTO public.match_results (home_team, home_team_id, full_time_home_team_goals, 
				half_time_home_team_goals, home_team_shots, home_team_shots_on_target, home_team_corners,
				home_team_fouls_committed, home_team_yellow_cards, home_team_red_cards, away_team, away_team_id,
				full_time_away_team_goals, half_time_away_team_goals, away_team_shots, away_team_shots_on_target, 
				away_team_corners, away_team_fouls_committed, away_team_yellow_cards, away_team_red_cards,
				 full_time_result, half_time_result, date_start, league_id) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)`,
			matchResult.HomeTeam,
			matchResult.HomeTeamId,
			matchResult.FullTimeHomeTeamGoals,
			matchResult.HalfTimeHomeTeamGoals,
			matchResult.HomeTeamShots,
			matchResult.HomeTeamShotsOnTarget,
			matchResult.HomeTeamCorners,
			matchResult.HomeTeamFoulsCommitted,
			matchResult.HomeTeamYellowCards,
			matchResult.HomeTeamRedCards,
			matchResult.AwayTeam,
			matchResult.AwayTeamId,
			matchResult.FullTimeAwayTeamGoals,
			matchResult.HalfTimeAwayTeamGoals,
			matchResult.AwayTeamShots,
			matchResult.AwayTeamShotsOnTarget,
			matchResult.AwayTeamCorners,
			matchResult.AwayTeamFoulsCommitted,
			matchResult.AwayTeamYellowCards,
			matchResult.AwayTeamRedCards,
			matchResult.FullTimeResult,
			matchResult.HalfTimeResult,
			matchResult.DateStart,
			matchResult.LeagueId,
		)

		if err != nil {
			fmt.Println(err)
		}
	}
}
