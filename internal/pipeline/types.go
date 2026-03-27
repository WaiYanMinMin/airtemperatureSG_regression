package pipeline

import "time"

type TemperatureRow struct {
	City      string
	District  string
	Year      int
	Month     int
	Latitude  float64
	Longitude float64
	ValueC    float64
}

type DemographicRow struct {
	City         string
	District     string
	Year         int
	Population   float64
	ElderlyShare float64
	IncomeIndex  float64
}

type HealthRow struct {
	City             string
	District         string
	Year             int
	Latitude         float64
	Longitude        float64
	HeatCasesPer100K float64
}

type MergedRow struct {
	City             string
	District         string
	Year             int
	Latitude         float64
	Longitude        float64
	HeatCasesPer100K float64
	TemperatureC     float64
	Population       float64
	ElderlyShare     float64
	IncomeIndex      float64
	DistanceKM       float64
}

type PredictionRow struct {
	MergedRow
	PredictedHeatCasesPer100K float64
	Residual                  float64
}

type ModelResult struct {
	Coefficients map[string]float64
	RMSE         float64
	MAE          float64
	R2           float64
	TrainedAt    time.Time
}
