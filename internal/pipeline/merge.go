package pipeline

import (
	"fmt"
	"math"
	"strings"
)

const EarthRadiusKM = 6371.0

func BuildMergedDataset(
	healthRows []HealthRow,
	tempRows []TemperatureRow,
	demoRows []DemographicRow,
	maxDistanceKM float64,
) ([]MergedRow, error) {
	if len(healthRows) == 0 || len(tempRows) == 0 || len(demoRows) == 0 {
		return nil, fmt.Errorf("health, temperature and demographics datasets must all be non-empty")
	}

	demoIndex := make(map[string]DemographicRow, len(demoRows))
	for _, d := range demoRows {
		key := districtYearKey(d.City, d.District, d.Year)
		demoIndex[key] = d
	}

	tempsByCityYear := make(map[string][]TemperatureRow)
	for _, t := range tempRows {
		key := cityYearKey(t.City, t.Year)
		tempsByCityYear[key] = append(tempsByCityYear[key], t)
	}

	merged := make([]MergedRow, 0, len(healthRows))
	for _, h := range healthRows {
		demoKey := districtYearKey(h.City, h.District, h.Year)
		demo, ok := demoIndex[demoKey]
		if !ok {
			continue
		}

		tempCandidates := tempsByCityYear[cityYearKey(h.City, h.Year)]
		if len(tempCandidates) == 0 {
			continue
		}

		bestTemp, bestDistance, ok := nearestTemperature(h, tempCandidates)
		if !ok {
			continue
		}
		if maxDistanceKM > 0 && bestDistance > maxDistanceKM {
			continue
		}

		merged = append(merged, MergedRow{
			City:             h.City,
			District:         h.District,
			Year:             h.Year,
			Latitude:         h.Latitude,
			Longitude:        h.Longitude,
			HeatCasesPer100K: h.HeatCasesPer100K,
			TemperatureC:     bestTemp.ValueC,
			Population:       demo.Population,
			ElderlyShare:     demo.ElderlyShare,
			IncomeIndex:      demo.IncomeIndex,
			DistanceKM:       bestDistance,
		})
	}

	if len(merged) == 0 {
		return nil, fmt.Errorf("no merged rows found; check city/district/year keys and coordinate coverage")
	}
	return merged, nil
}

func nearestTemperature(healthRow HealthRow, candidates []TemperatureRow) (TemperatureRow, float64, bool) {
	bestDistance := math.MaxFloat64
	var best TemperatureRow
	found := false

	for _, c := range candidates {
		d := HaversineKM(healthRow.Latitude, healthRow.Longitude, c.Latitude, c.Longitude)
		if d < bestDistance {
			bestDistance = d
			best = c
			found = true
		}
	}
	return best, bestDistance, found
}

func HaversineKM(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := toRadians(lat2 - lat1)
	dLon := toRadians(lon2 - lon1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(toRadians(lat1))*math.Cos(toRadians(lat2))*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return EarthRadiusKM * c
}

func toRadians(v float64) float64 {
	return v * math.Pi / 180
}

func cityYearKey(city string, year int) string {
	return strings.ToLower(strings.TrimSpace(city)) + "|" + fmt.Sprintf("%d", year)
}

func districtYearKey(city, district string, year int) string {
	return strings.ToLower(strings.TrimSpace(city)) + "|" + strings.ToLower(strings.TrimSpace(district)) + "|" + fmt.Sprintf("%d", year)
}
