package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const exampleRSSFeedSingleItem = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
	<channel>
		<title>Foo Technology</title>
		<item>
			<title>Strongly Lettered Word</title>
			<link>http://johnj.com/bazzy.html</link>
			<guid isPermaLink="true">http://johnj.com/bazzy.html</guid>
		</item>
	</channel>
</rss>`

func TestRSSFeedItems(t *testing.T) {
	rawData := []byte(exampleRSSFeedSingleItem)
	items := RSSFeedItems(rawData)
	assert.Equal(t, len(items), 1)
	assert.Equal(t, items[0].Title, "Strongly Lettered Word")
	assert.Equal(t, items[0].URL, "http://johnj.com/bazzy.html")
}
