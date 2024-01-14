package main

import (
	"path/filepath"
	"testing"

	"adventofcode23/internal/projectpath"
)

func TestExamplePart1(t *testing.T) {
	t.Parallel()
	d := NewDay17b(filepath.Join(projectpath.Root, "cmd", "day17", "example.txt"))

	want := 102
	got := d.Part1()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestExamplePart2(t *testing.T) {
	t.Parallel()
	d := NewDay17b(filepath.Join(projectpath.Root, "cmd", "day17", "example.txt"))

	want := 94
	got := d.Part2()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestExample2Part2(t *testing.T) {
	t.Parallel()
	d := NewDay17b(filepath.Join(projectpath.Root, "cmd", "day17", "example2_part2.txt"))

	want := 71
	got := d.Part2()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}
