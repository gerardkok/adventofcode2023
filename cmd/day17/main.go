package main

import (
	"math"
	"path/filepath"

	pq "github.com/emirpasic/gods/queues/priorityqueue"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

type Day17 struct {
	day.DayInput
}

func NewDay17(inputFile string) Day17 {
	return Day17{day.DayInput(inputFile)}
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

type queue struct {
	*pq.Queue
}

type visited map[state]int

func (v visited) get(s state) int {
	if v, ok := v[s]; ok {
		return v
	}
	return math.MaxInt
}

func newQueue() queue {
	q := pq.NewWith(func(a, b any) int {
		return a.(node).cost - b.(node).cost
	})
	return queue{q}
}

func (q *queue) enqueue(s state, cost int) {
	i := node{s, cost}
	q.Enqueue(i)
}

func (q *queue) dequeue() (state, int) {
	i, _ := q.Dequeue()
	return i.(node).state, i.(node).cost
}

func findMinCost(network map[state][]node, endRow, endColumn int) int {
	q := newQueue()
	q.enqueue(state{0, 0, north}, 0)
	q.enqueue(state{0, 0, west}, 0)
	v := make(visited)

	for {
		s, cost := q.dequeue()
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

func (h heatMap) edges(start state, minSteps, maxSteps int) []node {
	result := make([]node, 0)

	for s := minSteps; s <= maxSteps; s++ {
		for _, turn := range turnMap[start.entrance] {
			r, c := start.row+turn.dRow*s, start.column+turn.dColumn*s
			if r < 0 || r > len(h)-1 || c < 0 || c > len(h[0])-1 {
				continue
			}

			end := state{r, c, turn.nextEntrance}

			cost := 0
			for i := 1; i <= s; i++ {
				cost += h[start.row+turn.dRow*i][start.column+turn.dColumn*i]
			}
			result = append(result, node{end, cost})
		}
	}

	return result
}

func (h heatMap) makeNetwork(minSteps, maxSteps int) map[state][]node {
	result := make(map[state][]node)
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
	network := heatMap.makeNetwork(1, 3)
	return findMinCost(network, len(heatMap)-1, len(heatMap[0])-1)
}

func (d Day17) Part2() int {
	lines, _ := d.ReadLines()
	heatMap := makeHeatMap(lines)
	network := heatMap.makeNetwork(4, 10)
	return findMinCost(network, len(heatMap)-1, len(heatMap[0])-1)
}

func main() {
	d := NewDay17(filepath.Join(projectpath.Root, "cmd", "day17", "input.txt"))

	day.Solve(d)
}
