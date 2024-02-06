package main

import (
	"flag"
	"time"
)

var (
	// Pause DF-AI when the queue length gets above this threshold.
	maxQueuedLines = 1000

	// Unpause DF-AI when the queue length gets below this threshold.
	minQueuedLines = 500

	// Minimum number of lines before a duplicate is allowed.
	minLinesBeforeDuplicate = 500

	// Maximum number of "fuzzy" (some words changed) duplicates allowed.
	maxFuzzyDuplicates = 5

	// Number of lines to remember for "fuzzy" duplicate checking.
	fuzzyDuplicateWindow = 10

	// Maximum number of words that can differ in a "fuzzy" duplicate.
	maxFuzzyDifferentWords = 2

	// Make a post this often.
	postInterval = 2 * time.Hour
)

func parseFlags() {
	flag.IntVar(&maxQueuedLines, "max-queued-lines", maxQueuedLines, "Pause DF-AI when the queue length gets above this threshold.")
	flag.IntVar(&minQueuedLines, "min-queued-lines", minQueuedLines, "Unpause DF-AI when the queue length gets below this threshold.")
	flag.IntVar(&minLinesBeforeDuplicate, "min-lines-before-duplicate", minLinesBeforeDuplicate, "Minimum number of lines before a duplicate is allowed.")
	flag.IntVar(&maxFuzzyDuplicates, "max-fuzzy-duplicates", maxFuzzyDuplicates, "Maximum number of “fuzzy” (some words changed) duplicates allowed.")
	flag.IntVar(&fuzzyDuplicateWindow, "fuzzy-duplicate-window", fuzzyDuplicateWindow, "Number of lines to remember for “fuzzy” duplicate checking.")
	flag.IntVar(&maxFuzzyDifferentWords, "max-fuzzy-different-words", maxFuzzyDifferentWords, "Maximum number of words that can differ in a “fuzzy” duplicate.")
	flag.DurationVar(&postInterval, "post-interval", postInterval, "Make a post this often.")

	flag.Parse()
}
