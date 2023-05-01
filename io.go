package main

import (
	"fmt"
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

func spit(fileName string, contents string) error {
	dirName := path.Dir(fileName)
	mkdirIfNotExists(dirName)
	ioutil.WriteFile(fileName, []byte(contents), 0644)
	return nil
}

func slurp(fileName string) (string, error) {
	body, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(body), nil
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

func httpGetBytes(url string) ([]byte, error) {
	res, err := http.Get(url)
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

func getFeedItems(fs Feed, verbose bool) ([]FeedEntry, error) {
	body, err := httpGetBytes(fs.URL)
	if err != nil {
		return nil, err
	}
	switch fs.FeedType {
	case rssType:
		return rssFeedItems(fs, body), nil
	case atomType:
		return atomFeedItems(fs, body), nil
	default:
		return nil, fmt.Errorf("bad feed type, %v", fs.FeedType)
	}
}
