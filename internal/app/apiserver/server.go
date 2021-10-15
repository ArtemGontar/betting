package apiserver

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ArtemGontar/betting/internal/app/model"
	"github.com/ArtemGontar/betting/internal/app/reader"
	"github.com/ArtemGontar/betting/internal/app/service"
	"github.com/ArtemGontar/betting/internal/app/store"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	ctxKeyUser ctxKey = iota
	ctxKeyRequestID
)

type server struct {
	router *mux.Router
	logger *logrus.Logger
	store  store.Store
}

type ctxKey int8

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
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.HandleFunc("/hello", s.handleHello()).Methods("GET")
	s.router.HandleFunc("/match-results/batch", s.FillFromDatasetHandler()).Methods("POST")
	s.router.HandleFunc("/leagues/{id:[0-9]+}/statistic", s.GetAvgLeagueStatisticHandler()).Methods("GET")
	s.router.HandleFunc("/teams/{name:\\w+}/statistic", s.GetTeamStatisticHandler()).Methods("GET")
	s.router.HandleFunc("/match/statistic", s.GetMatchStatisticHandler()).Methods("GET")
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		next.ServeHTTP(rw, r)
		logger.Infof("completed in %v", time.Now().Sub(start))

	})
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

		s.logger.Info("Fill from dataset with name", req.DatasetName)

		matchesResults, err := reader.ReadMatchResultsFromDataset("datasets/" + req.DatasetName)
		if err != nil {
			s.respond(rw, r, http.StatusCreated, err)
		}
		s.logger.Info("Successfully readed from dataset", req.DatasetName)
		s.store.MatchResult().InsertMatchResults(matchesResults)
		s.logger.Info("Successfully insert match results")
		s.respond(rw, r, http.StatusCreated, matchesResults)
	}
}

func (s *server) GetAvgLeagueStatisticHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idString := vars["id"]
		id, e := strconv.Atoi(idString)
		if e != nil {
			s.logger.Error(e)
		}
		s.logger.Info("GetAvgLeagueStatistic for league", id)
		resp := service.LeagueStatistics(s.store, id)
		s.respond(rw, r, http.StatusCreated, resp)
	}
}

func (s *server) GetTeamStatisticHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		s.logger.Info("GetTeamStatisticHandler for team", name)
		teamStat := model.TeamStatistic{
			TeamName: name,
		}
		league := 1
		leagueStatistic := service.LeagueStatistics(s.store, league)
		service.TeamStatistics(s.store, &teamStat, leagueStatistic.AvgHomeScoredGoals, leagueStatistic.AvgHomeConcededGoals)
		s.respond(rw, r, http.StatusCreated, teamStat)
	}
}

func (s *server) GetMatchStatisticHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		league := 1
		homeTeam := r.URL.Query().Get("homeTeam")
		awayTeam := r.URL.Query().Get("awayTeam")
		s.logger.Info("GetMatchStatisticHandler for %s - %s", homeTeam, awayTeam)
		homeTeamStat := model.TeamStatistic{
			TeamName: homeTeam,
			IsHome:   true,
		}
		awayTeamStat := model.TeamStatistic{
			TeamName: awayTeam,
			IsHome:   false,
		}

		leagueStat := service.LeagueStatistics(s.store, league)
		service.TeamStatistics(s.store, &homeTeamStat, leagueStat.AvgHomeScoredGoals, leagueStat.AvgHomeConcededGoals)
		service.TeamStatistics(s.store, &awayTeamStat, leagueStat.AvgAwayScoredGoals, leagueStat.AvgAwayConcededGoals)
		service.PoissonDistribution(&homeTeamStat, awayTeamStat, leagueStat.AvgHomeScoredGoals)
		service.PoissonDistribution(&awayTeamStat, homeTeamStat, leagueStat.AvgAwayScoredGoals)
		// matches against each other (last 5)
		eachOtherResult := service.AgainstEachOtherResults(s.store, &homeTeamStat, &awayTeamStat)
		ms := model.MatchStatistic{
			HomeTeamName:            homeTeam,
			AwayTeamName:            awayTeam,
			HomePredictScore:        homeTeamStat.PredictScore,
			AwayPredictScore:        awayTeamStat.PredictScore,
			HomePoissonDistribution: homeTeamStat.PoissonDistribution,
			AwayPoissonDistribution: awayTeamStat.PoissonDistribution,
			AgainstEachOtherResult:  eachOtherResult,
		}
		s.respond(rw, r, http.StatusCreated, ms)
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
