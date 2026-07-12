package main

import (
	"context"
	"database/sql"
	"fmt"

	"log"
	"strconv"
	"time"

	"github.com/thenitverse/gator/internal/database"

	"github.com/google/uuid"
	"github.com/lib/pq"
	// "errors")
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("no username was provided")
	}
	_, err := s.db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}
	err = s.cfg.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Printf("User has been set to: <%s>", cmd.Args[0])
	return nil

}
func handlerRegister(s *state, cmds command) error {
	if len(cmds.Args) != 1 {
		return fmt.Errorf("no username was provided")
	}
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmds.Args[0],
	})
	if err != nil {
		return err

	}

	err = s.cfg.SetUser(cmds.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User has been set to: <%s>", cmds.Args[0])
	fmt.Println(user.Name)
	return nil

}
func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("not enough arguments")
	}
	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %v\n", timeBetweenReqs)

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		if err := scraperFeeds(s); err != nil {
			fmt.Println("error scraping feeds:", err)
		}
	}
	return nil

}
func scraperFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get next feed to fetch: %w", err)
	}
	err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return err
	}
	fetchedFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	for _, item := range fetchedFeed.Channel.Item {
		publishedAt, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			log.Printf("couldn't parse publishedAt %q: %v", item.PubDate, err)
			continue
		}

		post, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			FeedID:      feed.ID,
			PublishedAt: publishedAt,
		})
		if err != nil {
			pqErr, ok := err.(*pq.Error)
			if ok && pqErr.Code == "23505" {
				continue

			}
			log.Println(err)
			continue

		}
		fmt.Println("Found post:", post.Title)

	}

	return nil

}
func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.Args) == 1 {
		if specifiedLimit, err := strconv.Atoi(cmd.Args[0]); err == nil {
			limit = specifiedLimit
		} else {
			return fmt.Errorf("invalid limit: %w", err)
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts for user: %w", err)
	}

	fmt.Printf("Found %d posts for user %s:\n", len(posts), user.Name)
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}

	return nil
}
