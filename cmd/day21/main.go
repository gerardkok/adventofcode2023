package main

import (
	"fmt"
	"path/filepath"
	"strings"

	pq "github.com/emirpasic/gods/queues/priorityqueue"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

type Day21 struct {
	day.DayInput
	stepsPart1, stepsPart2 int
}

func NewDay21(inputFile string, stepsPart1, stepsPart2 int) Day21 {
	return Day21{day.DayInput(inputFile), stepsPart1, stepsPart2}
}

type garden struct {
	plots [][]byte
	start plot
}

type plot struct {
	row, column int
}

type move struct {
	dRow, dColumn int
}

type node struct {
	plot     plot
	distance int
}

type queue struct {
	*pq.Queue
}

type visited map[plot]struct{}

var neigbourMoves = []move{{0, -1}, {-1, 0}, {0, 1}, {1, 0}}

func newQueue() queue {
	q := pq.NewWith(func(a, b any) int {
		return a.(node).distance - b.(node).distance
	})
	return queue{q}
}

func (q *queue) enqueue(p plot, distance int) {
	i := node{p, distance}
	q.Enqueue(i)
}

func (q *queue) dequeue() (plot, int) {
	i, _ := q.Dequeue()
	return i.(node).plot, i.(node).distance
}

func (q *queue) empty() bool {
	return q.Empty()
}

func odd(n int) bool {
	return n%2 == 1
}

func (g garden) findReachable(start plot, steps, nCycles int) []int {
	height := len(g.plots)
	cycle := steps % height
	parity := cycle % 2

	q := newQueue()
	q.enqueue(start, 0)
	v := make(visited)
	reachable := 0

	result := make([]int, 0)

	for !q.empty() {
		p, distance := q.dequeue()
		if distance > cycle {
			cycle += height
			r := reachable
			if odd(len(result)) {
				r = len(v) - reachable
			}
			result = append(result, r)
			if len(result) == nCycles {
				return result
			}
		}

		if _, ok := v[p]; ok {
			continue
		}

		v[p] = struct{}{}
		if distance%2 == parity {
			reachable++
		}

		for _, m := range neigbourMoves {
			newPlot := plot{p.row + m.dRow, p.column + m.dColumn}
			r := (newPlot.row%height + height) % height
			c := (newPlot.column%height + height) % height
			if g.plots[r][c] == '#' {
				continue
			}

			q.enqueue(newPlot, distance+1)
		}
	}
	return nil
}

func makeGarden(lines []string) garden {
	plots := make([][]byte, len(lines))
	var start plot
	for r, line := range lines {
		plots[r] = []byte(line)
		c := strings.Index(line, "S")
		if c != -1 {
			start = plot{r, c}
		}
	}
	return garden{plots, start}
}

func (n plot) String() string {
	return fmt.Sprintf("(%d, %d)", n.row, n.column)
}

func (d Day21) Part1() int {
	lines, _ := d.ReadLines()
	garden := makeGarden(lines)

	return garden.findReachable(garden.start, d.stepsPart1, 1)[0]
}

func lagrangeInterpolation(y0, y1, y2 int) (int, int, int) {
	a := y0/2 - y1 + y2/2
	b := -3*(y0/2) + 2*y1 - y2/2
	c := y0
	return a, b, c
}

func (d Day21) Part2() int {
	lines, _ := d.ReadLines()
	garden := makeGarden(lines)

	iterations := garden.findReachable(garden.start, d.stepsPart2, 3)

	a, b, c := lagrangeInterpolation(iterations[2], iterations[1], iterations[0])
	x := d.stepsPart2 / len(garden.plots)

	return a*x*x + b*x + c
}

func main() {
	d := NewDay21(filepath.Join(projectpath.Root, "cmd", "day21", "input.txt"), 64, 26501365)

	day.Solve(d)
}
