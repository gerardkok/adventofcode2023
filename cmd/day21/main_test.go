package main

import (
	"path/filepath"
	"testing"

	"adventofcode23/internal/projectpath"
)

func TestExamplePart1(t *testing.T) {
	t.Parallel()
	d := NewDay21(filepath.Join(projectpath.Root, "cmd", "day21", "example.txt"))

	want := 16
	got := d.Part1()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}
