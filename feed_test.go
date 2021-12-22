package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFeedItems(t *testing.T) {
	rawData := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
	<channel>
		<title>NYT &gt; Technology</title>
		<item>
			<title>Apple Pushes Return to Office Date Back to T.B.D.</title>
			<link>https://www.nytimes.com/2021/12/15/technology/apple-return-to-work.html</link>
			<guid isPermaLink="true">https://www.nytimes.com/2021/12/15/technology/apple-return-to-work.html</guid>
		</item>
	</channel>
</rss>`)
	assert.Equal(t, len(FeedItems(rawData)), 1)
}
