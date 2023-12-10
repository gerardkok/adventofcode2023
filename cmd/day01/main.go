package main

import (
	"maps"
	"path/filepath"
	"strings"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

type Day01 struct {
	day.DayInput
}

var (
	digits = map[string]int{
		"1": 1,
		"2": 2,
		"3": 3,
		"4": 4,
		"5": 5,
		"6": 6,
		"7": 7,
		"8": 8,
		"9": 9,
	}
	words = map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
		"five":  5,
		"six":   6,
		"seven": 7,
		"eight": 8,
		"nine":  9,
	}
)

func NewDay01(inputFile string) Day01 {
	return Day01{day.DayInput(inputFile)}
}

func firstDigit(s string, digitValues map[string]int) int {
	resultIndex := len(s)
	resultValue := 0

	for key, value := range digitValues {
		index := strings.Index(s, key)
		if index != -1 && index < resultIndex {
			resultIndex = index
			resultValue = value
		}
	}

	return resultValue
}

func lastDigit(s string, digitValues map[string]int) int {
	resultIndex := -1
	resultValue := 0

	for key, value := range digitValues {
		index := strings.LastIndex(s, key)
		if index != -1 && index > resultIndex {
			resultIndex = index
			resultValue = value
		}
	}

	return resultValue
}

func (d Day01) Part1() int {
	lines, _ := d.ReadLines()

	sum := 0
	for _, line := range lines {
		first := firstDigit(line, digits)
		last := lastDigit(line, digits)
		combined := 10*first + last
		sum += combined
	}

	return sum
}

func (d Day01) Part2() int {
	lines, _ := d.ReadLines()

	digitsAndWords := digits
	maps.Copy(digitsAndWords, words)
	sum := 0
	for _, line := range lines {
		first := firstDigit(line, digitsAndWords)
		last := lastDigit(line, digitsAndWords)
		combined := 10*first + last
		sum += combined
	}

	return sum
}

func main() {
	d := NewDay01(filepath.Join(projectpath.Root, "cmd", "day01", "input.txt"))

	day.Solve(d)
}
