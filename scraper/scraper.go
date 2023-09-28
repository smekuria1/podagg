package scraper

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/smekuria1/podagg/internal/db"
	"github.com/smekuria1/podagg/parser"
)

func StartScraping(
	db *db.Queries,
	goroutines int,
	timeBetweenRequest time.Duration,
) {
	l := log.New(os.Stdout, "podagg-scraper", log.LstdFlags)

	l.Printf("Scraping on %v gorouting every %s duration", goroutines, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(goroutines),
		)

		if err != nil {
			l.Println("error fetching feeds:", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go ScrapeFeed(wg, db, feed)
		}

		wg.Wait()

	}
}

func ScrapeFeed(wg *sync.WaitGroup, database *db.Queries, feed db.Feed) {
	defer wg.Done()
	l := log.New(os.Stdout, "podagg-scraper", log.LstdFlags)
	_, err := database.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		l.Println("Error marking feed as feteched", err)
	}

	rssFeed, err := parser.URLToFeed(feed.Url)
	if err != nil {
		l.Println("Error fetching/scraping feed:", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}
		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			l.Printf("Error parsing Pub Date format should be %s skipping post %v", time.RFC1123Z, item.Title)
			continue
		}
		_, err = database.CreatePosts(context.Background(), db.CreatePostsParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			PublishedAt: pubAt,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			l.Println("Failed to create post", err)
		}
	}

	l.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))

}
