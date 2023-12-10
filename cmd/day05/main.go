package main

import (
	"bufio"
	"math"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

var almanacRE = regexp.MustCompile(`(?ms)^seeds: (.*)$\s*^seed-to-soil map:\s*([\d\s]+)$\s*^soil-to-fertilizer map:\s*([\d\s]+)$\s*^fertilizer-to-water map:\s*([\d\s]+)$\s*^water-to-light map:\s*([\d\s]+)$\s*^light-to-temperature map:\s*([\d\s]+)$\s*^temperature-to-humidity map:\s*([\d\s]+)$\s*^humidity-to-location map:\s*([\d\s]+)$`)

type numberRange struct {
	destination, source, length int
}

type Day05 struct {
	day.DayInput
}

func NewDay05(inputFile string) Day05 {
	return Day05{day.DayInput(inputFile)}
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

	nMappings := len(matches[0]) - 2

	mappings = make([][]numberRange, nMappings)
	for i := 0; i < nMappings; i++ {
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

func (d Day05) Part1() int {
	input, _ := d.ReadFile()
	seeds, mappings := parseInput(string(input))

	locations := make([]int, len(seeds))
	for i, seed := range seeds {
		locations[i] = findLocation(seed, mappings)
	}

	location := slices.Min(locations)
	return location
}

func (d Day05) Part2() int {
	input, _ := d.ReadFile()
	seeds, mappings := parseInput(string(input))

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
	d := NewDay05(filepath.Join(projectpath.Root, "cmd", "day05", "input.txt"))

	day.Solve(d)
}
