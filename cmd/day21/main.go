package main

import (
	"path/filepath"
	"strings"

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

var moves = []move{{0, -1}, {-1, 0}, {0, 1}, {1, 0}}

func (g garden) isRock(p plot) bool {
	r := (p.row%len(g.plots) + len(g.plots)) % len(g.plots)
	c := (p.column%len(g.plots[0]) + len(g.plots[0])) % len(g.plots[0])
	return g.plots[r][c] == '#'
}

func (p plot) neighbours() []plot {
	result := make([]plot, 0)
	for _, m := range moves {
		q := plot{p.row + m.dRow, p.column + m.dColumn}
		result = append(result, q)
	}
	return result
}

func (g garden) countReachable(start plot, cycles []int) []int {
	result := make([]int, len(cycles))

	seen := [2]map[plot]struct{}{{}, {}} // seen plots per parity
	startStep := 0
	todo := map[plot]struct{}{start: {}}

	for i := 0; i < len(cycles); i++ {
		var parity int

		for step := startStep; step <= cycles[i]; step++ {
			parity = step % 2
			todoNextStep := make(map[plot]struct{})

			for p := range todo {
				if _, ok := seen[parity][p]; ok {
					continue
				}

				seen[parity][p] = struct{}{}

				for _, q := range p.neighbours() {
					if g.isRock(q) {
						continue
					}
					todoNextStep[q] = struct{}{}
				}
			}

			todo = todoNextStep
		}

		result[i] = len(seen[parity])
		startStep = cycles[i] + 1
	}

	return result
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

func (d Day21) Part1() int {
	lines, _ := d.ReadLines()
	garden := makeGarden(lines)

	return garden.countReachable(garden.start, []int{d.stepsPart1})[0]
}

func lagrangeInterpolation(y0, y1, y2 int) (int, int, int) {
	a := (y2 - (2 * y1) + y0) / 2
	b := y1 - y0 - a
	c := y0
	return a, b, c
}

func (d Day21) Part2() int {
	lines, _ := d.ReadLines()
	garden := makeGarden(lines)

	nCycles := 3
	cycles := make([]int, nCycles)
	for i := 0; i < nCycles; i++ {
		cycles[i] = d.stepsPart2%len(garden.plots) + i*len(garden.plots)
	}

	iterations := garden.countReachable(garden.start, cycles)

	a, b, c := lagrangeInterpolation(iterations[0], iterations[1], iterations[2])
	x := d.stepsPart2 / len(garden.plots)

	return a*x*x + b*x + c
}

func main() {
	d := NewDay21(filepath.Join(projectpath.Root, "cmd", "day21", "input.txt"), 64, 26501365)

	day.Solve(d)
}
