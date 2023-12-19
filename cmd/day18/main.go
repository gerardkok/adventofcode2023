package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

type Day18 struct {
	day.DayInput
}

func NewDay18(inputFile string) Day18 {
	return Day18{day.DayInput(inputFile)}
}

type action struct {
	direction byte
	steps     int
}

type terrain [][]byte

type turn struct {
	dRow, dColumn int
	corner, edge  byte
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

func makePlan(lines []string) []action {
	result := make([]action, len(lines))
	for i, line := range lines {
		a, _, _ := strings.Cut(line, " (")
		dStr, sStr, _ := strings.Cut(a, " ")
		d := byte(dStr[0])
		s, _ := strconv.Atoi(sStr)
		result[i] = action{d, s}
	}
	return result
}

func dimensions(plan []action) (rows, columns, startRow, startColumn int) {
	var up, right, down, left, maxUp, maxRight, maxDown, maxLeft int

	for _, a := range plan {
		switch a.direction {
		case 'U':
			up += a.steps
			down -= a.steps
			maxUp = max(up, maxUp)
		case 'R':
			right += a.steps
			left -= a.steps
			maxRight = max(right, maxRight)
		case 'D':
			up -= a.steps
			down += a.steps
			maxDown = max(down, maxDown)
		default:
			right -= a.steps
			left += a.steps
			maxLeft = max(left, maxLeft)
		}
	}

	rows, columns = maxUp+maxDown+1, maxRight+maxLeft+1
	startRow, startColumn = maxUp, maxLeft

	return
}

func makeTerrain(rows, columns int) terrain {
	result := make([][]byte, rows)
	for i := range result {
		result[i] = []byte(strings.Repeat(".", columns))
	}
	return result
}

func (t terrain) dig(startRow, startColumn int, plan []action) {
	r := startRow
	c := startColumn
	prevDirection := plan[len(plan)-1].direction
	fmt.Printf("prev direction: %c\n", prevDirection)
	for _, a := range plan {
		turn := turnMap[[2]byte{prevDirection, a.direction}]
		fmt.Printf("turn: %d %d, %c, %c\n", turn.dRow, turn.dColumn, turn.corner, turn.edge)
		prevDirection = a.direction
		t[r][c] = turn.corner
		r, c = r+turn.dRow, c+turn.dColumn
		for s := 1; s < a.steps; s++ {
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

func (t terrain) countDug() int {
	result := 0
	for i := range t {
		for j := range t[i] {
			if t[i][j] != '.' || t.isInside(i, j) {
				// part of the trench
				result++
			}
		}
	}
	return result
}

func (d Day18) Part1() int {
	lines, _ := d.ReadLines()
	plan := makePlan(lines)
	rows, columns, startRow, startColumn := dimensions(plan)
	terrain := makeTerrain(rows, columns)
	terrain.dig(startRow, startColumn, plan)

	for _, line := range terrain {
		fmt.Println(string(line))
	}

	return terrain.countDug()
}

func (d Day18) Part2() int {
	return 0
}

func main() {
	d := NewDay18(filepath.Join(projectpath.Root, "cmd", "day18", "input.txt"))

	fmt.Println(d.Part1())
	// day.Solve(d)
}
