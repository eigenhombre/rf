# rf


<img src="/nightscan.jpeg" width="400">

`rf` is a terminal-based RSS/Atom feed reader, slightly reminiscent of [a certain '80s-era Usenet news reader](https://en.wikipedia.org/wiki/Rn_(newsreader)).

I wrote this small feed reader in Go because I was unhappy with my Web-based options.  I am still adding features but am posting it on GitHub to document my progress (I am new to Go).

# Implementation

Implemented in Go.  Currently the `o` and `P` commands only work on MacOS.

# Build

Check out this repo, then

    go build .

# Usage

## Command Line Options

    rf -h         # Show help
    rf -verbose   # Scan feeds verbosely
    rf            # Scan feeds; press `-h` for help

## Example

    $  rf -verbose
    Got 24910 bytes for PG.
    Got 52183 bytes for PLISP.
    Got 6004 bytes for PGO.
    Got 37464 bytes for NYTTECH.
    Got 192070 bytes for PCLOJURE.
    Got 10022567 bytes for MATT.
            PG    SEEN: Is There Such a Thing as Good Taste?
            PG    SEEN: Beyond Smart
            PG    SEEN: Weird Languages
            PG    SEEN: How to Work Hard
        ...
            PG    SEEN: Beating the Averages
            PG    SEEN: Lisp for Web-Based Applications
            PG    SEEN: Chapter 1 of Ansi Common Lisp
            PG    SEEN: Chapter 2 of Ansi Common Lisp
            PG    SEEN: Programming Bottom-Up
            PG    SEEN: This Year We Can End the Death Penalty in California
        PLISP    SEEN: Quicklisp news: December 2021 Quicklisp dist update now available
        ...
        PLISP    SEEN: vindarel: Lisp for the web: pagination and cleaning up HTML with LQuery
        PLISP     NEW: Stelian Ionescu: On New IDEs
                        https://blog.cddr.org/posts/2021-11-23-on-new-ides/

## Press `?` for available keystrokes / commands:
    ? ?
				N next feed
				B bottom of feed
				P post
				p prev article
				s skip article for now
				n skip article for now
				x mark article read
				X mark all articles in feed as read
				o open
				q quit program

## Press `x` to mark an article read:
    ? x
        PLISP    SEEN: Tim Bradshaw: The endless droning
        PGO    SEEN: The most popular Go items of 2021
        ...
        PGO    SEEN: Why region based memory allocation help with fragmentation
    NYTTECH    SEEN: Amazon Reaches Labor Deal, Giving Workers More Power to Organize
        ...
    NYTTECH    SEEN: Navigational Apps for the Blind Could Have a Broader Appeal
    PCLOJURE     NEW: Advent of Code: Day 25 - Sea Cucumber
                        https://andreyorst.gitlab.io/posts/2021-12-25-advent-of-code-day-25/

## Press `o` to open an article in your browser:

    ? o
    PCLOJURE     NEW: Advent of Code: Day 21 - Dirac Dice
                        https://andreyorst.gitlab.io/posts/2021-12-21-advent-of-code-day-21/

(Will use your default Web browser -- **only works on MacOS** currently.)

## `q` quits:
    ? q
