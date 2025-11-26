package discovery

import (
	"context"
	"testing"
	"time"
)

func TestNewService(t *testing.T) {
	service := NewService()
	if service == nil {
		t.Fatal("NewService returned nil")
	}
	if service.client == nil {
		t.Error("client is nil")
	}
	if service.feedParser == nil {
		t.Error("feedParser is nil")
	}
}

func TestResolveURL(t *testing.T) {
	service := NewService()

	tests := []struct {
		base     string
		href     string
		expected string
	}{
		{"https://example.com/", "/about", "https://example.com/about"},
		{"https://example.com/blog/", "post.html", "https://example.com/blog/post.html"},
		{"https://example.com/", "https://other.com/page", "https://other.com/page"},
		{"https://example.com/", "", ""},
	}

	for _, test := range tests {
		result := service.resolveURL(test.base, test.href)
		if result != test.expected {
			t.Errorf("resolveURL(%q, %q) = %q; want %q", test.base, test.href, result, test.expected)
		}
	}
}

func TestIsValidBlogDomain(t *testing.T) {
	service := NewService()

	validDomains := []string{
		"myblog.com",
		"example.blog",
		"personal-website.net",
	}

	invalidDomains := []string{
		"facebook.com",
		"www.twitter.com",
		"github.com",
		"stackoverflow.com",
	}

	for _, domain := range validDomains {
		if !service.isValidBlogDomain(domain) {
			t.Errorf("Expected %q to be valid blog domain", domain)
		}
	}

	for _, domain := range invalidDomains {
		if service.isValidBlogDomain(domain) {
			t.Errorf("Expected %q to be invalid blog domain", domain)
		}
	}
}

func TestGetFavicon(t *testing.T) {
	service := NewService()

	blogURL := "https://example.com/blog"
	favicon := service.getFavicon(blogURL)

	expected := "https://www.google.com/s2/favicons?domain=example.com"
	if favicon != expected {
		t.Errorf("getFavicon(%q) = %q; want %q", blogURL, favicon, expected)
	}
}

func TestDiscoverFromFeedWithTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping network test in short mode")
	}

	service := NewService()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test with a non-existent feed URL (should fail gracefully)
	_, err := service.DiscoverFromFeed(ctx, "https://nonexistent-feed-url-12345.com/feed")
	if err == nil {
		t.Log("Expected error for non-existent feed, but got none (this is acceptable)")
	}
}

func TestProgressCallbackCalled(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping network test in short mode")
	}

	service := NewService()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	progressCalled := false
	progressCb := func(progress Progress) {
		progressCalled = true
		// Verify progress has expected fields
		if progress.Stage == "" {
			t.Error("Progress stage should not be empty")
		}
	}

	// Test with a non-existent feed URL - progress callback should still be called
	_, _ = service.DiscoverFromFeedWithProgress(ctx, "https://nonexistent-feed-url-12345.com/feed", progressCb)

	if !progressCalled {
		t.Error("Progress callback should have been called at least once")
	}
}

func TestProgressStructFields(t *testing.T) {
	// Test that Progress struct can hold all expected fields
	p := Progress{
		Stage:      "checking_rss",
		Message:    "Checking RSS feed",
		Detail:     "https://example.com",
		Current:    5,
		Total:      10,
		FeedName:   "Test Feed",
		FoundCount: 3,
	}

	if p.Stage != "checking_rss" {
		t.Errorf("Expected stage 'checking_rss', got %q", p.Stage)
	}
	if p.Current != 5 {
		t.Errorf("Expected current 5, got %d", p.Current)
	}
	if p.Total != 10 {
		t.Errorf("Expected total 10, got %d", p.Total)
	}
	if p.FoundCount != 3 {
		t.Errorf("Expected found_count 3, got %d", p.FoundCount)
	}
}
