package main

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

type Feed struct {
	RSSVersion xml.Name `xml:"rss"`
	Channel    Channel  `xml:"channel"`
}

type Channel struct {
	XMLName xml.Name `xml:"channel"`
	Title   string   `xml:"title"`
	Items   []Item   `xml:"item"`
}

type Item struct {
	XMLName xml.Name `xml:"item"`
	Title   string   `xml:"title"`
	URL     string   `xml:"guid"`
}

func RawFeedData(endpoint string) ([]byte, error) {
	res, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func FeedItems(rawFeedData []byte) []Item {
	feed := Feed{}
	xml.Unmarshal(rawFeedData, &feed)
	return feed.Channel.Items
}
