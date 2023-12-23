package main

import (
	"fmt"
	"path/filepath"
	"strings"

	pq "github.com/emirpasic/gods/queues/priorityqueue"
	"golang.org/x/exp/maps"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

type Day21 struct {
	day.DayInput
}

func NewDay21(inputFile string) Day21 {
	return Day21{day.DayInput(inputFile)}
}

type garden struct {
	plots [][]byte
	start node
}

type node struct {
	row, column int
}

type move struct {
	dRow, dColumn int
}

type entry struct {
	node node
	cost int
}

type queue struct {
	*pq.Queue
}

type visited map[node]bool

var neigbourMoves = []move{{0, -1}, {-1, 0}, {0, 1}, {1, 0}}

func newQueue() queue {
	q := pq.NewWith(func(a, b any) int {
		return a.(entry).cost - b.(entry).cost
	})
	return queue{q}
}

func (q *queue) enqueue(n node, cost int) {
	i := entry{n, cost}
	q.Enqueue(i)
}

func (q *queue) dequeue() (node, int) {
	i, _ := q.Dequeue()
	return i.(entry).node, i.(entry).cost
}

func (q *queue) empty() bool {
	return q.Empty()
}

func findReachable(network map[node][]node, start node, limit int) int {
	q := newQueue()
	q.enqueue(start, 0)
	v := make(visited)
	reachable := 0

	for !q.empty() {
		n, cost := q.dequeue()
		if cost > limit {
			return reachable
		}

		if v[n] {
			continue
		}

		v[n] = true
		fmt.Printf("increasing reachable for node: %s\n", n)
		reachable++

		for _, newNode := range network[n] {
			q.enqueue(newNode, cost+2)
		}
	}
	return reachable
}

func makeGarden(lines []string) garden {
	plots := make([][]byte, len(lines))
	var start node
	for r, line := range lines {
		plots[r] = []byte(line)
		c := strings.Index(line, "S")
		if c != -1 {
			start = node{r, c}
		}
	}
	return garden{plots, start}
}

func (g garden) neighbours(n node) map[node]struct{} {
	result := make(map[node]struct{})
	for _, m := range neigbourMoves {
		r, c := n.row+m.dRow, n.column+m.dColumn
		if r < 0 || r > len(g.plots)-1 || c < 0 || c > len(g.plots[0])-1 {
			continue
		}
		if g.plots[r][c] == '#' {
			continue
		}
		neighbour := node{r, c}
		result[neighbour] = struct{}{}
	}
	return result
}

func (g garden) nextNeighbours(n node) map[node]struct{} {
	result := make(map[node]struct{})
	neighbours := g.neighbours(n)
	for m := range neighbours {
		nextNeighbours := g.neighbours(m)
		for k := range nextNeighbours {
			result[k] = struct{}{}
		}
	}
	delete(result, n)
	return result
}

func (g garden) makeNetwork(start node) map[node][]node {
	result := make(map[node][]node, 0)
	alt := ((start.row + start.column) % 2) + 1
	for r, row := range g.plots {
		alt = (alt + 1) % 2
		for c := alt; c < len(row); c += 2 {
			if g.plots[r][c] == '#' {
				continue
			}
			n := node{r, c}
			neighbours := g.nextNeighbours(n)
			result[n] = maps.Keys(neighbours)
		}
	}
	return result
}

func (n node) String() string {
	return fmt.Sprintf("(%d, %d)", n.row, n.column)
}

func (d Day21) Part1() int {
	lines, _ := d.ReadLines()
	garden := makeGarden(lines)
	network := garden.makeNetwork(garden.start)

	for k, v := range network {
		neighbours := make([]string, len(v))
		for i, n := range v {
			neighbours[i] = fmt.Sprint(n)
		}
		fmt.Printf("%s: [%s]\n", k, strings.Join(neighbours, ", "))
	}
	return findReachable(network, garden.start, 64)
}

func (d Day21) Part2() int {
	return 0
}

func main() {
	d := NewDay21(filepath.Join(projectpath.Root, "cmd", "day21", "input.txt"))

	day.Solve(d)
}
