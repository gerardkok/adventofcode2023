package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)

func readInput() ([]string, error) {
	file, err := os.Open("./cmd/puzzle01/input")
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

func firstDigit(s string) int {
	f := func(c rune) bool {
		return unicode.IsDigit(c)
	}
	index := strings.IndexFunc(s, f)
	if index == -1 {
		return -1
	}
	return int(s[index]) - '0'
}

func reverse(s string) string {
	chars := []rune(s)
	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}
	return string(chars)
}

func main() {
	lines, err := readInput()
	if err != nil {
		log.Fatalf("can't read input: %v\n", err)
	}

	sum := 0
	for _, line := range lines {
		first := firstDigit(line)
		last := firstDigit(reverse(line))
		combined := 10*first + last
		sum += combined
	}
	fmt.Println(sum)
}
