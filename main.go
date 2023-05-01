package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	_ "github.com/glebarez/go-sqlite"
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

// FeedEntry is a generic RSS/Atom feed post type.
type FeedEntry interface {
	EntryTitle() string
	EntryURL() string
	EntryDate() string
	Feed() Feed
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

func showPaginated(theTTY *tty.TTY, body string) {
	_, ySize, err := theTTY.Size()
	if err != nil {
		fmt.Println(err)
		return
	}
	blocks := chunks(body, ySize-3)
	for i, chunk := range blocks {
		fmt.Printf("%s", chunk)
		if i < len(blocks)-1 {
			fmt.Print("--More--")
			c := readChar(theTTY)
			if c == "q" {
				break
			}
		}
	}
}

func fetchAndShowArticle(theTTY *tty.TTY, item FeedEntry) {
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
	showPaginated(theTTY, wrapText(text, 80))
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

func interactWithItemsOld(feedspecs []Feed, items []FeedEntry, theTTY *tty.TTY, verbose, repl bool) error {
	fmt.Println("")
	i := 0
	i, done := nextItem(i, dirForward, nextUnread, items, verbose)
	if done && !repl {
		return nil
	}
	numToScroll := 0
	for {
		item := items[i]
		fmt.Printf("%3d", i)
		showItem(item)
		fmt.Print("? ")
		if numToScroll > 0 {
			numToScroll--
			i, _ = nextItem(i, dirForward, nextUnread, items, verbose)
			continue
		}
		c := readChar(theTTY)
		fmt.Println(c)
		switch c {
		case "E":
			for _, f := range feedspecs {
				fmt.Println(f.ShortName)
			}
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
			fetchAndShowArticle(theTTY, item)
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
		case "l":
			numToScroll = 5
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
			l list next five unread articles

			x mark article read
			X mark all articles in current feed read
			u mark article unread

			o open article in browser
			f fetch and show article online, in plain text (beta!)
			H post on Hacker News (must be logged in)

			A last article

			E list feeds
			q quit program

			h or ? this help message
			`)
		}
	}
}

func interactWithDbItems(db *sql.DB, stdin *tty.TTY, verbose, repl bool) {
	fmt.Println("")
	feedSpecs, err := getFeedSpecs(db)
	if err != nil {
		log.Fatal(err)
	}
	fsi := 0
	fmt.Printf("Current feed: %s\n", feedSpecs[fsi].ShortName)
	// // Get items from DB:
	// items, err := getFeedItemsFromDb(db)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%d items found in %d feeds.\n", len(items), len(feedSpecs))
	// for _, item := range items {
	// 	fmt.Printf("%12s (%d): %s ... %s.\n", item.FS.ShortName, item.FS.Id, item.Title, item.Date)
	// }
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

	// DB Setup:
	db, err := getDb(false)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	createTablesIfNotExist(db)
	feedSpecs, err := getFeedSpecs(db)
	if err != nil {
		log.Fatal(err)
	}
	for _, fs := range feedSpecs {
		if shouldFetchFeed(db, fs, 5*time.Minute) {
			fmt.Printf("Fetching items for '%s': %s\n", fs.ShortName, fs.URL)
			updateFeedItemsInDb(db, fs, verbose)
		}
	}
	// // ch := make(chan []FeedEntry, len(feedSpecs))
	// ch := make(chan []FeedEntry, 2000)
	// dones := make(chan bool, len(feedSpecs))
	// // ch := make(chan FeedEntry, len(feedSpecs))
	// var wg sync.WaitGroup
	// for _, fs := range feedSpecs {
	// 	wg.Add(1)
	// 	go func(fs FeedSpec) {
	// 		fmt.Printf("%12s: \n", fs.ShortName)
	// 		defer wg.Done()
	// 		items, err := getFeedItems(fs, verbose)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		if verbose {
	// 			fmt.Printf("%12s: %d items found\n", fs.ShortName, len(items))
	// 		} else {
	// 			fmt.Printf("%d ", len(items))
	// 		}
	// 		// for _, item := range items {
	// 		// 	ch <- item
	// 		// }
	// 		ch <- items
	// 	}(fs)
	// 	dones <- true
	// }
	// for i := 0; i < len(feedSpecs); i++ {
	// 	<-dones
	// }
	// // close(ch)
	// // wg.Wait()
	// // close(ch)

	stdin, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer stdin.Close()

	interactWithDbItems(db, stdin, verbose, repl)
	// // Consume and concatenate results:
	// var procItems []FeedEntry = []FeedEntry{}
	// fmt.Println("Waiting for all feeds to be processed...")
	// for items := range ch {
	// 	procItems = append(procItems, items...)
	// }
	// for item := range ch {
	// 	procItems = append(procItems, item)
	// }

	// interactWithItemsOld(feedSpecs, procItems, stdin, verbose, repl)
}
