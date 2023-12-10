package main

import (
	"path/filepath"
	"testing"

	"adventofcode23/internal/projectpath"
)

func TestExamplePart1(t *testing.T) {
	t.Parallel()
	d := NewDay01(filepath.Join(projectpath.Root, "cmd", "day01", "example-part1.txt"))

	want := 142
	got := d.Part1()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestExamplePart2(t *testing.T) {
	t.Parallel()
	d := NewDay01(filepath.Join(projectpath.Root, "cmd", "day01", "example-part2.txt"))

	want := 281
	got := d.Part2()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}
