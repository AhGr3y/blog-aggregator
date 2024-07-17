package main

import (
	"encoding/xml"
	"log"
	"net/http"
	"time"
)

func main() {

	httpClient := http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := httpClient.Get("https://wagslane.dev/index.xml")
	if err != nil {
		log.Printf("Error from httpClient.Get(): %s", err)
	}
	defer resp.Body.Close()

	decoder := xml.NewDecoder(resp.Body)
	params := RSS{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error from decoder.Decode(): %s", err)
	}

	log.Println(params.Channel.Items)
}

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
