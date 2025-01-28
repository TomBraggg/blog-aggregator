package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/tombraggg/blog-aggregator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.Args) == 1 {
		var err error
		limit, err = strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("invalid limit %w", err)
		}
	}

	posts, err := s.db.GetPostsByUserWithLimit(context.Background(), database.GetPostsByUserWithLimitParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts for user: %w", err)
	}

	fmt.Printf("Showing %d posts for user %s\n", len(posts), user.Name)
	for _, post := range posts {
		printPostInfo(database.Post{
			ID:          post.ID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       post.Title,
			Url:         post.Url,
			Description: post.Description,
			PublishedAt: post.PublishedAt,
			FeedID:      post.FeedID,
		})
	}

	return nil
}

func printPostInfo(post database.Post) {
	fmt.Println("=====================================")
	fmt.Printf("* Title:		%s\n", post.ID)
	fmt.Printf("* URL:			%v\n", post.Url)
	fmt.Printf("* Description:	%v\n", post.Description)
	fmt.Printf("* PublisheAt:	%s\n", post.PublishedAt)
	fmt.Println("=====================================")
}
