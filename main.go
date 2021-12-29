package main

import (
	"flag"
	"fmt"
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
	switch fs.FeedType {
	case RSSType:
		return RSSFeedItems(fs, body), nil
	case AtomType:
		return AtomFeedItems(fs, body), nil
	default:
		return nil, fmt.Errorf("bad feed type, %v", fs.FeedType)
	}
}

const (
	DIR_FORWARD = iota
	DIR_BACKWARD
)

func showSeenItem(item FeedEntry) {
	fmt.Printf("%12s %7s: %s\n", item.Feed().ShortName, "SEEN", item.EntryTitle())
}

func showNewItem(item FeedEntry) {
	fmt.Printf("%12s %7s: %s\n", item.Feed().ShortName, "NEW", item.EntryTitle())
	fmt.Printf("%12s %7s  %s\n", "", "", item.EntryURL())
}

func showItem(item FeedEntry) {
	if urlWasSeen(item.EntryURL()) {
		showSeenItem(item)
	} else {
		showNewItem(item)
	}
}

func scanItems(pos, dir int, items []FeedEntry, verbose bool) (int, bool) {
	for {
		if pos >= len(items) {
			pos = len(items) - 1
			// showItem(items[pos])
			return pos, true
		}
		if pos < 0 {
			pos = 0
			// showItem(items[pos])
			return 0, true
		}
		item := items[pos]
		if urlWasSeen(item.EntryURL()) {
			if dir == DIR_FORWARD {
				pos++
			} else {
				pos--
			}
		} else {
			showNewItem(item)
			return pos, false
		}
	}
}

func InteractWithItems(items []FeedEntry, theTTY *tty.TTY, verbose, repl bool) error {
	fmt.Println("")
	i := 0
	i, done := scanItems(i, DIR_FORWARD, items, verbose)
	if done && !repl {
		return nil
	}
	for {
		item := items[i]
		fmt.Print("? ")
		c := readChar(theTTY)
		fmt.Println(c)
		switch c {
		case "H":
			postItem(item, theTTY)
			recordURL(item.EntryURL())
			i++
		case "n":
			i++
			if i >= len(items) {
				i = len(items) - 1
			}
			showItem(items[i])
		case "N":
			i++
			i, _ = scanItems(i, DIR_FORWARD, items, verbose)
			showItem(items[i])
		case "p":
			i--
			if i < 0 {
				i = 0
			}
			showItem(items[i])
		case "P":
			i--
			i, _ = scanItems(i, DIR_BACKWARD, items, verbose)
			showItem(items[i])
		case "F":
			i = 0
			showItem(items[i])
		case "A":
			i = len(items) - 1
			showItem(items[i])
		case "x":
			recordURL(item.EntryURL())
			i, _ = scanItems(i, DIR_FORWARD, items, verbose)
			showItem(items[i])
		case "u":
			unRecordURL(item.EntryURL())
			showItem(item)
		case "o":
			macOpen(item.EntryURL())
		case "q":
			if verbose {
				fmt.Println("\n\nOK, See ya!")
			}
			return nil
		case "?":
			fmt.Println("USAGE:")
			fmt.Println(`
			F first article

			p prev article (read or unread)
			P prev unread article

			n next article (read or unread)
			N next unread article

			x mark article read
			u mark article unread
			o open article in browser
			H post on Hacker News (must be logged in)

			A last article
			q quit program

			? this help message
			`)
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
	var verbose, repl bool
	flagSet.BoolVar(&verbose, "verbose", false, "verbose output")
	flagSet.BoolVar(&repl, "repl", false, "always provide REPL")
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
			if verbose {
				fmt.Printf("%12s: %d items found\n", fs.ShortName, len(items))
			} else {
				fmt.Print(".")
			}
			ch <- items
		}(fs)
	}
	wg.Wait()
	close(ch)

	// Consume and concatenate results:
	var procItems []FeedEntry = []FeedEntry{}
	for items := range ch {
		procItems = append(procItems, items...)
	}

	InteractWithItems(procItems, stdin, verbose, repl)

	if verbose {
		fmt.Println("OK")
	}
}
