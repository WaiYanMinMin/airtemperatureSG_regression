package pipeline

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func CountCSVDataRows(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	rows, err := r.ReadAll()
	if err != nil {
		return 0, err
	}
	if len(rows) == 0 {
		return 0, fmt.Errorf("csv %s is empty", path)
	}
	if len(rows) == 1 {
		return 0, nil
	}
	return len(rows) - 1, nil
}

func LoadTemperatureCSV(path string) ([]TemperatureRow, error) {
	records, err := readCSV(path)
	if err != nil {
		return nil, err
	}

	out := make([]TemperatureRow, 0, len(records))
	for _, r := range records {
		row, err := parseTemperature(r)
		if err != nil {
			continue
		}
		out = append(out, row)
	}
	return out, nil
}

func LoadDemographicsCSV(path string) ([]DemographicRow, error) {
	records, err := readCSV(path)
	if err != nil {
		return nil, err
	}

	out := make([]DemographicRow, 0, len(records))
	for _, r := range records {
		row, err := parseDemographics(r)
		if err != nil {
			continue
		}
		out = append(out, row)
	}
	return out, nil
}

func LoadHealthCSV(path string) ([]HealthRow, error) {
	records, err := readCSV(path)
	if err != nil {
		return nil, err
	}

	out := make([]HealthRow, 0, len(records))
	for _, r := range records {
		row, err := parseHealth(r)
		if err != nil {
			continue
		}
		out = append(out, row)
	}
	return out, nil
}

func WriteMergedCSV(path string, rows []MergedRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	header := []string{
		"city", "district", "year", "lat", "lon", "heat_cases_per_100k",
		"temperature_c", "population", "elderly_share", "income_index", "distance_km_to_temp",
	}
	if err := w.Write(header); err != nil {
		return err
	}

	for _, r := range rows {
		record := []string{
			r.City, r.District, strconv.Itoa(r.Year),
			fmtFloat(r.Latitude), fmtFloat(r.Longitude), fmtFloat(r.HeatCasesPer100K),
			fmtFloat(r.TemperatureC), fmtFloat(r.Population), fmtFloat(r.ElderlyShare),
			fmtFloat(r.IncomeIndex), fmtFloat(r.DistanceKM),
		}
		if err := w.Write(record); err != nil {
			return err
		}
	}
	return nil
}

func WritePredictionsCSV(path string, rows []PredictionRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	header := []string{
		"city", "district", "year", "heat_cases_per_100k", "predicted_heat_cases_per_100k", "residual",
	}
	if err := w.Write(header); err != nil {
		return err
	}

	for _, r := range rows {
		record := []string{
			r.City, r.District, strconv.Itoa(r.Year),
			fmtFloat(r.HeatCasesPer100K), fmtFloat(r.PredictedHeatCasesPer100K), fmtFloat(r.Residual),
		}
		if err := w.Write(record); err != nil {
			return err
		}
	}
	return nil
}

func readCSV(path string) ([]map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.TrimLeadingSpace = true

	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) < 2 {
		return nil, fmt.Errorf("csv %s has no data rows", path)
	}

	header := make([]string, len(rows[0]))
	for i, h := range rows[0] {
		header[i] = normalizeKey(h)
	}

	out := make([]map[string]string, 0, len(rows)-1)
	for _, row := range rows[1:] {
		record := make(map[string]string, len(header))
		for i, key := range header {
			if i < len(row) {
				record[key] = strings.TrimSpace(row[i])
			}
		}
		out = append(out, record)
	}
	return out, nil
}

func parseTemperature(r map[string]string) (TemperatureRow, error) {
	year, err := parseIntField(r, "year")
	if err != nil {
		return TemperatureRow{}, err
	}
	month, err := parseOptionalIntField(r, "month", 1)
	if err != nil {
		return TemperatureRow{}, err
	}
	lat, err := parseFloatField(r, "lat")
	if err != nil {
		return TemperatureRow{}, err
	}
	lon, err := parseFloatField(r, "lon")
	if err != nil {
		return TemperatureRow{}, err
	}
	temp, err := parseFloatField(r, "temperature_c")
	if err != nil {
		return TemperatureRow{}, err
	}

	return TemperatureRow{
		City:      r["city"],
		District:  r["district"],
		Year:      year,
		Month:     month,
		Latitude:  lat,
		Longitude: lon,
		ValueC:    temp,
	}, nil
}

func parseDemographics(r map[string]string) (DemographicRow, error) {
	year, err := parseIntField(r, "year")
	if err != nil {
		return DemographicRow{}, err
	}
	pop, err := parseFloatField(r, "population")
	if err != nil {
		return DemographicRow{}, err
	}
	elderlyShare, err := parseFloatField(r, "elderly_share")
	if err != nil {
		return DemographicRow{}, err
	}
	incomeIndex, err := parseOptionalFloatField(r, "income_index", 0)
	if err != nil {
		return DemographicRow{}, err
	}

	return DemographicRow{
		City:         r["city"],
		District:     r["district"],
		Year:         year,
		Population:   pop,
		ElderlyShare: elderlyShare,
		IncomeIndex:  incomeIndex,
	}, nil
}

func parseHealth(r map[string]string) (HealthRow, error) {
	year, err := parseIntField(r, "year")
	if err != nil {
		return HealthRow{}, err
	}
	lat, err := parseFloatField(r, "lat")
	if err != nil {
		return HealthRow{}, err
	}
	lon, err := parseFloatField(r, "lon")
	if err != nil {
		return HealthRow{}, err
	}
	heatCases, err := parseFloatField(r, "heat_cases_per_100k")
	if err != nil {
		return HealthRow{}, err
	}

	return HealthRow{
		City:             r["city"],
		District:         r["district"],
		Year:             year,
		Latitude:         lat,
		Longitude:        lon,
		HeatCasesPer100K: heatCases,
	}, nil
}

func parseIntField(r map[string]string, key string) (int, error) {
	value := strings.TrimSpace(r[normalizeKey(key)])
	if value == "" {
		return 0, fmt.Errorf("missing %s", key)
	}
	v, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func parseOptionalIntField(r map[string]string, key string, defaultValue int) (int, error) {
	value := strings.TrimSpace(r[normalizeKey(key)])
	if value == "" {
		return defaultValue, nil
	}
	v, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func parseFloatField(r map[string]string, key string) (float64, error) {
	value := strings.TrimSpace(r[normalizeKey(key)])
	if value == "" {
		return 0, fmt.Errorf("missing %s", key)
	}
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func parseOptionalFloatField(r map[string]string, key string, defaultValue float64) (float64, error) {
	value := strings.TrimSpace(r[normalizeKey(key)])
	if value == "" {
		return defaultValue, nil
	}
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func normalizeKey(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func fmtFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', 6, 64)
}
