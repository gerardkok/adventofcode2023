package main

import (
	"path/filepath"
	"strings"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

type Puzzle10 struct {
	day.Day
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

type Diagram []string

func (d Diagram) value(t Tile) byte {
	return d[t.row][t.column]
}

func (d Diagram) connectsNorth(t Tile) bool {
	b := d.value(t)
	return b == '|' || b == 'L' || b == 'J'
}

func (d Diagram) connectsSouth(t Tile) bool {
	b := d.value(t)
	return b == '|' || b == '7' || b == 'F'
}

func (d Diagram) connectsEast(t Tile) bool {
	b := d.value(t)
	return b == '-' || b == 'L' || b == 'F'
}

func (d Diagram) connectsWest(t Tile) bool {
	b := d.value(t)
	return b == '-' || b == '7' || b == 'J'
}

func NewPuzzle10(inputFile string) Puzzle10 {
	d := day.Day{
		InputFile: inputFile,
	}

	return Puzzle10{d}
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
	switch d.value(current) {
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

func (p Puzzle10) Part1() int {
	input, _ := p.ReadLines()
	diagram := Diagram(input)
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
	result := make([]string, len(input)+2)
	result[0] = strings.Repeat(".", len(input[0])+2)
	result[len(input)+1] = result[0]
	for i, line := range input {
		result[i+1] = "." + line + "."
	}
	return result
}

func (p Puzzle10) Part2() int {
	input, _ := p.ReadLines()
	diagram := makeDiagram(input)
	mainLoop := make([][]byte, len(diagram))
	for i := range mainLoop {
		mainLoop[i] = make([]byte, len(diagram[0]))
	}

	S, left, right, valueS := diagram.findS()
	mainLoop[S.row][S.column] = valueS
	leftPrevious := S
	rightPrevious := S

	for left != right {
		mainLoop[left.row][left.column] = diagram[left.row][left.column]
		mainLoop[right.row][right.column] = diagram[right.row][right.column]
		left, leftPrevious = diagram.nextTile(left, leftPrevious), left
		right, rightPrevious = diagram.nextTile(right, rightPrevious), right
	}
	mainLoop[left.row][left.column] = diagram[left.row][left.column]

	return countInside(mainLoop)
}

func main() {
	p := NewPuzzle10(filepath.Join(projectpath.Root, "cmd", "day10", "input.txt"))

	day.Solve(p)
}
