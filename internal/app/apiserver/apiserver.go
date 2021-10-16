package apiserver

import (
	"database/sql"
	"net/http"

	league_service "github.com/ArtemGontar/betting/internal/app/service/league"
	statistic_service "github.com/ArtemGontar/betting/internal/app/service/statistic"
	match_results_repository "github.com/ArtemGontar/betting/internal/app/store/match_results"
)

func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()
	matchResultsRepository := match_results_repository.New(db)
	statisticService := statistic_service.New(matchResultsRepository)
	leagueService := league_service.New(matchResultsRepository)
	srv := newServer(leagueService, statisticService)
	return http.ListenAndServe(config.BindAddr, srv)
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
