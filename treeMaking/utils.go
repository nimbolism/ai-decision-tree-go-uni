package treemaking

import (
	"errors"
	"math"
	"sort"
)

// GetBestSplit finds the best attribute and value to split the data.
func GetBestSplit(data [][]string) (int, string, float64, bool, error) {
	// the variables to store the best split information
	maxIG := 0.0
	var bestSplit string
	var bestSplitVar int
	var isNumeric bool

	// looping through each column except the last one - the class label
	for i := 0; i < len(data[0])-1; i++ {
		// values of the current column - attribute
		x := make([]string, len(data))
		for j, row := range data {
			x[j] = row[i]
		}

		// information gain for splitting on the current attribute
		ig, split, numeric, valid, err := MaxInformationGainSplit(x, data, CalculateEntropy)
		if err != nil {
			return 0, "", 0, false, err
		}

		// update the best split information if the current split has higher information gain
		if valid && ig > maxIG {
			maxIG = ig
			bestSplit = split
			bestSplitVar = i
			isNumeric = numeric
		}
	}

	// return zero information gain when no valid split is found
	if maxIG == 0 {
		return 0, "", 0, false, nil
	}
	// return the best attribute
	return bestSplitVar, bestSplit, maxIG, isNumeric, nil
}

func MakeSplit(variable int, value string, data [][]string, isNumeric bool) ([][]string, [][]string, error) {
	// the variables to split the data onto
	var data1 [][]string
	var data2 [][]string

	for _, row := range data {
		val := row[variable]

		// applying the split logic to the data values
		if (isNumeric && val < value) || (!isNumeric && val == value) {
			data1 = append(data1, row)
		} else {
			data2 = append(data2, row)
		}
	}

	return data1, data2, nil
}

func CalculateEntropy(column []string) (float64, error) {
	// first we check wether the column is empty
	if len(column) == 0 {
		return 0, errors.New("input must be a non-empty slice")
	}

	// counting occurrences of each unique value in the column
	valueCounts := make(map[string]float64)
	for _, value := range column {
		valueCounts[value]++
	}

	for key := range valueCounts {
		valueCounts[key] /= float64(len(column))
	}

	// entropy calculation
	var entropy float64
	for _, value := range valueCounts {
		entropy -= value * math.Log2(value+1e-10)
	}

	return entropy, nil
}

// MakePrediction makes a prediction based on the majority class in the data.
func MakePrediction(data [][]string) string {
	// valueCounts is for counting occurrences of each class label
	valueCounts := make(map[string]int)
	for _, row := range data {
		valueCounts[row[len(row)-1]]++
	}

	// class label with the highest count
	maxCount := 0
	var prediction string
	for value, count := range valueCounts {
		if count > maxCount {
			maxCount = count
			prediction = value
		}
	}

	// returning the best class!
	return prediction
}

// MaxInformationGainSplit calculates the maximum information gain and the corresponding split value.
func MaxInformationGainSplit(x []string, data [][]string, lossFunction func([]string) (float64, error)) (float64, string, bool, bool, error) {
	splitValues := make([]string, 0)
	informationGains := make([]float64, 0)

	uniqueValues := make(map[string]bool)
	for _, value := range x {
		uniqueValues[value] = true
	}

	var uniqueValueList []string
	for key := range uniqueValues {
		uniqueValueList = append(uniqueValueList, key)
	}

	sort.Strings(uniqueValueList)

	options := uniqueValueList[1:]

	for _, val := range options {
		mask := make([]bool, len(x))
		for i, value := range x {
			mask[i] = value < val
		}

		// extract the last column of data as y - label class
		y := make([]string, len(data))
		for i := range data {
			y[i] = data[i][len(data[i])-1]
		}

		ig, err := InformationGain(y, mask, lossFunction)
		if err != nil {
			return 0, "", false, false, err
		}

		informationGains = append(informationGains, ig)
		splitValues = append(splitValues, val)
	}

	if len(informationGains) == 0 {
		return 0, "", false, false, nil
	}

	maxIG := informationGains[0]
	maxIGIndex := 0
	for i, ig := range informationGains {
		if ig > maxIG {
			maxIG = ig
			maxIGIndex = i
		}
	}

	return maxIG, splitValues[maxIGIndex], true, true, nil
}

// InformationGain calculates the information gain.
func InformationGain(y []string, mask []bool, lossFunction func([]string) (float64, error)) (float64, error) {
	numPositive := 0
	for _, value := range mask {
		if value {
			numPositive++
		}
	}
	numNegative := len(mask) - numPositive

	if numPositive == 0 || numNegative == 0 {
		return 0, nil
	}

	lossFull, err := lossFunction(y)
	if err != nil {
		return 0, err
	}

	yPositive := make([]string, 0, numPositive)
	yNegative := make([]string, 0, numNegative)
	for i, value := range mask {
		if value {
			yPositive = append(yPositive, y[i])
		} else {
			yNegative = append(yNegative, y[i])
		}
	}

	lossPositive, err := lossFunction(yPositive)
	if err != nil {
		return 0, err
	}

	lossNegative, err := lossFunction(yNegative)
	if err != nil {
		return 0, err
	}

	informationGain := lossFull - (float64(numPositive)/(float64(numPositive)+float64(numNegative)))*lossPositive -
		(float64(numNegative)/(float64(numPositive)+float64(numNegative)))*lossNegative

	return informationGain, nil
}
