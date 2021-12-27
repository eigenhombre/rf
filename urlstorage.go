package main

import (
	"errors"
	"os"
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

func recordURL(url string) {
	spit(metaDataFilePath(url), "PLACEHOLDER")
}
