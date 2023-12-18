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
	row, column   int
	entrance      direction
	directionKept int
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
	north: {{0, -1, east}, {1, 0, north}, {0, 1, west}},
	east:  {{-1, 0, south}, {0, -1, east}, {1, 0, north}},
	south: {{0, 1, west}, {-1, 0, south}, {0, -1, east}},
	west:  {{1, 0, north}, {0, 1, west}, {-1, 0, south}},
}

type item struct {
	state state
	loss  int
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
		return a.(item).loss - b.(item).loss
	})
	return queue{q}
}

func (q *queue) enqueue(s state, loss int) {
	i := item{s, loss}
	q.Enqueue(i)
}

func (q *queue) dequeue() (state, int) {
	entry, _ := q.Dequeue()
	i := entry.(item)
	return i.state, i.loss
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

func (h heatMap) findMinLoss() int {
	q := newQueue()
	v := make(visited)
	q.enqueue(state{
		row:           1,
		column:        0,
		entrance:      north,
		directionKept: 1,
	}, 0)
	q.enqueue(state{
		row:           0,
		column:        1,
		entrance:      west,
		directionKept: 1,
	}, 0)
	for !q.isEmpty() {
		s, l := q.dequeue()
		// fmt.Printf("considering: [%d, %d] entrace: %s, direction kept: %d, l: %d\n", s.row, s.column, s.entrance.print(), s.directionKept, l)
		loss := l + h[s.row][s.column]
		// fmt.Printf("heat from block: %d, new loss: %d\n", h[s.row][s.column], loss)
		if s.row == len(h)-1 && s.column == len(h[0])-1 {
			return loss
		}

		// if v, ok := visited[s]; ok && v <= loss {
		// 	continue
		// }
		if v.get(s) <= loss {
			continue
		}

		v[s] = loss
		for _, d := range turnMap[s.entrance] {
			nextDirectionKept := 1
			if d.nextEntrance == s.entrance {
				if s.directionKept > 2 {
					continue
				}
				nextDirectionKept = s.directionKept + 1
			}

			r := s.row + d.dRow
			c := s.column + d.dColumn

			if r < 0 || r > len(h)-1 || c < 0 || c > len(h[0])-1 {
				continue
			}

			s := state{
				row:           r,
				column:        c,
				entrance:      d.nextEntrance,
				directionKept: nextDirectionKept,
			}
			q.enqueue(s, loss)
		}
	}
	return -1
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
	return heatMap.findMinLoss()
}

func (d Day17) Part2() int {
	return 0
}

func main() {
	d := NewDay17(filepath.Join(projectpath.Root, "cmd", "day17", "input.txt"))

	day.Solve(d)
}
