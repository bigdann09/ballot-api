package models

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"time"

	"encoding/json"

	"github.com/ballot/internals/utils"
	"github.com/gocolly/colly/v2"
)

type Article struct {
	Title       string            `json:"title"`
	Slug        string            `json:"slug"`
	PublishedBy string            `json:"published_by"`
	Thumbnail   string            `json:"thumb"`
	Metadata    string            `json:"metadata"`
	Contents    []string          `json:"contents"`
	Socials     map[string]string `json:"socials"`
}

func GetArticles() []*Article {
	var articles []*Article
	c := colly.NewCollector(
		colly.AllowedDomains("https://edition.cnn.com", "www.edition.cnn.com", "edition.cnn.com"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
		colly.CacheDir("./news_cache"),
	)

	c.SetRequestTimeout(120 * time.Second)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Got a response from", r.Request.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)
	})

	c.OnHTML("article.live-story-post.liveStoryPost", func(e *colly.HTMLElement) {
		var article Article
		article.Title = template.HTMLEscapeString(e.ChildText("h2.live-story-post__headline.inline-placeholder"))
		article.PublishedBy = e.ChildText("span.live-story-post__byline.inline-placeholder")
		article.Thumbnail = e.ChildAttr("div.image_live-story.image_live-story__hide-placeholder", "data-url")
		article.Metadata = e.ChildText("span[data-editable=metaCaption]")
		article.Slug = utils.ParseToSlug(article.Title)

		// get social links
		socials := make(map[string]string)
		e.ForEach("button.social-share_compact__share", func(_ int, f *colly.HTMLElement) {
			attr := f.Attr("data-url")
			if strings.Contains(attr, "twitter") {
				socials["twitter"] = attr
			} else if strings.Contains(attr, "facebook") {
				socials["facebook"] = attr
			} else if strings.Contains(attr, "cnn.com/politics") {
				socials["website"] = attr
			}
		})

		e.ForEach("p.paragraph.inline-placeholder.vossi-paragraph", func(_ int, f *colly.HTMLElement) {
			content := f.Text
			article.Contents = append(article.Contents, strings.TrimSpace(content))
		})

		if article.Title == "" {
			article.Title = template.HTMLEscapeString(e.ChildText("h2.live-story-post__headline.inline-placeholder > strong"))
		}

		article.Socials = socials
		articles = append(articles, &article)
	})

	c.OnScraped(func(r *colly.Response) {

		// convert articles to byte
		parsed, err := json.MarshalIndent(&articles, "", "  ")
		if err != nil {
			log.Print(err)
		}

		// create file to store articles
		os.WriteFile("articles.json", parsed, 0666)
	})

	day := time.Now().Day()
	url := fmt.Sprintf("https://edition.cnn.com/politics/live-news/trump-harris-election-10-%d-24/index.html", day)

	c.Visit(url)

	return articles
}
