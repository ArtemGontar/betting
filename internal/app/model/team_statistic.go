package model

type TeamStatistic struct {
	TeamName            string
	IsHome              bool
	AttackPower         float64
	DefencePower        float64
	AvgScoredGoals      float64
	AvgConcededGoals    float64
	PredictScore        float64
	FullTimeResults     []string
	PoissonDistribution []float64
}
