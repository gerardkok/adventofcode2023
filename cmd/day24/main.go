package main

import (
	"fmt"
	"math"
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

func projectionYZ(h hailstone) equation {
	a := float64(h.vz) / float64(h.vy)
	b := float64(h.pz) - (a * float64(h.py))
	return equation{a, b}
}

func parallel(f, g equation) bool {
	return f.a == g.a
}

func (h hailstone) inPast(x float64) bool {
	return (h.vx > 0 && x < float64(h.px)) || (h.vx < 0 && x > float64(h.px))
}

func (h hailstone) inPastYZ(x float64) bool {
	return (h.vy > 0 && x < float64(h.py)) || (h.vy < 0 && x > float64(h.py))
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
	if math.IsNaN(x) {
		return -1, -1
	}

	if a.inPast(x) || b.inPast(x) {
		return -1, -1
	}

	return x, fa.apply(x)
}

func intersectionYZ(a, b hailstone) (float64, float64) {
	fa := projectionYZ(a)
	fb := projectionYZ(b)

	//fmt.Printf("fa: %v, fb: %v\n", fa, fb)

	if parallel(fa, fb) {
		return -1, -1
	}

	// a and b intersect at x
	x := (fb.b - fa.b) / (fa.a - fb.a)
	if math.IsNaN(x) {
		return -1, -1
	}

	if a.inPastYZ(x) || b.inPastYZ(x) {
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

func intersectingVelocitiesXY(rvx, rvy int, hailstones []hailstone) bool {
	type coord struct {
		x, y int
	}
	intersectionMap := make(map[coord]int)
	for i, a := range hailstones {
		relA := hailstone{a.px, a.py, a.pz, a.vx + rvx, a.vy + rvy, a.vz}
		for j := i + 1; j < len(hailstones); j++ {
			b := hailstones[j]
			relB := hailstone{b.px, b.py, b.pz, b.vx + rvx, b.vy + rvy, b.vz}
			x, y := intersection(relA, relB)
			if x == -1 && y == -1 {
				continue
			}
			//fmt.Printf("intersection: %d, %d\n", int(x), int(y))
			intersectionMap[coord{int(x), int(y)}]++
			if intersectionMap[coord{int(x), int(y)}] > 5 {
				return true
			} else if len(intersectionMap) > 5 {
				return false
			}
		}
	}
	return false
}

func intersectingVelocitiesYZ(rvy, rvz int, hailstones []hailstone) bool {
	type coord struct {
		x, y int
	}
	intersectionMap := make(map[coord]int)
	for i, a := range hailstones {
		relA := hailstone{a.px, a.py, a.pz, a.vx, a.vy + rvy, a.vz + rvz}
		for j := i + 1; j < len(hailstones); j++ {
			b := hailstones[j]
			relB := hailstone{b.px, b.py, b.pz, b.vx, b.vy + rvy, b.vz + rvz}
			x, y := intersectionYZ(relA, relB)
			if x == -1 && y == -1 {
				continue
			}
			//fmt.Printf("intersection: %d, %d\n", int(x), int(y))
			intersectionMap[coord{int(x), int(y)}]++
			if intersectionMap[coord{int(x), int(y)}] > 5 {
				fmt.Printf("x: %d, y: %d\n", int(x), int(y))
				return true
			} else if len(intersectionMap) > 5 {
				return false
			}
		}
	}
	return false
}

func findRock(hailstones []hailstone) hailstone {
	for rvx := -400; rvx < 400; rvx++ {
		for rvy := -400; rvy < 400; rvy++ {
			if intersectingVelocitiesXY(rvx, rvy, hailstones) {
				fmt.Printf("rvx: %d, rvy: %d\n", rvx, rvy)
				//return hailstone{}
			}
		}
	}
	for rvy := -400; rvy < 400; rvy++ {
		for rvz := -400; rvz < 400; rvz++ {
			if intersectingVelocitiesYZ(rvy, rvz, hailstones) {
				fmt.Printf("rvy: %d, rvz: %d\n", rvy, rvz)
				//return hailstone{}
			}
		}
	}
	return hailstone{}
}

func (d Day24) Part1() int {
	lines, _ := d.ReadLines()
	hailstones := parseLines(lines)

	return countIntersections(hailstones, d.lower, d.upper)
}

func (d Day24) Part2() int {
	lines, _ := d.ReadLines()
	hailstones := parseLines(lines)

	_ = findRock(hailstones)

	return 0
}

func main() {
	d := NewDay24(filepath.Join(projectpath.Root, "cmd", "day24", "input.txt"), 2e14, 4e14)

	//day.Solve(d)
	fmt.Println(d.Part2())
}
