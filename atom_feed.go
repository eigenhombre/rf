package main

import (
	"encoding/xml"
)

type AtomFeed struct {
	AtomTag xml.Name `xml:"feed"`
	Title   string   `xml:"title"`
	Items   []Entry  `xml:"entry"`
}

type Link struct {
	Href string `xml:"href,attr"`
}

type Entry struct {
	XMLName xml.Name `xml:"entry"`
	Title   string   `xml:"title"`
	URL     Link     `xml:"link"`
}

func atomToGeneric(entry Entry) GenericFeedEntry {
	return GenericFeedEntry{entry.Title, entry.URL.Href}
}

func AtomFeedItems(rawFeedData []byte) []GenericFeedEntry {
	feed := AtomFeed{}
	xml.Unmarshal(rawFeedData, &feed)
	ret := []GenericFeedEntry{}
	for _, item := range feed.Items {
		ret = append(ret, atomToGeneric(item))
	}
	return ret
}
