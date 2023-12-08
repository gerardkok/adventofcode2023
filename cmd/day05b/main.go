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

func combineMappings(a, b numberRange) numberRange {
	aStartSource := a.source
	aEndSource := a.source + a.length
	aStartDestination := a.destination
	aEndDestination := a.destination + a.length

	bStartSource := b.source
	bEndSource := b.source + b.length
	bStartDestination := b.destination
	bEndDestination := b.destination + b.length

	if aStartDestination >= bEndSource || aEndDestination <= bStartSource {
		// no overlap
		return numberRange{}
	}

	cStartSource := aStartSource
	cStartDestination := bStartDestination
	if aStartDestination > bStartSource {
		cStartDestination += (aStartDestination - bStartSource)
	} else if aStartDestination < bStartSource {
		cStartSource += (bStartSource - aStartDestination)
	}

	cEndSource := aEndSource
	cEndDestination := bEndDestination
	if aEndDestination > bEndSource {
		cEndSource -= (aEndDestination - bEndSource)
	} else if aEndDestination < bEndSource {
		cEndDestination -= (bEndSource - aEndDestination)
	}

	length := cEndDestination - cStartDestination
	return numberRange{cStartDestination, cStartSource, length}
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
		seedMapping := numberRange{seeds[i], seeds[i], seeds[i+1]}
		seedMappings = append(seedMappings, seedMapping)
	}

	result := seedMappings
	for i := range mappings {
		combinedMappings := make([]numberRange, 0)
		for _, combinedMapping := range result {
			for _, m := range mappings[i] {
				combined := combineMappings(combinedMapping, m)
				if combined.length > 0 {
					combinedMappings = append(combinedMappings, combined)
				}
			}
		}
		result = combinedMappings
	}

	// for _, firstMapping := range mappings[0] {
	// 	for _, secondMapping := range mappings[1] {
	// 		combined := combineMappings(firstMapping, secondMapping)
	// 		if combined.length > 0 {
	// 			combinedMappings = append(combinedMappings, combined)
	// 		}
	// 		//fmt.Printf("first mapping: %v\nsecond mapping: %v\ncombined: %v\n", firstMapping, secondMapping, combined)
	// 	}
	// }
	min := math.MaxInt
	for _, r := range result {
		if r.destination < min {
			min = r.destination
		}
	}
	fmt.Printf("result mappings: %v\nmin: %d\n", result, min)
}
