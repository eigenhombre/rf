package main

import (
	"fmt"

	"github.com/mattn/go-tty"
)

// postItem "posts" an item on Hacker News, using the clipboard
// to send values to the Submit page (you must be logged into
// the site).
func postItem(item FeedEntry, theTTY *tty.TTY) {
	fmt.Printf("Posting %q...\n", item)
	macOpen("https://news.ycombinator.com/submit")
	pbCopy(item.EntryTitle())
	fmt.Println("ANY KEY TO COPY URL...")
	_ = readChar(theTTY)
	pbCopy(item.EntryURL())
}
