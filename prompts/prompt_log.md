# Gen AI Prompt Log (Phased Learning Plan)

This log is structured to show that I used AI as a tutor while learning Go from a Python background.
I used prompts to understand concepts, then implemented and validated code myself.

---

## Phase 1: Go Fundamentals (Python -> Go mindset)

### Entry 1 - Basics & Syntax
- Date: 2026-03-25
- Prompt:
  "I have a Python background. Teach me Go basics (variables, types, loops, conditionals, functions) by comparing each concept with Python and giving examples."
- Learn Focus:
  - Static typing vs Python dynamic typing
  - `for` loops as the main Go loop
  - Functions and multiple return values
- Why:
  - Foundation before writing any data pipeline code in Go.
- Output Summary:
  - Received side-by-side Python vs Go examples and beginner syntax explanations.
- Validation:
  - Rewrote simple Python-style snippets into Go.
  - Compiled examples and fixed type errors manually.
  - Confirmed I can explain `:=`, type declarations, and return signatures.

### Entry 2 - Data Structures
- Date: 2026-03-25
- Prompt:
  "Explain Go arrays, slices, maps, and structs with examples. Compare them with Python lists, dictionaries, and classes."
- Learn Focus:
  - Slices (`[]T`) for dynamic collections
  - Maps for key-value indexing
  - Structs for typed row models
- Why:
  - These are core tools replacing most Pandas-like thinking in raw Go.
- Output Summary:
  - Received examples for converting list/dict/class patterns into Go slices/maps/structs.
- Validation:
  - Built row structs (`TemperatureRow`, `HealthRow`, etc.).
  - Used slices + append in parsing loops.
  - Used maps for keyed merges and grouped summaries.

### Entry 3 - Error Handling
- Date: 2026-03-25
- Prompt:
  "Teach me Go error handling patterns and how they differ from Python exceptions. Give real examples."
- Learn Focus:
  - `if err != nil`
  - Explicit error returns instead of try/except
- Why:
  - Critical for file reading and safe data processing.
- Output Summary:
  - Learned idiomatic pattern of returning `(value, error)` and handling errors at call sites.
- Validation:
  - Added `if err != nil` checks at every I/O and model step.
  - Tested with invalid/missing file paths to verify failure behavior.

---

## Phase 2: File Handling & Data Loading

### Entry 4 - Read CSV Files
- Date: 2026-03-25
- Prompt:
  "Show me how to read and parse CSV files in Go using built-in packages. Then compare with pandas.read_csv."
- Learn Focus:
  - `encoding/csv`
  - File I/O and row parsing
- Why:
  - First requirement for assignment datasets.
- Output Summary:
  - Got a reusable pattern for header normalization + field parsing.
- Validation:
  - Implemented CSV loaders for temperature, demographics, and health.
  - Checked row counts and malformed-row skipping behavior.

### Entry 5 - Work with JSON Data
- Date: 2026-03-25
- Prompt:
  "Teach me how to read, parse, and write JSON data in Go using structs. Compare with Python json and pandas."
- Learn Focus:
  - Struct tags (e.g., `json:"stationId"`)
  - Unmarshal/decode flow
- Why:
  - APIs + climate datasets are often JSON before conversion to CSV.
- Output Summary:
  - Learned JSON decoding into nested structs and converting to flat CSV rows.
- Validation:
  - Built converter from Singapore JSON to `temperature_singapore.csv`.
  - Verified generated columns match pipeline requirements.

---

## Phase 3: Data Cleaning & Preprocessing

### Entry 6 - Data Cleaning
- Date: 2026-03-25
- Prompt:
  "Show me how to clean a dataset in Go: handling missing values, filtering rows, and transforming columns. Use CSV data examples."
- Learn Focus:
  - Manual cleaning logic
  - Reusable parsing helper functions
- Why:
  - No Pandas shortcut in base Go; cleaning must be explicit.
- Output Summary:
  - Used field-by-field parsing with required/optional value checks.
- Validation:
  - Tested missing fields and parse errors.
  - Confirmed dropped rows are tracked in exploration report.

