package pipeline

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

func WriteExplorationPlots(outputDir string, tempRows []TemperatureRow, healthRows []HealthRow) error {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return err
	}

	yearCounts := map[string]float64{}
	for _, r := range tempRows {
		yearCounts[fmt.Sprintf("%d", r.Year)]++
	}
	labels, values := sortedPairs(yearCounts)
	if err := writeBarChartSVG(
		outputDir+"/temperature_rows_by_year.svg",
		"Temperature Rows by Year",
		labels,
		values,
		"Year",
		"Rows",
	); err != nil {
		return err
	}

	type agg struct {
		sum   float64
		count int
	}
	heatByDistrict := map[string]agg{}
	for _, r := range healthRows {
		key := r.City + "::" + r.District
		a := heatByDistrict[key]
		a.sum += r.HeatCasesPer100K
		a.count++
		heatByDistrict[key] = a
	}
	means := map[string]float64{}
	for k, a := range heatByDistrict {
		if a.count > 0 {
			means[k] = a.sum / float64(a.count)
		}
	}
	dLabels, dValues := sortedPairs(means)
	if len(dLabels) > 20 {
		dLabels = dLabels[:20]
		dValues = dValues[:20]
	}
	if err := writeBarChartSVG(
		outputDir+"/mean_heat_cases_by_district.svg",
		"Mean Heat Cases per 100k by District (Top 20)",
		dLabels,
		dValues,
		"District",
		"Mean Heat Cases per 100k",
	); err != nil {
		return err
	}

	tempValues := make([]float64, 0, len(tempRows))
	for _, r := range tempRows {
		tempValues = append(tempValues, r.ValueC)
	}
	hLabels, hValues := histogram(tempValues, 10)
	if err := writeBarChartSVG(
		outputDir+"/temperature_histogram.svg",
		"Temperature Distribution",
		hLabels,
		hValues,
		"Temperature bins (deg C)",
		"Count",
	); err != nil {
		return err
	}

	return nil
}

func sortedPairs(m map[string]float64) ([]string, []float64) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	values := make([]float64, 0, len(keys))
	for _, k := range keys {
		values = append(values, m[k])
	}
	return keys, values
}

func histogram(values []float64, bins int) ([]string, []float64) {
	if len(values) == 0 || bins <= 0 {
		return []string{}, []float64{}
	}
	minV := math.MaxFloat64
	maxV := -math.MaxFloat64
	for _, v := range values {
		if v < minV {
			minV = v
		}
		if v > maxV {
			maxV = v
		}
	}
	if minV == maxV {
		return []string{fmt.Sprintf("%.2f-%.2f", minV, maxV)}, []float64{float64(len(values))}
	}

	width := (maxV - minV) / float64(bins)
	counts := make([]float64, bins)
	for _, v := range values {
		idx := int((v - minV) / width)
		if idx == bins {
			idx = bins - 1
		}
		counts[idx]++
	}

	labels := make([]string, bins)
	for i := 0; i < bins; i++ {
		low := minV + float64(i)*width
		high := low + width
		labels[i] = fmt.Sprintf("%.2f-%.2f", low, high)
	}
	return labels, counts
}

func writeBarChartSVG(path, title string, labels []string, values []float64, xAxis, yAxis string) error {
	if len(labels) == 0 || len(values) == 0 || len(labels) != len(values) {
		return os.WriteFile(path, []byte("<svg xmlns='http://www.w3.org/2000/svg' width='800' height='200'><text x='20' y='40'>No data available</text></svg>"), 0o644)
	}

	const (
		width  = 1200.0
		height = 700.0
		left   = 90.0
		right  = 30.0
		top    = 70.0
		bottom = 180.0
	)
	plotW := width - left - right
	plotH := height - top - bottom

	maxV := 0.0
	for _, v := range values {
		if v > maxV {
			maxV = v
		}
	}
	if maxV == 0 {
		maxV = 1
	}

	barW := plotW / float64(len(values))

	var b strings.Builder
	b.WriteString(fmt.Sprintf("<svg xmlns='http://www.w3.org/2000/svg' width='%.0f' height='%.0f'>", width, height))
	b.WriteString("<rect width='100%' height='100%' fill='white'/>")
	b.WriteString(fmt.Sprintf("<text x='%.0f' y='35' font-size='22' font-family='Arial'>%s</text>", left, escapeXML(title)))

	x0 := left
	y0 := top + plotH
	b.WriteString(fmt.Sprintf("<line x1='%.2f' y1='%.2f' x2='%.2f' y2='%.2f' stroke='black'/>", x0, y0, x0+plotW, y0))
	b.WriteString(fmt.Sprintf("<line x1='%.2f' y1='%.2f' x2='%.2f' y2='%.2f' stroke='black'/>", x0, top, x0, y0))

	for i, v := range values {
		h := (v / maxV) * plotH
		x := left + float64(i)*barW + barW*0.12
		y := y0 - h
		w := barW * 0.76
		b.WriteString(fmt.Sprintf("<rect x='%.2f' y='%.2f' width='%.2f' height='%.2f' fill='#4f81bd'/>", x, y, w, h))

		labelX := x + w/2
		labelY := y0 + 18
		b.WriteString(fmt.Sprintf("<g transform='translate(%.2f,%.2f) rotate(45)'><text font-size='10' text-anchor='start' font-family='Arial'>%s</text></g>",
			labelX, labelY, escapeXML(labels[i])))
	}

	for i := 0; i <= 5; i++ {
		tickV := maxV * float64(i) / 5.0
		tickY := y0 - (tickV/maxV)*plotH
		b.WriteString(fmt.Sprintf("<line x1='%.2f' y1='%.2f' x2='%.2f' y2='%.2f' stroke='#dddddd'/>", x0, tickY, x0+plotW, tickY))
		b.WriteString(fmt.Sprintf("<text x='%.2f' y='%.2f' font-size='11' text-anchor='end' font-family='Arial'>%.2f</text>", x0-8, tickY+4, tickV))
	}

	b.WriteString(fmt.Sprintf("<text x='%.2f' y='%.2f' font-size='14' font-family='Arial'>%s</text>", left+plotW/2-80, height-40, escapeXML(xAxis)))
	b.WriteString(fmt.Sprintf("<g transform='translate(25,%.2f) rotate(-90)'><text font-size='14' font-family='Arial'>%s</text></g>", top+plotH/2+70, escapeXML(yAxis)))
	b.WriteString("</svg>")

	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func escapeXML(s string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
		"'", "&apos;",
	)
	return replacer.Replace(s)
}
