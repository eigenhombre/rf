package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"

	"github.com/mattn/go-tty"
	"jaytaylor.com/html2text"
)

const (
	feedStateDir = "/tmp/rss.feeds"
)

const (
	rssType  = "rss"
	atomType = "atom"
)

// FeedSpec is a generic RSS/Atom feed specifier.
type FeedSpec struct {
	ShortName string `json:"name"`
	URL       string `json:"url"`
	FeedType  string `json:"type"`
}

// FeedEntry is a generic RSS/Atom feed post type.
type FeedEntry interface {
	EntryTitle() string
	EntryURL() string
	Feed() FeedSpec
}

func getFeedItems(fs FeedSpec, verbose bool) ([]FeedEntry, error) {
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
			}
			return pos + 1, false
		}
		if pos == 0 {
			return pos, true
		}
		return pos - 1, false
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

func fetchAndShowArticle(item FeedEntry) {
	fmt.Printf("Fetching %s... ", item.EntryTitle())
	body, err := httpGetBytes(item.EntryURL())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("(%d bytes)\n", len(body))
	text, err := html2text.FromString(string(body),
		html2text.Options{PrettyTables: false, OmitLinks: true, TextOnly: true})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(wrapText(text, 80))
}

func unreadItemIndices(items []FeedEntry) []int {
	ret := []int{}
	for i, item := range items {
		if !urlWasSeen(item.EntryURL()) {
			ret = append(ret, i)
		}
	}
	return ret
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
			if urlWasSeen(items[i].EntryURL()) && !repl {
				return nil
			}
		case "N":
			i, _ = nextItem(i, dirForward, nextAny, items, verbose)
		case "p":
			i, _ = nextItem(i, dirBackward, nextUnread, items, verbose)
		case "P":
			i, _ = nextItem(i, dirBackward, nextAny, items, verbose)
		case "R":
			uis := unreadItemIndices(items)
			if len(uis) == 0 {
				fmt.Println("No unread items!")
			} else {
				i = uis[rand.Intn(len(uis))]
			}
		case "F":
			i = 0
		case "f":
			fetchAndShowArticle(item)
		case "A":
			i = len(items) - 1
		case "x":
			recordURL(item.EntryURL())
			i, _ = nextItem(i, dirForward, nextUnread, items, verbose)
			if urlWasSeen(items[i].EntryURL()) && !repl {
				return nil
			}
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
			R random unread article

			x mark article read
			X mark all articles in current feed read
			u mark article unread

			o open article in browser
			f fetch and show article online, in plain text (beta!)
			H post on Hacker News (must be logged in)

			A last article
			q quit program

			h or ? this help message
			`)
		}
	}
}

func main() {
	err := mkdirIfNotExists(feedStateDir)
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
	feedSpecs, err := serializedFeeds()
	if err != nil {
		fmt.Println(`There was a problem reading your feed configuration file.

Create a file $HOME/.rffeeds.json as follows:

[
	{
		"name": "PG",
		"url": "http://www.aaronsw.com/2002/feeds/pgessays.rss",
		"type": "rss"
	  }
]

Then restart the program.
		
		`)
		os.Exit(-1)
	}
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

	stdin, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer stdin.Close()

	// Consume and concatenate results:
	var procItems []FeedEntry = []FeedEntry{}
	for items := range ch {
		procItems = append(procItems, items...)
	}

	interactWithItems(procItems, stdin, verbose, repl)
}
