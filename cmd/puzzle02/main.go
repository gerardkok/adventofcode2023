package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)

type digitWord struct {
	word  string
	value int
}

type match struct {
	index, value int
}

var digitWords = []digitWord{
	{"zero", 0},
	{"one", 1},
	{"two", 2},
	{"three", 3},
	{"four", 4},
	{"five", 5},
	{"six", 6},
	{"seven", 7},
	{"eight", 8},
	{"nine", 9},
}

func readInput() ([]string, error) {
	file, err := os.Open("./cmd/puzzle02/input")
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

func firstDigitWord(s string) match {
	result := match{len(s), len(s)}
	for _, digitWord := range digitWords {
		index := strings.Index(s, digitWord.word)
		//		fmt.Printf("%d, %s\n", index, digitWord.word)
		if index != -1 && index < result.index {
			result = match{index, digitWord.value}
		}
	}
	return result
}

func lastDigitWord(s string) match {
	result := match{0, -1}
	for _, digitWord := range digitWords {
		index := strings.LastIndex(s, digitWord.word)
		if index != -1 && index > result.index {
			result = match{index, digitWord.value}
		}
	}
	return result
}

func firstDigitChar(s string) match {
	f := func(c rune) bool {
		return unicode.IsDigit(c)
	}
	index := strings.IndexFunc(s, f)
	if index == -1 {
		return match{len(s), len(s)}
	}
	return match{index, int(s[index]) - '0'}
}

func lastDigitChar(s string) match {
	f := func(c rune) bool {
		return unicode.IsDigit(c)
	}
	index := strings.LastIndexFunc(s, f)
	if index == -1 {
		return match{0, -1}
	}
	return match{index, int(s[index]) - '0'}
}

func firstDigit(s string) int {
	c := firstDigitChar(s)
	w := firstDigitWord(s)
	result := c.value
	if w.index < c.index {
		result = w.value
	}
	//fmt.Printf("c: %+v, d: %+v, result: %d\n", c, w, result)
	return result
}

func lastDigit(s string) int {
	c := lastDigitChar(s)
	w := lastDigitWord(s)
	result := c.value
	if w.index > c.index {
		result = w.value
	}
	return result
}

func main() {
	lines, err := readInput()
	if err != nil {
		log.Fatalf("can't read input: %v\n", err)
	}

	sum := 0
	for _, line := range lines {
		first := firstDigit(line)
		last := lastDigit(line)
		combined := 10*first + last
		sum += combined
		//		fmt.Printf("%s, %d%d = %d\n", line, first, last, combined)
	}
	fmt.Println(sum)
}
