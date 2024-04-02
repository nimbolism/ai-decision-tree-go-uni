package treemaking

import (
	"fmt"
	"math/rand"
)

var (
	maxDepth = 8
)

type DecisionTreeNode struct {
	AttributeIndex int
	AttributeValue string
	IsLeaf         bool
	ClassLabel     string
	Children       map[string]*DecisionTreeNode
}

func ConstructDecisionTree(data [][]string) (*DecisionTreeNode, error) {
	fmt.Printf("Length of data BEFORE oversampling: %d \n", len(data))

	data = oversampleRecords(data)

	fmt.Printf("Length of data AFTER oversampling: %d \n", len(data))

	// shuffling the data
	rand.Shuffle(len(data), func(i, j int) {
		data[i], data[j] = data[j], data[i]
	})

	// the recursive function to create the tree
	return constructDecisionTreeRecursive(data, maxDepth)
}

func oversampleRecords(data [][]string) [][]string {
	yesRecords := make([][]string, 0)
	noRecords := make([][]string, 0)

	// separate 'yes' and 'no' records
	for _, record := range data {
		if record[len(record)-1] == "true" {
			yesRecords = append(yesRecords, record)
		} else {
			noRecords = append(noRecords, record)
		}
	}

	// oversampling the 'yes' records to balance the dataset - it seems like the 'no' records are way less
	for len(yesRecords) < len(noRecords) {
		index := rand.Intn(len(yesRecords))
		yesRecords = append(yesRecords, yesRecords[index])
	}

	// combining the oversampled 'yes' records with the 'no' records
	data = append(yesRecords, noRecords...)
	return data
}

// constructDecisionTreeRecursive constructs a decision tree recursively.
func constructDecisionTreeRecursive(data [][]string, depth int) (*DecisionTreeNode, error) {
	// base case for recursive function - when depth is 0 or data is empty and return a leaf node
	if depth == 0 || len(data) == 0 {
		return &DecisionTreeNode{
			IsLeaf:     true,
			ClassLabel: MakePrediction(data), // class label based on majority class
		}, nil
	}

	// calculate the best attribute and value to the data split on
	bestAttribute, bestValue, maxIG, isNumeric, err := GetBestSplit(data)
	if err != nil {
		return nil, err // return error when issue occurs with split calculation
	}
	fmt.Printf("Best attribute: %d, Best value: %s, Maximum information gain: %.2f\n", bestAttribute, bestValue, maxIG)

	// return a leaf node when no valid split is found
	if maxIG == 0 {
		return &DecisionTreeNode{
			IsLeaf:     true,
			ClassLabel: MakePrediction(data), // Class label based on majority class
		}, nil
	}

	// create a new decision tree node with the best attribute and value
	node := &DecisionTreeNode{
		AttributeIndex: bestAttribute,
		AttributeValue: bestValue,
		IsLeaf:         false,
		Children:       make(map[string]*DecisionTreeNode),
	}
	fmt.Printf("Created node with attribute %d and value %s.\n", bestAttribute, bestValue)

	// splitting data into leftData and rightData based on the best attribute and value
	leftData, rightData, err := MakeSplit(bestAttribute, bestValue, data, isNumeric)
	if err != nil {
		return nil, err // return error when there's an issue with data splitting
	}

	// construct the left and right subtrees
	leftSubtree, err := constructDecisionTreeRecursive(leftData, depth-1)
	if err != nil {
		return nil, err // return error when there's an issue with left subtree construction
	}
	rightSubtree, err := constructDecisionTreeRecursive(rightData, depth-1)
	if err != nil {
		return nil, err // return error when there's an issue with right subtree construction
	}

	// assign left and right subtrees to the current node
	node.Children["left"] = leftSubtree
	node.Children["right"] = rightSubtree

	return node, nil
}
