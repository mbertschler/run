package main

import "testing"

func TestRootPath(t *testing.T) {
	args := []string{"git.exahome.net/tools/run"}
	out := Main(args)
	if out != 0 {
		t.Error("non 0 exit status", out)
	}
}

func TestExahome(t *testing.T) {
	args := []string{"hello"}
	out := Main(args)
	if out != 0 {
		t.Error("non 0 exit status", out)
	}
}

func TestBuilder(t *testing.T) {
	args := []string{"builder"}
	out := Main(args)
	if out != 0 {
		t.Error("non 0 exit status", out)
	}
}
