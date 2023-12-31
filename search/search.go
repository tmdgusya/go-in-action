package search

import (
	"log"
	"sync"
)

// This variable is considered a package-level variable
// unexported identifiers start with a lowercase letter
var matchers = make(map[string]Matcher)

func Run(searchTerm string) {
	// Retrieve the list of feeds to search through
	feeds, err := RetrieveFeeds()

	if err != nil {
		log.Fatal(err)
	}

	// create an unbuffered channel to receive match results
	// A good rule of thumb when declaring variables is to use the keyword var when declaring variables that will be initialized to their zero value,
	// and to use the short variable declaration operator when you're providing extra initialization or making a function call.
	results := make(chan *Result)

	// Setup a wait group so we can process all the feeds
	var waitGroup sync.WaitGroup

	// Set the number of goroutines we need to wait for while they process the individual feeds
	waitGroup.Add(len(feeds))

	// Launch a goroutine for each feed to find the results
	for _, feed := range feeds {
		// Retrieve a matcher for the search
		matcher, exists := matchers[feed.Type]

		if !exists {
			matcher = matchers["default"]
		}

		// Launch the goroutine to perform the search
		go func(matcher Matcher, feed *Feed) {
			Match(matcher, feed, searchTerm, results)
			waitGroup.Done()
		}(matcher, feed)
	}

	// Launch a goroutine to monitor when all the work is done
	go func() {
		// Wait for everything to be processed
		waitGroup.Wait()

		// Close the channel to signal to the Display function that we can exit the program
		close(results)
	}()

	// Start displaying results as they are available and return after the final result is displayed
	Display(results)
}

// Register is called to register a matcher for use by the program.
func Register(feedType string, matcher Matcher) {
	if _, exists := matchers[feedType]; exists {
		log.Fatalln(feedType, "Matcher already registered")
	}

	log.Println("Register", feedType, "matcher")
	matchers[feedType] = matcher
}
