# Urban Heat and Health Regression (Go)

This project is a Go-based data pipeline for analyzing the relationship between urban heat exposure and health risk.
It includes:
- loading temperature, demographics, and health datasets
- cleaning and merging data
- spatial matching by latitude/longitude (nearest temperature point)
- training a multiple linear regression model
- comparing results across years, cities, and districts
- generating exploration reports and plots
- exporting processed CSVs, predictions, and model summaries

## Project Structure

- `cmd/prepare_singapore_airtemp`: converts provided Singapore JSON into `temperature` CSV format
- `cmd/heatmodel`: main pipeline and regression runner
- `internal/pipeline`: reusable loading, matching, exploration, model, plotting, and report code
- `dataset`: sample datasets and raw JSON
- `output`: generated processed files (created at runtime)
- `prompts/prompt_log.md`: phased AI prompt log used during development

## Dataset Column Requirements

### Temperature CSV
Required columns:
- `city`
- `district`
- `year`
- `lat`
- `lon`
- `temperature_c`

Optional:
- `month` (defaults to `1` if missing)

### Demographics CSV
Required columns:
- `city`
- `district`
- `year`
- `population`
- `elderly_share`

Optional:
- `income_index` (defaults to `0`)

### Health CSV (target variable)
Required columns:
- `city`
- `district`
- `year`
- `lat`
- `lon`
- `heat_cases_per_100k`

## Quick Start

Generate a temperature CSV from the provided Singapore JSON:

```bash
go run ./cmd/prepare_singapore_airtemp \
  -input dataset/NGDS_AirTemperatureacrossSingapore.json \
  -output dataset/temperature_singapore.csv
```

Run the pipeline:

```bash
go run ./cmd/heatmodel \
  -temperature dataset/temperature_singapore.csv \
  -demographics dataset/demographics_sample.csv \
  -health dataset/health_sample.csv \
  -out output \
  -max-distance-km 10 \
  -train-ratio 0.8 \
  -seed 42
```

## Outputs

After running `cmd/heatmodel`, these files are created:
- `output/processed_merged.csv`
- `output/predictions.csv`
- `output/model_report.txt`
- `output/data_exploration.txt`
- `output/plots/temperature_rows_by_year.svg`
- `output/plots/temperature_histogram.svg`
- `output/plots/mean_heat_cases_by_district.svg`

`model_report.txt` includes:
- regression metrics (`RMSE`, `MAE`, `R^2`)
- learned coefficients
- grouped summaries by year, city, district
- short summary and conclusion

`data_exploration.txt` includes:
- parsed vs dropped row counts per dataset
- coverage by city/year/district
- feature statistics (count/min/max/mean)

## How the Pipeline Works

1. **Load data** from CSV files (`internal/pipeline/csvio.go`)
2. **Clean rows** with strict field parsing and optional defaults
3. **Spatially match** health rows to nearest temperature points (`internal/pipeline/merge.go`)
4. **Explore data** via textual stats (`internal/pipeline/explore.go`) and SVG plots (`internal/pipeline/plot.go`)
5. **Train model** with OLS regression (`internal/pipeline/regression.go`)
6. **Evaluate and summarize** with metrics and grouped comparisons (`internal/pipeline/report.go`)

