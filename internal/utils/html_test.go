package utils

import (
	"strings"
	"testing"
)

func TestCleanHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "No malformed tags",
			input:    "<p>Hello world</p>",
			expected: "<p>Hello world</p>",
		},
		{
			name:     "Malformed opening tag <p-->",
			input:    "<p-->Hello world</p>",
			expected: "<p>Hello world</p>",
		},
		{
			name:     "Malformed div tag <div-->",
			input:    "<div-->Content here</div>",
			expected: "<div>Content here</div>",
		},
		{
			name:     "Multiple malformed tags",
			input:    "<p-->First paragraph</p><div-->Second content</div>",
			expected: "<p>First paragraph</p><div>Second content</div>",
		},
		{
			name:     "Malformed tag with extra dashes",
			input:    "<p---->Too many dashes</p>",
			expected: "<p>Too many dashes</p>",
		},
		{
			name:     "Mixed valid and malformed tags",
			input:    "<div><p-->Malformed inside</p></div>",
			expected: "<div><p>Malformed inside</p></div>",
		},
		{
			name:     "Malformed self-closing img tag",
			input:    `<img src="test.png"-->`,
			expected: `<img src="test.png">`,
		},
		{
			name:     "Complex CoolShell-like content",
			input:    `<p--><script async src="https://example.com/ad.js"></script><img src="test.png">Content here</p>`,
			expected: `<p><script async src="https://example.com/ad.js"></script><img src="test.png">Content here</p>`,
		},
		{
			name:     "Content with nested tags",
			input:    "<div--><p>Nested <strong>bold</strong> text</p></div>",
			expected: "<div><p>Nested <strong>bold</strong> text</p></div>",
		},
		{
			name:     "Content with links",
			input:    `<p-->Visit <a href="https://example.com">this link</a></p>`,
			expected: `<p>Visit <a href="https://example.com">this link</a></p>`,
		},
		{
			name:     "Malformed br tag",
			input:    `Line 1<br-->Line 2`,
			expected: `Line 1<br>Line 2`,
		},
		{
			name:     "Malformed hr tag without attributes",
			input:    `Section 1<hr-->Section 2`,
			expected: `Section 1<hr>Section 2`,
		},
		{
			name:     "Real CoolShell sample",
			input:    `<p--><img decoding="async" loading="lazy" class="alignright" src="test.png" alt="" width="300">这两天技术圈里热议的一件事</p>`,
			expected: `<p><img decoding="async" loading="lazy" class="alignright" src="test.png" alt="" width="300">这两天技术圈里热议的一件事</p>`,
		},
		{
			name:     "Multiple malformed tags in sequence",
			input:    `<p-->Text<img src="a.png"--><br-->More text</p>`,
			expected: `<p>Text<img src="a.png"><br>More text</p>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanHTML(tt.input)
			if result != tt.expected {
				t.Errorf("CleanHTML() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestCleanHTML_Whitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Leading whitespace",
			input:    "   <p>Content</p>",
			expected: "<p>Content</p>",
		},
		{
			name:     "Trailing whitespace",
			input:    "<p>Content</p>   ",
			expected: "<p>Content</p>",
		},
		{
			name:     "Both leading and trailing whitespace",
			input:    "  \n\t  <p>Content</p>  \n\t  ",
			expected: "<p>Content</p>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanHTML(tt.input)
			if result != tt.expected {
				t.Errorf("CleanHTML() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestCleanHTML_PreservesValidHTML(t *testing.T) {
	validHTML := `<div class="content">
		<h1>Title</h1>
		<p>First paragraph with <strong>bold</strong> and <em>italic</em> text.</p>
		<p>Second paragraph with <a href="https://example.com">a link</a>.</p>
		<img src="image.jpg" alt="Description">
	</div>`

	result := CleanHTML(validHTML)

	// Should preserve all content and fix whitespace
	if !strings.Contains(result, "<h1>Title</h1>") {
		t.Error("Valid heading tag was modified")
	}
	if !strings.Contains(result, `<a href="https://example.com">a link</a>`) {
		t.Error("Valid link tag was modified")
	}
	if !strings.Contains(result, `<img src="image.jpg" alt="Description">`) {
		t.Error("Valid img tag was modified")
	}
}
