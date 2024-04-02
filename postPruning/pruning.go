package postpruning

import (
	"fmt"

	"github.com/nimbolism/decision-tree/post"
	treemaking "github.com/nimbolism/decision-tree/treeMaking"
)

func Prune(root *treemaking.DecisionTreeNode, validationData [][]string) {
	if root == nil || root.IsLeaf {
		return
	}

	for _, child := range root.Children {
		Prune(child, validationData)
	}

	if allChildrenAreLeaves(root) {
		oldAccuracy := post.Evaluate(root, validationData)

		// replace the subtree with a leaf node
		originalChildren := root.Children
		root.Children = make(map[string]*treemaking.DecisionTreeNode)
		root.IsLeaf = true
		root.ClassLabel = getMajorityClass(validationData)

		newAccuracy := post.Evaluate(root, validationData)

		// Revert to the original subtree if pruning doesn't improve accuracy
		if newAccuracy <= oldAccuracy {
			root.IsLeaf = false
			root.Children = originalChildren
		}

		fmt.Printf("Pruning subtree at node %v\n", root.AttributeIndex)
		fmt.Printf("Old Accuracy on validationData: %.2f%%\n", oldAccuracy*100)
		fmt.Printf("New Accuracy on validationData: %.2f%%\n", newAccuracy*100)
		fmt.Println(newAccuracy <= oldAccuracy)
	}
}
func allChildrenAreLeaves(node *treemaking.DecisionTreeNode) bool {
	for _, child := range node.Children {
		if !child.IsLeaf {
			return false
		}
	}
	return true
}

func getMajorityClass(data [][]string) string {
	classCounts := make(map[string]int)

	for _, record := range data {
		classLabel := record[len(record)-1]
		classCounts[classLabel]++
	}

	maxCount := 0
	majorityClass := ""

	for classLabel, count := range classCounts {
		if count > maxCount {
			maxCount = count
			majorityClass = classLabel
		}
	}

	return majorityClass
}
