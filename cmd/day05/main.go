package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
)

var almanacRE = regexp.MustCompile(`(?ms)^seeds: (.*)$\s*^seed-to-soil map:\s*([\d\s]+)$\s*^soil-to-fertilizer map:\s*([\d\s]+)$\s*^fertilizer-to-water map:\s*([\d\s]+)$\s*^water-to-light map:\s*([\d\s]+)$\s*^light-to-temperature map:\s*([\d\s]+)$\s*^temperature-to-humidity map:\s*([\d\s]+)$\s*^humidity-to-location map:\s*([\d\s]+)$`)

type numberRange struct {
	destination, source, length int
}

func makeMapping(s string) []numberRange {
	ranges := make([]numberRange, 0)
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		destination, _ := strconv.Atoi(fields[0])
		source, _ := strconv.Atoi(fields[1])
		length, _ := strconv.Atoi(fields[2])
		r := numberRange{destination, source, length}
		ranges = append(ranges, r)
	}
	return ranges
}

func parseInput(input string) (seeds []int, mappings [][]numberRange) {
	matches := almanacRE.FindAllStringSubmatch(string(input), -1)

	s := strings.Fields(matches[0][1])
	seeds = make([]int, len(s))
	for i, m := range s {
		s, _ := strconv.Atoi(m)
		seeds[i] = s
	}

	mappings = make([][]numberRange, 7)
	for i := 0; i < 7; i++ {
		mappings[i] = makeMapping(matches[0][i+2])
	}
	return seeds, mappings
}

func findNext(seed int, mapping []numberRange) int {
	for _, n := range mapping {
		s := seed - n.source
		if s >= 0 && s < n.length {
			return n.destination + s
		}
	}
	return seed
}

func findLocation(seed int, mappings [][]numberRange) int {
	result := seed
	for _, mapping := range mappings {
		result = findNext(result, mapping)
	}
	return result
}

func part1(seeds []int, mappings [][]numberRange) int {
	locations := make([]int, len(seeds))
	for i, seed := range seeds {
		locations[i] = findLocation(seed, mappings)
	}

	location := slices.Min(locations)
	return location
}

func part2(seeds []int, mappings [][]numberRange) int {
	min := math.MaxInt
	for i := 0; i < len(seeds); i += 2 {
		maxSeed := seeds[i] + seeds[i+1]
		for j := seeds[i]; j < maxSeed; j++ {
			location := findLocation(j, mappings)
			if location < min {
				min = location
			}
		}
	}
	return min
}

func part3(seeds []int, mappings [][]numberRange) int {
	var wg sync.WaitGroup

	min := math.MaxInt
	minMtx := &sync.Mutex{}

	wg.Add(len(seeds) / 2)

	for i := 0; i < len(seeds); i += 2 {
		start := seeds[i]
		maxSeed := start + seeds[i+1]
		go func() {
			defer wg.Done()

			for seed := start; seed < maxSeed; seed++ {
				loc := findLocation(seed, mappings)
				if loc < min {
					minMtx.Lock()
					min = loc
					minMtx.Unlock()
				}
			}
		}()
	}

	wg.Wait()

	return min
}

func main() {
	input, err := os.ReadFile("./cmd/day05/input")
	if err != nil {
		log.Fatalf("can't read input: %v\n", err)
	}

	seeds, mappings := parseInput(string(input))

	//fmt.Println(part1(seeds, mappings))
	//fmt.Println(part2(seeds, mappings))
	fmt.Println(part3(seeds, mappings))
}
