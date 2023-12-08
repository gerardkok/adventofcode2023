package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type handType int

const (
	highCard handType = iota
	onePair
	twoPair
	threeOfAKind
	fullHouse
	fourOfAKind
	fiveOfAKind

	faceOrder1 = "23456789TJQKA"
	faceOrder2 = "J23456789TQKA"
)

var (
	highCardPattern     = []int{1, 1, 1, 1, 1}
	onePairPattern      = []int{1, 1, 1, 2}
	twoPairPattern      = []int{1, 2, 2}
	threeOfAKindPattern = []int{1, 1, 3}
	fullHousePattern    = []int{2, 3}
	fourOfAKindPattern  = []int{1, 4}
	fiveOfAKindPattern  = []int{5}
)

type handBid struct {
	hand     string
	bid      int
	handType handType
}

func readInput() ([]string, error) {
	file, err := os.Open("./cmd/day07/input")
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

func equalPatterns(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func pattern1(hand string) []int {
	cardValues := make(map[rune]int)
	for _, c := range hand {
		cardValues[c]++
	}
	pattern := make([]int, 0)
	for _, v := range cardValues {
		pattern = append(pattern, v)
	}
	sort.Ints(pattern)
	return pattern
}

func pattern2(hand string) []int {
	cardValues := make(map[rune]int)
	for _, c := range hand {
		cardValues[c]++
	}
	pattern := make([]int, 0)
	for k, v := range cardValues {
		if k != 'J' {
			pattern = append(pattern, v)
		}
	}
	sort.Ints(pattern)
	if len(pattern) == 0 {
		// JJJJJ
		return fiveOfAKindPattern
	}
	pattern[len(pattern)-1] += cardValues['J']
	return pattern
}

func getHandType(hand string, patternFn func(string) []int) handType {
	pattern := patternFn(hand)
	switch {
	case equalPatterns(pattern, highCardPattern):
		return highCard
	case equalPatterns(pattern, onePairPattern):
		return onePair
	case equalPatterns(pattern, twoPairPattern):
		return twoPair
	case equalPatterns(pattern, threeOfAKindPattern):
		return threeOfAKind
	case equalPatterns(pattern, fullHousePattern):
		return fullHouse
	case equalPatterns(pattern, fourOfAKindPattern):
		return fourOfAKind
	case equalPatterns(pattern, fiveOfAKindPattern):
		return fiveOfAKind
	}
	return highCard
}

func newHandBid(hand string, bid int, patternFn func(string) []int) handBid {
	handType := getHandType(hand, patternFn)
	return handBid{hand, bid, handType}
}

func part1(input []string) int {
	handBids := make([]handBid, len(input))
	for i, line := range input {
		hand, b, _ := strings.Cut(line, " ")
		bid, _ := strconv.Atoi(b)
		handBids[i] = newHandBid(hand, bid, pattern1)
	}

	sort.Slice(handBids, func(i, j int) bool {
		if handBids[i].handType < handBids[j].handType {
			return true
		}
		if handBids[i].handType > handBids[j].handType {
			return false
		}

		for k := 0; k < 5; k++ {
			faceValueI := strings.IndexByte(faceOrder1, handBids[i].hand[k])
			faceValueJ := strings.IndexByte(faceOrder1, handBids[j].hand[k])
			if faceValueI < faceValueJ {
				return true
			}
			if faceValueI > faceValueJ {
				return false
			}
		}
		return false
	})

	totalWinnings := 0
	for i, h := range handBids {
		w := (i + 1) * h.bid
		totalWinnings += w
	}
	return totalWinnings
}

func part2(input []string) int {
	handBids := make([]handBid, len(input))
	for i, line := range input {
		hand, b, _ := strings.Cut(line, " ")
		bid, _ := strconv.Atoi(b)
		handBids[i] = newHandBid(hand, bid, pattern2)
	}

	sort.Slice(handBids, func(i, j int) bool {
		if handBids[i].handType < handBids[j].handType {
			return true
		}
		if handBids[i].handType > handBids[j].handType {
			return false
		}

		for k := 0; k < 5; k++ {
			faceValueI := strings.IndexByte(faceOrder2, handBids[i].hand[k])
			faceValueJ := strings.IndexByte(faceOrder2, handBids[j].hand[k])
			if faceValueI < faceValueJ {
				return true
			}
			if faceValueI > faceValueJ {
				return false
			}
		}
		return false
	})

	totalWinnings := 0
	for i, h := range handBids {
		w := (i + 1) * h.bid
		totalWinnings += w
	}
	return totalWinnings
}

func main() {
	input, err := readInput()
	if err != nil {
		log.Fatalf("can't read input: %v\n", err)
	}

	fmt.Println(part1(input))
	fmt.Println(part2(input))
}
