package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"hello-go/internal/pipeline"
)

func main() {
	temperaturePath := flag.String("temperature", "", "Path to temperature CSV")
	demographicsPath := flag.String("demographics", "", "Path to demographics CSV")
	healthPath := flag.String("health", "", "Path to health CSV")
	outputDir := flag.String("out", "output", "Directory for outputs")
	maxDistanceKM := flag.Float64("max-distance-km", 10, "Maximum distance for spatial matching")
	trainRatio := flag.Float64("train-ratio", 0.8, "Train set ratio between 0 and 1")
	seed := flag.Int64("seed", 42, "Random seed for reproducible split")
	flag.Parse()

	if *temperaturePath == "" || *demographicsPath == "" || *healthPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := os.MkdirAll(*outputDir, 0o755); err != nil {
		log.Fatalf("create output dir: %v", err)
	}

	tempRows, err := pipeline.LoadTemperatureCSV(*temperaturePath)
	if err != nil {
		log.Fatalf("load temperature csv: %v", err)
	}
	demoRows, err := pipeline.LoadDemographicsCSV(*demographicsPath)
	if err != nil {
		log.Fatalf("load demographics csv: %v", err)
	}
	healthRows, err := pipeline.LoadHealthCSV(*healthPath)
	if err != nil {
		log.Fatalf("load health csv: %v", err)
	}

	mergedRows, err := pipeline.BuildMergedDataset(healthRows, tempRows, demoRows, *maxDistanceKM)
	if err != nil {
		log.Fatalf("build merged dataset: %v", err)
	}

	model, predictions, err := pipeline.TrainAndEvaluate(mergedRows, *trainRatio, *seed)
	if err != nil {
		log.Fatalf("train and evaluate model: %v", err)
	}

	mergedPath := filepath.Join(*outputDir, "processed_merged.csv")
	predPath := filepath.Join(*outputDir, "predictions.csv")
	reportPath := filepath.Join(*outputDir, "model_report.txt")
	explorePath := filepath.Join(*outputDir, "data_exploration.txt")
	plotDir := filepath.Join(*outputDir, "plots")

	if err := pipeline.WriteMergedCSV(mergedPath, mergedRows); err != nil {
		log.Fatalf("write merged csv: %v", err)
	}
	if err := pipeline.WritePredictionsCSV(predPath, predictions); err != nil {
		log.Fatalf("write predictions csv: %v", err)
	}
	if err := pipeline.WriteDataExplorationReport(
		explorePath,
		*temperaturePath,
		*demographicsPath,
		*healthPath,
		tempRows,
		demoRows,
		healthRows,
	); err != nil {
		log.Fatalf("write data exploration report: %v", err)
	}
	if err := pipeline.WriteExplorationPlots(plotDir, tempRows, healthRows); err != nil {
		log.Fatalf("write exploration plots: %v", err)
	}

	if err := pipeline.WriteModelReport(
		reportPath,
		model,
		pipeline.SummariesByYear(predictions),
		pipeline.SummariesByCity(predictions),
		pipeline.SummariesByDistrict(predictions),
	); err != nil {
		log.Fatalf("write model report: %v", err)
	}

	fmt.Printf("Done. Processed rows: %d\n", len(mergedRows))
	fmt.Printf("RMSE: %.4f | MAE: %.4f | R^2: %.4f\n", model.RMSE, model.MAE, model.R2)
	fmt.Printf("Output files:\n- %s\n- %s\n- %s\n- %s\n- %s\n", mergedPath, predPath, reportPath, explorePath, plotDir)
}
