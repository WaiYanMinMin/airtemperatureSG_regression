# Urban Heat and Health Regression (Go)

This project gives you an assignment-ready Go workflow for:
- loading temperature, demographics, and health datasets
- cleaning and merging data
- spatial matching by latitude/longitude (nearest temperature point)
- training a multiple linear regression model
- analysing differences across years, cities, and districts
- exporting processed CSVs and model outputs

## Project Structure

- `cmd/prepare_singapore_airtemp`: converts provided Singapore JSON into `temperature` CSV format
- `cmd/heatmodel`: main pipeline and regression runner
- `internal/pipeline`: reusable loading, matching, model, and report code
- `dataset`: sample datasets and raw JSON
- `output`: generated processed files (created at runtime)
- `prompts/prompt_log_template.md`: Gen AI prompt log template

## 1) Dataset Column Requirements

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

## 2) Quick Start with Included Files

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

## 3) Output Files

After running `cmd/heatmodel`, these files are created:
- `output/processed_merged.csv`
- `output/predictions.csv`
- `output/model_report.txt`

`model_report.txt` includes:
- regression metrics (`RMSE`, `MAE`, `R^2`)
- learned coefficients
- grouped summaries by year, city, district

## 4) Mapping to Assignment Tasks

1. **Load and explore datasets**: CSV loaders in `internal/pipeline/csvio.go`
2. **Clean and prepare**: invalid rows are dropped during parsing
3. **Spatial matching**: Haversine nearest-neighbour in `internal/pipeline/merge.go`
4. **Regression model in Go**: OLS in `internal/pipeline/regression.go`
5. **Differences across year/district/city**: grouped summaries in report
6. **AI prompts and validation**: fill `prompts/prompt_log_template.md`
7. **Presentation insights**: derive charts/tables from output CSV/report

## 5) Replace with Lecturer Datasets

Use the same required column names in your real CSV files. If your source files use different names, rename the headers before running.

For multiple cities (Singapore, Bangkok, Jakarta, Shanghai), keep city names in the `city` column and place all rows in the same CSV files.
