package main

import (
	"path/filepath"
	"strings"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

type Day10 struct {
	day.DayInput
}

func NewDay10(inputFile string) Day10 {
	return Day10{day.DayInput(inputFile)}
}

type Tile struct {
	row, column int
}

func (t Tile) north() Tile {
	return Tile{t.row - 1, t.column}
}

func (t Tile) south() Tile {
	return Tile{t.row + 1, t.column}
}

func (t Tile) east() Tile {
	return Tile{t.row, t.column + 1}
}

func (t Tile) west() Tile {
	return Tile{t.row, t.column - 1}
}

type Diagram [][]byte

func (d Diagram) get(t Tile) byte {
	return d[t.row][t.column]
}

func (d Diagram) set(t Tile, b byte) {
	d[t.row][t.column] = b
}

func (d Diagram) connectsNorth(t Tile) bool {
	b := d.get(t)
	return b == '|' || b == 'L' || b == 'J'
}

func (d Diagram) connectsSouth(t Tile) bool {
	b := d.get(t)
	return b == '|' || b == '7' || b == 'F'
}

func (d Diagram) connectsEast(t Tile) bool {
	b := d.get(t)
	return b == '-' || b == 'L' || b == 'F'
}

func (d Diagram) connectsWest(t Tile) bool {
	b := d.get(t)
	return b == '-' || b == '7' || b == 'J'
}

func (d Diagram) findS() (S, left, right Tile, valueS byte) {
	S, left, right = Tile{}, Tile{}, Tile{}
	for i := range d {
		for j := range d[i] {
			if d[i][j] == 'S' {
				S = Tile{i, j}
			}
		}
	}

	n := S.north()
	s := S.south()
	e := S.east()
	w := S.west()
	switch {
	case d.connectsSouth(n) && d.connectsNorth(s):
		return S, n, s, '|'
	case d.connectsSouth(n) && d.connectsEast(w):
		return S, n, w, 'J'
	case d.connectsSouth(n) && d.connectsWest(e):
		return S, n, e, 'L'
	case d.connectsNorth(s) && d.connectsEast(w):
		return S, s, w, '7'
	case d.connectsNorth(s) && d.connectsWest(e):
		return S, s, e, 'F'
	default:
		// d.connectsEast(w) && d.connectsWest(e)
		return S, w, e, '-'
	}
}

func (d Diagram) nextTile(current, previous Tile) Tile {
	var r Tile
	var l Tile
	switch d.get(current) {
	case '|':
		l = current.north()
		r = current.south()
	case '-':
		l = current.east()
		r = current.west()
	case 'L':
		l = current.north()
		r = current.east()
	case '7':
		l = current.south()
		r = current.west()
	case 'F':
		l = current.south()
		r = current.east()
	default:
		l = current.north()
		r = current.west()
	}

	if l == previous {
		return r
	}
	return l
}

func (d Day10) Part1() int {
	input, _ := d.ReadLines()
	diagram := makeDiagram(input)
	S, left, right, _ := diagram.findS()
	leftPrevious := S
	rightPrevious := S
	distance := 1

	for left != right {
		left, leftPrevious = diagram.nextTile(left, leftPrevious), left
		right, rightPrevious = diagram.nextTile(right, rightPrevious), right
		distance++
	}
	return distance
}

func odd(num int) bool {
	return num%2 == 1
}

func isInside(mainLoop [][]byte, row, column int) bool {
	crossed := 0
	m := min(row, column)
	for i := 0; i < m; i++ {
		v := mainLoop[row-1-i][column-1-i]
		if v == '|' || v == '-' || v == 'J' || v == 'F' {
			crossed++
		}
	}

	return odd(crossed)
}

func countInside(mainLoop [][]byte) int {
	result := 0
	for i := range mainLoop {
		for j := range mainLoop[i] {
			if mainLoop[i][j] != 0 {
				// part of the main loop
				continue
			}

			if isInside(mainLoop, i, j) {
				result++
			}
		}
	}
	return result
}

func makeDiagram(input []string) Diagram {
	result := make(Diagram, len(input)+2)
	result[0] = []byte(strings.Repeat(".", len(input[0])+2))
	result[len(input)+1] = result[0]
	for i, line := range input {
		result[i+1] = []byte("." + line + ".")
	}
	return result
}

func (d Day10) Part2() int {
	input, _ := d.ReadLines()
	diagram := makeDiagram(input)
	mainLoop := make(Diagram, len(diagram))
	for i := range mainLoop {
		mainLoop[i] = make([]byte, len(diagram[0]))
	}

	S, left, right, valueS := diagram.findS()
	mainLoop.set(S, valueS)
	leftPrevious := S
	rightPrevious := S

	for left != right {
		mainLoop.set(left, diagram.get(left))
		mainLoop.set(right, diagram.get(right))
		left, leftPrevious = diagram.nextTile(left, leftPrevious), left
		right, rightPrevious = diagram.nextTile(right, rightPrevious), right
	}
	mainLoop.set(left, diagram.get(left))

	return countInside(mainLoop)
}

func main() {
	d := NewDay10(filepath.Join(projectpath.Root, "cmd", "day10", "input.txt"))

	day.Solve(d)
}
