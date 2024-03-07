package apigateway

import "strings"

type Node struct {
	Method  string
	Tag     string
	Resolve string
	Nodes   []*Node
}

var resources []*Node = []*Node{}

func newNode(method, tag string) *Node {
	return &Node{
		Method:  method,
		Tag:     tag,
		Resolve: "",
		Nodes:   []*Node{},
	}
}

func findNode(method, tag string, nodes []*Node) *Node {
	for _, node := range nodes {
		if node.Method == method && node.Tag == tag {
			return node
		}
	}

	return nil
}

func AddNode(method, path, resolve string) {
	var node *Node
	tags := strings.Split(path, "/")

	for _, tag := range tags {
		if node == nil {
			node = findNode(method, tag, resources)
			if node == nil {
				node = newNode(method, tag)
				resources = append(resources, node)
			}
		} else {
			main := *node
			node = findNode(method, tag, main.Nodes)
			if node == nil {
				node = newNode(method, tag)
				main.Nodes = append(main.Nodes, node)
			}
		}
	}

	if node != nil {
		node.Resolve = resolve
	}
}
