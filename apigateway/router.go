package apigateway

import (
	"regexp"
	"strings"

	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/strs"
)

type Node struct {
	Tag     string
	Resolve et.Json
	Nodes   []*Node
}

type Nodes struct {
	Routes []*Node
}

type Resolve struct {
	Node    *Node
	Params  []et.Json
	Resolve string
}

// List of routes
var routes *Nodes

// Load routes from file
func load() error {

	return nil
}

// Save routes to file
func save() error {
	// Convertion struct to json
	jsonData, err := et.Marshal(routes)
	if err != nil {
		return err
	}

	console.Log(jsonData.ToString())

	return nil
}

// Create a new node from routes
func newNode(tag string, nodes []*Node) (*Node, []*Node) {
	result := &Node{
		Tag:     tag,
		Resolve: et.Json{},
		Nodes:   []*Node{},
	}

	nodes = append(nodes, result)

	return result, nodes
}

// Find a node from routes
func findNode(tag string, nodes []*Node) *Node {
	for _, node := range nodes {
		if node.Tag == tag {
			return node
		}
	}

	return nil
}

func findResolve(tag string, nodes []*Node, route *Resolve) (*Node, *Resolve) {
	node := findNode(tag, nodes)
	if node == nil {
		// Define regular expression
		regex := regexp.MustCompile(`^\{.*\}$`)
		// Find node by regular expression
		for _, n := range nodes {
			if regex.MatchString(n.Tag) {
				if route == nil {
					route = &Resolve{
						Params: []et.Json{},
					}
				}
				route.Node = n
				route.Params = append(route.Params, et.Json{n.Tag: tag})
				return n, route
			}
		}
	} else if route == nil {
		route = &Resolve{
			Node:   node,
			Params: []et.Json{},
		}
	} else {
		route.Node = node
	}

	return node, route
}

// Add a route to the list
func AddRoute(method, path, resolve string) {
	node := findNode(method, routes.Routes)
	if node == nil {
		node, routes.Routes = newNode(method, routes.Routes)
	}

	tags := strings.Split(path, "/")
	for _, tag := range tags {
		if len(tag) > 0 {
			find := findNode(tag, node.Nodes)
			if find == nil {
				node, node.Nodes = newNode(tag, node.Nodes)
			} else {
				node = find
			}
		}
	}

	if node != nil {
		node.Resolve = et.Json{
			"method":  method,
			"resolve": resolve,
		}

		save()
	}
}

// Get a route from the list
func GetResolve(method, path string) *Resolve {
	node := findNode(method, routes.Routes)
	if node == nil {
		return nil
	}

	var result *Resolve
	tags := strings.Split(path, "/")
	for _, tag := range tags {
		if len(tag) > 0 {
			node, result = findResolve(tag, node.Nodes, result)
			if node == nil {
				return nil
			}
		}
	}

	if result != nil {
		result.Resolve = node.Resolve.Str("resolve")
		for _, param := range result.Params {
			for key, value := range param {
				result.Resolve = strings.Replace(result.Resolve, key, "%v", -1)
				result.Resolve = strs.Format(result.Resolve, value)
			}
		}
	}

	return result
}

// Init routes
func init() {
	routes = &Nodes{
		Routes: []*Node{},
	}

	load()
}
