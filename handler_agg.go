package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tombraggg/blog-aggregator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <time between request>", cmd.Name)
	}
	frequency, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	log.Printf("Collecting feeds every %s...", frequency)

	ticker := time.NewTicker(frequency)
	for ; ; <-ticker.C {
		scrapeFeeds(*s)
	}
}

func scrapeFeeds(s state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Println("Could not get next feed to fetch: %w", err)
		return
	}
	log.Println("Found a feed to fetch!")
	scrapeFeed(s.db, feed)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedFetchedByID(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Could not mark feed as fetched: %v", err)
		return
	}

	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Could not fetch feed: %v", err)
		return
	}

	for _, item := range feedData.Channel.Item {
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: item.PubDate,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			fmt.Printf("Could not create post: %v", err)
			continue
		}
	}
}
