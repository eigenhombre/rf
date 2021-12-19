package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
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
	GUID    string   `xml:"guid"`
}

func RawFeedData(endpoint string) ([]byte, error) {
	res, err := http.Get("https://rss.nytimes.com/services/xml/rss/nyt/Technology.xml")
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

func main() {
	body, err := RawFeedData("https://rss.nytimes.com/services/xml/rss/nyt/Technology.xml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Got %d bytes in XML body.\n", len(body))

	items := FeedItems(body)
	for _, item := range items {
		fmt.Println(item.Title)
	}
	fmt.Println("OK")
}
