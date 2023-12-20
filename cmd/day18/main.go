package main

import (
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"

	"golang.org/x/exp/maps"
)

type Day18 struct {
	day.DayInput
}

func NewDay18(inputFile string) Day18 {
	return Day18{day.DayInput(inputFile)}
}

type action struct {
	direction    byte
	steps        int
	start, compr coord
	corner, edge byte
}

type terrain [][]byte

type turn struct {
	dRow, dColumn int
	corner, edge  byte
}

type coord struct {
	row, column int
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

func makePlan1(lines []string) []action {
	result := make([]action, len(lines))
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

func makePlan2(lines []string) []action {
	result := make([]action, len(lines))
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

func addCoords(plan []action) {
	r := 0
	c := 0
	prevDirection := plan[len(plan)-1].direction
	for i, a := range plan {
		turn := turnMap[[2]byte{prevDirection, a.direction}]
		prevDirection = a.direction
		plan[i].corner, plan[i].edge = turn.corner, turn.edge
		plan[i].start = coord{r, c}
		r, c = r+a.steps*turn.dRow, c+a.steps*turn.dColumn
	}
}

func compressionMap(m map[int]struct{}) (map[int]int, []int) {
	n := maps.Keys(m)
	sort.Ints(n)
	translations := make(map[int]int, len(n))

	for i, j := range n {
		translations[j] = i
	}
	widths := make([]int, len(n)-1)
	for i, j := 0, 1; i < len(widths); i, j = i+1, j+1 {
		widths[i] = abs(n[j] - n[i])
	}
	return translations, widths
}

func addComprCoords(plan []action) ([]int, []int) {
	rows := make(map[int]struct{})
	columns := make(map[int]struct{})
	for _, a := range plan {
		rows[a.start.row] = struct{}{}
		rows[a.start.row+1] = struct{}{}
		columns[a.start.column] = struct{}{}
		columns[a.start.column+1] = struct{}{}
	}

	comprRows, rowWidths := compressionMap(rows)
	comprColumns, columnWidths := compressionMap(columns)

	for i, a := range plan {
		plan[i].compr = coord{comprRows[a.start.row], comprColumns[a.start.column]}
	}

	return rowWidths, columnWidths
}

func makeTerrain(rows, columns int) terrain {
	result := make([][]byte, rows)
	for i := range result {
		result[i] = []byte(strings.Repeat(".", columns))
	}
	return result
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func (t terrain) dig(plan []action) {
	r := plan[0].compr.row
	c := plan[0].compr.column
	prevDirection := plan[len(plan)-1].direction
	for i, j := 0, 1; i < len(plan); i, j = i+1, (j+1)%len(plan) {
		a := plan[i]
		b := plan[j]
		turn := turnMap[[2]byte{prevDirection, a.direction}]
		prevDirection = a.direction
		t[r][c] = turn.corner
		r, c = r+turn.dRow, c+turn.dColumn
		for !(r == b.compr.row && c == b.compr.column) {
			t[r][c] = turn.edge
			r, c = r+turn.dRow, c+turn.dColumn
		}
	}
}

func (t terrain) isInside(row, column int) bool {
	inside := false
	m := min(row, column)
	for i := 0; i < m; i++ {
		v := t[row-1-i][column-1-i]
		if v == '|' || v == '-' || v == 'J' || v == 'F' {
			inside = !inside
		}
	}

	return inside
}

func (t terrain) countDugOut(rows, columns []int) int {
	result := 0
	for i := range t {
		for j := range t[i] {
			if t[i][j] != '.' || t.isInside(i, j) {
				result += rows[i] * columns[j]
			}
		}
	}
	return result
}

func (d Day18) Part1() int {
	lines, _ := d.ReadLines()
	plan := makePlan1(lines)
	addCoords(plan)
	rows, columns := addComprCoords(plan)
	t := makeTerrain(len(rows), len(columns))
	t.dig(plan)

	return t.countDugOut(rows, columns)
}

func (d Day18) Part2() int {
	lines, _ := d.ReadLines()
	plan := makePlan2(lines)
	addCoords(plan)
	rows, columns := addComprCoords(plan)
	t := makeTerrain(len(rows), len(columns))
	t.dig(plan)

	return t.countDugOut(rows, columns)
}

func main() {
	d := NewDay18(filepath.Join(projectpath.Root, "cmd", "day18", "input.txt"))

	day.Solve(d)
}
