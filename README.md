# rf


<img src="/nightscan.jpeg" width="400">

`rf` is a terminal-based RSS/Atom feed reader, slightly reminiscent of [a certain '80s-era Usenet news reader](https://en.wikipedia.org/wiki/Rn_(newsreader)).

I wrote this small feed reader because I was unhappy with my Web-based options.

# Features
## Current

- Completely text-based (except for actual viewing / posting of articles);
- Cross-post to Hacker News;
- Single-keystroke-driven;
- Fetch all feeds in parallel to minimize latency;
- Workflow otherwise similar in feel to [`rn`](https://en.wikipedia.org/wiki/Rn_(newsreader)).  (You may have to be over 50 and/or have other problems for this to have any appeal.)

## Planned

- Configurable list of feeds (JSON and/or YAML);
- Maybe: provide in-band, plain-text preview of text contents.

# Implementation

Implemented in Go.  Currently the `o` and `H` commands only work on MacOS.  PRs for other platforms will be enthusiastically reviewed.

# Build

Check out this repo, then

    go build .

# Usage

## Command Line Options

    rf -h         # Show help
    rf -verbose   # Read feeds verbosely
    rf            # Read feeds

## Keyboard Commands

    F first article

    p prev article (read or unread)
    P prev unread article

    n next article (read or unread)
    N next unread article

    x mark article read
    o open article in browser
    H post on Hacker News (must be logged in)

    A last article
    q quit program

