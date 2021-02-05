package main

import (
	"fmt"
	"strings"
)

func cut(s string, count int, sep string) string {
	if count < 0 {
		return s
	}
	if len(s) > count {
		return s[:count] + sep
	}
	return s
}

var sizes []string

func init() {
	sizes = []string{
		"B ",
		"KB",
		"MB",
		"GB",
		"TB",
	}
}

func toSize(size int64) string {
	for i := 0; i < len(sizes); i++ {
		if size < pow(1024, i+1) {
			return fmt.Sprintf("%4d%s", size/pow(1024, i), sizes[i])
		}
	}
	return "   big"
}

func pow(x, exp int) int64 {
	ret := int64(1)
	for i := 0; i < exp; i++ {
		ret *= int64(x)
	}
	return ret
}

func splitter(cmd string) []string {
	ret := []string{}
	for _, s := range strings.Fields(cmd) {
		ret = append(ret, strings.TrimSpace(s))
	}
	return ret
}
