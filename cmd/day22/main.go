package main

import (
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
	vertical bool
	occupies []coord
}

type zBuffer [][]int

func parseCoord(s string) coord {
	c := strings.Split(s, ",")
	xValue, _ := strconv.Atoi(c[0])
	yValue, _ := strconv.Atoi(c[1])
	zValue, _ := strconv.Atoi(c[2])
	return coord{x: xValue, y: yValue, z: zValue}
}

func newBrick(start, end coord) brick {
	o := orientation(start, end)
	vertical := o == z
	if start[o] > end[o] {
		start, end = end, start
	}
	length := end[o] - start[o] + 1
	occupies := make([]coord, length)

	for i := 0; i < length; i++ {
		occupies[i] = make(coord)
		for a := range start {
			occupies[i][a] = start[a]
			if a == o {
				occupies[i][a] = start[a] + i
			}
		}
	}

	return brick{vertical, occupies}
}

func parseBrick(line string) brick {
	// all bricks seem to be 1x1xh
	ends := strings.Split(line, "~")
	start := parseCoord(ends[0])
	end := parseCoord(ends[1])
	return newBrick(start, end)
}

func parseBricks(lines []string) []brick {
	result := make([]brick, len(lines))
	for i, line := range lines {
		result[i] = parseBrick(line)
	}
	return result
}

func orientation(a, b coord) axis {
	for ax := range a {
		if a[ax] != b[ax] {
			return ax
		}
	}
	return x
}

func (b brick) start() coord {
	return b.occupies[0]
}

func (b brick) end() coord {
	return b.occupies[len(b.occupies)-1]
}

func (b brick) top() []coord {
	if b.vertical {
		return []coord{b.end()}
	}

	return b.occupies
}

func (b brick) lower(to int) {
	for i := range b.occupies {
		b.occupies[i][z] = to
		if b.vertical {
			b.occupies[i][z] = to + i
		}
	}
}

func (b brick) clone() brick {
	occ := make([]coord, len(b.occupies))
	for i, o := range b.occupies {
		occ[i] = coord{x: o[x], y: o[y], z: o[z]}
	}
	return brick{b.vertical, occ}
}

func (zb *zBuffer) maxZ(b brick) int {
	result := 0
	for _, c := range b.occupies {
		result = max(result, (*zb)[c[x]][c[y]])
	}
	return result
}

func compact(bricks []brick) {
	zb := makeZBuffer(maxAxis(bricks, x)+1, maxAxis(bricks, y)+1)
	for i, b := range bricks {
		lowerableTo := zb.maxZ(b) + 1
		bricks[i].lower(lowerableTo)
		for _, c := range bricks[i].top() {
			zb[c[x]][c[y]] = c[z]
		}
	}
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
	localZb := zb.clone()

	for i := layer + 1; i < len(bricks); i++ {
		b := bricks[i].clone()
		lowerableTo := localZb.maxZ(b) + 1
		if lowerableTo < b.start()[z] {
			b.lower(lowerableTo)
			result++
		}

		for _, c := range b.top() {
			localZb[c[x]][c[y]] = c[z]
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
		result = max(result, b.end()[a])
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
		return bricks[i].start()[z] < bricks[j].start()[z]
	})

	compact(bricks)

	return countDisintegratable(bricks)
}

func (d Day22) Part2() int {
	lines, _ := d.ReadLines()
	bricks := parseBricks(lines)

	// sort on z
	sort.Slice(bricks, func(i, j int) bool {
		return bricks[i].start()[z] < bricks[j].start()[z]
	})

	compact(bricks)

	return countFalling(bricks)
}

func main() {
	d := NewDay22(filepath.Join(projectpath.Root, "cmd", "day22", "input.txt"))

	day.Solve(d)
}
