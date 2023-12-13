package main

import (
	"path/filepath"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

type Day13 struct {
	day.DayInput
}

func NewDay13(inputFile string) Day13 {
	return Day13{day.DayInput(inputFile)}
}

func readBlocks(lines []string) [][]string {
	result := make([][]string, 0)
	result = append(result, make([]string, 0))

	i := 0
	for _, line := range lines {
		if len(line) == 0 {
			i++
			result = append(result, make([]string, 0))
			continue
		}

		result[i] = append(result[i], line)
	}

	return result
}

func reflect(pattern [][]byte, r int, nSmudges int) int {
	smudgesFound := 0
	for f, b := r, r-1; b >= 0 && f < len(pattern); b, f = b-1, f+1 {
		for i := range pattern[0] {
			if pattern[b][i] != pattern[f][i] {
				smudgesFound++
			}
		}
	}
	return smudgesFound
}

func reflection(pattern [][]byte, nSmudges int) int {
	for i := 1; i < len(pattern); i++ {
		if reflect(pattern, i, nSmudges) == nSmudges {
			return i
		}
	}

	return 0
}

func transpose(slice [][]byte) [][]byte {
	xl := len(slice[0])
	yl := len(slice)
	result := make([][]byte, xl)
	for i := range result {
		result[i] = make([]byte, yl)
	}
	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = slice[j][i]
		}
	}
	return result
}

func sumNotes(lines []string, nSmudges int) int {
	blocks := readBlocks(lines)

	sum := 0
	for _, block := range blocks {
		pattern := make([][]byte, len(block))

		for r, line := range block {
			pattern[r] = []byte(line)
		}

		r := reflection(pattern, nSmudges)
		sum += 100 * r

		transposed := transpose(pattern)
		s := reflection(transposed, nSmudges)
		sum += s
	}

	return sum
}

func (d Day13) Part1() int {
	lines, _ := d.ReadLines()

	return sumNotes(lines, 0)
}

func (d Day13) Part2() int {
	lines, _ := d.ReadLines()

	return sumNotes(lines, 1)
}

func main() {
	d := NewDay13(filepath.Join(projectpath.Root, "cmd", "day13", "input.txt"))

	day.Solve(d)
}
