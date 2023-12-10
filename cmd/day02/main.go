package main

import (
	"path/filepath"
	"strconv"
	"strings"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

type Day02 struct {
	day.DayInput
}

func NewDay02(inputFile string) Day02 {
	return Day02{day.DayInput(inputFile)}
}

func isGamePossible(game string) bool {
	rounds := strings.Split(game, "; ")
	for _, round := range rounds {
		scoreMap := scoreMap(round)
		if !isRoundPossible(scoreMap) {
			return false
		}
	}
	return true
}

func isRoundPossible(scoreMap map[string]int) bool {
	return scoreMap["red"] <= 12 && scoreMap["green"] <= 13 && scoreMap["blue"] <= 14
}

func power(game string) int {
	rounds := strings.Split(game, "; ")
	required := map[string]int{
		"red":   0,
		"green": 0,
		"blue":  0,
	}

	for _, round := range rounds {
		scoreMap := scoreMap(round)
		for k, v := range scoreMap {
			if v > required[k] {
				required[k] = v
			}
		}
	}

	return required["red"] * required["green"] * required["blue"]
}

func scoreMap(round string) map[string]int {
	scoreMap := make(map[string]int)
	scores := strings.Split(round, ", ")
	for _, s := range scores {
		c, color, _ := strings.Cut(s, " ")
		score, _ := strconv.Atoi(c)
		scoreMap[color] += score
	}
	return scoreMap
}

func (d Day02) Part1() int {
	lines, _ := d.ReadLines()

	sum := 0

	for i, line := range lines {
		_, game, _ := strings.Cut(line, ": ")
		index := i + 1
		if isGamePossible(game) {
			sum += index
		}
	}

	return sum
}

func (d Day02) Part2() int {
	lines, _ := d.ReadLines()

	sum := 0

	for _, line := range lines {
		_, game, _ := strings.Cut(line, ": ")
		power := power(game)
		sum += power
	}

	return sum
}

func main() {
	d := NewDay02(filepath.Join(projectpath.Root, "cmd", "day02", "input.txt"))

	day.Solve(d)
}
