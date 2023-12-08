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
	file, err := os.Open("./cmd/puzzle05/input")
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
	// fmt.Printf("checking part: %+v\n", part)
	for i := part.left - 1; i < part.right+1; i++ {
		// fmt.Printf("[%d, %d] = %c\n", part.line-1, i, schema[part.line-1][i])
		if schema[part.line-1][i] != '.' {
			return false
		}
		// fmt.Printf("[%d, %d] = %c\n", part.line+1, i, schema[part.line+1][i])
		if schema[part.line+1][i] != '.' {
			return false
		}
	}
	// fmt.Printf("[%d, %d] = %c\n", part.line, part.left-1, schema[part.line][part.left-1])
	// fmt.Printf("[%d, %d] = %c\n", part.line, part.right, schema[part.line][part.right])
	// fmt.Println("---")
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
				// fmt.Printf("part %+v not surrounded by dots\n", part)
				result = append(result, part)
			}
		}
	}
	return result
}

func main() {
	input, err := readInput()
	if err != nil {
		log.Fatalf("can't read input: %v\n", err)
	}

	schema := makeSchema(input)
	// for _, s := range schema {
	// 	fmt.Println(s)
	// }
	partNumbers := partNumbers(schema)
	sum := 0
	for _, p := range partNumbers {
		sum += p.partNumber
	}
	fmt.Println(sum)
}
