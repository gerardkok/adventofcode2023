package main

import (
	"path/filepath"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

type Day11 struct {
	day.DayInput
	expansionPart1, expansionPart2 int
}

type galaxy struct {
	row, column int
}

type space struct {
	galaxies                []galaxy
	rowWidths, columnWidths []int
}

func NewDay11(inputFile string, expansionPart1, expansionPart2 int) Day11 {
	return Day11{day.DayInput(inputFile), expansionPart1, expansionPart2}
}

func parseInput(lines []string, expansion int) space {
	galaxies := make([]galaxy, 0)
	rowWidths := make([]int, len(lines))
	for i := range rowWidths {
		rowWidths[i] = expansion
	}
	columnWidths := make([]int, len(lines[0]))
	for i := range columnWidths {
		columnWidths[i] = expansion
	}

	for row, line := range lines {
		for column, c := range line {
			if c == '.' {
				continue
			}
			g := galaxy{row, column}
			galaxies = append(galaxies, g)
			rowWidths[row] = 1
			columnWidths[column] = 1
		}
	}

	return space{galaxies, rowWidths, columnWidths}
}

func expand(indices []int) []int {
	result := make([]int, len(indices))

	sumPrevious := 0
	for i, index := range indices {
		result[i] = sumPrevious
		sumPrevious += index
	}

	return result
}

func abs(num int) int {
	if num < 0 {
		return -num
	}
	return num
}

func distance(a, b galaxy) int {
	return abs(a.row-b.row) + abs(a.column-b.column)
}

func sumDistances(galaxies []galaxy) int {
	sum := 0
	for i, a := range galaxies {
		for _, b := range galaxies[:i+1] {
			sum += distance(a, b)
		}
	}
	return sum
}

func (s space) expand() []galaxy {
	expandedRows := expand(s.rowWidths)
	expandedColumns := expand(s.columnWidths)

	result := make([]galaxy, len(s.galaxies))

	for i, g := range s.galaxies {
		result[i] = galaxy{expandedRows[g.row], expandedColumns[g.column]}
	}

	return result
}

func (d Day11) Part1() int {
	lines, _ := d.ReadLines()
	space := parseInput(lines, d.expansionPart1)

	return sumDistances(space.expand())
}

func (d Day11) Part2() int {
	lines, _ := d.ReadLines()
	space := parseInput(lines, d.expansionPart2)

	return sumDistances(space.expand())
}

func main() {
	d := NewDay11(filepath.Join(projectpath.Root, "cmd", "day11", "input.txt"), 2, 1000000)

	day.Solve(d)
}
