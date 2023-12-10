package main

import (
	"path/filepath"
	"strconv"
	"strings"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

type Day09 struct {
	day.DayInput
}

func NewDay09(inputFile string) Day09 {
	return Day09{day.DayInput(inputFile)}
}

func binomialCoefficients(n int) []int {
	result := make([]int, n+1)
	result[0] = 1
	for i := 1; i <= n; i++ {
		result[i] = result[i-1] * (n + 1 - i) / i
	}
	return result
}

func even(n int) bool {
	return n%2 == 0
}

func sign(n int) int {
	if even(n) {
		return 1
	}
	return -1
}

func extrapolate(row []int) int {
	n := len(row)
	sign := sign(n + 1)
	binomialCoefficients := binomialCoefficients(n)
	result := 0
	for i := 0; i < n; i++ {
		r := sign * row[i] * binomialCoefficients[i]
		result += r
		sign = -sign
	}
	return result
}

func reverse(r []int) []int {
	n := len(r)
	result := make([]int, n)
	for i, e := range r {
		result[n-1-i] = e
	}
	return result
}

func (d Day09) Part1() int {
	input, _ := d.ReadLines()
	sum := 0
	for _, line := range input {
		fields := strings.Fields(line)
		sequence := make([]int, len(fields))
		for i, field := range fields {
			f, _ := strconv.Atoi(field)
			sequence[i] = f
		}
		s := extrapolate(sequence)
		sum += s
	}
	return sum
}

func (d Day09) Part2() int {
	input, _ := d.ReadLines()
	sum := 0
	for _, line := range input {
		fields := strings.Fields(line)
		sequence := make([]int, len(fields))
		for i, field := range fields {
			f, _ := strconv.Atoi(field)
			sequence[i] = f
		}
		s := extrapolate(reverse(sequence))
		sum += s
	}
	return sum
}

func main() {
	d := NewDay09(filepath.Join(projectpath.Root, "cmd", "day09", "input"))

	day.Solve(d)
}
