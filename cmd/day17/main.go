package main

import (
	"math"
	"path/filepath"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

type Day17 struct {
	day.DayInput
}

func NewDay17(inputFile string) Day17 {
	return Day17{day.DayInput(inputFile)}
}

type direction bool

type heatMap [][]int

type state struct {
	row, column int
	entrance    direction
}

const (
	vertical   = direction(true)
	horizontal = direction(false)

	maxCostPerStep = 9
)

type turn [2]int

var (
	directions = []direction{vertical, horizontal}
	turnMap    = map[direction][]turn{
		vertical:   {{0, -1}, {0, 1}},
		horizontal: {{-1, 0}, {1, 0}},
	}
)

type node struct {
	state state
	cost  int
}

type queue [][]node

type visited map[state]int

type network map[state][]node

func (v visited) get(s state) int {
	if v, ok := v[s]; ok {
		return v
	}
	return math.MaxInt
}

func (q *queue) enqueue(s state, cost int) {
	i := cost % len(*q)
	(*q)[i] = append((*q)[i], node{s, cost})
}

func (network network) dijkstra(endRow, endColumn, maxSteps int) int {
	buckets := maxSteps*maxCostPerStep + 1
	q := make(queue, buckets)
	q.enqueue(state{0, 0, vertical}, 0)
	q.enqueue(state{0, 0, horizontal}, 0)
	v := make(visited)

	for index := 0; ; index = (index + 1) % len(q) {
		for len(q[index]) > 0 {
			n := q[index][0]
			q[index] = q[index][1:]
			s, cost := n.state, n.cost

			if s.row == endRow && s.column == endColumn {
				return cost
			}

			if v.get(s) <= cost {
				continue
			}

			v[s] = cost

			for _, newNode := range network[s] {
				q.enqueue(newNode.state, cost+newNode.cost)
			}
		}
	}
}

func (h heatMap) cumulativeCost(row, column, steps int, t turn) int {
	cost := 0
	for i := 1; i <= steps; i++ {
		cost += h[row+t[0]*i][column+t[1]*i]
	}
	return cost
}

func (h heatMap) outside(row, column int) bool {
	return row < 0 || row > len(h)-1 || column < 0 || column > len(h[0])-1
}

func (h heatMap) edges(start state, minSteps, maxSteps int) []node {
	var result []node

	for _, turn := range turnMap[start.entrance] {
		for s := minSteps; s <= maxSteps; s++ {
			r, c := start.row+turn[0]*s, start.column+turn[1]*s
			if h.outside(r, c) {
				// next step will also be outside map, so no need to 'continue'
				break
			}

			end := state{r, c, !start.entrance}

			cost := h.cumulativeCost(start.row, start.column, s, turn)
			result = append(result, node{end, cost})
		}
	}

	return result
}

func (h heatMap) makeNetwork(minSteps, maxSteps int) network {
	result := make(network)
	for r, row := range h {
		for c := range row {
			for _, d := range directions {
				start := state{r, c, d}
				result[start] = h.edges(start, minSteps, maxSteps)
			}
		}
	}
	return result
}

func makeHeatMap(lines []string) heatMap {
	result := make(heatMap, len(lines))
	for i, line := range lines {
		result[i] = make([]int, len(line))
		for j, ch := range line {
			result[i][j] = int(ch - '0')
		}
	}
	return result
}

func (d Day17) Part1() int {
	lines, _ := d.ReadLines()
	heatMap := makeHeatMap(lines)
	minSteps := 1
	maxSteps := 3
	network := heatMap.makeNetwork(minSteps, maxSteps)
	return network.dijkstra(len(heatMap)-1, len(heatMap[0])-1, maxSteps)
}

func (d Day17) Part2() int {
	lines, _ := d.ReadLines()
	heatMap := makeHeatMap(lines)
	minSteps := 4
	maxSteps := 10
	network := heatMap.makeNetwork(minSteps, maxSteps)
	return network.dijkstra(len(heatMap)-1, len(heatMap[0])-1, maxSteps)
}

func main() {
	d := NewDay17(filepath.Join(projectpath.Root, "cmd", "day17", "input.txt"))

	day.Solve(d)
}
