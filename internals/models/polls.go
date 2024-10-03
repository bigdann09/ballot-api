package models

import (
	"log"
	"reflect"

	"github.com/gocolly/colly/v2"
)

type Poll struct {
	Harris string `json:"harris"`
	Trump  string `json:"trump"`
}

type PollResult struct {
	Candidate string `json:"candidate"`
	Result    string `json:"result"`
}

// scrape polls
func ScrapedPolls() []Poll {
	var polls = make([]Poll, 1)
	c := colly.NewCollector(
		colly.AllowedDomains("https://www.270towin.com", "www.270towin.com", "270towin.com"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
		colly.CacheDir("./polls_cache"),
	)

	c.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)
	})

	c.OnHTML("table#polls", func(e *colly.HTMLElement) {
		var poll Poll
		poll.Harris = e.ChildText("tr#poll_avg_row > td:nth-child(2)")
		poll.Trump = e.ChildText("tr#poll_avg_row > td:nth-child(3)")
		polls = append(polls, poll)
	})

	c.Visit("https://www.270towin.com/2024-presidential-election-polls/national")

	var modpolls []Poll
	for key, poll := range polls {
		if key == 1 {
			pollValue := reflect.ValueOf(poll)
			if !pollValue.IsZero() {
				modpolls = append(modpolls, poll)
			}
		}
	}

	return modpolls
}