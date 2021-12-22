package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

const (
	feedStateDir = "/tmp/rss.feeds"
)

func postItem(item Item) {
	fmt.Printf("Posting %q...\n", item)
	cmd := exec.Command("open", "https://news.ycombinator.com/submit")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	PbCopy(item.Title)
	fmt.Println("ANY KEY TO COPY URL...")
	_ = readChar()
	PbCopy(item.URL)
}

func getRssFeedURLs() []string {
	return []string{
		"https://rss.nytimes.com/services/xml/rss/nyt/Technology.xml",
		"http://planet.lisp.org/rss20.xml",
		// FIXME: Handle Atom posts:
		// "http://planet.clojure.in/atom.xml",
		// "https://planetgolang.dev/index.xml",
	}
}

func HandleFeed(url string) error {
	fmt.Printf("Handling feed %s....\n", url)
	body, err := RawFeedData(url)
	if err != nil {
		return err
	}
	fmt.Printf("Got %d bytes in XML body.\n", len(body))

	items := FeedItems(body)
articles:
	for _, item := range items {
		if urlWasSeen(item.URL) {
			fmt.Println("REPEAT: " + item.Title)
		} else {
			fmt.Println("   NEW: " + item.Title)
			fmt.Println("        " + item.URL)
			fmt.Print("Post article? ")
			c := readChar()
			fmt.Println("")
			switch strings.ToLower(c) {
			case "y":
				postItem(item)
				recordURL(item.URL)
			case "q":
				fmt.Println("\nQuitting....")
				break articles
			default:
				recordURL(item.URL)
			}
		}
	}
	return nil
}

func main() {
	err := mkdirIfNotExists(feedStateDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, feed := range getRssFeedURLs() {
		err = HandleFeed(feed)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("OK")
}
