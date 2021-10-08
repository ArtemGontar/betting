package models

import "time"

type MatchResults struct {
	Id                     int
	HomeTeam               string    `json:"home_team"`
	HomeTeamId             int       `json:"home_team_id"`
	FullTimeHomeTeamGoals  int       `json:"full_time_home_team_goals"`
	HalfTimeHomeTeamGoals  int       `json:"half_time_home_team_goals"`
	HomeTeamShots          int       `json:"home_team_shots"`
	HomeTeamShotsOnTarget  int       `json:"home_team_shots_on_target"`
	HomeTeamCorners        int       `json:"home_team_corners"`
	HomeTeamFoulsCommitted int       `json:"home_team_fouls_committed"`
	HomeTeamYellowCards    int       `json:"home_team_yellow_cards"`
	HomeTeamRedCards       int       `json:"home_team_red_cards"`
	AwayTeam               string    `json:"away_team"`
	AwayTeamId             int       `json:"away_team_id"`
	FullTimeAwayTeamGoals  int       `json:"full_time_away_team_goals"`
	HalfTimeAwayTeamGoals  int       `json:"half_time_away_team_goals"`
	AwayTeamShots          int       `json:"away_team_shots"`
	AwayTeamShotsOnTarget  int       `json:"away_team_shots_on_target"`
	AwayTeamCorners        int       `json:"away_team_corners"`
	AwayTeamFoulsCommitted int       `json:"away_team_fouls_committed"`
	AwayTeamYellowCards    int       `json:"away_team_yellow_cards"`
	AwayTeamRedCards       int       `json:"away_team_red_cards"`
	FullTimeResult         string    `json:"full_time_result"`
	HalfTimeResult         string    `json:"half_time_result"`
	DateStart              time.Time `json:"date_start"`
	LeagueId               int       `json:"league_id"`
}
