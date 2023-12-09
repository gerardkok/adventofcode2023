package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func readInput() ([]string, error) {
	file, err := os.Open("./cmd/day09/input")
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

func allZeroes(s []int) bool {
	for _, i := range s {
		if i != 0 {
			return false
		}
	}
	return true
}

func makeTriangle(sequence []int) [][]int {
	result := make([][]int, 0)
	n := 0
	result = append(result, sequence)
	for !allZeroes(result[n]) {
		s := make([]int, len(result[n])-1)
		for i := range s {
			s[i] = result[n][i+1] - result[n][i]
		}
		n++
		result = append(result, s)
	}
	return result
}

func expandTriangleForward(triangle [][]int) [][]int {
	depth := len(triangle)
	result := make([][]int, depth)
	result[depth-1] = make([]int, len(triangle[depth-1])+1)
	for i := range result[depth-1] {
		result[depth-1][i] = 0
	}

	for i := depth - 2; i >= 0; i-- {
		result[i] = make([]int, len(triangle[i])+1)
		_ = copy(result[i], triangle[i])
		result[i][len(triangle[i])] = result[i][len(triangle[i])-1] + result[i+1][len(result[i+1])-1]
	}

	return result
}

func extrapolateForward(sequence []int) int {
	triangle := makeTriangle(sequence)
	expandedTriangle := expandTriangleForward(triangle)
	return expandedTriangle[0][len(expandedTriangle[0])-1]
}

func part1(input []string) int {
	sum := 0
	for _, line := range input {
		fields := strings.Fields(line)
		sequence := make([]int, len(fields))
		for i, field := range fields {
			f, _ := strconv.Atoi(field)
			sequence[i] = f
		}
		nextValue := extrapolateForward(sequence)
		sum += nextValue
	}
	return sum
}

func expandTriangleBackward(triangle [][]int) [][]int {
	depth := len(triangle)
	result := make([][]int, depth)
	result[depth-1] = make([]int, len(triangle[depth-1])+1)
	for i := range result[depth-1] {
		result[depth-1][i] = 0
	}

	for i := depth - 2; i >= 0; i-- {
		result[i] = make([]int, len(triangle[i])+1)
		_ = copy(result[i][1:], triangle[i])
		result[i][0] = result[i][1] - result[i+1][0]
	}

	return result
}

func extrapolateBackward(sequence []int) int {
	triangle := makeTriangle(sequence)
	fmt.Printf("triangle: %v\n", triangle)
	expandedTriangle := expandTriangleBackward(triangle)
	fmt.Printf("expanded triangle: %v\n", expandedTriangle)
	return expandedTriangle[0][0]
}

func part2(input []string) int {
	sum := 0
	for _, line := range input {
		fields := strings.Fields(line)
		sequence := make([]int, len(fields))
		for i, field := range fields {
			f, _ := strconv.Atoi(field)
			sequence[i] = f
		}
		nextValue := extrapolateBackward(sequence)
		sum += nextValue
	}
	return sum
}

func main() {
	input, err := readInput()
	if err != nil {
		log.Fatalf("can't read input: %v\n", err)
	}

	fmt.Println(part1(input))
	fmt.Println(part2(input))
}
