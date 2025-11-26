package discovery

import "errors"

var (
	errTooManyRedirects       = errors.New("too many redirects")
	errFriendLinkPageNotFound = errors.New("friend link page not found")
	errRSSFeedNotFound        = errors.New("RSS feed not found")
)
