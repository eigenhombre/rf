package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/mattn/go-tty"
)

const (
	feedStateDir = "/tmp/rss.feeds"
)

const (
	RSSType = iota
	AtomType
)

type FeedSpec struct {
	ShortName string
	URL       string
	FeedType  int
}

type FeedEntry interface {
	EntryTitle() string
	EntryURL() string
	Feed() FeedSpec
}

func allFeedSpecs() []FeedSpec {
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
		{"PG", "http://www.aaronsw.com/2002/feeds/pgessays.rss", RSSType},
	}
}

func getFeedItems(fs FeedSpec, verbose bool) ([]FeedEntry, error) {
	body, err := RawFeedData(fs.URL)
	if err != nil {
		return nil, err
	}
	if verbose {
		fmt.Printf("Got %d bytes for %s.\n", len(body), fs.ShortName)
	}
	switch fs.FeedType {
	case RSSType:
		return RSSFeedItems(fs, body), nil
	case AtomType:
		return AtomFeedItems(fs, body), nil
	default:
		return nil, fmt.Errorf("bad feed type, %v", fs.FeedType)
	}
}

func HandleFeed(items []FeedEntry, theTTY *tty.TTY, verbose bool) error {
	i := 0
	for {
		if i >= len(items) {
			return nil
		}
		item := items[i]
		if urlWasSeen(item.EntryURL()) {
			if verbose {
				fmt.Printf("%10s %7s: %s\n", item.Feed().ShortName, "SEEN", item.EntryTitle())
			}
			i++
		} else {
			fmt.Printf("%10s %7s: %s\n", item.Feed().ShortName, "NEW", item.EntryTitle())
			fmt.Printf("%10s %7s  %s\n", "", "", item.EntryURL())
			fmt.Print("? ")
			c := readChar(theTTY)
			fmt.Println(c)
			switch c {
			case "P":
				postItem(item, theTTY)
				recordURL(item.EntryURL())
				i++
			case "s":
				i++
			case "n":
				i++
			case "x":
				i++
				recordURL(item.EntryURL())
			case "o":
				macOpen(item.EntryURL())
			case "X":
				if verbose {
					fmt.Println("\nMarking all articles in feed as read...")
				}
				for _, ir := range items {
					recordURL(ir.EntryURL())
				}
				return nil
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
					if !urlWasSeen(items[i].EntryURL()) {
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
				x mark article read
				X mark all articles in feed as read
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
	flagSet := flag.NewFlagSet("args", flag.ContinueOnError)
	var verbose bool
	flagSet.BoolVar(&verbose, "verbose", false, "verbose output")
	err = flagSet.Parse(os.Args[1:])
	if err != nil {
		// Usage() is called inside Parse
		return
	}
	feedSpecs := allFeedSpecs()
	ch := make(chan []FeedEntry, len(feedSpecs))
	var wg sync.WaitGroup
	for _, fs := range feedSpecs {
		wg.Add(1)
		go func(fs FeedSpec) {
			defer wg.Done()
			items, err := getFeedItems(fs, verbose)
			if err != nil {
				log.Fatal(err)
			}
			ch <- items
		}(fs)
	}
	wg.Wait()
	close(ch)
	for items := range ch {
		err = HandleFeed(items, stdin, verbose)
		if err != nil {
			log.Fatal(err)
		}
	}
	if verbose {
		fmt.Println("OK")
	}
}
