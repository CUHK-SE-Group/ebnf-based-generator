package main

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
)

type Node struct {
	Name     string
	Children []*Node
}

func NewNode(name string) *Node {
	return &Node{Name: name}
}

func (n *Node) AddChild(child *Node) {
	n.Children = append(n.Children, child)
}

func DFS(n *Node, depth int, f func(*Node, int)) {
	f(n, depth)
	if depth > 0 {
		for _, child := range n.Children {
			DFS(child, depth-1, f)
		}
	}
}

func buildTree(name string, grammar map[string][]string, depth int) *Node {
	if depth == 0 || !strings.Contains(name, "<") && !strings.Contains(name, ">") {
		return NewNode(name)
	}

	node := NewNode(name)
	for _, prod := range grammar[name] {
		re := regexp.MustCompile(`(<[^>]*>|[^<>]+)`)
		parts := re.FindAllString(prod, -1)
		virtualChild := NewNode(prod)
		for _, symbol := range parts {
			child := buildTree(symbol, grammar, depth-1)
			virtualChild.AddChild(child)
		}
		node.AddChild(virtualChild)
	}

	return node
}

func main() {
	grammar := map[string][]string{
		"<start>":   {"<expr>"},
		"<expr>":    {"<term> + <expr>", "<term> - <expr>", "<term>"},
		"<term>":    {"<term> * <factor>", "<term> / <factor>", "<factor>"},
		"<factor>":  {"+<factor>", "-<factor>", "(<expr>)", "<integer>", "<integer>.<integer>"},
		"<integer>": {"<digit><integer>", "<digit>"},
		"<digit>":   {"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
	}

	root := buildTree("<start>", grammar, 8)

	file, _ := os.Create("tree.dot")
	defer file.Close()

	file.WriteString("digraph G {\n")
	edges := make(map[string]bool)
	DFS(root, 10, func(n *Node, depth int) {
		if len(n.Children) == 0 {
			return
		}
		selectedChild := rand.Int() % len(n.Children)
		child := n.Children[selectedChild]
		//for _, child := range n.Children {
		edge := fmt.Sprintf("\"%s\" -> \"%s\"", n.Name, child.Name)
		if _, exists := edges[edge]; !exists {
			file.WriteString(edge + ";\n")
			edges[edge] = true
		}
		//}
	})
	file.WriteString("}")
}
