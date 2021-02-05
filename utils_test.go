package main

import (
	"reflect"
	"testing"
)

// workaround for error because of the use of flag.Parse()
// https://stackoverflow.com/a/58192326/1420669
var _ = func() bool {
	testing.Init()
	return true
}()

func TestSplitter(t *testing.T) {
	for _, test := range []struct {
		in  string
		out []string
	}{
		{" foo      ", []string{"foo"}},
		{" foo  bar    ", []string{"foo", "bar"}},
		{"foo baz", []string{"foo", "baz"}},
	} {
		if s := splitter(test.in); !reflect.DeepEqual(s, test.out) {
			t.Error("Expected:", s, "to be:", test.out, "with:", test.in)
		}
	}
}

func TestToSize(t *testing.T) {
	for _, test := range []struct {
		in  int64
		out string
	}{
		{1023, "1023B "},
		{1024, "   1KB"},
		{pow(1024, 2), "   1MB"},
		{pow(1024, 3), "   1GB"},
		{pow(1024, 4), "   1TB"},
		{pow(1024, 5), "   big"},
	} {
		if s := toSize(test.in); !reflect.DeepEqual(s, test.out) {
			t.Error("Expected:", s, "to be:", test.out, "with:", test.in)
		}
	}
}

func TestPow(t *testing.T) {
	for _, test := range []struct {
		in  []int
		out int64
	}{
		{[]int{2, 2}, 4},
		{[]int{2, 0}, 1},
		{[]int{0, 0}, 1},
		{[]int{0, 1}, 0},
		{[]int{10, 1}, 10},
		{[]int{10, 10}, 10000000000},
	} {
		if s := pow(test.in[0], test.in[1]); !reflect.DeepEqual(s, test.out) {
			t.Error("Expected:", s, "to be:", test.out, "with:", test.in)
		}
	}
}

func TestCut(t *testing.T) {
	for _, test := range []struct {
		in  string
		cnt int
		sep string
		out string
	}{
		{"pippo", 3, "--", "pip--"},
		{"pippo", 10, "--", "pippo"},
		{"pippo", 1, "", "p"},
		{"pippo", 0, "-", "-"},
		{"pippo", -1, "ciao", "pippo"},
	} {
		if s := cut(test.in, test.cnt, test.sep); !reflect.DeepEqual(s, test.out) {
			t.Error("Expected:", s, "to be:", test.out, "with:", test.in)
		}
	}
}
