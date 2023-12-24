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

type zBuffer [][]int

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

func (b brick) top() []coord {
	if b.orientation() == z {
		return []coord{b.end}
	}

	return b.occupies()
}

func (b brick) lower(to int) (brick, []coord) {
	start := coord{x: b.start[x], y: b.start[y], z: to}
	end := coord{x: b.end[x], y: b.end[y], z: to}

	if b.orientation() == z {
		start = coord{x: b.start[x], y: b.start[y], z: to}
		end = coord{x: b.end[x], y: b.end[y], z: to + b.length() - 1}
	}

	n := brick{start, end}
	return n, n.top()
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
		result[i] = newB
		for _, c := range newZ {
			zbuffer[c[x]][c[y]] = c[z]
		}
	}
	return result
}

func makeZBuffer(xMax, yMax int) zBuffer {
	result := make(zBuffer, xMax)
	for i := range result {
		result[i] = make([]int, yMax)
	}
	return result
}

func (zb zBuffer) clone() zBuffer {
	result := make(zBuffer, len(zb))
	for i, zbx := range zb {
		result[i] = make([]int, len(zbx))
		copy(result[i], zb[i])
	}
	return result
}

func (zb *zBuffer) countFalling(bricks []brick, layer int) int {
	result := 0
	localZ := zb.clone()

	for i := layer + 1; i < len(bricks); i++ {
		b := bricks[i]
		highestZ := 0
		for _, c := range b.occupies() {
			highestZ = max(highestZ, localZ[c[x]][c[y]])
		}

		newB, newZ := b.lower(highestZ + 1)

		if !equalBricks(b, newB) {
			result++
		}
		for _, c := range newZ {
			localZ[c[x]][c[y]] = c[z]
		}
	}

	for _, c := range bricks[layer].top() {
		(*zb)[c[x]][c[y]] = c[z]
	}

	return result
}

func (zb *zBuffer) compactable(bricks []brick, layer int) bool {
	return zb.countFalling(bricks, layer) > 0
}

func maxAxis(bricks []brick, a axis) int {
	result := 0
	for _, b := range bricks {
		result = max(result, b.end[a])
	}
	return result
}

func countDisintegratable(bricks []brick) int {
	zb := makeZBuffer(maxAxis(bricks, x)+1, maxAxis(bricks, y)+1)
	result := 0
	for i := 0; i < len(bricks); i++ {
		if !zb.compactable(bricks, i) {
			result++
		}
	}
	return result
}

func countFalling(bricks []brick) int {
	zb := makeZBuffer(maxAxis(bricks, x)+1, maxAxis(bricks, y)+1)
	result := 0
	for i := 0; i < len(bricks); i++ {
		falling := zb.countFalling(bricks, i)
		result += falling
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

	compacted := compact(bricks)

	// sort on z
	sort.Slice(compacted, func(i, j int) bool {
		return compacted[i].start[z] < compacted[j].start[z]
	})

	return countDisintegratable(compacted)
}

func (d Day22) Part2() int {
	lines, _ := d.ReadLines()
	bricks := parseBricks(lines)

	// sort on z
	sort.Slice(bricks, func(i, j int) bool {
		return bricks[i].start[z] < bricks[j].start[z]
	})

	compacted := compact(bricks)

	// sort on z
	sort.Slice(compacted, func(i, j int) bool {
		return compacted[i].start[z] < compacted[j].start[z]
	})

	return countFalling(compacted)
}

func main() {
	d := NewDay22(filepath.Join(projectpath.Root, "cmd", "day22", "input.txt"))

	day.Solve(d)
}
