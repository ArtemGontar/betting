package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ArtemGontar/betting/internal/app/model"
	"github.com/ArtemGontar/betting/internal/app/reader"
	"github.com/ArtemGontar/betting/internal/app/service"
	"github.com/ArtemGontar/betting/internal/app/store"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type server struct {
	router *mux.Router
	logger *logrus.Logger
	store  store.Store
}

func newServer(store store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/hello", s.handleHello()).Methods("GET")
	s.router.HandleFunc("/leagues/{id:[0-9]+}/stat", s.GetAvgLeagueStatisticHandler()).Methods("GET")
	s.router.HandleFunc("/teams/{name:\\w+}/stat", s.GetTeamStatisticHandler()).Methods("GET")
	s.router.HandleFunc("/match-results/batch", s.FillFromDatasetHandler()).Methods("POST")
}

func (s *server) handleHello() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		s.respond(rw, r, http.StatusCreated, "hello")
	}
}

func (s *server) FillFromDatasetHandler() http.HandlerFunc {
	type request struct {
		DatasetName string `json:"datasetName"`
	}
	return func(rw http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(rw, r, http.StatusBadRequest, err)
			return
		}

		matchesResults, err := reader.ReadMatchResultsFromDataset("datasets/" + req.DatasetName)
		if err != nil {
			s.respond(rw, r, http.StatusCreated, err)
		}
		s.store.MatchResult().InsertMatchResults(matchesResults)
		s.respond(rw, r, http.StatusCreated, "Success")
	}
}

func (s *server) GetAvgLeagueStatisticHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idString := vars["id"]
		id, e := strconv.Atoi(idString)
		if e == nil {
			fmt.Printf("%T \n %v", id, id)
		}
		resp := service.LeagueAvgStatistics(s.store, id)
		s.respond(rw, r, http.StatusCreated, resp)
	}
}

func (s *server) GetTeamStatisticHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		teamStat := model.TeamStatistic{
			TeamName: name,
		}
		league := 1
		leagueStatistic := service.LeagueAvgStatistics(s.store, league)
		service.TeamStatistics(s.store, &teamStat, leagueStatistic.AvgHomeScoredGoals, leagueStatistic.AvgHomeConcededGoals)
		s.respond(rw, r, http.StatusCreated, teamStat)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
