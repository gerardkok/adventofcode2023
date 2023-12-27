package main

import (
	"path/filepath"
	"strings"

	"golang.org/x/exp/maps"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

type Day25 struct {
	day.DayInput
}

type edge [2]string

type graph struct {
	vertices map[string]struct{}
	edges    map[edge]struct{}
}

type subset struct {
	parent string
	rank   int
}

func NewDay25(inputFile string) Day25 {
	return Day25{day.DayInput(inputFile)}
}

func findRoot(subsets map[string]subset, s string) string {
	if subsets[s].parent == s {
		return s
	}

	a := subsets[s]
	a.parent = findRoot(subsets, subsets[s].parent)
	subsets[s] = a

	return a.parent
}

func union(subsets map[string]subset, x, y string) {
	xroot := findRoot(subsets, x)
	yroot := findRoot(subsets, y)

	a := subsets[xroot]
	b := subsets[yroot]

	switch {
	case subsets[xroot].rank < subsets[yroot].rank:
		a.parent = yroot
	case subsets[xroot].rank > subsets[yroot].rank:
		b.parent = xroot
	default:
		b.parent = xroot
		a.rank++
	}

	subsets[xroot] = a
	subsets[yroot] = b
}

func parseGraph(lines []string) graph {
	result := graph{
		vertices: make(map[string]struct{}),
		edges:    make(map[edge]struct{}),
	}
	for _, line := range lines {
		component, c, _ := strings.Cut(line, ": ")
		result.vertices[component] = struct{}{}
		connections := strings.Split(c, " ")
		for _, connection := range connections {
			result.vertices[connection] = struct{}{}
			result.edges[edge{component, connection}] = struct{}{}
		}
	}
	return result
}

func randomEdge(edges map[edge]struct{}) edge {
	var result edge
	for e := range edges {
		result = e
	}
	delete(edges, result)
	return result
}

func (g graph) minCut(subsets map[string]subset) int {
	result := 0
	for e := range g.edges {
		root0 := findRoot(subsets, e[0])
		root1 := findRoot(subsets, e[1])
		if root0 == root1 {
			continue
		}
		result++
	}
	return result
}

func sizes(subsets map[string]subset) []int {
	result := make(map[string]int)
	for s := range subsets {
		root := findRoot(subsets, s)
		result[root]++
	}

	return maps.Values(result)
}

func (g graph) kargers() (int, int, int) {
	subsets := make(map[string]subset)
	for v := range g.vertices {
		subsets[v] = subset{v, 0}
	}

	n := len(g.vertices)
	edges := make(map[edge]struct{})
	maps.Copy(edges, g.edges)

	for n > 2 {
		e := randomEdge(edges)
		subset0 := findRoot(subsets, e[0])
		subset1 := findRoot(subsets, e[1])

		if subset0 == subset1 {
			continue
		}

		union(subsets, subset0, subset1)
		n--
	}

	cuts := g.minCut(subsets)

	s := sizes(subsets)

	return cuts, s[0], s[1]
}

func (g graph) findCut() int {
	for {
		cut, s0, s1 := g.kargers()
		if cut == 3 {
			return s0 * s1
		}
	}
}

func (d Day25) Part1() int {
	lines, _ := d.ReadLines()
	graph := parseGraph(lines)

	return graph.findCut()
}

func (d Day25) Part2() int {
	return 0
}

func main() {
	d := NewDay25(filepath.Join(projectpath.Root, "cmd", "day25", "input.txt"))

	day.Solve(d)
}
