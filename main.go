package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

const (
	feedStateDir = "/tmp/rss.feeds"
)

func postItem(item GenericFeedEntry) {
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

const (
	RSSType = iota
	AtomType
)

func getRssFeedURLs() map[string]int {
	return map[string]int{
		"https://rss.nytimes.com/services/xml/rss/nyt/Technology.xml": RSSType,
		"http://planet.lisp.org/rss20.xml":                            RSSType,
		"https://rss.nytimes.com/services/xml/rss/nyt/Science.xml":    RSSType,
		"http://planet.clojure.in/atom.xml":                           AtomType,
		"https://planetgolang.dev/index.xml":                          AtomType,
		"https://matthewrocklin.com/blog/atom.xml":                    AtomType,
	}
}

func HandleFeed(url string, feedType int, verbose bool) error {
	if verbose {
		fmt.Printf("Handling feed %s....\n", url)
	}
	body, err := RawFeedData(url)
	if err != nil {
		return err
	}
	if verbose {
		fmt.Printf("Got %d bytes in XML body.\n", len(body))
	}
	var items []GenericFeedEntry
	switch feedType {
	case RSSType:
		items = RSSFeedItems(body)
	case AtomType:
		items = AtomFeedItems(body)
	default:
		log.Fatal(fmt.Sprintf("Bad feed type, %d!", feedType))
	}
articles:
	for _, item := range items {
		if urlWasSeen(item.URL) {
			if verbose {
				fmt.Println("    REPEAT: " + item.Title)
			}
		} else {
			fmt.Println("       NEW: " + item.Title)
			fmt.Println("            " + item.URL)
			fmt.Print("Post article? ")
			c := readChar()
			fmt.Println("")
			switch strings.ToLower(c) {
			case "y":
				postItem(item)
				recordURL(item.URL)
			case "q":
				if verbose {
					fmt.Println("\nWill stop processing articles in this feed....")
				}
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
	verbose := flag.Bool("verbose", false, "verbose output")
	flag.Parse()
	for feed, feedType := range getRssFeedURLs() {
		err = HandleFeed(feed, feedType, *verbose)
		if err != nil {
			log.Fatal(err)
		}
	}
	if *verbose {
		fmt.Println("OK")
	}
}
