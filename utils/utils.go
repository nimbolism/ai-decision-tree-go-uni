package utils

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

func ReadCSV(filename string) ([][]string, error) {
	csvfile, err := os.Open("data/" + filename)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// reading the file
	r := csv.NewReader(csvfile)
	records, err := r.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return nil, nil
	}
	return records, nil
}

func SaveCSV(records [][]string, filename string) error {
	outputFile, err := os.Create("data/" + filename)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return err
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	for _, record := range records {
		if err := writer.Write(record); err != nil {
			fmt.Println("Error writing record to output file:", err)
			return err
		}
	}
	return nil
}

func SplitData(records [][]string, testRatio float64) ([][]string, [][]string) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	rng.Shuffle(len(records), func(i, j int) {
		records[i], records[j] = records[j], records[i]
	})

	splitIndex := int(float64(len(records)) * (1 - testRatio))

	firstData := records[:splitIndex]
	secondData := records[splitIndex:]

	return firstData, secondData
}
