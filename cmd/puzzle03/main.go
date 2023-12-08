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

func isGamePossible(game string) bool {
	rounds := strings.Split(game, "; ")
	for _, round := range rounds {
		if !isRoundPossible(round) {
			return false
		}
	}
	return true
}

func isRoundPossible(round string) bool {
	scoreMap := make(map[string]int)
	scores := strings.Split(round, ", ")
	for _, s := range scores {
		cubeScore := strings.Split(s, " ")
		score, _ := strconv.Atoi(cubeScore[0])
		color := cubeScore[1]
		scoreMap[color] += score
	}
	return isScoreMapPossible(scoreMap)
}

func isScoreMapPossible(scoreMap map[string]int) bool {
	return scoreMap["red"] <= 12 && scoreMap["green"] <= 13 && scoreMap["blue"] <= 14
}

func main() {
	lines, err := readInput()
	if err != nil {
		log.Fatalf("can't read input: %v\n", err)
	}

	sum := 0

	for i, line := range lines {
		parts := strings.Split(line, ": ")
		index := i + 1
		game := parts[1]
		if isGamePossible(game) {
			sum += index
		}
	}
	fmt.Println(sum)
}
