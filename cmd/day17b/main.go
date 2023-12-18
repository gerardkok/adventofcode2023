package main

import (
	"math"
	"path/filepath"

	pq "github.com/emirpasic/gods/queues/priorityqueue"

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

var turnMap = map[direction][]turn{
	north: {{0, -1, east}, {0, -2, east}, {0, -3, east}, {0, 1, west}, {0, 2, west}, {0, 3, west}},
	east:  {{-1, 0, south}, {-2, 0, south}, {-3, 0, south}, {1, 0, north}, {2, 0, north}, {3, 0, north}},
	south: {{0, 1, west}, {0, 2, west}, {0, 3, west}, {0, -1, east}, {0, -2, east}, {0, -3, east}},
	west:  {{1, 0, north}, {2, 0, north}, {3, 0, north}, {-1, 0, south}, {-2, 0, south}, {-3, 0, south}},
}

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
	entry, _ := q.Dequeue()
	i := entry.(node)
	return i.state, i.cost
}

func (q *queue) isEmpty() bool {
	return q.Empty()
}

func (d direction) print() string {
	switch d {
	case north:
		return "north"
	case east:
		return "east"
	case south:
		return "south"
	default:
		return "west"
	}
}

// func (h heatMap) findMinLoss() int {
// 	q := newQueue()
// 	v := make(visited)
// 	q.enqueue(state{
// 		row:           1,
// 		column:        0,
// 		entrance:      north,
// 		directionKept: 1,
// 	}, 0)
// 	q.enqueue(state{
// 		row:           0,
// 		column:        1,
// 		entrance:      west,
// 		directionKept: 1,
// 	}, 0)
// 	for !q.isEmpty() {
// 		s, l := q.dequeue()
// 		// fmt.Printf("considering: [%d, %d] entrace: %s, direction kept: %d, l: %d\n", s.row, s.column, s.entrance.print(), s.directionKept, l)
// 		loss := l + h[s.row][s.column]
// 		// fmt.Printf("heat from block: %d, new loss: %d\n", h[s.row][s.column], loss)
// 		if s.row == len(h)-1 && s.column == len(h[0])-1 {
// 			return loss
// 		}

// 		// if v, ok := visited[s]; ok && v <= loss {
// 		// 	continue
// 		// }
// 		if v.get(s) <= loss {
// 			continue
// 		}

// 		v[s] = loss
// 		for _, d := range turnMap[s.entrance] {
// 			nextDirectionKept := 1
// 			if d.nextEntrance == s.entrance {
// 				if s.directionKept > 2 {
// 					continue
// 				}
// 				nextDirectionKept = s.directionKept + 1
// 			}

// 			r := s.row + d.dRow
// 			c := s.column + d.dColumn

// 			if r < 0 || r > len(h)-1 || c < 0 || c > len(h[0])-1 {
// 				continue
// 			}

// 			s := state{
// 				row:           r,
// 				column:        c,
// 				entrance:      d.nextEntrance,
// 				directionKept: nextDirectionKept,
// 			}
// 			q.enqueue(s, loss)
// 		}
// 	}
// 	return -1
// }

func findMinCost(network map[state][]node, endRow, endColumn int) int {
	q := newQueue()
	v := make(visited)
	q.enqueue(state{
		row:      0,
		column:   0,
		entrance: north,
	}, 0)
	q.enqueue(state{
		row:      0,
		column:   0,
		entrance: west,
	}, 0)

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

func (h heatMap) makeNetwork() map[state][]node {
	result := make(map[state][]node)
	for r, row := range h {
		for c := range row {
			for _, d := range []direction{north, east, south, west} {
				start := state{r, c, d}
				result[start] = make([]node, 0)

				for _, turn := range turnMap[d] {
					if r+turn.dRow < 0 || r+turn.dRow > len(h)-1 || c+turn.dColumn < 0 || c+turn.dColumn > len(h[0])-1 {
						continue
					}
					end := state{r + turn.dRow, c + turn.dColumn, turn.nextEntrance}
					cost := 0
					if turn.dColumn == 0 {
						if turn.dRow > 0 {
							for i := r + 1; i <= r+turn.dRow; i++ {
								cost += h[i][c]
							}
						} else {
							for i := r - 1; i >= r+turn.dRow; i-- {
								cost += h[i][c]
							}
						}
					} else {
						if turn.dColumn > 0 {
							for i := c + 1; i <= c+turn.dColumn; i++ {
								cost += h[r][i]
							}
						} else {
							for i := c - 1; i >= c+turn.dColumn; i-- {
								cost += h[r][i]
							}
						}
					}
					result[start] = append(result[start], node{end, cost})
				}
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
	network := heatMap.makeNetwork()
	// for k, v := range network {
	// 	fmt.Printf("start: %v -> %v\n", k, v)
	// }
	return findMinCost(network, len(heatMap)-1, len(heatMap[0])-1)
}

func (d Day17b) Part2() int {
	return 0
}

func main() {
	d := NewDay17b(filepath.Join(projectpath.Root, "cmd", "day17", "input.txt"))

	day.Solve(d)
}
