package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var partNumberRE = regexp.MustCompile(`\d+`)

type partNumber struct {
	line, left, right, partNumber int
}

func readInput() ([]string, error) {
	file, err := os.Open("./cmd/day03/input")
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

func makeSchema(input []string) []string {
	result := make([]string, len(input)+2)
	result[0] = strings.Repeat(".", len(input[0])+2)
	result[len(input)+1] = result[0]
	for i, line := range input {
		result[i+1] = "." + line + "."
	}
	return result
}

func surroundedByDots(part partNumber, schema []string) bool {
	for i := part.left - 1; i < part.right+1; i++ {
		if schema[part.line-1][i] != '.' {
			return false
		}
		if schema[part.line+1][i] != '.' {
			return false
		}
	}
	return schema[part.line][part.left-1] == '.' && schema[part.line][part.right] == '.'
}

func partNumbers(schema []string) []partNumber {
	result := make([]partNumber, 0)
	for i, line := range schema {
		matches := partNumberRE.FindAllStringIndex(line, -1)
		if matches == nil {
			continue
		}

		for _, match := range matches {
			number, _ := strconv.Atoi(line[match[0]:match[1]])
			part := partNumber{i, match[0], match[1], number}
			if !surroundedByDots(part, schema) {
				result = append(result, part)
			}
		}
	}
	return result
}

func part1(input []string) int {
	schema := makeSchema(input)
	partNumbers := partNumbers(schema)
	sum := 0
	for _, p := range partNumbers {
		sum += p.partNumber
	}
	return sum
}

func gear(line, col int) string {
	return fmt.Sprintf("%d/%d", line, col)
}

func attachedToGears(part partNumber, schema []string) []string {
	result := make([]string, 0)
	for i := part.left - 1; i < part.right+1; i++ {
		if schema[part.line-1][i] == '*' {
			result = append(result, gear(part.line-1, i))
		}
		if schema[part.line+1][i] == '*' {
			result = append(result, gear(part.line+1, i))
		}
	}
	if schema[part.line][part.left-1] == '*' {
		result = append(result, gear(part.line, part.left-1))
	}
	if schema[part.line][part.right] == '*' {
		result = append(result, gear(part.line, part.right))
	}
	return result
}

func gearMap(parts []partNumber, schema []string) map[string][]partNumber {
	result := make(map[string][]partNumber)
	for _, part := range parts {
		gears := attachedToGears(part, schema)
		for _, gear := range gears {
			if _, ok := result[gear]; !ok {
				result[gear] = make([]partNumber, 0)
			}
			result[gear] = append(result[gear], part)
		}
	}
	return result
}

func part2(input []string) int {
	schema := makeSchema(input)
	partNumbers := partNumbers(schema)
	gearMap := gearMap(partNumbers, schema)
	sum := 0
	for _, gearList := range gearMap {
		if len(gearList) == 2 {
			gear := gearList[0].partNumber * gearList[1].partNumber
			sum += gear
		}
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
