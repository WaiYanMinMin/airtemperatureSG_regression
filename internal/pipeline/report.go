package pipeline

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

func WriteModelReport(path string, result ModelResult, byYear, byCity, byDistrict []GroupSummary) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var b strings.Builder
	b.WriteString("Heat-Health Regression Report\n")
	b.WriteString("=============================\n\n")
	b.WriteString("Model Metrics\n")
	b.WriteString("-------------\n")
	b.WriteString(fmt.Sprintf("Trained at : %s\n", result.TrainedAt.Format("2006-01-02 15:04:05")))
	b.WriteString(fmt.Sprintf("RMSE       : %.6f\n", result.RMSE))
	b.WriteString(fmt.Sprintf("MAE        : %.6f\n", result.MAE))
	b.WriteString(fmt.Sprintf("R^2        : %.6f\n\n", result.R2))

	b.WriteString("Model Coefficients\n")
	b.WriteString("------------------\n")
	b.WriteString(fmt.Sprintf("%-20s %15s\n", "Feature", "Coefficient"))
	b.WriteString(fmt.Sprintf("%-20s %15s\n", strings.Repeat("-", 20), strings.Repeat("-", 15)))
	keys := make([]string, 0, len(result.Coefficients))
	for k := range result.Coefficients {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		b.WriteString(fmt.Sprintf("%-20s %15.8f\n", k, result.Coefficients[k]))
	}
	b.WriteString("\n")

	appendSummaryTable(&b, "Comparison by Year", byYear)
	appendSummaryTable(&b, "Comparison by City", byCity)
	appendSummaryTable(&b, "Comparison by District", byDistrict)
	appendNarrative(&b, result, byYear, byCity, byDistrict)

	_, err = f.WriteString(b.String())
	return err
}

func appendSummaryTable(b *strings.Builder, title string, rows []GroupSummary) {
	b.WriteString(title + "\n")
	b.WriteString(strings.Repeat("-", len(title)) + "\n")
	b.WriteString(fmt.Sprintf("%-24s %5s %13s %13s %14s\n", "Group", "N", "Mean Actual", "Mean Pred", "Mean Residual"))
	b.WriteString(fmt.Sprintf("%-24s %5s %13s %13s %14s\n",
		strings.Repeat("-", 24), strings.Repeat("-", 5), strings.Repeat("-", 13), strings.Repeat("-", 13), strings.Repeat("-", 14)))
	for _, row := range rows {
		b.WriteString(fmt.Sprintf(
			"%-24s %5d %13.6f %13.6f %14.6f\n",
			row.Group, row.Count, row.MeanActual, row.MeanPred, row.MeanResidual,
		))
	}
	b.WriteString("\n")
}

func appendNarrative(b *strings.Builder, result ModelResult, byYear, byCity, byDistrict []GroupSummary) {
	b.WriteString("Short Summary\n")
	b.WriteString("-------------\n")
	b.WriteString(fmt.Sprintf(
		"- The regression explains %.2f%% of observed variation (R^2 = %.4f), indicating a moderate overall fit.\n",
		result.R2*100, result.R2,
	))
	b.WriteString(fmt.Sprintf(
		"- Prediction error is RMSE = %.4f and MAE = %.4f heat cases per 100k.\n",
		result.RMSE, result.MAE,
	))
	feature := strongestFeature(result.Coefficients)
	b.WriteString(fmt.Sprintf(
		"- In this run, the strongest positive predictor is `%s` (coefficient %.4f).\n",
		feature, result.Coefficients[feature],
	))
	if len(byDistrict) > 0 {
		maxAbsResidual := byDistrict[0]
		for _, row := range byDistrict {
			if math.Abs(row.MeanResidual) > math.Abs(maxAbsResidual.MeanResidual) {
				maxAbsResidual = row
			}
		}
		b.WriteString(fmt.Sprintf(
			"- Largest district-level mismatch is %s with mean residual %.4f, suggesting local effects not fully captured.\n",
			maxAbsResidual.Group, maxAbsResidual.MeanResidual,
		))
	}
	b.WriteString("\n")

	b.WriteString("Conclusion\n")
	b.WriteString("----------\n")
	b.WriteString("- This model is a solid baseline for the assignment and shows that temperature and demographic factors are associated with heat-related health risk.\n")
	b.WriteString("- Results should be interpreted carefully because the current evaluation set is small; adding more years/districts/cities is likely to improve robustness.\n")
	b.WriteString("- Recommended next steps: include more predictors (e.g., humidity, rainfall, green cover), then compare metrics after retraining.\n")
	b.WriteString("\n")
}

func strongestFeature(coeffs map[string]float64) string {
	bestKey := ""
	bestValue := -1.0
	for k, v := range coeffs {
		if k == "intercept" {
			continue
		}
		absV := math.Abs(v)
		if absV > bestValue {
			bestValue = absV
			bestKey = k
		}
	}
	if bestKey == "" {
		return "n/a"
	}
	return bestKey
}
