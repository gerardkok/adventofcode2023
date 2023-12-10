package day

import (
	"bufio"
	"fmt"
	"os"
)

type Day interface {
	ReadLines() ([]string, error)
	Part1() int
	Part2() int
}

type DayInput string

func (d DayInput) ReadLines() ([]string, error) {
	file, err := os.Open(string(d))
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

func Solve(p Day) {
	fmt.Println(p.Part1())
	fmt.Println(p.Part2())
}
