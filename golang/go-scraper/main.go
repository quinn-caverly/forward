package main

import (
	"log"
	"time"

	"github.com/quinn-caverly/go-scraper/stussy"
)

func main() {
	est, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatal("Was not able to load EST timezone: ", err)
	}
	date_scraped := time.Now().In(est).Format("2006-01-02")

    stussy.Scrape(date_scraped)
}
