package models

type League struct {
	LeagueId  int    `json:"league_id"`
	CountryId int    `json:"country_id"`
	Title     string `json:"title"`
}
