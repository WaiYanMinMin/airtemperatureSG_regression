package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
	"time"
)

type apiPayload struct {
	Data struct {
		Stations []struct {
			ID       string `json:"id"`
			Location struct {
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			} `json:"location"`
		} `json:"stations"`
		Readings []struct {
			Timestamp string `json:"timestamp"`
			Data      []struct {
				StationID string  `json:"stationId"`
				Value     float64 `json:"value"`
			} `json:"data"`
		} `json:"readings"`
	} `json:"data"`
}

func main() {
	input := flag.String("input", "dataset/NGDS_AirTemperatureacrossSingapore.json", "Input JSON path")
	output := flag.String("output", "dataset/temperature_singapore.csv", "Output CSV path")
	flag.Parse()

	in, err := os.Open(*input)
	if err != nil {
		log.Fatalf("open input json: %v", err)
	}
	defer in.Close()

	var payload apiPayload
	if err := json.NewDecoder(in).Decode(&payload); err != nil {
		log.Fatalf("decode input json: %v", err)
	}

	stationCoords := map[string][2]float64{}
	for _, s := range payload.Data.Stations {
		stationCoords[s.ID] = [2]float64{s.Location.Latitude, s.Location.Longitude}
	}

	out, err := os.Create(*output)
	if err != nil {
		log.Fatalf("create output csv: %v", err)
	}
	defer out.Close()

	w := csv.NewWriter(out)
	defer w.Flush()

	if err := w.Write([]string{"city", "district", "year", "month", "lat", "lon", "temperature_c"}); err != nil {
		log.Fatalf("write header: %v", err)
	}

	for _, reading := range payload.Data.Readings {
		t, err := time.Parse(time.RFC3339, reading.Timestamp)
		if err != nil {
			continue
		}
		year := strconv.Itoa(t.Year())
		month := strconv.Itoa(int(t.Month()))

		for _, d := range reading.Data {
			coords, ok := stationCoords[d.StationID]
			if !ok {
				continue
			}
			record := []string{
				"Singapore",
				d.StationID,
				year,
				month,
				floatToString(coords[0]),
				floatToString(coords[1]),
				floatToString(d.Value),
			}
			if err := w.Write(record); err != nil {
				log.Fatalf("write record: %v", err)
			}
		}
	}
}

func floatToString(v float64) string {
	return strconv.FormatFloat(v, 'f', 6, 64)
}
