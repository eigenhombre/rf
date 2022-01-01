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
	rssType = iota
	atomType
)

// FeedSpec is a generic RSS/Atom feed specifier.
type FeedSpec struct {
	ShortName string
	URL       string
	FeedType  int
}

// FeedEntry is a generic RSS/Atom feed post type.
type FeedEntry interface {
	EntryTitle() string
	EntryURL() string
	Feed() FeedSpec
}

func allFeedSpecs() []FeedSpec {
	return []FeedSpec{
		{"NYTTECH", "https://rss.nytimes.com/services/xml/rss/nyt/Technology.xml", rssType},
		{"NYTSCI", "https://rss.nytimes.com/services/xml/rss/nyt/Science.xml", rssType},
		{"PLISP", "http://planet.lisp.org/rss20.xml", rssType},
		{"PCLOJURE", "http://planet.clojure.in/atom.xml", atomType},
		{"PGO", "https://planetgolang.dev/index.xml", atomType},
		{"MATT", "https://matthewrocklin.com/blog/atom.xml", atomType},
		{"PP", "https://paintingperceptions.com/feed", rssType},
		{"ILLUS", "http://illustrationart.blogspot.com/feeds/posts/default", atomType},
		{"MUDDY", "https://www.muddycolors.com/feed/", rssType},
		{"GURNEY", "http://gurneyjourney.blogspot.com/feeds/posts/default", atomType},
		{"PG", "http://www.aaronsw.com/2002/feeds/pgessays.rss", rssType},
		{"FOGUS", "http://blog.fogus.me/feed", rssType},
		{"KLIPSE", "https://blog.klipse.tech/feed.xml", atomType},
		{"BORK", "https://blog.michielborkent.nl/atom.xml", atomType},
	}
}

func getFeedItems(fs FeedSpec, verbose bool) ([]FeedEntry, error) {
	body, err := rawFeedData(fs.URL)
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

func markAllItemsInFeedRead(i int, items []FeedEntry, verbose bool) []FeedEntry {
	thisFeed := items[i].Feed()
	for _, item := range items {
		if item.Feed() == thisFeed && !urlWasSeen(item.EntryURL()) {
			recordURL(item.EntryURL())
			if verbose {
				showItem(item)
			}
		}
	}
	return items
}

const (
	dirForward = iota
	dirBackward
)

const (
	nextUnread = iota
	nextAny
)

func nextItem(pos, dir, changeKind int, items []FeedEntry, verbose bool) (int, bool) {
	// Move forward or backwards by a single item:
	if changeKind == nextAny {
		if dir == dirForward {
			if pos == len(items)-1 {
				return pos, true
			} else {
				return pos + 1, false
			}
		} else {
			if pos == 0 {
				return pos, true
			} else {
				return pos - 1, false
			}
		}
	}
	// Scan forward or backwards for unread items:
	scanStarted := false
	for {
		if pos >= len(items) {
			pos = len(items) - 1
			return pos, true
		}
		if pos < 0 {
			pos = 0
			return 0, true
		}
		item := items[pos]
		if urlWasSeen(item.EntryURL()) || !scanStarted {
			if dir == dirForward {
				pos++
			} else {
				pos--
			}
			scanStarted = true
		} else {
			return pos, false
		}
	}
}

func interactWithItems(items []FeedEntry, theTTY *tty.TTY, verbose, repl bool) error {
	fmt.Println("")
	i := 0
	i, done := nextItem(i, dirForward, nextUnread, items, verbose)
	if done && !repl {
		return nil
	}
	for {
		item := items[i]
		fmt.Printf("%3d", i)
		showItem(item)
		fmt.Print("? ")
		c := readChar(theTTY)
		fmt.Println(c)
		switch c {
		case "H":
			postItem(item, theTTY)
			recordURL(item.EntryURL())
			i, _ = nextItem(i, dirForward, nextUnread, items, verbose)
		case "n":
			i, _ = nextItem(i, dirForward, nextUnread, items, verbose)
		case "N":
			i++
			i, _ = nextItem(i, dirForward, nextAny, items, verbose)
		case "p":
			i, _ = nextItem(i, dirBackward, nextUnread, items, verbose)
		case "P":
			i, _ = nextItem(i, dirBackward, nextAny, items, verbose)
		case "F":
			i = 0
		case "A":
			i = len(items) - 1
		case "x":
			recordURL(item.EntryURL())
			i, _ = nextItem(i, dirForward, nextUnread, items, verbose)
		case "X":
			items = markAllItemsInFeedRead(i, items, verbose)
			recordURL(item.EntryURL())
			i, _ = nextItem(i, dirForward, nextUnread, items, verbose)
		case "u":
			// NOTE: Can return non-nil error if file doesn't exist, but we ignore it:
			unRecordURL(item.EntryURL())
		case "o":
			err := macOpen(item.EntryURL())
			if err != nil {
				fmt.Println(err)
			}
		case "q":
			if verbose {
				fmt.Println("\n\nOK, See ya!")
			}
			return nil
		case "?", "h":
			fmt.Println("USAGE:")
			fmt.Println(`
			F first article

			n next unread article
			N next article (read or unread)
			p prev unread article
			P prev article (read or unread)

			x mark article read
			X mark all articles in current feed read
			u mark article unread

			o open article in browser
			H post on Hacker News (must be logged in)

			A last article
			q quit program

			h or ? this help message
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

	interactWithItems(procItems, stdin, verbose, repl)
}
