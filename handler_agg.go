package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/tombraggg/blog-aggregator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <time between request>", cmd.Name)
	}
	frequency, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("could not parse frequency string to duration: %w", err)
	}

	ticker := time.NewTicker(frequency)
	for ; ; <-ticker.C {
		go scrapeFeeds(*s)
	}
}

func scrapeFeeds(s state) error {
	seenFeeds := make(map[database.Feed]struct{})

	for {
		feedToFetch, err := s.db.GetNextFeedToFetch(context.Background())
		if err != nil {
			return fmt.Errorf("could not get next feed to fetch: %w", err)
		}

		if _, exitsts := seenFeeds[feedToFetch]; exitsts {
			return nil
		}

		seenFeeds[feedToFetch] = struct{}{}

		_, err = s.db.MarkFeedFetchedByID(context.Background(), database.MarkFeedFetchedByIDParams{
			LastFetchedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			UpdatedAt: time.Now(),
			ID:        feedToFetch.ID,
		})
		if err != nil {
			return fmt.Errorf("could not mark feed as fetched")
		}

		feed, err := fetchFeed(context.Background(), feedToFetch.Url)
		if err != nil {
			return fmt.Errorf("could not fetch feed %w", err)
		}

		printRSSFeed(*feed)
	}
}
