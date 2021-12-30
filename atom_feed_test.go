package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const exampleAtomFeedSingleItem = `<?xml version="1.0" encoding="utf-8" standalone="yes" ?>
<feed xmlns="http://www.w3.org/2005/Atom">
	<title>Planet Arrakis</title>
	<link rel="self" href="http://johnj.com/atom.xml"/>
	<link href="http://johnj.com/"/>
	<updated>2021-12-22T18:04:16+00:00</updated>
	<generator uri="http://johnj.com/">http://johnj.com</generator>
	<entry>
		<title type="html">Another Atomic Article</title>
		<link href="http://johnj.com/aaa.html"/>
		<updated>2021-12-22T16:44:37+00:00</updated>
		<content type="html">Some <em>HTML</em> content</content>
		<author>
			<name>Eig N. Hombre</name>
			<email>noreply@eigenhombre.com</email>
			<uri>http://johnj.com/</uri>
		</author>
		<source>
			<title type="html">AAA</title>
			<link rel="self" href="http://johnj.com/xyz.html"/>
		</source>
	</entry>
</feed>`

func TestAtomFeedItems(t *testing.T) {
	rawData := []byte(exampleAtomFeedSingleItem)
	// FIXME: Add/improve test coverage around FeedSpec part:
	items := atomFeedItems(FeedSpec{}, rawData)
	assert.Equal(t, len(items), 1)
	assert.Equal(t, items[0].EntryTitle(), "Another Atomic Article")
	assert.Equal(t, items[0].EntryURL(), "http://johnj.com/aaa.html")
}
