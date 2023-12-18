package main

import (
	"path/filepath"
	"strings"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

type Day16 struct {
	day.DayInput
}

type direction int

type tile map[direction]struct{}

const (
	north direction = iota
	east
	south
	west
)

type grid [][]byte

type bounce struct {
	dRow, dColumn int
	entrance      direction
}

var bounceMap = map[byte]map[direction][]bounce{
	' ': {}, // edge, beam dims
	'.': {
		north: {{1, 0, north}},
		east:  {{0, -1, east}},
		south: {{-1, 0, south}},
		west:  {{0, 1, west}},
	},
	'/': {
		north: {{0, -1, east}},
		east:  {{1, 0, north}},
		south: {{0, 1, west}},
		west:  {{-1, 0, south}},
	},
	'\\': {
		north: {{0, 1, west}},
		east:  {{-1, 0, south}},
		south: {{0, -1, east}},
		west:  {{1, 0, north}},
	},
	'|': {
		north: {{1, 0, north}},
		east:  {{1, 0, north}, {-1, 0, south}},
		south: {{-1, 0, south}},
		west:  {{1, 0, north}, {-1, 0, south}},
	},
	'-': {
		north: {{0, -1, east}, {0, 1, west}},
		east:  {{0, -1, east}},
		south: {{0, -1, east}, {0, 1, west}},
		west:  {{0, 1, west}},
	},
}

func (t tile) isEnergized() bool {
	return len(t) > 0
}

func (t tile) seen(from direction) bool {
	_, ok := t[from]
	return ok
}

func NewDay16(inputFile string) Day16 {
	return Day16{day.DayInput(inputFile)}
}

func makeGrid(lines []string) grid {
	// add edge where beam dims
	result := make(grid, len(lines)+2)
	result[0] = []byte(strings.Repeat(" ", len(lines[0])+2))
	for i, line := range lines {
		result[i+1] = []byte(" " + line + " ")
	}
	result[len(lines)+1] = []byte(strings.Repeat(" ", len(lines[0])+2))
	return result
}

func makeTiles(height, width int) [][]tile {
	result := make([][]tile, height)
	for i := range result {
		result[i] = make([]tile, width)
		for j := range result[i] {
			result[i][j] = make(map[direction]struct{})
		}
	}
	return result
}

func (g grid) beam(tiles [][]tile, row, column int, entrance direction) {
	if tiles[row][column].seen(entrance) {
		return
	}

	tiles[row][column][entrance] = struct{}{}
	bounceTo := bounceMap[g[row][column]][entrance]
	for _, b := range bounceTo {
		g.beam(tiles, row+b.dRow, column+b.dColumn, b.entrance)
	}
}

func (g grid) countEnergized(row, column int, entrance direction) int {
	tiles := makeTiles(len(g), len(g[0]))

	g.beam(tiles, row, column, entrance)

	result := 0

	// exclude edge
	for i := 1; i < len(tiles)-1; i++ {
		for j := 1; j < len(tiles[i])-1; j++ {
			if tiles[i][j].isEnergized() {
				result++
			}
		}
	}

	return result
}

func (d Day16) Part1() int {
	lines, _ := d.ReadLines()
	grid := makeGrid(lines)

	return grid.countEnergized(1, 1, west)
}

func (d Day16) Part2() int {
	lines, _ := d.ReadLines()
	grid := makeGrid(lines)

	maxEnergized := 0
	for row := 1; row < len(grid)-1; row++ {
		maxEnergized = max(maxEnergized, grid.countEnergized(row, 1, west), grid.countEnergized(row, len(grid[0])-1, east))
	}
	for column := 1; column < len(grid[0])-1; column++ {
		maxEnergized = max(maxEnergized, grid.countEnergized(1, column, north), grid.countEnergized(len(grid)-1, column, south))
	}

	return maxEnergized
}

func main() {
	d := NewDay16(filepath.Join(projectpath.Root, "cmd", "day16", "input.txt"))

	day.Solve(d)
}
