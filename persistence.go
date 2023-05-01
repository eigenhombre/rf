package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

// Feed is a generic RSS/Atom feed specifier.
type Feed struct {
	Id        int
	ShortName string
	URL       string
	FeedType  string
}

type Entry struct {
	Id    int
	Title string
	URL   string
	Date  string
	FS    Feed
}

func getDb(isMem bool) (*sql.DB, error) {
	if isMem {
		db, err := sql.Open("sqlite", ":memory:")
		if err != nil {
			return nil, err
		}
		return db, nil
	}
	feedDb := os.Getenv("HOME") + "/.rffeeds.sqlite"
	fmt.Println("Feed DB:", feedDb)
	db, err := sql.Open("sqlite", feedDb)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createTablesIfNotExist(db *sql.DB) {
	creates := []string{`
		CREATE TABLE IF NOT EXISTS feeds (
		id INTEGER NOT NULL PRIMARY KEY,
		shortdesc TEXT UNIQUE NOT NULL,
		url TEXT NOT NULL,
		type TEXT NOT NULL,
		fetched TEXT NULL
		);`,
		`
		CREATE TABLE IF NOT EXISTS entries (
		id INTEGER NOT NULL PRIMARY KEY,
		title TEXT NOT NULL,
		url TEXT NOT NULL,
		date TEXT NULL,
		feed_id INTEGER NOT NULL,
		-- boolean read:
		read INTEGER DEFAULT 0,
		UNIQUE (title, url, feed_id)
		);`,
	}
	for _, create := range creates {
		if _, err := db.Exec(create); err != nil {
			panic(err)
		}
	}
}

func addFeedIfNotExists(db *sql.DB, fs Feed) {
	_, err := db.Exec(`INSERT OR IGNORE INTO feeds (shortdesc, url, type) VALUES (?, ?, ?)`,
		fs.ShortName,
		fs.URL,
		fs.FeedType)
	if err != nil {
		panic(err)
	}
}

// FIXME: unused
func populateExistingFeedspecs(db *sql.DB, feedSpecs []Feed) {
	for _, fs := range feedSpecs {
		fmt.Printf("%12s: %s (%s)\n", fs.ShortName, fs.URL, fs.FeedType)
		// Check if the row already exists:
		var id int
		err := db.QueryRow(`SELECT id FROM feeds WHERE shortdesc = ?`, fs.ShortName).Scan(&id)
		if err == nil {
			continue
		}
		_, err = db.Exec(`INSERT INTO feeds (shortdesc, url, type) VALUES (?, ?, ?)`,
			fs.ShortName,
			fs.URL,
			fs.FeedType)
		if err != nil {
			panic(err)
		}
	}
}

func showFeedSpecsFromDb(db *sql.DB) {
	// Look up feeds in DB:
	rows, err := db.Query(`SELECT id, shortdesc, url, type FROM feeds`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var shortdesc string
		var url string
		var feedtype string
		if err := rows.Scan(&id, &shortdesc, &url, &feedtype); err != nil {
			panic(err)
		}
		fmt.Println(id, shortdesc, url, feedtype)
	}
}

func getFeedSpecs(db *sql.DB) ([]Feed, error) {
	// Look up feeds in DB:
	rows, err := db.Query(`SELECT id, shortdesc, url, type FROM feeds`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ret := []Feed{}
	for rows.Next() {
		var id int
		var shortdesc string
		var url string
		var feedtype string
		if err := rows.Scan(&id, &shortdesc, &url, &feedtype); err != nil {
			panic(err)
		}
		ret = append(ret, Feed{Id: id, ShortName: shortdesc, URL: url, FeedType: feedtype})
	}
	return ret, nil
}

func updateFeedItemsInDb(db *sql.DB, fs Feed, verbose bool) {
	items, err := getFeedItems(fs, verbose)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%12s: %d items found\n", fs.ShortName, len(items))
	for _, item := range items {
		time, err := parseTime(item.EntryDate())
		if err != nil {
			fmt.Println("NOT ADDING ITEM with invalid time: ", item.EntryDate())
			continue
		}
		fmt.Printf("%12s (%d): %s ... %s.\n", fs.ShortName, fs.Id, item.EntryTitle(), item.EntryDate())
		_, err = db.Exec(`INSERT OR IGNORE INTO entries (title, url, date, feed_id) VALUES (?, ?, ?, ?)`,
			item.EntryTitle(),
			item.EntryURL(),
			time,
			fs.Id)
		if err != nil {
			panic(err)
		}
	}
	// Set feed as read (fetched == now UTC):
	_, err = db.Exec(`UPDATE feeds SET fetched = datetime('now', 'utc') WHERE id = ?`, fs.Id)
	if err != nil {
		panic(err)
	}
}

func shouldFetchFeed(db *sql.DB, fs Feed, duration time.Duration) bool {
	var fetchedStr string
	err := db.QueryRow(`SELECT fetched FROM feeds WHERE id = ?`, fs.Id).Scan(&fetchedStr)
	if err != nil {
		panic(err)
	}
	if fetchedStr == "" {
		return true
	}
	fetched, err := time.Parse("2006-01-02 15:04:05", fetchedStr)
	if err != nil {
		panic(err)
	}
	if time.Since(fetched) >= duration {
		return true
	}
	return false
}

func getFeedItemsFromDb(db *sql.DB) ([]Entry, error) {
	// Look up feeds in DB:
	rows, err := db.Query(`SELECT id, title, url, date, feed_id FROM entries`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ret := []Entry{}
	for rows.Next() {
		var id int
		var title string
		var url string
		var date string
		var feedId int
		if err := rows.Scan(&id, &title, &url, &date, &feedId); err != nil {
			panic(err)
		}
		entry := Entry{Id: id, Title: title, URL: url, Date: date, FS: Feed{Id: feedId}}
		ret = append(ret, entry)
	}
	return ret, nil
}
