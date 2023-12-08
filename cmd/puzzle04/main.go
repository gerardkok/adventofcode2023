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
	file, err := os.Open("./cmd/puzzle03/input")
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

func power(game string) int {
	rounds := strings.Split(game, "; ")
	minimumNeeded := map[string]int{
		"red":   0,
		"green": 0,
		"blue":  0,
	}
	for _, round := range rounds {
		scoreMap := scoreMap(round)
		for k, v := range scoreMap {
			if v > minimumNeeded[k] {
				minimumNeeded[k] = v
			}
		}
	}
	power := 1
	for _, v := range minimumNeeded {
		power *= v
	}
	return power
}

func scoreMap(round string) map[string]int {
	scoreMap := make(map[string]int)
	scores := strings.Split(round, ", ")
	for _, s := range scores {
		cubeScore := strings.Split(s, " ")
		score, _ := strconv.Atoi(cubeScore[0])
		color := cubeScore[1]
		scoreMap[color] += score
	}
	return scoreMap
}

func main() {
	lines, err := readInput()
	if err != nil {
		log.Fatalf("can't read input: %v\n", err)
	}

	sum := 0

	for _, line := range lines {
		parts := strings.Split(line, ": ")
		game := parts[1]
		power := power(game)
		sum += power
	}
	fmt.Println(sum)
}
