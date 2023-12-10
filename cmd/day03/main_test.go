package main

import (
	"path/filepath"
	"testing"

	"adventofcode23/internal/projectpath"
)

func TestExamplePart1(t *testing.T) {
	t.Parallel()
	d := NewDay03(filepath.Join(projectpath.Root, "cmd", "day03", "example.txt"))

	want := 4361
	got := d.Part1()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestExamplePart2(t *testing.T) {
	t.Parallel()
	d := NewDay03(filepath.Join(projectpath.Root, "cmd", "day03", "example.txt"))

	want := 467835
	got := d.Part2()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}
