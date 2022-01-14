# rf

<img src="nightscan.jpeg" width="400">

![build](https://github.com/eigenhombre/rf/actions/workflows/build.yml/badge.svg)

`rf` is a terminal-based RSS/Atom feed reader, slightly reminiscent of [a certain '80s-era Usenet news reader](https://en.wikipedia.org/wiki/Rn_(newsreader)).

I wrote this small feed reader because I was unhappy with my Web-based options.

# Features

- Completely text-based (except for actual viewing / posting of articles)
- Cross-post to Hacker News
- Single-keystroke-driven
- Fetch all feeds in parallel to minimize latency
- Workflow similar in feel to [`rn`](https://en.wikipedia.org/wiki/Rn_(newsreader)).  (you may have to be over 50 and/or have other problems for this to have any appeal)
- Configurable JSON list of feeds
- Provide plain-text preview of text contents inline, similar to a venerable Usenet news reader (*IN PROGRESS / BETA*).

# Implementation

Implemented in Go.  Currently the `o` and `H` commands only work on MacOS.  PRs for other platforms will be enthusiastically reviewed.

# Build

Check out this repo, then

    go build .

I'll set up `go get` eventually.

# Usage

## Command Line Options

    rf -h         # Show help
    rf -verbose   # Read feeds verbosely
    rf            # Read feeds

## Keyboard Commands

			F first article

			n next unread article
			N next article (read or unread)
			p prev unread article
			P prev article (read or unread)
			R random unread article

			x mark article read
			X mark all articles in current feed read
			u mark article unread

			o open article in browser
			f fetch and show article online, in plain text (beta!)
			H post on Hacker News (must be logged in)

			A last article
			q quit program

			h or ? this help message
