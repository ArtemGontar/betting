package model

type MatchStatistic struct {
	HomeTeamName            string
	HomePredictScore        float64
	HomePoissonDistribution []float64
	AwayTeamName            string
	AwayPredictScore        float64
	AwayPoissonDistribution []float64
	AgainstEachOtherResult  []Result
}
