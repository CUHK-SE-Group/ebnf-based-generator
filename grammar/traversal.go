package grammar

type TraversalAlgo func(node *Node, f func(node *Node))
type Traversal interface {
	func(*Node, func(node *Node))
}

func BFS(node *Node, f func(node *Node)) {
	// We use a slice as a queue in this function
	var queue []*Node

	// Enqueue the root node
	queue = append(queue, node)

	for len(queue) > 0 {
		// Dequeue a node
		node := queue[0]
		queue = queue[1:]

		// Print the node's value (or do something else with it)
		f(node)

		// Enqueue all children
		for _, child := range node.Children {
			queue = append(queue, child)
		}
	}
}

func DFS(node *Node, f func(node *Node)) {
	dfs(node, f)
}

// dfs is a helper function for the DFS method
func dfs(node *Node, f func(node *Node)) {
	// Print the node's value (or do something else with it)
	f(node)

	// Call dfs for all children
	for _, child := range node.Children {
		dfs(child, f)
	}
}
