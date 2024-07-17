package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/ahgr3y/blog-aggregator/internal/database"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Items       []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

func fetchRSSFromFeed(feedUrl string) (*RSS, error) {

	httpClient := http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := httpClient.Get(feedUrl)
	if err != nil {
		return &RSS{}, err
	}
	defer resp.Body.Close()

	decoder := xml.NewDecoder(resp.Body)
	params := RSS{}
	err = decoder.Decode(&params)
	if err != nil {
		return &RSS{}, err
	}

	return &params, nil
}

func parseDateString(dateString string) time.Time {

	formats := []string{
		time.RFC1123Z,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC3339,
	}

	for _, format := range formats {
		pubDate, err := time.Parse(format, dateString)
		if err != nil {
			log.Printf("error parsing string to time: %s", err)
			continue
		} else {
			return pubDate
		}
	}

	return time.Now().UTC()
}

func savePost(db *database.Queries, feed database.Feed, rssItem RSSItem) error {

	_, err := db.CreatePost(context.Background(), database.CreatePostParams{
		ID:        database.GenerateUUID(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Title:     rssItem.Title,
		Url:       rssItem.Link,
		Description: sql.NullString{
			String: rssItem.Description,
			Valid:  true,
		},
		PublishedAt: parseDateString(rssItem.PubDate),
		FeedID:      feed.ID,
	})
	if err != nil {
		return err
	}

	return nil
}

// scrapFeed - Print the title of each post in each feed,
// save the each post in the database,
// then mark each feed as fetched.
func scrapFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed, report *[]string) {

	defer wg.Done()

	// Fetch data from feed and print title of posts,
	// then save the post to the database.
	rssData, err := fetchRSSFromFeed(feed.Url)
	if err != nil {
		log.Printf("Error from fetchRSSFromFeed: %s", err)
		return
	}
	for _, item := range rssData.Channel.Items {
		println(item.Title)
		err = savePost(db, feed, item)
		if err != nil {
			continue
		}
	}

	// Mark feed as fetched
	err = db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Error from MarkFeedFetched: %s", err)
		return
	}

	*report = append(*report, fmt.Sprintf("%s processed with %v posts", feed.Name, len(rssData.Channel.Items)))
}

// startScraper - For every 10 minutes, the next 10 feeds from the
// database will be fetched and processed.
func startScraper(db *database.Queries) {

	ticker := time.NewTicker(time.Minute * 10)
	summaryReport := &[]string{}

	log.Println("Scraper started scraping...")
	for ; ; <-ticker.C {
		wg := &sync.WaitGroup{}

		log.Println("Getting next feeds to fetch...")
		feeds, err := db.GetNextFeedsToFetch(context.Background(), 10)
		if err != nil {
			log.Printf("Error from GetNextFeedsToFetch: %s", err)
			continue
		}
		log.Println("Finished acquiring feeds to fetch...")

		log.Println("Start processing feeds...")
		for _, feed := range feeds {
			log.Printf("Processing %v...", feed.Name)
			wg.Add(1)
			go scrapFeed(db, wg, feed, summaryReport)
		}

		wg.Wait()

		log.Println("Feeds processed successfully...")
		for _, report := range *summaryReport {
			log.Println(report)
		}
	}

}
