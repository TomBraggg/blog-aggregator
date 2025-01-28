package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tombraggg/blog-aggregator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}
	url := cmd.Args[0]

	feedToFollow, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("could not find feed: %w", err)
	}

	ffRow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feedToFollow.ID,
	})
	if err != nil {
		return fmt.Errorf("could not follow feed: %w", err)
	}

	fmt.Println("Feed follow created:")
	printFeedFollow(ffRow.UserName, ffRow.FeedName)
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}
	url := cmd.Args[0]

	_, err := s.db.DeleteFeedFollowByURL(context.Background(), database.DeleteFeedFollowByURLParams{
		UserID: user.ID,
		Url:    url,
	})
	if err != nil {
		return fmt.Errorf("could not get delete feed by URL: %w", err)
	}

	fmt.Printf("Successfully unfollowed feed: %s", url)

	return nil
}

func handlerListFeedFollows(s *state, cmd command, user database.User) error {
	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("could not get followed feeds: %w", err)
	}

	fmt.Print("Followed Feeds:")
	for _, follow := range feedFollows {
		feed, err := s.db.GetFeedByID(context.Background(), follow.FeedID)
		if err != nil {
			return fmt.Errorf("could not get feed by ID: %w", err)
		}
		printFeed(feed, user)
	}

	return nil
}

func printFeedFollow(username, feedname string) {
	fmt.Println("=====================================")
	fmt.Printf("* User:		%s\n", username)
	fmt.Printf("* Feed:		%s\n", feedname)
	fmt.Println("=====================================")
}
