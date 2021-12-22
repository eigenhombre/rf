package main

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/mattn/go-tty"
)

func mkdirIfNotExists(dirName string) error {
	return os.MkdirAll(dirName, 0755)
}

func removeFileExtension(fileName string) string {
	return strings.TrimSuffix(
		strings.TrimSuffix(fileName, path.Ext(fileName)),
		"/",
	)
}

func readChar() string {
	tty, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer tty.Close()

	r, err := tty.ReadRune()
	if err != nil {
		log.Fatal(err)
	}
	return string(r)
}
