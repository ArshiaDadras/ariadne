package internal

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/ArshiaDadras/Ariadne/pkg"
)

func SaveObject(obj interface{}, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}

	_, err = file.Write(b)
	if err != nil {
		return err
	}

	return nil
}

func ParseCSV(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	reader.Comma = '\t'

	data, err := reader.ReadAll()
	return data[1:], err
}

func ParseGPSData(path string) ([]GPSPoint, error) {
	data, err := ParseCSV(path)
	if err != nil {
		return nil, err
	}

	points := make([]GPSPoint, 0, len(data))
	for _, row := range data {
		latitude, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, err
		}

		longitude, err := strconv.ParseFloat(row[3], 64)
		if err != nil {
			return nil, err
		}

		dateTime, err := time.Parse(TimeFormat, row[0]+" "+row[1])
		if err != nil {
			return nil, err
		}

		points = append(points, GPSPoint{
			Location: pkg.Point{Longitude: longitude, Latitude: latitude},
			Time:     dateTime,
		})
	}

	slices.SortFunc(points, func(a, b GPSPoint) int {
		if a.Time.Before(b.Time) {
			return -1
		} else if a.Time.After(b.Time) {
			return 1
		}
		return 0
	})

	return points, nil
}
