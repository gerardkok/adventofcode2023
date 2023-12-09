package main

import (
	"path/filepath"
	"testing"

	"adventofcode23/internal/projectpath"
)

func TestExamplePart1(t *testing.T) {
	t.Parallel()
	p, err := NewPuzzle09(filepath.Join(projectpath.Root, "cmd", "day09", "testinput"))
	if err != nil {
		t.Fatal(err)
	}
	want := 114
	got := p.Part1()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}
