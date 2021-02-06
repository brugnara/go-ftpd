package main

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

func TestNewFtp(t *testing.T) {
	for _, test := range []struct {
		in  string
		out *ftpd
	}{
		{" foo      ", &ftpd{"foo", "foo"}},
		{" ./foo/", &ftpd{"foo", "foo"}},
		{"./foo", &ftpd{"foo", "foo"}},
		{"../foo/baz", &ftpd{"../foo/baz", "../foo/baz"}},
		{"   ", &ftpd{"./public", "./public"}},
		{"", &ftpd{"./public", "./public"}},
	} {
		if s := newFtp(test.in); !reflect.DeepEqual(s, test.out) {
			t.Error("Expected:", s, "to be:", test.out, "with:", test.in)
		}
	}
}

func TestHello(t *testing.T) {
	f := &ftpd{"", ""}
	var buff bytes.Buffer
	f.hello(&buff)

	if buff.String() == "" {
		t.Error("Something should be written to the writer")
	}

	if !strings.Contains(buff.String(), "Welcome") {
		t.Error("The FTPD should be more kind and say Welcome")
	}
}

func TestHelp(t *testing.T) {
	f := &ftpd{"", ""}
	var buff bytes.Buffer
	f.help(&buff)

	if buff.String() == "" {
		t.Error("Something should be written to the writer")
	}

	if !strings.Contains(buff.String(), "Available commands:") {
		t.Error("What an unuseful help is this?")
	}
}

func TestCursor(t *testing.T) {
	for _, test := range []struct {
		in  string
		out string
	}{
		{"/public", "$ / > "},
		{"////ciao///", "$ / > "},
		{"", "$ / > "},
		{"../", "$ / > "},
		{"../   /", "$ / > "},
	} {
		f := newFtp(test.in)
		var b bytes.Buffer
		f.cursor(&b)
		if s := b.String(); s != test.out {
			t.Error("Expected:", s, "to be:", test.out, "with:", test.in)
		}
	}
}

func TestCommandGenerics(t *testing.T) {
	for _, test := range []struct {
		in  []string
		out bool
	}{
		// invalids
		{[]string{"foo", ""}, false},
		{[]string{"foo", "baz"}, false},
		{[]string{"foo", "cd"}, false},
		{[]string{"foo", "ls"}, false},
		{[]string{"foo", "cat"}, false},
		{[]string{"foo", "cat baz"}, false},
		// valids
		{[]string{"public", "quit"}, true},
		{[]string{"public", "ls"}, true},
		{[]string{"public", "cd tt"}, true},
		{[]string{"public", "cat test.txt"}, true},
	} {
		f := &ftpd{test.in[0], test.in[0]}
		if o := f.command(ioutil.Discard, test.in[1]); o != test.out {
			t.Error("Expected:", o, "to be:", test.out, "with:", test.in)
		}
	}
}

func TestCommandCursor(t *testing.T) {
	for _, test := range []struct {
		in  []string
		out string
	}{
		// only valids here!
		{[]string{"public", "ls"}, "$ / > "},
		{[]string{"public", "cd tt"}, "$ /tt > "},
		{[]string{"public", "cd blabla"}, "$ / > "},
		{[]string{"public", "cat test.txt"}, "$ / > "},
	} {
		f := &ftpd{test.in[0], test.in[0]}
		var b bytes.Buffer
		f.command(&b, test.in[1])
		lines := strings.Split(b.String(), "\n")
		if len(lines) == 0 {
			t.Error("nothing here but at least the cursor should be there..")
			continue
		}
		cursor := lines[len(lines)-1]
		if cursor != test.out {
			t.Error("Expected:", cursor, "to be:", test.out, "with:", test.in)
		}
	}
}
