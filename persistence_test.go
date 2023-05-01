package main

import "testing"

var exampleFeed = Feed{
	ShortName: "test",
	URL:       "http://example.com",
	FeedType:  "rss",
}

func TestPersistFeed(t *testing.T) {
	db, err := getDb(true)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	createTablesIfNotExist(db)

	addFeedIfNotExists(db, exampleFeed)
	feedSpecs, err := getFeedSpecs(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(feedSpecs) != 1 {
		t.Fatalf("expected 1 feedSpec, got %d", len(feedSpecs))
	}
	if feedSpecs[0].ShortName != exampleFeed.ShortName {
		t.Fatalf("expected %s, got %s", exampleFeed.ShortName, feedSpecs[0].ShortName)
	}
	if feedSpecs[0].URL != exampleFeed.URL {
		t.Fatalf("expected %s, got %s", exampleFeed.URL, feedSpecs[0].URL)
	}
	if feedSpecs[0].FeedType != exampleFeed.FeedType {
		t.Fatalf("expected %s, got %s", exampleFeed.FeedType, feedSpecs[0].FeedType)
	}
}

func TestPersistFeedItem(t *testing.T) {
	db, err := getDb(true)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	createTablesIfNotExist(db)
	addFeedIfNotExists(db, exampleFeed)
	// Add a feed item:
}
