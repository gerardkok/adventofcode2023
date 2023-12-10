package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type race struct {
	time     int
	distance int
}

func readInput() ([]string, error) {
	file, err := os.Open("./cmd/day06/testinput")
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

func part1(input []string) int {
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

func part2(input []string) int {
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
	input, err := readInput()
	if err != nil {
		log.Fatalf("can't read input: %v\n", err)
	}

	fmt.Println(part1(input))
	fmt.Println(part2(input))
}
