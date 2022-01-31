# rf

<img src="nightscan.jpeg" width="400">

![build](https://github.com/eigenhombre/rf/actions/workflows/build.yml/badge.svg)

`rf` is a terminal-based RSS/Atom feed reader, slightly reminiscent of [a certain '80s-era Usenet news reader](https://en.wikipedia.org/wiki/Rn_(newsreader)).

I wrote this small feed reader because I was unhappy with my Web-based options.

# Features

- Largely terminal-based; no Web bloat, though you can open articles in the browser;
- Cross-post to Hacker News;
- Single-keystroke-driven;
- Fetch all feeds in parallel to minimize latency;
- Workflow similar in feel to [`rn`](https://en.wikipedia.org/wiki/Rn_(newsreader)) (you may have to be over 50 and/or have other problems for this to have any appeal);
- Configurable JSON list of feeds (coming soon: add/remove subscriptions through the program itself);
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

## Feed Configuration

Currently, you have to add feeds manually by creating `.rffeeds.json` 
in your home directory.  Example:

	[
		{
			"name": "PG",
			"url": "http://www.aaronsw.com/2002/feeds/pgessays.rss",
			"type": "rss"
		},
		{
			"name": "MATT",
			"url": "https://matthewrocklin.com/blog/atom.xml",
			"type": "atom"
		}
	]

# License

Copyright © 2022, John Jacobsen. MIT License.

# Disclaimer

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