### Entry 7 - Data Transformation Pipeline
- Date: 2026-03-25
- Prompt:
  "Help me design a reusable data preprocessing pipeline in Go (like pandas pipelines), including filtering, mapping, and feature transformation."
- Learn Focus:
  - Modular function design
  - Pipeline sequencing
- Why:
  - Needed a maintainable workflow for assignment tasks.
- Output Summary:
  - Defined staged pipeline: load -> clean -> spatial match -> model -> report.
- Validation:
  - Confirmed each step can run independently.
  - Verified final merged output for expected columns.

---

## Phase 4: Basic Data Analysis

### Entry 8 - Descriptive Statistics
- Date: 2026-03-25
- Prompt:
  "Show me how to compute mean, median, variance, and standard deviation in Go. Compare with numpy."
- Learn Focus:
  - Math/stat helpers in Go
  - Writing statistical functions manually
- Why:
  - Core exploration step before model training.
- Output Summary:
  - Learned manual stat calculations and summary table patterns.
- Validation:
  - Implemented and reviewed count/min/max/mean summaries in exploration output.
  - Cross-checked values with quick manual checks on sample ranges.

### Entry 9 - Data Exploration
- Date: 2026-03-25
- Prompt:
  "Help me perform exploratory data analysis (EDA) in Go, including summary stats and basic insights. Suggest libraries if needed."
- Learn Focus:
  - Coverage summaries by city/year/district
  - Lightweight charting options
  - Awareness of `gonum` / `gota` ecosystem
- Why:
  - Assignment explicitly needs exploration and insights.
- Output Summary:
  - Added explicit exploration report and generated SVG plots.
- Validation:
  - Confirmed `output/data_exploration.txt` and plot files are created every run.
  - Checked labels and group counts are correct.

---

##  Phase 5: Regression Model

### Entry 10 - Linear Regression
- Date: 2026-03-25
- Prompt:
  "Teach me how to implement linear regression in Go from scratch, then using a library like gonum. Compare with scikit-learn."
- Learn Focus:
  - OLS math intuition
  - Matrix-based training workflow
- Why:
  - Core requirement: build regression model in Go.
- Output Summary:
  - Implemented a from-scratch OLS baseline (normal equation + linear solver).
- Validation:
  - Checked coefficients are finite.
  - Confirmed reproducible results with fixed seed.

### Entry 11 - Model Evaluation
- Date: 2026-03-25
- Prompt:
  "Show me how to evaluate a regression model in Go using metrics like RMSE and R-squared."
- Learn Focus:
  - RMSE, MAE, R^2
  - Train/test split evaluation
- Why:
  - Needed objective quality checks for report discussion.
- Output Summary:
  - Implemented metric functions and prediction residual output.
- Validation:
  - Verified metrics print in CLI and report.
  - Reviewed residuals in `output/predictions.csv`.

---

## Phase 6: Project Integration

### Entry 12 - End-to-End Project
- Date: 2026-03-25
- Prompt:
  "Guide me step-by-step to build a Go project that loads CSV, cleans data, performs EDA, trains regression, and evaluates model. Explain each step and structure like a real project."
- Learn Focus:
  - Real-world project layout
  - End-to-end reproducible workflow
- Why:
  - Needed to connect all assignment tasks in one runnable pipeline.
- Output Summary:
  - Produced an integrated CLI workflow and documented execution commands.
- Validation:
  - Ran project end-to-end:
    - `go run ./cmd/prepare_singapore_airtemp`
    - `go run ./cmd/heatmodel ...`
  - Confirmed all outputs generated successfully.

---

## Responsible AI Use Statement
- AI was used as a tutor for concepts, syntax translation (Python -> Go), and debugging guidance.
- I validated outputs by compiling, running, checking generated files, and reviewing metrics manually.
- Final code decisions, integration, and assignment interpretation were done by me.
- I did not request AI to generate final presentation conclusions without checking results myself.
# Gen AI Prompt Log (Heat-Health Regression in Go)

This log documents how I used AI support responsibly while still writing, testing, and understanding the project myself.

---
