package feed

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"time"

	"github.com/mmcdole/gofeed"

	"MrRSS/backend"
	"MrRSS/backend/history"
)

type FeedOriginContentsInfo struct {
	Title string
	Image string
	Item  gofeed.Item
}

func GetFeedContent(db *sql.DB) []backend.FeedContentsInfo {
	feedList := GetFeedList(db)

	var feedOriginContent []FeedOriginContentsInfo

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	feedParser := gofeed.NewParser()

	for _, feed := range feedList {
		link := feed.Link
		feed, err := feedParser.ParseURLWithContext(link, ctx)
		if err != nil {
			continue
		}

		for _, item := range feed.Items {
			// Get the URL of the feed
			u, err := url.Parse(feed.Link)
			if err != nil {
				panic(err)
			}

			// Get the image URL for the feed
			imageUrl := "https://www.google.com/s2/favicons?sz=16&domain=" + u.Host
			if feed.Image != nil {
				imageUrl = feed.Image.URL
			}

			// Append the item to the feedOriginContent
			feedOriginContent = append(feedOriginContent, FeedOriginContentsInfo{
				Title: feed.Title,
				Image: imageUrl,
				Item:  *item,
			})
		}
	}

	sort.Slice(feedOriginContent, func(i, j int) bool {
		return feedOriginContent[i].Item.PublishedParsed.After(*feedOriginContent[j].Item.PublishedParsed)
	})

	var feedContent []backend.FeedContentsInfo

	for _, item := range feedOriginContent {
		// Get the image URL
		imageURL := ""
		filterImageUrl := filterImage(item.Item.Content)
		if item.Item.Image != nil {
			imageURL = item.Item.Image.URL
		} else if filterImageUrl != nil {
			imageURL = *filterImageUrl
		}

		// Get the time since the item was published
		timeSinceStr := getTimeSince(item.Item.PublishedParsed)

		feedContentItem := backend.FeedContentsInfo{
			FeedTitle: item.Title,
			FeedImage: item.Image,
			Title:     item.Item.Title,
			Link:      item.Item.Link,
			TimeSince: timeSinceStr,
			Time:      item.Item.PublishedParsed.Format("2006-01-02 15:04"),
			Image:     imageURL,
			Content:   item.Item.Content,
			Readed:    false,
		}

		// Check if the item is in the history
		if history.CheckInHistory(db, feedContentItem) {
			feedContentItem.Readed = history.GetHistoryReaded(db, feedContentItem)
		}

		// Append the item to the feedContent
		feedContent = append(feedContent, feedContentItem)
	}

	return feedContent
}

func filterImage(content string) *string {
	imgRegex := regexp.MustCompile(`img[^>]*src="([^"]*)`)

	var firstImageURL *string
	imgMatches := imgRegex.FindStringSubmatch(content)
	if len(imgMatches) > 1 {
		firstImageURL = &imgMatches[1]
	}

	return firstImageURL
}

func getTimeSince(t *time.Time) string {
	timeSince := time.Since(*t)
	timeStr := "now"

	if timeSince > 0 {
		if timeSince < time.Hour {
			minutes := int(timeSince.Minutes())
			if minutes == 1 {
				timeStr = fmt.Sprintf("%d minute ago", minutes)
			} else {
				timeStr = fmt.Sprintf("%d minutes ago", minutes)
			}
		} else if timeSince < 24*time.Hour {
			hours := int(timeSince.Hours())
			if hours == 1 {
				timeStr = fmt.Sprintf("%d hour ago", hours)
			} else {
				timeStr = fmt.Sprintf("%d hours ago", hours)
			}
		} else if timeSince < 365*24*time.Hour {
			days := int(timeSince.Hours() / 24)
			if days == 1 {
				timeStr = fmt.Sprintf("%d day ago", days)
			} else {
				timeStr = fmt.Sprintf("%d days ago", days)
			}
		} else {
			years := int(timeSince.Hours() / (365 * 24))
			if years == 1 {
				timeStr = fmt.Sprintf("%d year ago", years)
			} else {
				timeStr = fmt.Sprintf("%d years ago", years)
			}
		}
	}

	return timeStr
}
