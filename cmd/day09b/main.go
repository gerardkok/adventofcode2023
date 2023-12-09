package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

type Puzzle09 struct {
	day.Day
}

type command struct {
	inputFile string
	output    io.Writer
}

type option func(*command) error

func WithInputFile(file string) option {
	return func(c *command) error {
		c.inputFile = file
		return nil
	}
}

func WithOutput(output io.Writer) option {
	return func(c *command) error {
		if output == nil {
			return errors.New("nil output writer")
		}
		c.output = output
		return nil
	}
}

func NewCommand(opts ...option) (command, error) {
	c := command{
		output: os.Stdout,
	}
	for _, opt := range opts {
		err := opt(&c)
		if err != nil {
			return command{}, err
		}
	}
	return c, nil
}

func (c command) readInput() ([]string, error) {
	file, err := os.Open(c.inputFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := make([]string, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func NewPuzzle09(inputFile string) (Puzzle09, error) {
	// d, err := day.NewDay(
	// 	day.WithInputFile(inputFile),
	// )
	// if err != nil {
	// 	return Puzzle09{}, nil
	// }
	d := day.Day{
		InputFile: inputFile,
	}

	return Puzzle09{d}, nil
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

func (p Puzzle09) Part1() int {
	input, _ := p.ReadLines()
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

func (p Puzzle09) Part2() int {
	input, _ := p.ReadLines()
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
	p, err := NewPuzzle09(filepath.Join(projectpath.Root, "cmd", "day09", "input"))
	if err != nil {
		log.Fatalf("can't create command: %v\n", err)
	}

	day.Solve(p)
}
