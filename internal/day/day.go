package day

import (
	"bufio"
	"fmt"
	"os"
)

type Puzzle interface {
	ReadLines() ([]string, error)
	Part1() int
	Part2() int
}

type Day struct {
	InputFile string
}

// type Option func(*Day) error

// func WithInputFile(file string) Option {
// 	return func(d *Day) error {
// 		d.InputFile = file
// 		return nil
// 	}
// }

// func NewDay(opts ...Option) (Day, error) {
// 	c := Day{}
// 	for _, opt := range opts {
// 		err := opt(&c)
// 		if err != nil {
// 			return Day{}, err
// 		}
// 	}
// 	return c, nil
// }

func (d Day) ReadLines() ([]string, error) {
	file, err := os.Open(d.InputFile)
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

func Solve(p Puzzle) {
	fmt.Println(p.Part1())
	fmt.Println(p.Part2())
}
