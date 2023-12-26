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

// type equationOldStyle struct {
// 	// y = ax + b
// 	a, b float64
// }

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

// func (h hailstone) projectionOldStyle(p plane) equationOldStyle {
// 	ax1 := projectionMap[p][0]
// 	ax2 := projectionMap[p][1]
// 	a := float64(h.velocity[ax2]) / float64(h.velocity[ax1])
// 	b := float64(h.position[ax2]) - (a * float64(h.position[ax1]))
// 	return equationOldStyle{a, b}
// }

func d(f, g equation) int {
	return g.b*f.a - f.b*g.a
}

func parallel(f, g equation) bool {
	return d(f, g) == 0
}

// func parallelOldStyle(f, g equationOldStyle) bool {
// 	return f.a == g.a
// }

func (h hailstone) positionAtTime(nanosec int) (int, int, int) {
	xPos := h.velocity[x]*nanosec + h.position[x]
	yPos := h.velocity[y]*nanosec + h.position[y]
	zPos := h.velocity[z]*nanosec + h.position[z]
	return xPos, yPos, zPos
}

func (h hailstone) timeAtPosition(pos float64, p plane) float64 {
	return (pos - float64(h.position[p])) / float64(h.velocity[p])
}

func (h hailstone) inPast(xIntersect float64, p plane) bool {
	ax1 := projectionMap[p][0]
	return h.timeAtPosition(xIntersect, ax1) < 0
}

// func (h hailstone) inPastOldStyle(xIntersect float64, p plane) bool {
// 	ax1 := projectionMap[p][0]
// 	return (h.velocity[ax1] > 0 && xIntersect < float64(h.position[ax1])) ||
// 		(h.velocity[ax1] < 0 && xIntersect > float64(h.position[ax1]))
// }

// func (h hailstone) inPastYZ(x float64) bool {
// 	return (h.vy > 0 && x < float64(h.py)) || (h.vy < 0 && x > float64(h.py))
// }

// func (f equation) apply(x float64) float64 {
// 	return f.a*x + f.b
// }

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

// func (f equationOldStyle) apply(x float64) float64 {
// 	return f.a*x + f.b
// }

// func intersectionOldStyle(a, b hailstone, p plane) (float64, float64) {
// 	fa := a.projectionOldStyle(p)
// 	fb := b.projectionOldStyle(p)

// 	if parallelOldStyle(fa, fb) {
// 		fmt.Println("old parallel")
// 		return -1, -1
// 	}

// 	// a and b intersect at x
// 	x := (fb.b - fa.b) / (fa.a - fb.a)

// 	if a.inPastOldStyle(x, p) || b.inPastOldStyle(x, p) {
// 		fmt.Println("old in past")
// 		return -1, -1
// 	}

// 	return x, fa.apply(x)
// }

// func intersectionYZ(a, b hailstone) (float64, float64) {
// 	fa := projectionYZ(a)
// 	fb := projectionYZ(b)

// 	//fmt.Printf("fa: %v, fb: %v\n", fa, fb)

// 	if parallel(fa, fb) {
// 		return -1, -1
// 	}

// 	// a and b intersect at x
// 	x := (fb.b - fa.b) / (fa.a - fb.a)
// 	if math.IsNaN(x) {
// 		return -1, -1
// 	}

// 	if a.inPastYZ(x) || b.inPastYZ(x) {
// 		return -1, -1
// 	}

// 	return x, fa.apply(x)
// }

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

func intersectingVelocities(rvx, rvy int, p plane, hailstones []hailstone) bool {
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
			//fmt.Printf("intersection: %d, %d\n", int(x), int(y))
			intersectionMap[coord{int(x), int(y)}]++
			if intersectionMap[coord{int(x), int(y)}] > 5 {
				fmt.Printf("rvx: %d, x: %f, y: %f\n", -rvx, x, y)
				return true
			} else if len(intersectionMap) > 5 {
				return false
			}
		}
	}
	return false
}

// func intersectingVelocitiesYZ(rvy, rvz int, hailstones []hailstone) bool {
// 	type coord struct {
// 		x, y int
// 	}
// 	intersectionMap := make(map[coord]int)
// 	for i, a := range hailstones {
// 		relA := hailstone{a.px, a.py, a.pz, a.vx, a.vy + rvy, a.vz + rvz}
// 		for j := i + 1; j < len(hailstones); j++ {
// 			b := hailstones[j]
// 			relB := hailstone{b.px, b.py, b.pz, b.vx, b.vy + rvy, b.vz + rvz}
// 			x, y := intersectionYZ(relA, relB)
// 			if x == -1 && y == -1 {
// 				continue
// 			}
// 			//fmt.Printf("intersection: %d, %d\n", int(x), int(y))
// 			intersectionMap[coord{int(x), int(y)}]++
// 			if intersectionMap[coord{int(x), int(y)}] > 5 {
// 				fmt.Printf("x: %d, y: %d\n", int(x), int(y))
// 				return true
// 			} else if len(intersectionMap) > 5 {
// 				return false
// 			}
// 		}
// 	}
// 	return false
// }

func findRock(hailstones []hailstone) hailstone {

	for rvx := -400; rvx < 400; rvx++ {
		for rvy := -400; rvy < 400; rvy++ {
			if intersectingVelocities(rvx, rvy, z, hailstones) {
				fmt.Printf("rvx: %d, rvy: %d\n", rvx, rvy)
				//return hailstone{}
			}
		}
	}
	for rvy := -400; rvy < 400; rvy++ {
		for rvz := -400; rvz < 400; rvz++ {
			if intersectingVelocities(rvy, rvz, x, hailstones) {
				fmt.Printf("rvy: %d, rvz: %d\n", rvy, rvz)
				//return hailstone{}
			}
		}
	}
	for rvx := -400; rvx < 400; rvx++ {
		for rvz := -400; rvz < 400; rvz++ {
			if intersectingVelocities(rvx, rvz, y, hailstones) {
				fmt.Printf("rvx: %d, rvz: %d\n", rvx, rvz)
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
