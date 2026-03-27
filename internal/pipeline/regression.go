package pipeline

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)
type datasetSplit struct {
	Train []MergedRow
	Test  []MergedRow
}

func TrainAndEvaluate(rows []MergedRow, trainRatio float64, seed int64) (ModelResult, []PredictionRow, error) {
	if len(rows) < 10 {
		return ModelResult{}, nil, fmt.Errorf("need at least 10 merged rows; got %d", len(rows))
	}

	split := splitDataset(rows, trainRatio, seed)
	if len(split.Train) == 0 || len(split.Test) == 0 {
		return ModelResult{}, nil, fmt.Errorf("train/test split failed")
	}

	featureNames := []string{"intercept", "temperature_c", "population", "elderly_share", "income_index"}
	trainX, trainY := matrixFromRows(split.Train)
	coeffs, err := fitOLS(trainX, trainY)
	if err != nil {
		return ModelResult{}, nil, err
	}

	preds := make([]PredictionRow, 0, len(split.Test))
	actual := make([]float64, 0, len(split.Test))
	predValues := make([]float64, 0, len(split.Test))
	for _, r := range split.Test {
		p := predictOne(coeffs, r)
		residual := r.HeatCasesPer100K - p
		preds = append(preds, PredictionRow{
			MergedRow:                 r,
			PredictedHeatCasesPer100K: p,
			Residual:                  residual,
		})
		actual = append(actual, r.HeatCasesPer100K)
		predValues = append(predValues, p)
	}

	result := ModelResult{
		Coefficients: make(map[string]float64, len(featureNames)),
		RMSE:         rmse(actual, predValues),
		MAE:          mae(actual, predValues),
		R2:           r2(actual, predValues),
		TrainedAt:    time.Now(),
	}
	for i, name := range featureNames {
		result.Coefficients[name] = coeffs[i]
	}

	return result, preds, nil
}

func splitDataset(rows []MergedRow, trainRatio float64, seed int64) datasetSplit {
	if trainRatio <= 0 || trainRatio >= 1 {
		trainRatio = 0.8
	}
	data := append([]MergedRow(nil), rows...)

	rng := rand.New(rand.NewSource(seed))
	rng.Shuffle(len(data), func(i, j int) {
		data[i], data[j] = data[j], data[i]
	})

	cut := int(float64(len(data)) * trainRatio)
	if cut < 1 {
		cut = 1
	}
	if cut >= len(data) {
		cut = len(data) - 1
	}

	return datasetSplit{
		Train: data[:cut],
		Test:  data[cut:],
	}
}

func matrixFromRows(rows []MergedRow) ([][]float64, []float64) {
	x := make([][]float64, len(rows))
	y := make([]float64, len(rows))
	for i, r := range rows {
		x[i] = []float64{1.0, r.TemperatureC, r.Population, r.ElderlyShare, r.IncomeIndex}
		y[i] = r.HeatCasesPer100K
	}
	return x, y
}

func fitOLS(x [][]float64, y []float64) ([]float64, error) {
	if len(x) == 0 || len(x) != len(y) {
		return nil, fmt.Errorf("invalid matrix dimensions")
	}
	k := len(x[0])

	xtx := make([][]float64, k)
	for i := range xtx {
		xtx[i] = make([]float64, k)
	}
	xty := make([]float64, k)

	for i := range x {
		for row := 0; row < k; row++ {
			xty[row] += x[i][row] * y[i]
			for col := 0; col < k; col++ {
				xtx[row][col] += x[i][row] * x[i][col]
			}
		}
	}

	return solveLinearSystem(xtx, xty)
}

