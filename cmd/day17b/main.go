package main

import (
	"fmt"
	"math"
	"path/filepath"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

type Day17b struct {
	day.DayInput
}

func NewDay17b(inputFile string) Day17b {
	return Day17b{day.DayInput(inputFile)}
}

type direction int

type heatMap [][]int

type state struct {
	row, column int
	entrance    direction
}

const (
	north direction = iota
	east
	south
	west
)

type turn struct {
	dRow, dColumn int
	nextEntrance  direction
}

var (
	directions = []direction{north, east, south, west}
	turnMap    = map[direction][]turn{
		north: {{0, -1, east}, {0, 1, west}},
		east:  {{-1, 0, south}, {1, 0, north}},
		south: {{0, 1, west}, {0, -1, east}},
		west:  {{1, 0, north}, {-1, 0, south}},
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

func (q *queue) enqueue(s state, heuristic, cost int) {
	i := (heuristic + cost) % len(*q)
	(*q)[i] = append((*q)[i], node{s, cost})
}

func (network network) aStar(endRow, endColumn, buckets int) int {
	q := make(queue, buckets)
	q.enqueue(state{0, 0, north}, endRow+endColumn, 0)
	q.enqueue(state{0, 0, west}, endRow+endColumn, 0)
	v := make(visited)
	count := 0

	for index := (endRow + endColumn) % len(q); ; index = (index + 1) % len(q) {
		for len(q[index]) > 0 {
			n := q[index][0]
			count++
			q[index] = q[index][1:]
			s, cost := n.state, n.cost

			if s.row == endRow && s.column == endColumn {
				fmt.Printf("count: %d\n", count)
				return cost
			}

			if v.get(s) <= cost {
				continue
			}

			v[s] = cost

			for _, newNode := range network[s] {
				heuristic := endRow - newNode.state.row + endColumn - newNode.state.column
				q.enqueue(newNode.state, heuristic, cost+newNode.cost)
			}
		}
	}
}

func (h heatMap) cumulativeCost(row, column, steps int, t turn) int {
	cost := 0
	for i := 1; i <= steps; i++ {
		cost += h[row+t.dRow*i][column+t.dColumn*i]
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
			r, c := start.row+turn.dRow*s, start.column+turn.dColumn*s
			if h.outside(r, c) {
				// next step will also be outside map, so no need to 'continue'
				break
			}

			end := state{r, c, turn.nextEntrance}

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

func (d Day17b) Part1() int {
	lines, _ := d.ReadLines()
	heatMap := makeHeatMap(lines)
	network := heatMap.makeNetwork(1, 3)
	return network.aStar(len(heatMap)-1, len(heatMap[0])-1, 3*(9+1)+1)
}

func (d Day17b) Part2() int {
	lines, _ := d.ReadLines()
	heatMap := makeHeatMap(lines)
	network := heatMap.makeNetwork(4, 10)
	return network.aStar(len(heatMap)-1, len(heatMap[0])-1, 10*(9+1)+1)
}

func main() {
	d := NewDay17b(filepath.Join(projectpath.Root, "cmd", "day17", "input.txt"))

	day.Solve(d)
}
