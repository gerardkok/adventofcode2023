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
)

var almanacRE = regexp.MustCompile(`(?ms)^seeds: (.*)$\s*^seed-to-soil map:\s*([\d\s]+)$\s*^soil-to-fertilizer map:\s*([\d\s]+)$\s*^fertilizer-to-water map:\s*([\d\s]+)$\s*^water-to-light map:\s*([\d\s]+)$\s*^light-to-temperature map:\s*([\d\s]+)$\s*^temperature-to-humidity map:\s*([\d\s]+)$\s*^humidity-to-location map:\s*([\d\s]+)$`)

type numberRange struct {
	destination, source, length int
}

func makeMapping(s string) []numberRange {
	// assumes ranges are consecutive
	ranges := make([]numberRange, 0)
	scanner := bufio.NewScanner(strings.NewReader(s))
	starts := make([]int, 0)
	ends := make([]int, 0)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		destination, _ := strconv.Atoi(fields[0])
		source, _ := strconv.Atoi(fields[1])
		starts = append(starts, source)
		length, _ := strconv.Atoi(fields[2])
		ends = append(ends, destination+length)
		r := numberRange{destination, source, length}
		ranges = append(ranges, r)
	}
	minSource := slices.Min(starts)
	if minSource > 0 {
		firstRange := numberRange{0, 0, minSource}
		ranges = append(ranges, firstRange)
	}
	maxSource := slices.Max(ends)
	if maxSource < math.MaxInt {
		lastRange := numberRange{maxSource, maxSource, math.MaxInt - maxSource}
		ranges = append(ranges, lastRange)
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

func mergeRanges(a, b numberRange) numberRange {
	if a.destination >= b.source+b.length || a.destination+a.length <= b.source {
		// no overlap
		return numberRange{}
	}

	resultSource := a.source
	resultDestination := b.destination
	if a.destination < b.source {
		resultSource += (b.source - a.destination)
	} else if a.destination > b.source {
		resultDestination += (a.destination - b.source)
	}

	resultLength := a.source + a.length - resultSource
	if a.destination+a.length > b.source+b.length {
		resultLength -= (a.destination + a.length - b.source - b.length)
	}

	return numberRange{resultDestination, resultSource, resultLength}
}

func mergeMappings(a, b []numberRange) []numberRange {
	result := make([]numberRange, 0)
	for _, m := range a {
		for _, n := range b {
			combined := mergeRanges(m, n)
			if combined.length > 0 {
				result = append(result, combined)
			}
		}
	}
	return result
}

func main() {
	input, err := os.ReadFile("./cmd/day05b/input")
	if err != nil {
		log.Fatalf("can't read input: %v\n", err)
	}

	seeds, mappings := parseInput(string(input))
	for i, m := range mappings {
		fmt.Printf("[%d] %v\n", i, m)
	}

	// fmt.Println(part1(seeds, mappings))
	// fmt.Println(part2(seeds, mappings))
	seedMappings := make([]numberRange, len(seeds)/2)
	for i := 0; i < len(seeds); i += 2 {
		seedMappings[i/2] = numberRange{seeds[i], seeds[i], seeds[i+1]}
	}

	// result := seedMappings
	// for i := range mappings {
	// 	combinedMappings := make([]numberRange, 0)
	// 	for _, combinedMapping := range result {
	// 		for _, m := range mappings[i] {
	// 			combined := mergeRanges(combinedMapping, m)
	// 			if combined.length > 0 {
	// 				combinedMappings = append(combinedMappings, combined)
	// 			}
	// 		}
	// 	}
	// 	result = combinedMappings
	// }

	result := seedMappings
	for _, mapping := range mappings {
		result = mergeMappings(result, mapping)
	}

	min := math.MaxInt
	for _, r := range result {
		if r.destination < min {
			min = r.destination
		}
	}
	fmt.Printf("result mappings: %v\nmin: %d\n", result, min)
}
