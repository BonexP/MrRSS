package utils

import (
	"regexp"
	"strings"
)

// selfClosingTags is the list of HTML self-closing tags to handle
const selfClosingTags = "img|br|hr|input|meta|link"

// Compile regex patterns once at package initialization for better performance
var (
	// Matches malformed opening tags like <p-->, <div-->
	malformedTagRegex = regexp.MustCompile(`<([a-zA-Z][a-zA-Z0-9]*)\s*--+>`)
	
	// Matches malformed self-closing tags with attributes like <img src="..." -->
	malformedSelfClosingWithAttrs = regexp.MustCompile(`<(` + selfClosingTags + `)\s+([^<>]+?)--+>`)
	
	// Matches malformed self-closing tags without attributes like <br-->
	malformedSelfClosingNoAttrs = regexp.MustCompile(`<(` + selfClosingTags + `)\s*--+>`)
)

// CleanHTML sanitizes HTML content by fixing common malformed patterns
// that can cause rendering issues.
func CleanHTML(html string) string {
	if html == "" {
		return html
	}

	// Fix malformed opening tags like <p--> to <p>
	// This pattern matches tags like <p-->, <div-->, etc.
	html = malformedTagRegex.ReplaceAllString(html, "<$1>")

	// Fix malformed self-closing tags like <img-->, <br--> to <img>, <br>
	// Some feeds have broken self-closing tags with or without attributes
	// Pattern 1: Tags with attributes (e.g., <img src="..." -->)
	// Use [^<>]+ to avoid matching angle brackets and nested tags
	html = malformedSelfClosingWithAttrs.ReplaceAllString(html, "<$1 $2>")
	
	// Pattern 2: Tags without attributes (e.g., <br-->)
	html = malformedSelfClosingNoAttrs.ReplaceAllString(html, "<$1>")

	// Trim any leading/trailing whitespace
	html = strings.TrimSpace(html)

	return html
}
