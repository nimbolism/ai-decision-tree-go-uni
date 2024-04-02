package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func DataClean(records [][]string) ([][]string, error) {
	// removing first row - column names!
	records = removeFirstRow(records)

	// removing evero mising value - we cant approximate!
	records = removeMissing(records)

	// removing the id column - first column
	for i := range records {
		records[i] = append(records[i][:0], records[i][1:]...)
	}

	// using numerical values instead of texts for education and maritial
	records = mapTextToNumeric(records)

	// converting numerical columns to integers
	intColumns := []int{0, 5, 9, 11, 12, 13, 14}
	for _, record := range records {
		for _, colIdx := range intColumns {
			val, err := strconv.Atoi(record[colIdx])
			if err != nil {
				fmt.Printf("Error converting column %d to integer: %v\n", colIdx, err)
				return nil, err
			}
			record[colIdx] = strconv.Itoa(val)
		}
	}

	// adding age restrictions
	records = filterRecordsAge(records)

	// converting text to boolean values
	boolColumns := map[int]bool{4: true, 6: true, 7: true, 16: true}
	for _, record := range records {
		for colIdx, isBool := range boolColumns {
			if isBool {
				if record[colIdx] == "yes" {
					record[colIdx] = "true"
				} else {
					record[colIdx] = "false"
				}
			}
		}
	}

	// standardizing balance column and scaling its values!
	records = standardizeBalance(records)

	// doing one-hot encoding for job attribute
	records = addJobFeatures(records)

	// removing the original job column
	for i := range records {
		records[i] = append(records[i][:1], records[i][2:]...)
	}

	// writing cleaned data to a CSV file for later reviews!
	err := SaveCSV(records, "cleaned_data.csv")
	if err != nil {
		return nil, err
	}

	fmt.Println("Data cleaning completed successfully.")
	return records, nil
}

func filterRecordsAge(records [][]string) [][]string {
	var filteredRecords [][]string

	for _, record := range records {
		val, err := strconv.Atoi(record[0])
		if err != nil {
			fmt.Printf("Error converting column 0 to integer: %v\n", err)
			continue // skipping record if conversion to integer fails
		}

		// limit age between 0 and 120
		if val >= 0 && val <= 120 {
			filteredRecords = append(filteredRecords, record)
		}
	}

	return filteredRecords
}

func removeFirstRow(records [][]string) [][]string {
	if len(records) > 0 {
		return records[1:]
	}
	return records
}

func removeMissing(records [][]string) [][]string {
	cleanRecords := make([][]string, 0)
	for i, record := range records {
		hasMissing := false
		for _, value := range record {
			if value == "" {
				hasMissing = true
				break
			}
		}
		if !hasMissing {
			cleanRecords = append(cleanRecords, record)
		} else {
			fmt.Printf("Row %d has missing values and will be removed.\n", i+1)
		}
	}
	return cleanRecords
}

func mapTextToNumeric(records [][]string) [][]string {
	for _, record := range records {
		switch strings.ToLower(record[3]) {
		case "tertiary":
			record[3] = "3"
		case "secondary":
			record[3] = "2"
		case "primary":
			record[3] = "1"
		default:
			record[3] = "0"
		}

		switch strings.ToLower(record[2]) {
		case "married":
			record[2] = "1"
		case "single":
			record[2] = "2"
		case "divorced":
			record[2] = "3"
		default:
			record[2] = "0"
		}
	}
	return records
}
func standardizeBalance(records [][]string) [][]string {
	sum := 0.0
	for _, record := range records {
		balance, _ := strconv.ParseFloat(record[5], 64)
		sum += balance
	}
	mean := sum / float64(len(records))

	sum = 0.0
	for _, record := range records {
		balance, _ := strconv.ParseFloat(record[5], 64)
		variance := balance - mean
		sum += variance * variance
	}
	stddev := math.Sqrt(sum / float64(len(records)))

	for i, record := range records {
		balance, _ := strconv.ParseFloat(record[5], 64)
		scaled := (balance - mean) / stddev
		records[i][5] = fmt.Sprintf("%f", scaled)
	}
	return records
}

func addJobFeatures(records [][]string) [][]string {
	jobCategories := make(map[string]bool)
	for _, record := range records {
		job := record[1]
		jobCategories[job] = true
	}
	for job := range jobCategories {
		newFeature := make([]string, len(records))
		for i, record := range records {
			if record[1] == job {
				newFeature[i] = "1"
			} else {
				newFeature[i] = "0"
			}
		}
		for i := range records {
			records[i] = append(records[i][:2], append([]string{newFeature[i]}, records[i][2:]...)...)
		}
	}
	return records
}
