package main

import (
	"encoding/xml"
)

type atomFeed struct {
	AtomTag xml.Name    `xml:"feed"`
	Title   string      `xml:"title"`
	Items   []AtomEntry `xml:"entry"`
}

type atomLink struct {
	Href string `xml:"href,attr"`
}

// AtomEntry is an entry / blog post for an RSS feed.
type AtomEntry struct {
	XMLName xml.Name `xml:"entry"`
	Title   string   `xml:"title"`
	URL     atomLink `xml:"link"`
	fs      FeedSpec
}

// EntryTitle returns an Atom post's title.
func (e AtomEntry) EntryTitle() string {
	return e.Title
}

// EntryURL returns an Atom post's URL.
func (e AtomEntry) EntryURL() string {
	return e.URL.Href
}

// Feed returns the feed specifier for a given Atom feed item.
func (e AtomEntry) Feed() FeedSpec {
	return e.fs
}

func atomFeedItems(fs FeedSpec, rawFeedData []byte) []FeedEntry {
	feed := atomFeed{}
	xml.Unmarshal(rawFeedData, &feed)
	ret := []FeedEntry{}
	for _, item := range feed.Items {
		item.fs = fs
		ret = append(ret, item)
	}
	return ret
}
