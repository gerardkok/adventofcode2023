package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

func readInput() ([]string, error) {
	file, err := os.Open("./cmd/day04/input")
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

func numbers(s string) []int {
	fields := strings.Fields(s)
	result := make([]int, len(fields))
	for i, field := range fields {
		number, _ := strconv.Atoi(field)
		result[i] = number
	}
	return result
}

func countMatches(line string) int {
	_, card, _ := strings.Cut(line, ": ")
	winningNumbersStr, myNumbersStr, _ := strings.Cut(card, " | ")
	winningNumbers := numbers(winningNumbersStr)
	myNumbers := numbers(myNumbersStr)
	count := 0
	for _, myNumber := range myNumbers {
		if slices.Contains(winningNumbers, myNumber) {
			count++
		}
	}
	return count
}

func cardValue(line string) int {
	count := countMatches(line)
	if count == 0 {
		return 0
	}
	return 1 << (count - 1)
}

func part1(lines []string) int {
	sum := 0
	for _, line := range lines {
		value := cardValue(line)
		sum += value
	}
	return sum
}

func part2(lines []string) int {
	matchCount := make([]int, len(lines))
	for i, line := range lines {
		matchCount[i] = countMatches(line)
	}
	copies := make([]int, len(lines))
	for i := range copies {
		copies[i] = 1
	}
	for i, count := range matchCount {
		for j := 0; j < count; j++ {
			copies[i+j+1] += copies[i]
		}
	}
	sum := 0
	for _, copy := range copies {
		sum += copy
	}
	return sum
}

func main() {
	lines, err := readInput()
	if err != nil {
		log.Fatalf("can't read input: %v\n", err)
	}

	part1 := part1(lines)
	fmt.Println(part1)

	part2 := part2(lines)
	fmt.Println(part2)
}
