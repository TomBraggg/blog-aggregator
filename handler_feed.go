package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tombraggg/blog-aggregator/internal/database"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}
	name := cmd.Args[0]
	url := cmd.Args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("could not create feed: %w", err)
	}
	fmt.Println("Feed created successfully:")
	printFeed(feed, user)

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not create feed follow: %w", err)
	}
	fmt.Println("Feed followed successfully:")
	printFeedFollow(feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func handlerGetFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetAllFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("could not get all feeds: %w", err)
	}
	for _, feed := range feeds {
		userName, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("could not get user for feed: %w", err)
		}
		printFeed(feed, userName)
	}
	return nil
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Println("=====================================")
	fmt.Printf("* ID:			%s\n", feed.ID)
	fmt.Printf("* Created:		%v\n", feed.CreatedAt)
	fmt.Printf("* Updated:		%v\n", feed.UpdatedAt)
	fmt.Printf("* Name:			%s\n", feed.Name)
	fmt.Printf("* URL:			%s\n", feed.Url)
	fmt.Printf("* User:			%s\n", user.Name)
	fmt.Println("=====================================")
}
