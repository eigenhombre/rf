package main

import (
	"encoding/xml"
)

type RSSFeed struct {
	RSSVersion xml.Name   `xml:"rss"`
	Channel    RSSChannel `xml:"channel"`
}

type RSSChannel struct {
	XMLName xml.Name  `xml:"channel"`
	Title   string    `xml:"title"`
	Items   []RSSItem `xml:"item"`
}

type RSSItem struct {
	XMLName xml.Name `xml:"item"`
	Title   string   `xml:"title"`
	URL     string   `xml:"link"`
}

func rssToGeneric(entry RSSItem) GenericFeedEntry {
	return GenericFeedEntry{entry.Title, entry.URL}
}

func RSSFeedItems(rawFeedData []byte) []GenericFeedEntry {
	feed := RSSFeed{}
	xml.Unmarshal(rawFeedData, &feed)
	ret := []GenericFeedEntry{}
	for _, item := range feed.Channel.Items {
		ret = append(ret, rssToGeneric(item))
	}
	return ret
}
