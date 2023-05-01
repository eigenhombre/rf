package main

import (
	"encoding/xml"
	"strings"
)

type rssFeed struct {
	RSSVersion xml.Name   `xml:"rss"`
	Channel    rssChannel `xml:"channel"`
}

type rssChannel struct {
	XMLName xml.Name   `xml:"channel"`
	Title   string     `xml:"title"`
	Items   []RSSEntry `xml:"item"`
}

// RSSEntry is an entry / blog post for an RSS feed.
type RSSEntry struct {
	XMLName xml.Name `xml:"item"`
	Title   string   `xml:"title"`
	URL     string   `xml:"link"`
	GUID    string   `xml:"guid"`
	Date    string   `xml:"pubDate"`
	fs      Feed
}

// EntryTitle returns an RSS post's title.
func (r RSSEntry) EntryTitle() string {
	return r.Title
}

// EntryURL returns an RSS post's URL.
func (r RSSEntry) EntryURL() string {
	if len(r.URL) > 0 {
		return r.URL
	}
	return r.GUID
}

func (r RSSEntry) EntryDate() string {
	return r.Date
}

// Feed returns the feed specifier for a given RSS feed item.
func (r RSSEntry) Feed() Feed {
	return r.fs
}

func rssFeedItems(fs Feed, rawFeedData []byte) []FeedEntry {
	feed := rssFeed{}
	xml.Unmarshal(rawFeedData, &feed)
	ret := []FeedEntry{}
	for _, item := range feed.Channel.Items {
		item.URL = strings.TrimSpace(item.URL)
		item.GUID = strings.TrimSpace(item.GUID)
		item.fs = fs
		ret = append(ret, item)
	}
	return ret
}
