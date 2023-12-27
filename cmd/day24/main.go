package main

import (
	"fmt"
	"path/filepath"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

type Day24 struct {
	day.DayInput
	lower, upper float64
}

type hailstone struct {
	position, velocity map[plane]int
}

type equation struct {
	// linear equation  ax + by = c
	a, b, c int
}

type plane int

const (
	x plane = iota
	y
	z
)

var projectionMap = map[plane][2]plane{
	x: {y, z},
	y: {x, z},
	z: {x, y},
}

func NewDay24(inputFile string, lower, upper float64) Day24 {
	return Day24{day.DayInput(inputFile), lower, upper}
}

func parseLine(line string) hailstone {
	var px, py, pz, vx, vy, vz int
	fmt.Sscanf(line, "%d, %d, %d @ %d, %d, %d", &px, &py, &pz, &vx, &vy, &vz)
	positionMap := map[plane]int{
		x: px,
		y: py,
		z: pz,
	}
	velocityMap := map[plane]int{
		x: vx,
		y: vy,
		z: vz,
	}
	return hailstone{positionMap, velocityMap}
}

func parseLines(lines []string) []hailstone {
	result := make([]hailstone, len(lines))
	for i, line := range lines {
		result[i] = parseLine(line)
	}
	return result
}

func (p plane) String() string {
	switch p {
	case x:
		return "x"
	case y:
		return "y"
	default:
		return "z"
	}
}

func (h hailstone) projection(p plane) equation {
	ax1 := projectionMap[p][0]
	ax2 := projectionMap[p][1]
	a := h.velocity[ax2]
	b := -h.velocity[ax1]
	c := h.velocity[ax2]*h.position[ax1] - h.velocity[ax1]*h.position[ax2]
	return equation{a, b, c}
}

func d(f, g equation) int {
	return g.b*f.a - f.b*g.a
}

func parallel(f, g equation) bool {
	return d(f, g) == 0
}

func (h hailstone) timeAtPosition(pos float64, p plane) float64 {
	return (pos - float64(h.position[p])) / float64(h.velocity[p])
}

func (h hailstone) inPast(xIntersect float64, p plane) bool {
	ax1 := projectionMap[p][0]
	return h.timeAtPosition(xIntersect, ax1) < 0
}

func intersection(a, b hailstone, p plane) (float64, float64) {
	fa := a.projection(p)
	fb := b.projection(p)

	if parallel(fa, fb) {
		return -1, -1
	}

	// a and b intersect at xIntersect, yIntersect
	d := d(fa, fb)
	factorA := float64(fa.c) / float64(d)
	factorB := float64(fb.c) / float64(d)
	xIntersect := float64(fb.b)*factorA - float64(fa.b)*factorB
	yIntersect := float64(fa.a)*factorB - float64(fb.a)*factorA

	if a.inPast(xIntersect, p) || b.inPast(xIntersect, p) {
		return -1, -1
	}

	return xIntersect, yIntersect
}

func (h hailstone) String() string {
	return fmt.Sprintf("%d, %d, %d @ %d, %d, %d", h.position[x], h.position[y], h.position[z], h.velocity[x], h.velocity[y], h.velocity[z])
}

func countIntersections(hailstones []hailstone, lower, upper float64) int {
	result := 0
	for i := range hailstones {
		for j := i + 1; j < len(hailstones); j++ {
			x, y := intersection(hailstones[i], hailstones[j], z)
			if x >= lower && x <= upper && y >= lower && y <= upper {
				result++
			}
		}
	}
	return result
}

func intersectingVelocities(rvx, rvy int, p plane, hailstones []hailstone) (hailstone, bool) {
	type coord struct {
		x, y int
	}
	ax1 := projectionMap[p][0]
	ax2 := projectionMap[p][1]
	intersectionMap := make(map[coord]int)
	for i, a := range hailstones {
		relVelocityA := map[plane]int{
			p:   a.velocity[p],
			ax1: a.velocity[ax1] + rvx,
			ax2: a.velocity[ax2] + rvy,
		}
		relA := hailstone{a.position, relVelocityA}
		for j := i + 1; j < len(hailstones); j++ {
			b := hailstones[j]
			relVelocityB := map[plane]int{
				p:   b.velocity[p],
				ax1: b.velocity[ax1] + rvx,
				ax2: b.velocity[ax2] + rvy,
			}
			relB := hailstone{b.position, relVelocityB}
			x, y := intersection(relA, relB, p)
			if x == -1 && y == -1 {
				continue
			}
			intersectionMap[coord{int(x), int(y)}]++
			if intersectionMap[coord{int(x), int(y)}] > 5 {
				hPosition := map[plane]int{
					ax1: int(x),
					ax2: int(y),
				}
				hVelocity := map[plane]int{
					ax1: -rvx,
					ax2: -rvy,
				}
				return hailstone{hPosition, hVelocity}, true
			} else if len(intersectionMap) > 5 {
				return hailstone{}, false
			}
		}
	}
	return hailstone{}, false
}

func findRockProjection(hailstones []hailstone, p plane) hailstone {
	for rvx := -400; rvx < 400; rvx++ {
		for rvy := -400; rvy < 400; rvy++ {
			result, found := intersectingVelocities(rvx, rvy, p, hailstones)
			if found {
				return result
			}
		}
	}
	return hailstone{}
}

func findRock(hailstones []hailstone) hailstone {
	result := findRockProjection(hailstones, z)
	h := findRockProjection(hailstones, x)
	result.position[z] = h.position[z]
	result.velocity[z] = h.velocity[z]
	return result
}

func (d Day24) Part1() int {
	lines, _ := d.ReadLines()
	hailstones := parseLines(lines)

	return countIntersections(hailstones, d.lower, d.upper)
}

func (d Day24) Part2() int {
	lines, _ := d.ReadLines()
	hailstones := parseLines(lines)

	rock := findRock(hailstones)

	return rock.position[x] + rock.position[y] + rock.position[z]
}

func main() {
	d := NewDay24(filepath.Join(projectpath.Root, "cmd", "day24", "input.txt"), 2e14, 4e14)

	day.Solve(d)
}
