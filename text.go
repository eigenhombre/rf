package main

import (
	"strings"
)

func wrapText(input string, column int) string {
	splits := strings.Split(input, " ")
	ret := ""
	thisRow := 0
	for _, el := range splits {
		if thisRow+len(el) > column {
			ret += "\n" + el
			thisRow = len(el)
			continue
		}
		if thisRow == 0 || el == "" {
			ret += el
		} else {
			ret += " " + el
		}
		thisRow += len(el)
	}
	return ret
}

func chunks(input string, chunkHeight int) []string {
	splits := strings.Split(input, "\n")
	ret := []string{}
	cur := ""
	curLen := 0
	for i, line := range splits {
		cur += line + "\n"
		curLen++
		if curLen >= chunkHeight || i == len(splits)-1 {
			ret = append(ret, cur)
			curLen = 0
			cur = ""
		}
	}
	return ret
}
