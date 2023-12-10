package main

import (
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

type Day04 struct {
	day.DayInput
}

func NewDay04(inputFile string) Day04 {
	return Day04{day.DayInput(inputFile)}
}

func numbers(s string) []int {
	fields := strings.Fields(s)
	result := make([]int, len(fields))
	for i, field := range fields {
		number, _ := strconv.Atoi(field)
		result[i] = number
	}
	return result
}

func countMatches(line string) int {
	_, card, _ := strings.Cut(line, ": ")
	winningNumbersStr, myNumbersStr, _ := strings.Cut(card, " | ")
	winningNumbers := numbers(winningNumbersStr)
	myNumbers := numbers(myNumbersStr)
	count := 0
	for _, myNumber := range myNumbers {
		if slices.Contains(winningNumbers, myNumber) {
			count++
		}
	}
	return count
}

func cardValue(line string) int {
	count := countMatches(line)
	if count == 0 {
		return 0
	}
	return 1 << (count - 1)
}

func (d Day04) Part1() int {
	lines, _ := d.ReadLines()
	sum := 0
	for _, line := range lines {
		value := cardValue(line)
		sum += value
	}
	return sum
}

func (d Day04) Part2() int {
	lines, _ := d.ReadLines()
	matchCount := make([]int, len(lines))
	for i, line := range lines {
		matchCount[i] = countMatches(line)
	}
	copies := make([]int, len(lines))
	for i := range copies {
		copies[i] = 1
	}
	for i, count := range matchCount {
		for j := 0; j < count; j++ {
			copies[i+j+1] += copies[i]
		}
	}
	sum := 0
	for _, copy := range copies {
		sum += copy
	}
	return sum
}

func main() {
	d := NewDay04(filepath.Join(projectpath.Root, "cmd", "day04", "input.txt"))

	day.Solve(d)
}
