# fs


<img src="/nightscan.jpeg" width="400">

A terminal-based RSS/Atom "feed scanner," slightly reminiscent of [a certain '80s-era Usenet news reader](https://en.wikipedia.org/wiki/Rn_(newsreader)).

I wrote this small feed reader in Go because I was unhappy with my Web-based options.  I am still adding features but am posting it on GitHub to document my progress (I am new to Go).

# Usage

    fs -h         # Show help
    fs -verbose   # Scan feeds verbosely
    fs            # Scan feeds; press `-h` for help

    08:10:26 fs 37.7F     ≡ * ☐ ~ (master) >  fs -verbose
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
    ? x
    PCLOJURE    SEEN: Fermat's Christmas Theorem : Fixed Points Come In Pairs
    PCLOJURE    SEEN: Christmas Theorem: Some Early Orbits
    PCLOJURE    SEEN: Clojure CLI tools - which execution option to use
    PCLOJURE    SEEN: News
    PCLOJURE     NEW: Middle Business Analyst
                        https://agiliway.com/middle-business-analyst-2/
    ? x
    PCLOJURE    SEEN: Simple Component Driven ClojureScript
    PCLOJURE    SEEN: Clojure Deref (Dec 23, 2021)
    PCLOJURE    SEEN: Full Stack Software Engineer at Urbest
    PCLOJURE    SEEN: Best Python IDEs for Data Science!
    PCLOJURE    SEEN: Full Stack Developer at Edgewood Software Corp.
    PCLOJURE    SEEN: CIDER 1.2 (Nice)
    PCLOJURE    SEEN: A Christmas Card in Clojure
    PCLOJURE     NEW: Advent of Code: Day 21 - Dirac Dice
                        https://andreyorst.gitlab.io/posts/2021-12-21-advent-of-code-day-21/
    ? o
    PCLOJURE     NEW: Advent of Code: Day 21 - Dirac Dice
                        https://andreyorst.gitlab.io/posts/2021-12-21-advent-of-code-day-21/
    ? q


    OK, See ya!
    OK

