package main

import (
	"path/filepath"
	"testing"

	"adventofcode23/internal/projectpath"
)

func TestExample1Part1(t *testing.T) {
	t.Parallel()
	d := NewDay08(filepath.Join(projectpath.Root, "cmd", "day08", "example1-part1.txt"))

	want := 2
	got := d.Part1()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestExample2Part1(t *testing.T) {
	t.Parallel()
	d := NewDay08(filepath.Join(projectpath.Root, "cmd", "day08", "example2-part1.txt"))

	want := 6
	got := d.Part1()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestExamplePart2(t *testing.T) {
	t.Parallel()
	d := NewDay08(filepath.Join(projectpath.Root, "cmd", "day08", "example-part2.txt"))

	want := 6
	got := d.Part2()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}
