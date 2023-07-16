package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/jesiahharris/rss-agg/internal/database"
)

func startScraping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequest time.Duration,
) {
	log.Printf("Scraping on %v goroutines every %s \n", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	// this for loop runs every interval(timeBetweenRequest) on the ticker channel
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Println("Error fetching feeds: ", err)
			// function should always be running, continue on err
			continue
		}

		wg := &sync.WaitGroup{}
		// iterate over feeds on main goroutine
		// for each feed, add 1 to waitgroup
		for _, feed := range feeds {
			wg.Add(1)

			// spawn new goroutine per feed
			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched:", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed:", err)
	}

	for _, item := range rssFeed.Channel.Item {
		log.Println("Found post:", item.Title, "on feed", feed.Name)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
