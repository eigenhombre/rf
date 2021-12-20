package main

import (
	"bufio"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	tty "github.com/mattn/go-tty"
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

const (
	feedDataDir = "/tmp/nyt.feeds"
)

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

func mkdirIfNotExists(dirName string) error {
	return os.MkdirAll(dirName, 0755)
}

func removeFileExtension(fileName string) string {
	return strings.TrimSuffix(fileName, path.Ext(fileName))
}

func metaDataFilePath(url string) string {
	var stripFront = regexp.MustCompile(`^http(?:s)\:\/\/`)
	var feedFile = stripFront.ReplaceAllString(url, "")
	return removeFileExtension(feedDataDir + "/" + feedFile)
}

func recordURL(url string) {
	spit(metaDataFilePath(url), "PLACEHOLDER")
}

func urlWasSeen(url string) bool {
	path := metaDataFilePath(url)
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func pbcopy(data string) {
	pbcopyCmd := exec.Command("pbcopy")
	pbcopyIn, _ := pbcopyCmd.StdinPipe()
	pbcopyCmd.Start()
	pbcopyIn.Write([]byte(data))
	pbcopyIn.Close()
	pbcopyCmd.Wait()
}

func postItem(item Item) {
	fmt.Printf("Posting %q...\n", item)
	cmd := exec.Command("open", "https://news.ycombinator.com/submit")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	pbcopy(item.Title)
	fmt.Println("ANY KEY TO COPY URL...")
	_ = readChar()
	pbcopy(item.URL)
}

func spit(fileName string, content string) {
	dirName := path.Dir(fileName)
	mkdirIfNotExists(dirName)
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	writer := bufio.NewWriter(file)
	writer.WriteString(content)
	defer writer.Flush()
}

func readChar() string {
	tty, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer tty.Close()

	r, err := tty.ReadRune()
	if err != nil {
		log.Fatal(err)
	}
	return string(r)
}

func main() {
	err := mkdirIfNotExists(feedDataDir)
	if err != nil {
		log.Fatal(err)
	}

	body, err := RawFeedData("https://rss.nytimes.com/services/xml/rss/nyt/Technology.xml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Got %d bytes in XML body.\n", len(body))

	items := FeedItems(body)
articles:
	for _, item := range items {
		if urlWasSeen(item.URL) {
			fmt.Println("REPEAT: " + item.Title)
		} else {
			fmt.Println("   NEW: " + item.Title)
			fmt.Print("Post article? ")
			c := readChar()
			fmt.Println("")
			switch strings.ToLower(c) {
			case "y":
				postItem(item)
				recordURL(item.URL)
			case "q":
				fmt.Println("\nQuitting....")
				break articles
			default:
				recordURL(item.URL)
			}
		}
	}
	fmt.Println("OK")
}
