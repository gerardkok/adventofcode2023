package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

type Day22 struct {
	day.DayInput
}

func NewDay22(inputFile string) Day22 {
	return Day22{day.DayInput(inputFile)}
}

type axis int

const (
	x axis = iota
	y
	z
)

type coord map[axis]int

type brick struct {
	start, end coord
}

func parseCoord(s string) coord {
	c := strings.Split(s, ",")
	xValue, _ := strconv.Atoi(c[0])
	yValue, _ := strconv.Atoi(c[1])
	zValue, _ := strconv.Atoi(c[2])
	return coord{x: xValue, y: yValue, z: zValue}
}

func parseBrick(line string) brick {
	// all bricks seem to be 1x1xh
	ends := strings.Split(line, "~")
	start := parseCoord(ends[0])
	end := parseCoord(ends[1])
	if end[x] < start[x] || end[y] < start[y] || end[z] < start[z] {
		start, end = end, start
	}
	return brick{start, end}
}

func parseBricks(lines []string) []brick {
	result := make([]brick, len(lines))
	for i, line := range lines {
		result[i] = parseBrick(line)
	}
	return result
}

func (b brick) orientation() axis {
	for a := range b.start {
		if b.start[a] < b.end[a] {
			return a
		}
	}
	return x
}

func (b brick) length() int {
	o := b.orientation()
	return b.end[o] - b.start[o] + 1
}

func (b brick) occupies() []coord {
	result := make([]coord, b.length())
	orient := b.orientation()
	for i := 0; i < b.length(); i++ {
		result[i] = make(coord)
		for a := range b.start {
			result[i][a] = b.start[a]
			if a == orient {
				result[i][a] = b.start[a] + i
			}
		}
	}
	return result
}

func (c coord) String() string {
	return fmt.Sprintf("%d,%d,%d", c[x], c[y], c[z])
}

func (a axis) String() string {
	switch a {
	case x:
		return "x"
	case y:
		return "y"
	case z:
		return "z"
	default:
		return "unknown"
	}
}

func (b brick) String() string {
	occupies := b.occupies()
	occupiesStr := make([]string, len(occupies))
	for i, o := range occupies {
		occupiesStr[i] = fmt.Sprintf("(%s)", o)
	}
	return fmt.Sprintf("%s~%s, orientation: %s, length: %d, occupies: [%s]", b.start, b.end, b.orientation(), b.length(), strings.Join(occupiesStr, ","))
}

func (b brick) printShort() string {
	return fmt.Sprintf("%s~%s", b.start, b.end)
}

// func occupy(bricks []brick) []coord {

// }

func equalCoords(a, b coord) bool {
	for k, v := range a {
		if v != b[k] {
			return false
		}
	}
	return true
}

func equalBricks(a, b brick) bool {
	return equalCoords(a.start, b.start) && equalCoords(a.end, b.end)
}

func (b brick) lower(to int) (brick, []coord) {
	if b.orientation() == z {
		start := coord{x: b.start[x], y: b.start[y], z: to}
		end := coord{x: b.end[x], y: b.end[y], z: to + b.length() - 1}
		n := brick{start, end}
		return n, []coord{n.end}
	}

	start := coord{x: b.start[x], y: b.start[y], z: to}
	end := coord{x: b.end[x], y: b.end[y], z: to}
	n := brick{start, end}
	return n, n.occupies()
}

func compact(bricks []brick) []brick {
	zbuffer := make([][]int, maxAxis(bricks, x)+1)
	for i := range zbuffer {
		zbuffer[i] = make([]int, maxAxis(bricks, y)+1)
	}

	result := make([]brick, len(bricks))
	for i, b := range bricks {
		highestZ := 0
		for _, c := range b.occupies() {
			highestZ = max(highestZ, zbuffer[c[x]][c[y]])
		}
		newB, newZ := b.lower(highestZ + 1)
		// special := brick{coord{x: 1, y: 5, z: 3}, coord{x: 1, y: 7, z: 3}}
		// if equalBricks(b, special) {
		//fmt.Printf("b: %s\nnewB: %s\nhighest: %d, newZ: %v\n", b, newB, highestZ, newZ)
		//}
		result[i] = newB
		for _, c := range newZ {
			zbuffer[c[x]][c[y]] = c[z]
		}
	}
	return result
}

func compactable(bricks []brick, skipBrick int) bool {
	zbuffer := make([][]int, maxAxis(bricks, x)+1)
	for i := range zbuffer {
		zbuffer[i] = make([]int, maxAxis(bricks, y)+1)
	}

	for i, b := range bricks {
		if i == skipBrick {
			// fmt.Printf("skipping %s\n", bricks[skipBrick])
			continue
		}

		highestZ := 0
		for _, c := range b.occupies() {
			highestZ = max(highestZ, zbuffer[c[x]][c[y]])
		}

		newB, newZ := b.lower(highestZ + 1)

		if !equalBricks(b, newB) {
			// fmt.Printf("brick %s can be lowered, failing support:\n", b)
			// for _, c := range b.occupies() {
			// 	fmt.Printf("[%d,%d]: %d ", c[x], c[y], zbuffer[c[x]][c[y]])
			// }
			// fmt.Println()
			return true
		}
		for _, c := range newZ {
			zbuffer[c[x]][c[y]] = c[z]
		}
	}
	return false
}

func maxAxis(bricks []brick, a axis) int {
	result := 0
	for _, b := range bricks {
		result = max(result, b.end[a])
	}
	return result
}

func countDisintegratable(bricks []brick) int {
	result := 0
	for i := 0; i < len(bricks); i++ {
		if !compactable(bricks, i) {
			result++
		}
	}
	return result
}

func (d Day22) Part1() int {
	lines, _ := d.ReadLines()
	bricks := parseBricks(lines)

	// sort on z
	sort.Slice(bricks, func(i, j int) bool {
		return bricks[i].start[z] < bricks[j].start[z]
	})

	// for _, b := range bricks {
	// 	fmt.Printf("%s\n", b.printShort())
	// }

	// for i, brick := range bricks {
	// 	fmt.Printf("[%3d] %s\n", i, brick)
	// }

	// initialCompactable := compactable(bricks, len(bricks))

	// fmt.Printf("initial setup compactable: %t\n", initialCompactable)

	compacted := compact(bricks)

	// sort on z
	sort.Slice(compacted, func(i, j int) bool {
		return compacted[i].start[z] < compacted[j].start[z]
	})
	// fmt.Println("compacted:")

	// for i, brick := range compacted {
	// 	fmt.Printf("[%3d] %s\n", i, brick)
	// }

	// compactedCompactable := compactable(compacted, len(bricks))

	// fmt.Printf("compacted setup compactable: %t\n", compactedCompactable)

	return countDisintegratable(compacted)
}

func (d Day22) Part2() int {
	return 0
}

func main() {
	d := NewDay22(filepath.Join(projectpath.Root, "cmd", "day22", "input.txt"))

	day.Solve(d)
}
