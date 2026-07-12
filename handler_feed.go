package main

import (
	"context"
	"fmt"

	"time"

	"github.com/thenitverse/gator/internal/database"

	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("no user provided")
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}
	fmt.Println("Feed created successfully:")
	fmt.Printf(" * ID:      %v\n", feed.ID)
	fmt.Printf(" * Name:    %v\n", feed.Name)
	fmt.Printf(" *Url:    %v\n", feed.Url)
	fmt.Printf(" *CreatedAt:    %v\n", feed.CreatedAt)
	fmt.Printf(" *UpdatedAt:    %v\n", feed.UpdatedAt)
	return nil

}
func handlerFeeds(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Println(feed.Name)
		fmt.Println(feed.Url)
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Println(user.Name)
	}
	return nil

}
func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("no url provided")
	}

	feed, err := s.db.GetFeedByURL(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}
	fmt.Printf("%s is now following %s\n", feedFollow.UserName, feedFollow.FeedName)

	return nil

}
func handlerFollowing(s *state, cmd command, user database.User) error {

	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	for _, item := range feedFollows {
		fmt.Println(item.FeedName)
	}
	return nil
}
