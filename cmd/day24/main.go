package main

import (
	"path/filepath"
	"strconv"
	"strings"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

type Day24 struct {
	day.DayInput
	lower, upper float64
}

type hailstone struct {
	px, py, pz, vx, vy, vz int
}

type equation struct {
	// linear equation of the orthogonal projection on the xy-plane: y = ax + b
	a, b float64
}

func NewDay24(inputFile string, lower, upper float64) Day24 {
	return Day24{day.DayInput(inputFile), lower, upper}
}

func parseLine(line string) hailstone {
	position, velocity, _ := strings.Cut(line, " @ ")
	p := strings.Split(position, ", ")
	v := strings.Split(velocity, ", ")
	px, _ := strconv.Atoi(p[0])
	py, _ := strconv.Atoi(p[1])
	pz, _ := strconv.Atoi(p[2])
	vx, _ := strconv.Atoi(v[0])
	vy, _ := strconv.Atoi(v[1])
	vz, _ := strconv.Atoi(v[2])
	return hailstone{px, py, pz, vx, vy, vz}
}

func parseLines(lines []string) []hailstone {
	result := make([]hailstone, len(lines))
	for i, line := range lines {
		result[i] = parseLine(line)
	}
	return result
}

func projection(h hailstone) equation {
	a := float64(h.vy) / float64(h.vx)
	b := float64(h.py) - (a * float64(h.px))
	return equation{a, b}
}

func parallel(f, g equation) bool {
	return f.a == g.a
}

func (h hailstone) inPast(x float64) bool {
	return (h.vx > 0 && x < float64(h.px)) || (h.vx < 0 && x > float64(h.px))
}

func (f equation) apply(x float64) float64 {
	return f.a*x + f.b
}

func intersection(a, b hailstone) (float64, float64) {
	fa := projection(a)
	fb := projection(b)

	if parallel(fa, fb) {
		return -1, -1
	}

	// a and b intersect at x
	x := (fb.b - fa.b) / (fa.a - fb.a)

	if a.inPast(x) || b.inPast(x) {
		return -1, -1
	}

	return x, fa.apply(x)
}

func countIntersections(hailstones []hailstone, lower, upper float64) int {
	result := 0
	for i := range hailstones {
		for j := i + 1; j < len(hailstones); j++ {
			x, y := intersection(hailstones[i], hailstones[j])
			if x >= lower && x <= upper && y >= lower && y <= upper {
				result++
			}
		}
	}
	return result
}

func (d Day24) Part1() int {
	lines, _ := d.ReadLines()
	hailstones := parseLines(lines)

	return countIntersections(hailstones, d.lower, d.upper)
}

func (d Day24) Part2() int {
	return 0
}

func main() {
	d := NewDay24(filepath.Join(projectpath.Root, "cmd", "day24", "input.txt"), 2e14, 4e14)

	day.Solve(d)
}
