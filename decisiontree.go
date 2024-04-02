package main

import (
	"fmt"

	"github.com/nimbolism/decision-tree/post"
	postpruning "github.com/nimbolism/decision-tree/postPruning"
	treemaking "github.com/nimbolism/decision-tree/treeMaking"
	"github.com/nimbolism/decision-tree/utils"
)

func main() {
	// opening and reading csv file
	records, err := utils.ReadCSV("data.csv")
	if err != nil {
		fmt.Println("Error during reading CSV file:", err)
		return
	}

	// calling the function to clean the data - this is custom for this dataset only!
	records, err = utils.DataClean(records)
	if err != nil {
		fmt.Println("Error during data cleaning:", err)
		return
	}

	// spliting data into test-train and spliting trainData to train-validation
	trainData, testData := utils.SplitData(records, 0.2)
	trainData, validationData := utils.SplitData(trainData, 0.2)

	fmt.Printf("Training data size: %d\n", len(trainData))
	fmt.Printf("Testing data size: %d\n", len(testData))

	// tree construction
	root, err := treemaking.ConstructDecisionTree(trainData)
	if err != nil {
		fmt.Println("Error constructing decision tree:", err)
		return
	}

	// pruning based on validation set
	postpruning.Prune(root, validationData)

	// evaluating accuracy to check overfit
	accuracy := post.Evaluate(root, trainData)
	fmt.Printf("Accuracy on train data: %.2f%%\n", accuracy*100)

	accuracy = post.Evaluate(root, validationData)
	fmt.Printf("Accuracy on validation data: %.2f%%\n", accuracy*100)

	accuracy = post.Evaluate(root, testData)
	fmt.Printf("Accuracy on test data: %.2f%%\n", accuracy*100)
}
