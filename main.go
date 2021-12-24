package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mattn/go-tty"
)

const (
	feedStateDir = "/tmp/rss.feeds"
)

func postItem(item GenericFeedEntry, theTTY *tty.TTY) {
	fmt.Printf("Posting %q...\n", item)
	macOpen("https://news.ycombinator.com/submit")
	PbCopy(item.Title)
	fmt.Println("ANY KEY TO COPY URL...")
	_ = readChar(theTTY)
	PbCopy(item.URL)
}

const (
	RSSType = iota
	AtomType
)

type FeedSpec struct {
	ShortName string
	URL       string
	FeedType  int
}

func getRssFeedURLs() []FeedSpec {
	return []FeedSpec{
		{"NYTTECH", "https://rss.nytimes.com/services/xml/rss/nyt/Technology.xml", RSSType},
		{"NYTSCI", "https://rss.nytimes.com/services/xml/rss/nyt/Science.xml", RSSType},
		{"PLISP", "http://planet.lisp.org/rss20.xml", RSSType},
		{"PCLOJURE", "http://planet.clojure.in/atom.xml", AtomType},
		{"PGO", "https://planetgolang.dev/index.xml", AtomType},
		{"MATT", "https://matthewrocklin.com/blog/atom.xml", AtomType},
		{"PP", "https://paintingperceptions.com/feed", RSSType},
		{"ILLUS", "http://illustrationart.blogspot.com/feeds/posts/default", AtomType},
		{"MUDDY", "https://www.muddycolors.com/feed/", RSSType},
		{"GURNEY", "http://gurneyjourney.blogspot.com/feeds/posts/default", AtomType},
	}
}

func getRssFeedItems(fs FeedSpec, verbose bool) ([]GenericFeedEntry, error) {
	if verbose {
		fmt.Printf("Handling feed '%s' (%s)....\n", fs.ShortName, fs.URL)
	}
	body, err := RawFeedData(fs.URL)
	if err != nil {
		return nil, err
	}
	if verbose {
		fmt.Printf("Got %d bytes in XML body.\n", len(body))
	}
	switch fs.FeedType { // FIXME: make and use a method
	case RSSType:
		return RSSFeedItems(body), nil
	case AtomType:
		return AtomFeedItems(body), nil
	default:
		return nil, fmt.Errorf("bad feed type, %v", fs.FeedType)
	}
}

func HandleFeed(fs FeedSpec, items []GenericFeedEntry, theTTY *tty.TTY, verbose bool) error {
	i := 0
	for {
		if i >= len(items) {
			return nil
		}
		item := items[i]
		if urlWasSeen(item.URL) {
			if verbose {
				fmt.Printf("%10s %7s: %s\n", fs.ShortName, "SEEN", item.Title)
			}
			i++
		} else {
			fmt.Printf("%10s %7s: %s\n", fs.ShortName, "NEW", item.Title)
			fmt.Printf("%10s %7s: %s\n", fs.ShortName, "", item.URL)
			fmt.Print("? ")
			c := readChar(theTTY)
			fmt.Println("")
			switch c {
			case "P":
				postItem(item, theTTY)
				recordURL(item.URL)
				i++
			case "s":
				i++
			case "n":
				i++
			case "x":
				i++
				recordURL(item.URL)
			case "o":
				macOpen(item.URL)
			case "N":
				if verbose {
					fmt.Println("\nWill stop processing articles in this feed....")
				}
				return nil
			case "q":
				if verbose {
					fmt.Println("\n\nOK, See ya!")
				}
				return io.EOF
			case "p":
				if i > 0 {
					i--
				}
				for {
					if i == 0 {
						break
					}
					if !urlWasSeen(items[i].URL) {
						break
					}
					i--
				}
			case "B":
				i = len(items) - 1
			case "?":
				fmt.Println(`
				N next feed
				B bottom of feed
				P post
				p prev article
				s skip article for now
				n skip article for now
				x mark article done
				o open
				q quit program
				`)
			}
		}
	}
}

func main() {
	stdin, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer stdin.Close()

	err = mkdirIfNotExists(feedStateDir)
	if err != nil {
		log.Fatal(err)
	}
	fs := flag.NewFlagSet("args", flag.ContinueOnError)
	var verbose bool
	fs.BoolVar(&verbose, "verbose", false, "verbose output")
	err = fs.Parse(os.Args[1:])
	if err != nil {
		// Usage() is called inside Parse
		return
	}
	for _, fs := range getRssFeedURLs() {
		items, err := getRssFeedItems(fs, verbose)
		if err != nil {
			log.Fatal(err)
		}
		err = HandleFeed(fs, items, stdin, verbose)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
	}
	if verbose {
		fmt.Println("OK")
	}
}
