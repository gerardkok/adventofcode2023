package main

import (
	"fmt"
	"math"
	"path/filepath"
	"strings"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

type Day25b struct {
	day.DayInput
}

type matrix [][]int

type graph struct {
	vertices  map[string]int
	adjacency matrix
}

func NewDay25b(inputFile string) Day25b {
	return Day25b{day.DayInput(inputFile)}
}

func (g *graph) addVertex(v string) int {
	if n, ok := g.vertices[v]; ok {
		return n
	}

	n := len(g.vertices)
	g.vertices[v] = n
	return n
}

func parseGraph(lines []string) graph {
	result := graph{
		vertices: make(map[string]int),
	}
	edges := make([][2]int, 0)
	for _, line := range lines {
		component, c, _ := strings.Cut(line, ": ")
		f := result.addVertex(component)
		connections := strings.Split(c, " ")
		for _, connection := range connections {
			t := result.addVertex(connection)
			edges = append(edges, [2]int{f, t})
		}
	}
	adjacency := make(matrix, len(result.vertices))
	for i := range adjacency {
		adjacency[i] = make([]int, len(result.vertices))
	}
	for _, e := range edges {
		adjacency[e[0]][e[1]] = 1
		adjacency[e[1]][e[0]] = 1
	}
	result.adjacency = adjacency
	return result
}

func (g graph) String() string {
	result := ""
	for v, n := range g.vertices {
		result += fmt.Sprintf("[%s] %d\n", v, n)
	}
	return result
}

func max_element(s []int) int {
	m := math.MinInt
	result := 0
	for i, e := range s {
		if e > m {
			m = e
			result = i
		}
	}
	return result
}

func (m *matrix) globalMinCut() (int, []int) {
	bestCut := math.MaxInt
	bestPartition := make([]int, 0)
	n := len(*m)
	co := make(matrix, n)

	for i := 0; i < n; i++ {
		co[i] = []int{i}
	}

	for ph := 1; ph < n; ph++ {
		w := make([]int, n)
		copy(w, (*m)[0])
		s, t := 0, 0
		for it := 0; it < n-ph; it++ {
			w[t] = math.MinInt
			s, t = t, max_element(w)
			for i := 0; i < n; i++ {
				w[i] += (*m)[t][i]
			}
		}
		if w[t]-(*m)[t][t] < bestCut {
			bestCut = w[t] - (*m)[t][t]
			bestPartition = co[t]
		}

		co[s] = append(co[s], co[t]...)
		for i := 0; i < n; i++ {
			(*m)[s][i] += (*m)[t][i]
		}
		for i := 0; i < n; i++ {
			(*m)[i][s] = (*m)[s][i]
		}
		(*m)[0][t] = math.MinInt
	}

	return bestCut, bestPartition
}

func (d Day25b) Part1() int {
	lines, _ := d.ReadLines()
	graph := parseGraph(lines)

	_, bestPartition := graph.adjacency.globalMinCut()

	return len(bestPartition) * (len(graph.adjacency) - len(bestPartition))
}

func (d Day25b) Part2() int {
	return 0
}

func main() {
	d := NewDay25b(filepath.Join(projectpath.Root, "cmd", "day25", "input.txt"))

	day.Solve(d)
}
