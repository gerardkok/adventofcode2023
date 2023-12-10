package main

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

var nodeRE = regexp.MustCompile(`(\w+)\s*=\s*\((\w+),\s*(\w+)\)`)

type node map[byte]string

type Day08 struct {
	day.DayInput
}

func NewDay08(inputFile string) Day08 {
	return Day08{day.DayInput(inputFile)}
}

func readInput() ([]string, error) {
	file, err := os.Open("./cmd/day08/input")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := make([]string, 0)
	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func parseInput(input []string) (directions string, graph map[string]node) {
	directions = input[0]
	graph = make(map[string]node)

	for i := 2; i < len(input); i++ {
		matches := nodeRE.FindStringSubmatch(input[i])
		name := matches[1]
		left := matches[2]
		right := matches[3]
		n := node{
			'L': left,
			'R': right,
		}
		graph[name] = n
	}

	return directions, graph
}

func (d Day08) Part1() int {
	input, _ := d.ReadLines()
	directions, graph := parseInput(input)

	steps := 0
	current := "AAA"
	for current != "ZZZ" {
		directionIndex := steps % len(directions)
		direction := directions[directionIndex]
		current = graph[current][direction]
		steps++
	}

	return steps
}

func findCycle(start, directions string, graph map[string]node) int {
	steps := 0
	current := start
	for current[2] != 'Z' {
		directionIndex := steps % len(directions)
		direction := directions[directionIndex]
		current = graph[current][direction]
		steps++
	}
	return steps
}

func startState(graph map[string]node) []string {
	result := make([]string, 0)
	for k := range graph {
		if strings.HasSuffix(k, "A") {
			result = append(result, k)
		}
	}
	return result
}

func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}

func (d Day08) Part2() int {
	input, _ := d.ReadLines()
	directions, graph := parseInput(input)

	current := startState(graph)
	steps := make([]int, len(current))
	for i, n := range current {
		steps[i] = findCycle(n, directions, graph)
	}

	return LCM(steps[0], steps[1], steps[2:]...)
}

func main() {
	d := NewDay08(filepath.Join(projectpath.Root, "cmd", "day08", "input.txt"))

	day.Solve(d)
}
