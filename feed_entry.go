package main

type FeedEntry interface {
	EntryTitle() string
	EntryURL() string
}
