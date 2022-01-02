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
