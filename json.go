package main

import (
	"encoding/json"
	"os"
)

// // Commented out pending feed addition:
// func serializeFeedsAsJSON(feedSpecs []FeedSpec) (string, error) {
// 	theBytes, err := json.MarshalIndent(feedSpecs, "", "  ")
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(theBytes), nil
// }

func feedFileName() string {
	return os.Getenv("HOME") + "/.rffeeds.json"
}

// // Commented out pending feed addition:
// func writeFeeds(feedSpecs []FeedSpec) error {
// 	m, err := serializeFeedsAsJSON(feedSpecs)
// 	if err != nil {
// 		return err
// 	}
// 	err = spit(feedFileName(), m)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func serializedFeeds() ([]Feed, error) {
	feedFile := feedFileName()
	body, err := slurp(feedFile)
	if err != nil {
		return nil, err
	}
	feeds := []Feed{}
	err = json.Unmarshal([]byte(body), &feeds)
	if err != nil {
		return nil, err
	}
	return feeds, nil
}
