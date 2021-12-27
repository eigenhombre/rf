package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"path"
	"regexp"
)

func metaDataFilePath(url string) (x string) {
	var stripFront = regexp.MustCompile(`^http(?:s)\:\/\/`)
	var feedFile = stripFront.ReplaceAllString(url, "")
	return removeFileExtension(feedStateDir + "/" + feedFile)
}

func urlWasSeen(url string) bool {
	path := metaDataFilePath(url)
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
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

func recordURL(url string) {
	spit(metaDataFilePath(url), "PLACEHOLDER")
}
