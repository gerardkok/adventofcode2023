package main

import (
	"path/filepath"
	"testing"

	"adventofcode23/internal/projectpath"
)

func TestExamplePart1(t *testing.T) {
	t.Parallel()
	d := NewDay18(filepath.Join(projectpath.Root, "cmd", "day18", "example.txt"))

	want := 62
	got := d.Part1()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

// func TestExamplePart2(t *testing.T) {
// 	t.Parallel()
// 	d := NewDay17(filepath.Join(projectpath.Root, "cmd", "day17", "example.txt"))

// 	want := 51
// 	got := d.Part2()
// 	if want != got {
// 		t.Errorf("want %d, got %d", want, got)
// 	}
// }
