package pipeline

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

type DatasetQuality struct {
	Name        string
	RawRows     int
	ParsedRows  int
	DroppedRows int
}

type NumericStats struct {
	Count int
	Min   float64
	Max   float64
	Mean  float64
}

func WriteDataExplorationReport(
	path string,
	temperaturePath, demographicsPath, healthPath string,
	tempRows []TemperatureRow,
	demoRows []DemographicRow,
	healthRows []HealthRow,
) error {
	tempRaw, err := CountCSVDataRows(temperaturePath)
	if err != nil {
		return fmt.Errorf("count temperature rows: %w", err)
	}
	demoRaw, err := CountCSVDataRows(demographicsPath)
	if err != nil {
		return fmt.Errorf("count demographics rows: %w", err)
	}
	healthRaw, err := CountCSVDataRows(healthPath)
	if err != nil {
		return fmt.Errorf("count health rows: %w", err)
	}

	quality := []DatasetQuality{
		{
			Name:        "Temperature",
			RawRows:     tempRaw,
			ParsedRows:  len(tempRows),
			DroppedRows: tempRaw - len(tempRows),
		},
		{
			Name:        "Demographics",
			RawRows:     demoRaw,
			ParsedRows:  len(demoRows),
			DroppedRows: demoRaw - len(demoRows),
		},
		{
			Name:        "Health",
			RawRows:     healthRaw,
			ParsedRows:  len(healthRows),
			DroppedRows: healthRaw - len(healthRows),
		},
	}

	var b strings.Builder
	b.WriteString("Data Exploration Report\n")
	b.WriteString("=======================\n\n")

	b.WriteString("1) Data Quality (Raw vs Parsed)\n")
	b.WriteString("-------------------------------\n")
	b.WriteString("dataset,raw_rows,parsed_rows,dropped_rows\n")
	for _, q := range quality {
		b.WriteString(fmt.Sprintf("%s,%d,%d,%d\n", q.Name, q.RawRows, q.ParsedRows, q.DroppedRows))
	}
	b.WriteString("\n")

	b.WriteString("2) Coverage\n")
	b.WriteString("-----------\n")
	appendCountSection(&b, "Temperature by City", countsFromTemperature(tempRows, func(r TemperatureRow) string { return r.City }))
	appendCountSection(&b, "Demographics by City", countsFromDemographics(demoRows, func(r DemographicRow) string { return r.City }))
	appendCountSection(&b, "Health by City", countsFromHealth(healthRows, func(r HealthRow) string { return r.City }))

	appendCountSection(&b, "Temperature by Year", countsFromTemperature(tempRows, func(r TemperatureRow) string { return fmt.Sprintf("%d", r.Year) }))
	appendCountSection(&b, "Demographics by Year", countsFromDemographics(demoRows, func(r DemographicRow) string { return fmt.Sprintf("%d", r.Year) }))
	appendCountSection(&b, "Health by Year", countsFromHealth(healthRows, func(r HealthRow) string { return fmt.Sprintf("%d", r.Year) }))

	appendCountSection(&b, "Health by District", countsFromHealth(healthRows, func(r HealthRow) string { return r.City + "::" + r.District }))
	b.WriteString("\n")

	b.WriteString("3) Numeric Feature Stats\n")
	b.WriteString("------------------------\n")
	writeStatsLine(&b, "temperature_c", statsFromTemperature(tempRows))
	writeStatsLine(&b, "population", statsFromDemographics(demoRows, func(r DemographicRow) float64 { return r.Population }))
	writeStatsLine(&b, "elderly_share", statsFromDemographics(demoRows, func(r DemographicRow) float64 { return r.ElderlyShare }))
	writeStatsLine(&b, "income_index", statsFromDemographics(demoRows, func(r DemographicRow) float64 { return r.IncomeIndex }))
	writeStatsLine(&b, "heat_cases_per_100k", statsFromHealth(healthRows, func(r HealthRow) float64 { return r.HeatCasesPer100K }))
	b.WriteString("\n")

	b.WriteString("Notes:\n")
	b.WriteString("- dropped_rows = raw_rows - parsed_rows\n")
	b.WriteString("- parsed rows are rows that passed required field parsing\n")
	b.WriteString("- this exploration step supports assignment task #1 before modeling\n")

	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func appendCountSection(b *strings.Builder, title string, counts map[string]int) {
	b.WriteString(title + "\n")
	b.WriteString(strings.Repeat("-", len(title)) + "\n")
	b.WriteString("group,count\n")
	keys := make([]string, 0, len(counts))
	for k := range counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		b.WriteString(fmt.Sprintf("%s,%d\n", k, counts[k]))
	}
	b.WriteString("\n")
}

func countsFromTemperature(rows []TemperatureRow, keyFn func(TemperatureRow) string) map[string]int {
	out := map[string]int{}
	for _, r := range rows {
		out[keyFn(r)]++
	}
	return out
}

func countsFromDemographics(rows []DemographicRow, keyFn func(DemographicRow) string) map[string]int {
	out := map[string]int{}
	for _, r := range rows {
		out[keyFn(r)]++
	}
	return out
}

func countsFromHealth(rows []HealthRow, keyFn func(HealthRow) string) map[string]int {
	out := map[string]int{}
	for _, r := range rows {
		out[keyFn(r)]++
	}
	return out
}

func statsFromTemperature(rows []TemperatureRow) NumericStats {
	values := make([]float64, 0, len(rows))
	for _, r := range rows {
		values = append(values, r.ValueC)
	}
	return stats(values)
}

func statsFromDemographics(rows []DemographicRow, valueFn func(DemographicRow) float64) NumericStats {
	values := make([]float64, 0, len(rows))
	for _, r := range rows {
		values = append(values, valueFn(r))
	}
	return stats(values)
}

func statsFromHealth(rows []HealthRow, valueFn func(HealthRow) float64) NumericStats {
	values := make([]float64, 0, len(rows))
	for _, r := range rows {
		values = append(values, valueFn(r))
	}
	return stats(values)
}

func stats(values []float64) NumericStats {
	if len(values) == 0 {
		return NumericStats{}
	}
	minV := math.MaxFloat64
	maxV := -math.MaxFloat64
	sum := 0.0
	for _, v := range values {
		if v < minV {
			minV = v
		}
		if v > maxV {
			maxV = v
		}
		sum += v
	}
	return NumericStats{
		Count: len(values),
		Min:   minV,
		Max:   maxV,
		Mean:  sum / float64(len(values)),
	}
}

func writeStatsLine(b *strings.Builder, name string, s NumericStats) {
	b.WriteString(fmt.Sprintf(
		"%s,count=%d,min=%.6f,max=%.6f,mean=%.6f\n",
		name, s.Count, s.Min, s.Max, s.Mean,
	))
}
