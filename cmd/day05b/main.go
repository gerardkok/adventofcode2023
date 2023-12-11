package main

import (
	"bufio"
	"math"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

var almanacRE = regexp.MustCompile(`(?ms)^seeds: (.*)$\s*^seed-to-soil map:\s*([\d\s]+)$\s*^soil-to-fertilizer map:\s*([\d\s]+)$\s*^fertilizer-to-water map:\s*([\d\s]+)$\s*^water-to-light map:\s*([\d\s]+)$\s*^light-to-temperature map:\s*([\d\s]+)$\s*^temperature-to-humidity map:\s*([\d\s]+)$\s*^humidity-to-location map:\s*([\d\s]+)$`)

type numberRange struct {
	destination, source, length int
}

type Day05b struct {
	day.DayInput
}

func NewDay05b(inputFile string) Day05b {
	return Day05b{day.DayInput(inputFile)}
}

func parseRange(s string) numberRange {
	fields := strings.Fields(s)
	destination, _ := strconv.Atoi(fields[0])
	source, _ := strconv.Atoi(fields[1])
	length, _ := strconv.Atoi(fields[2])
	return numberRange{destination, source, length}
}

func parseMapping(s string) []numberRange {
	ranges := make([]numberRange, 0)
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		r := parseRange(scanner.Text())
		ranges = append(ranges, r)
	}

	return addMissingRanges(ranges)
}

func addMissingRanges(mapping []numberRange) []numberRange {
	sort.Slice(mapping, func(i, j int) bool {
		return mapping[i].source < mapping[j].source
	})

	bloated := make([]numberRange, 2*len(mapping)+1)
	source := 0
	for i, m := range mapping {
		bloated[2*i] = numberRange{source, source, m.source - source}
		bloated[2*i+1] = m
		source = m.source + m.length
	}
	bloated[2*len(mapping)] = numberRange{source, source, math.MaxInt - source}

	return slices.DeleteFunc(bloated, func(n numberRange) bool {
		return n.length == 0
	})
}

func parseInput(input string) (seeds []int, mappings [][]numberRange) {
	matches := almanacRE.FindAllStringSubmatch(string(input), -1)

	s := strings.Fields(matches[0][1])
	seeds = make([]int, len(s))
	for i, m := range s {
		seed, _ := strconv.Atoi(m)
		seeds[i] = seed
	}

	nMappings := len(matches[0]) - 2
	mappings = make([][]numberRange, nMappings)
	for i := 0; i < nMappings; i++ {
		mappings[i] = parseMapping(matches[0][i+2])
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
			merged := mergeRanges(m, n)
			if merged.length > 0 {
				result = append(result, merged)
			}
		}
	}
	return result
}

func minLocation(seedMappings []numberRange, mappings [][]numberRange) int {
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

	return min
}

func (d Day05b) Part1() int {
	input, _ := d.ReadFile()
	seeds, mappings := parseInput(string(input))

	seedMappings := make([]numberRange, len(seeds))
	for i := 0; i < len(seeds); i++ {
		seedMappings[i] = numberRange{seeds[i], seeds[i], 1}
	}

	return minLocation(seedMappings, mappings)
}

func (d Day05b) Part2() int {
	input, _ := d.ReadFile()
	seeds, mappings := parseInput(string(input))

	seedMappings := make([]numberRange, len(seeds)/2)
	for i := 0; i < len(seeds); i += 2 {
		seedMappings[i/2] = numberRange{seeds[i], seeds[i], seeds[i+1]}
	}

	return minLocation(seedMappings, mappings)
}

func main() {
	d := NewDay05b(filepath.Join(projectpath.Root, "cmd", "day05", "input.txt"))

	day.Solve(d)
}
