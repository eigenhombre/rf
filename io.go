package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/mattn/go-tty"
)

func mkdirIfNotExists(dirName string) error {
	return os.MkdirAll(dirName, 0755)
}

func spit(fileName string, content string) {
	dirName := path.Dir(fileName)
	mkdirIfNotExists(dirName)
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	writer := bufio.NewWriter(file)
	writer.WriteString(content)
	defer writer.Flush()
	defer file.Close()
}

func rm(fileName string) error {
	return os.Remove(fileName)
}

func removeFileExtension(fileName string) string {
	return strings.TrimSuffix(
		strings.TrimSuffix(fileName, path.Ext(fileName)),
		"/",
	)
}

// readChar reads a character from stdin using unbuffered I/O.
func readChar(theTTY *tty.TTY) string {
	r, err := theTTY.ReadRune()
	if err != nil {
		log.Fatal(err)
	}
	return string(r)
}

func RawFeedData(endpoint string) ([]byte, error) {
	res, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