func solveLinearSystem(a [][]float64, b []float64) ([]float64, error) {
	n := len(a)
	aug := make([][]float64, n)
	for i := 0; i < n; i++ {
		aug[i] = make([]float64, n+1)
		copy(aug[i], a[i])
		aug[i][n] = b[i]
	}

	for pivot := 0; pivot < n; pivot++ {
		maxRow := pivot
		for r := pivot + 1; r < n; r++ {
			if math.Abs(aug[r][pivot]) > math.Abs(aug[maxRow][pivot]) {
				maxRow = r
			}
		}
		if math.Abs(aug[maxRow][pivot]) < 1e-12 {
			return nil, fmt.Errorf("matrix is singular or near-singular")
		}
		aug[pivot], aug[maxRow] = aug[maxRow], aug[pivot]

		divisor := aug[pivot][pivot]
		for c := pivot; c <= n; c++ {
			aug[pivot][c] /= divisor
		}
		for r := 0; r < n; r++ {
			if r == pivot {
				continue
			}
			factor := aug[r][pivot]
			for c := pivot; c <= n; c++ {
				aug[r][c] -= factor * aug[pivot][c]
			}
		}
	}

	coeffs := make([]float64, n)
	for i := 0; i < n; i++ {
		coeffs[i] = aug[i][n]
	}
	return coeffs, nil
}

func predictOne(coeffs []float64, r MergedRow) float64 {
	return coeffs[0] +
		coeffs[1]*r.TemperatureC +
		coeffs[2]*r.Population +
		coeffs[3]*r.ElderlyShare +
		coeffs[4]*r.IncomeIndex
}

func rmse(actual, predicted []float64) float64 {
	if len(actual) == 0 {
		return 0
	}
	sum := 0.0
	for i := range actual {
		diff := actual[i] - predicted[i]
		sum += diff * diff
	}
	return math.Sqrt(sum / float64(len(actual)))
}

func mae(actual, predicted []float64) float64 {
	if len(actual) == 0 {
		return 0
	}
	sum := 0.0
	for i := range actual {
		sum += math.Abs(actual[i] - predicted[i])
	}
	return sum / float64(len(actual))
}

func r2(actual, predicted []float64) float64 {
	if len(actual) == 0 {
		return 0
	}
	mean := 0.0
	for _, v := range actual {
		mean += v
	}
	mean /= float64(len(actual))

	ssTot, ssRes := 0.0, 0.0
	for i := range actual {
		dTot := actual[i] - mean
		dRes := actual[i] - predicted[i]
		ssTot += dTot * dTot
		ssRes += dRes * dRes
	}
	if ssTot == 0 {
		return 0
	}
	return 1 - (ssRes / ssTot)
}

type GroupSummary struct {
	Group        string
	Count        int
	MeanActual   float64
	MeanPred     float64
	MeanResidual float64
}

func SummariesByYear(rows []PredictionRow) []GroupSummary {
	return summarize(rows, func(r PredictionRow) string { return fmt.Sprintf("%d", r.Year) })
}

func SummariesByCity(rows []PredictionRow) []GroupSummary {
	return summarize(rows, func(r PredictionRow) string { return r.City })
}

func SummariesByDistrict(rows []PredictionRow) []GroupSummary {
	return summarize(rows, func(r PredictionRow) string { return r.City + "::" + r.District })
}

func summarize(rows []PredictionRow, keyFn func(PredictionRow) string) []GroupSummary {
	type acc struct {
		count int
		a     float64
		p     float64
		r     float64
	}
	m := map[string]*acc{}
	for _, row := range rows {
		key := keyFn(row)
		if _, ok := m[key]; !ok {
			m[key] = &acc{}
		}
		m[key].count++
		m[key].a += row.HeatCasesPer100K
		m[key].p += row.PredictedHeatCasesPer100K
		m[key].r += row.Residual
	}

	out := make([]GroupSummary, 0, len(m))
	for key, v := range m {
		out = append(out, GroupSummary{
			Group:        key,
			Count:        v.count,
			MeanActual:   v.a / float64(v.count),
			MeanPred:     v.p / float64(v.count),
			MeanResidual: v.r / float64(v.count),
		})
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Group < out[j].Group })
	return out
}
