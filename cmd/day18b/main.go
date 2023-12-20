package main

import (
	"path/filepath"
	"strconv"
	"strings"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

type Day18b struct {
	day.DayInput
}

func NewDay18b(inputFile string) Day18b {
	return Day18b{day.DayInput(inputFile)}
}

type action struct {
	direction byte
	steps     int
}

type plan []action

type turn struct {
	dRow, dColumn int
	corner, edge  byte
}

type coord struct {
	row, column int
}

type polygon struct {
	boundary int
	coords   []coord
}

var directionMap = map[string]byte{
	"U": 'U',
	"R": 'R',
	"D": 'D',
	"L": 'L',
	"0": 'R',
	"1": 'D',
	"2": 'L',
	"3": 'U',
}

var turnMap = map[[2]byte]turn{
	{'U', 'R'}: {0, 1, 'F', '-'},
	{'U', 'L'}: {0, -1, '7', '-'},
	{'R', 'D'}: {1, 0, '7', '|'},
	{'R', 'U'}: {-1, 0, 'J', '|'},
	{'D', 'L'}: {0, -1, 'J', '-'},
	{'D', 'R'}: {0, 1, 'L', '-'},
	{'L', 'U'}: {-1, 0, 'L', '|'},
	{'L', 'D'}: {1, 0, 'F', '|'},
}

func makePlan1(lines []string) plan {
	result := make(plan, len(lines))
	for i, line := range lines {
		fields := strings.Fields(line)
		d := directionMap[fields[0]]
		s, _ := strconv.Atoi(fields[1])
		result[i] = action{
			direction: d,
			steps:     s,
		}
	}
	return result
}

func makePlan2(lines []string) plan {
	result := make(plan, len(lines))
	for i, line := range lines {
		fields := strings.Fields(line)
		d := directionMap[string(fields[2][7])]
		hexSteps, _ := strconv.ParseInt(fields[2][2:7], 16, 64)
		s := int(hexSteps)
		result[i] = action{
			direction: d,
			steps:     s,
		}
	}
	return result
}

func (p plan) makePolygon() polygon {
	coords := make([]coord, len(p))
	boundary := 0
	row, column := 0, 0

	prevDirection := p[len(p)-1].direction
	for i, a := range p {
		turn := turnMap[[2]byte{prevDirection, a.direction}]
		prevDirection = a.direction
		coords[i] = coord{row, column}
		boundary += a.steps
		row, column = row+a.steps*turn.dRow, column+a.steps*turn.dColumn
	}

	return polygon{boundary, coords}
}

func det(a, b coord) int {
	return a.column*b.row - a.row*b.column
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func (p polygon) countInside() int {
	// shoelace formula
	result := 0
	for i, j := 0, 1; j < len(p.coords); i, j = i+1, j+1 {
		result += det(p.coords[i], p.coords[i+1])
	}
	return abs(result) / 2
}

func (p polygon) capacity() int {
	// pick's theorem: A = i + b/2 - 1
	// A = p.countInside(), b = p.boundary, b + i = p.capacity
	// i = A - b/2 + 1 => b + i = A + b/2 + 1
	return p.boundary/2 + p.countInside() + 1
}

func (d Day18b) Part1() int {
	lines, _ := d.ReadLines()
	plan := makePlan1(lines)
	polygon := plan.makePolygon()

	return polygon.capacity()
}

func (d Day18b) Part2() int {
	lines, _ := d.ReadLines()
	plan := makePlan2(lines)
	polygon := plan.makePolygon()

	return polygon.capacity()
}

func main() {
	d := NewDay18b(filepath.Join(projectpath.Root, "cmd", "day18", "input.txt"))

	day.Solve(d)
}
