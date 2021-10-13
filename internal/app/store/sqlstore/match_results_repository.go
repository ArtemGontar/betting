package sqlstore

import (
	"fmt"

	"github.com/ArtemGontar/betting/internal/app/models"
	"github.com/ArtemGontar/betting/internal/app/store"
)

type MatchResultRepository struct {
	store *Store
}

func (r *MatchResultRepository) SelectLeagueAvgScoredGoals(league int) (float64, float64, error) {
	var avgFullTimeHomeLeagueGoals float64
	var avgFullTimeAwayLeagueGoals float64
	if err := r.store.db.QueryRow(`SELECT avg(full_time_home_team_goals), avg(full_time_away_team_goals)
		FROM public.match_results WHERE league_id = $1`, league).Scan(
		&avgFullTimeHomeLeagueGoals,
		&avgFullTimeAwayLeagueGoals,
	); err != nil {
		return 0, 0, err
	}

	return avgFullTimeHomeLeagueGoals, avgFullTimeAwayLeagueGoals, nil
}

func (r *MatchResultRepository) SelectHomeTeamAvgGoals(team string) (float64, float64, error) {
	var avgFullTimeHomeTeamGoals float64
	var avgFullTimeAwayTeamGoals float64
	if err := r.store.db.QueryRow(`SELECT avg(full_time_home_team_goals), avg(full_time_away_team_goals)
		FROM public.match_results WHERE home_team = $1`, team).Scan(
		&avgFullTimeHomeTeamGoals,
		&avgFullTimeAwayTeamGoals,
	); err != nil {
		return 0, 0, err
	}

	return avgFullTimeHomeTeamGoals, avgFullTimeAwayTeamGoals, nil
}

func (r *MatchResultRepository) SelectAwayTeamAvgGoals(team string) (float64, float64, error) {
	var avgFullTimeHomeTeamGoals float64
	var avgFullTimeAwayTeamGoals float64
	if err := r.store.db.QueryRow(`SELECT avg(full_time_away_team_goals), avg(full_time_home_team_goals)
		FROM public.match_results WHERE away_team = $1`, team).Scan(
		&avgFullTimeAwayTeamGoals,
		&avgFullTimeHomeTeamGoals,
	); err != nil {
		return 0, 0, err
	}

	return avgFullTimeAwayTeamGoals, avgFullTimeHomeTeamGoals, nil
}

func (r *MatchResultRepository) SelectLastFiveGamesByTeam(team string) ([]store.Result, error) {
	fullTimeResults := make([]store.Result, 0)
	rows, err := r.store.db.Query(`SELECT home_team, away_team, full_time_result
		FROM public.match_results
		WHERE home_team = $1 OR away_team = $1
		ORDER BY date_start desc
		LIMIT 5`, team)
	if err != nil {
		panic(err)
	}
	var result string
	var homeTeam string
	var awayTeam string
	for rows.Next() {
		rows.Scan(&homeTeam, &awayTeam, &result)
		fullTimeResults = append(fullTimeResults, store.Result{HomeTeam: homeTeam, AwayTeam: awayTeam, Result: result})
	}

	return fullTimeResults, nil
}

func (r *MatchResultRepository) SelectAgainstEachOtherResults(team1 string, team2 string) ([]store.Result, error) {
	fullTimeResults := make([]store.Result, 0)
	rows, err := r.store.db.Query(`SELECT home_team, away_team, full_time_home_team_goals, full_time_away_team_goals, full_time_result
	FROM public.match_results
	WHERE (home_team = $1 AND away_team = $2) OR (home_team = $2 AND away_team = $1)
	LIMIT 5`, team1, team2)
	if err != nil {
		panic(err)
	}
	var result string
	var homeTeam string
	var awayTeam string
	var homeGoals int
	var awayGoals int
	for rows.Next() {
		rows.Scan(&homeTeam, &awayTeam, &homeGoals, &awayGoals, &result)
		fullTimeResults = append(fullTimeResults, store.Result{
			HomeTeam:              homeTeam,
			AwayTeam:              awayTeam,
			Result:                result,
			FullTimeHomeTeamGoals: homeGoals,
			FullTimeAwayTeamGoals: awayGoals,
		})
	}

	return fullTimeResults, nil
}

func (r *MatchResultRepository) InsertMatchResults(matchResults []models.MatchResult) {
	for _, matchResult := range matchResults {
		fmt.Println(matchResult)
		_, err := r.store.db.Exec(`INSERT INTO public.match_results (home_team, home_team_id, full_time_home_team_goals, 
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
