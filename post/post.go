package post

import treemaking "github.com/nimbolism/decision-tree/treeMaking"

// Evaluate calculates the accuracy of the decision tree on the given data.
func Evaluate(node *treemaking.DecisionTreeNode, data [][]string) float64 {
	correct := 0

	for _, row := range data {
		prediction := Predict(node, row)
		if prediction == row[len(row)-1] {
			correct++
		}
	}

	return float64(correct) / float64(len(data))
}

// Predict makes a prediction for a single data point using the decision tree.
func Predict(node *treemaking.DecisionTreeNode, row []string) string {
	if node.IsLeaf {
		return node.ClassLabel
	}

	value := row[node.AttributeIndex]
	if value < node.AttributeValue {
		return Predict(node.Children["left"], row)
	} else {
		return Predict(node.Children["right"], row)
	}
}
