package main

import (
	"path/filepath"
	"strconv"
	"strings"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

type race struct {
	time     int
	distance int
}

type Day06 struct {
	day.DayInput
}

func NewDay06(inputFile string) Day06 {
	return Day06{day.DayInput(inputFile)}
}

func iSqrt(num int) int {
	result := 0
	for i := 1; i*i <= num; i++ {
		result = i
	}
	return result
}

func winRaceOptions(r race) int {
	// quadratic formula
	D := r.time*r.time - 4*r.distance
	s := iSqrt(D)
	perfectSquare := s*s == D
	if s%2 == r.time%2 {
		// if both are even or both are odd, you can fit one more win in
		s++
	}
	if perfectSquare {
		// tied with record, subtract both ties
		return s - 2
	}
	return s
}

func (day Day06) Part1() int {
	input, _ := day.ReadLines()
	times := strings.Fields(input[0])
	distances := strings.Fields(input[1])
	races := make([]race, len(times)-1)
	for i := 0; i < len(times)-1; i++ {
		t, _ := strconv.Atoi(times[i+1])
		d, _ := strconv.Atoi(distances[i+1])
		races[i] = race{t, d}
	}

	result := 1
	for _, r := range races {
		w := winRaceOptions(r)
		result *= w
	}
	return result
}

func (day Day06) Part2() int {
	input, _ := day.ReadLines()
	times := strings.Fields(input[0])
	distances := strings.Fields(input[1])

	t := ""
	d := ""
	for i := 1; i < len(times); i++ {
		t += times[i]
		d += distances[i]
	}
	time, _ := strconv.Atoi(t)
	distance, _ := strconv.Atoi(d)
	r := race{time, distance}
	return winRaceOptions(r)
}

func main() {
	d := NewDay06(filepath.Join(projectpath.Root, "cmd", "day06", "input.txt"))

	day.Solve(d)
}
