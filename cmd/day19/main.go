package main

import (
	"path/filepath"
	"strconv"
	"strings"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

type Day19 struct {
	day.DayInput
}

func NewDay19(inputFile string) Day19 {
	return Day19{day.DayInput(inputFile)}
}

const (
	accepted = "A"
	rejected = "R"
)

type part map[byte]int

type rule struct {
	category  byte
	target    string
	threshold int
	applies   func(value, threshold int) bool
}

type workflows map[string][]rule

func parseInput(lines []string) (workflows, []part) {
	result := make([][]string, 2)
	result = append(result, make([]string, 0))

	i := 0
	for _, line := range lines {
		if len(line) == 0 {
			i++
			result = append(result, make([]string, 0))
			continue
		}

		result[i] = append(result[i], line)
	}

	w := parseWorkflows(result[0])
	r := parseParts(result[1])

	return w, r
}

func parseCondition(condition, target string) rule {
	category := condition[0]
	applies := func(value, threshold int) bool {
		return value > threshold
	}
	if condition[1] == byte('<') {
		applies = func(value, threshold int) bool {
			return value < threshold
		}
	}
	threshold, _ := strconv.Atoi(condition[2:])

	return rule{
		category:  category,
		target:    target,
		threshold: threshold,
		applies:   applies,
	}
}

func parseRule(r string) rule {
	condition, target, ok := strings.Cut(r, ":")
	if !ok {
		// just a target
		return rule{
			target: r,
			applies: func(_, _ int) bool {
				return true
			},
		}
	}
	return parseCondition(condition, target)
}

func parseRules(rules string) []rule {
	s := strings.Split(rules, ",")
	result := make([]rule, len(s))

	for i, r := range s {
		result[i] = parseRule(r)
	}
	return result
}

func parseWorkflows(w []string) workflows {
	result := make(workflows, len(w))
	for _, workflow := range w {
		name, r, _ := strings.Cut(workflow[:len(workflow)-1], "{") // cut off }
		result[name] = parseRules(r)
	}
	return result
}

func parsePart(rating string) part {
	result := make(part, 4)
	split := strings.Split(rating[1:len(rating)-1], ",") // cut off { and }

	for _, r := range split {
		category := r[0]
		value, _ := strconv.Atoi(r[2:])
		result[category] = value
	}
	return result
}

func parseParts(ratings []string) []part {
	result := make([]part, len(ratings))
	for i, rating := range ratings {
		result[i] = parsePart(rating)
	}
	return result
}

func (w workflows) accept(p part, workflow string) bool {
	if workflow == accepted || workflow == rejected {
		return workflow == accepted
	}

	rules := w[workflow]
	for _, r := range rules {
		if r.applies(p[r.category], r.threshold) {
			return w.accept(p, r.target)
		}
	}
	return false
}

func (p part) sum() int {
	result := 0
	for _, v := range p {
		result += v
	}
	return result
}

func (d Day19) Part1() int {
	lines, _ := d.ReadLines()
	workflows, ratings := parseInput(lines)

	sum := 0

	for _, rating := range ratings {
		if workflows.accept(rating, "in") {
			sum += rating.sum()
		}
	}

	return sum
}

func (d Day19) Part2() int {
	return 0
}

func main() {
	d := NewDay19(filepath.Join(projectpath.Root, "cmd", "day19", "input.txt"))

	day.Solve(d)
}
