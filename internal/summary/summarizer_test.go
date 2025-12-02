package summary

import (
	"strings"
	"testing"
)

func TestNewSummarizer(t *testing.T) {
	s := NewSummarizer()
	if s == nil {
		t.Error("NewSummarizer should not return nil")
	}
}

func TestSummarize_ShortText(t *testing.T) {
	s := NewSummarizer()
	result := s.Summarize("Short text.", Short)

	if !result.IsTooShort {
		t.Error("Expected IsTooShort to be true for short text")
	}
}

func TestSummarize_MediumLength(t *testing.T) {
	s := NewSummarizer()

	// Test with a longer text
	text := `Natural language processing is a field of artificial intelligence. It focuses on the interaction between computers and humans using natural language. The ultimate goal is to enable computers to understand, interpret, and generate human language. NLP combines computational linguistics with machine learning and deep learning. Applications include machine translation, sentiment analysis, and text summarization. Modern NLP uses transformer models that have revolutionized the field. These models can process text more effectively than previous approaches. The field continues to advance rapidly with new techniques being developed.`

	result := s.Summarize(text, Medium)

	if result.IsTooShort {
		t.Error("Expected IsTooShort to be false for medium-length text")
	}

	if result.Summary == "" {
		t.Error("Expected non-empty summary")
	}

	if len(result.KeyPoints) == 0 {
		t.Error("Expected some key points")
	}

	// Summary should be shorter than original
	if len(result.Summary) >= len(text) {
		t.Error("Summary should be shorter than original text")
	}
}

func TestSummarize_DifferentLengths(t *testing.T) {
	s := NewSummarizer()

	text := `Natural language processing is a field of artificial intelligence. It focuses on the interaction between computers and humans using natural language. The ultimate goal is to enable computers to understand, interpret, and generate human language. NLP combines computational linguistics with machine learning and deep learning. Applications include machine translation, sentiment analysis, and text summarization. Modern NLP uses transformer models that have revolutionized the field. These models can process text more effectively than previous approaches. The field continues to advance rapidly with new techniques being developed. Research institutions and companies are investing heavily in NLP research. The technology has many practical applications in everyday life.`

	shortResult := s.Summarize(text, Short)
	mediumResult := s.Summarize(text, Medium)
	longResult := s.Summarize(text, Long)

	// Short should have fewer sentences than medium
	if shortResult.SentenceCount >= mediumResult.SentenceCount {
		t.Errorf("Short summary (%d sentences) should have fewer sentences than medium (%d sentences)",
			shortResult.SentenceCount, mediumResult.SentenceCount)
	}

	// Medium should have fewer sentences than long
	if mediumResult.SentenceCount >= longResult.SentenceCount {
		t.Errorf("Medium summary (%d sentences) should have fewer sentences than long (%d sentences)",
			mediumResult.SentenceCount, longResult.SentenceCount)
	}
}

func TestSummarize_HTMLContent(t *testing.T) {
	s := NewSummarizer()

	htmlText := `<p>Natural language processing is a field of artificial intelligence.</p>
		<p>It focuses on the interaction between computers and humans using natural language.</p>
		<p>The ultimate goal is to enable computers to understand, interpret, and generate human language.</p>
		<p>NLP combines computational linguistics with machine learning and deep learning.</p>
		<p>Applications include machine translation, sentiment analysis, and text summarization.</p>
		<p>Modern NLP uses transformer models that have revolutionized the field.</p>
		<p>These models can process text more effectively than previous approaches.</p>
		<p>The field continues to advance rapidly with new techniques being developed.</p>`

	result := s.Summarize(htmlText, Medium)

	// Summary should not contain HTML tags
	if strings.Contains(result.Summary, "<p>") || strings.Contains(result.Summary, "</p>") {
		t.Error("Summary should not contain HTML tags")
	}

	if result.Summary == "" {
		t.Error("Expected non-empty summary for HTML content")
	}
}

func TestSummarize_ChineseText(t *testing.T) {
	s := NewSummarizer()

	// Longer Chinese text with more sentences
	chineseText := `自然语言处理是人工智能的一个重要领域。它专注于计算机与人类使用自然语言进行交互。最终目标是使计算机能够理解、解释和生成人类语言。自然语言处理结合了计算语言学与机器学习和深度学习。应用包括机器翻译、情感分析和文本摘要。现代自然语言处理使用已经彻底改变该领域的Transformer模型。这些模型可以比以前的方法更有效地处理文本。该领域正在快速发展，不断有新技术被开发出来。研究人员不断探索新的算法和方法。企业和学术界都在大力投资这一领域。`

	result := s.Summarize(chineseText, Medium)

	// Chinese text should generate a summary (may or may not be marked as too short depending on sentence count)
	if result.Summary == "" {
		t.Error("Expected non-empty summary for Chinese content")
	}
}

func TestCleanText(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "<p>Hello</p> <strong>World</strong>",
			expected: "Hello World",
		},
		{
			input:    "Hello&nbsp;World",
			expected: "Hello World",
		},
		{
			input:    "a &amp; b",
			expected: "a & b",
		},
		{
			input:    "  multiple   spaces  ",
			expected: "multiple spaces",
		},
	}

	for _, tt := range tests {
		result := cleanText(tt.input)
		if result != tt.expected {
			t.Errorf("cleanText(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestSplitSentences(t *testing.T) {
	text := "First sentence here. Second sentence here! Third sentence here?"
	sentences := splitSentences(text)

	if len(sentences) < 1 {
		t.Errorf("Expected at least 1 sentence, got %d", len(sentences))
	}
}

func TestTokenize(t *testing.T) {
	text := "The quick brown fox jumps over the lazy dog"
	tokens := tokenize(text)

	// Should not contain common stopwords like "the"
	for _, token := range tokens {
		if token == "the" {
			t.Errorf("tokenize should remove stopword: %s", token)
		}
	}

	// Should contain meaningful words like "quick", "brown", "fox"
	found := make(map[string]bool)
	for _, token := range tokens {
		found[token] = true
	}

	if !found["quick"] || !found["brown"] || !found["fox"] {
		t.Error("tokenize should retain meaningful words")
	}
}

func TestIsStopWord(t *testing.T) {
	stopWords := []string{"the", "a", "an", "and", "or", "in", "on", "at", "的", "了", "和"}
	nonStopWords := []string{"computer", "algorithm", "processing", "模型", "技术"}

	for _, word := range stopWords {
		if !isStopWord(word) {
			t.Errorf("%q should be a stopword", word)
		}
	}

	for _, word := range nonStopWords {
		if isStopWord(word) {
			t.Errorf("%q should not be a stopword", word)
		}
	}
}

func TestCalculateNumSentences(t *testing.T) {
	tests := []struct {
		totalSentences int
		length         SummaryLength
		minExpected    int
		maxExpected    int
	}{
		{10, Short, 1, 3},
		{10, Medium, 3, 5},
		{10, Long, 4, 6},
		{1, Short, 1, 1},
		{100, Short, 15, 25},
	}

	for _, tt := range tests {
		result := calculateNumSentences(tt.totalSentences, tt.length)
		if result < tt.minExpected || result > tt.maxExpected {
			t.Errorf("calculateNumSentences(%d, %s) = %d, expected between %d and %d",
				tt.totalSentences, tt.length, result, tt.minExpected, tt.maxExpected)
		}
	}
}

func TestSentenceSimilarity(t *testing.T) {
	// Similar sentences should have higher similarity
	s1 := "Natural language processing is important"
	s2 := "Natural language processing is essential"
	s3 := "Cooking recipes are delicious"

	sim12 := sentenceSimilarity(s1, s2)
	sim13 := sentenceSimilarity(s1, s3)

	if sim12 <= sim13 {
		t.Errorf("Similar sentences should have higher similarity: sim(%q, %q)=%f <= sim(%q, %q)=%f",
			s1, s2, sim12, s1, s3, sim13)
	}
}

func TestCalculateTFIDF(t *testing.T) {
	sentences := []string{
		"Natural language processing is important.",
		"Machine learning is a key technology.",
		"Natural language processing uses machine learning.",
	}

	scores := calculateTFIDF(sentences)

	if len(scores) != len(sentences) {
		t.Errorf("Expected %d scores, got %d", len(sentences), len(scores))
	}

	// All scores should be between 0 and 1 (normalized)
	for i, score := range scores {
		if score < 0 || score > 1 {
			t.Errorf("Score %d should be between 0 and 1, got %f", i, score)
		}
	}
}

func TestCalculateTextRank(t *testing.T) {
	sentences := []string{
		"Natural language processing is important.",
		"Machine learning is a key technology.",
		"Natural language processing uses machine learning.",
	}

	scores := calculateTextRank(sentences)

	if len(scores) != len(sentences) {
		t.Errorf("Expected %d scores, got %d", len(sentences), len(scores))
	}

	// All scores should be between 0 and 1 (normalized)
	for i, score := range scores {
		if score < 0 || score > 1 {
			t.Errorf("Score %d should be between 0 and 1, got %f", i, score)
		}
	}
}

func TestSummarize_EmptyText(t *testing.T) {
	s := NewSummarizer()
	result := s.Summarize("", Short)

	if !result.IsTooShort {
		t.Error("Expected IsTooShort to be true for empty text")
	}
}

func TestSummarize_SingleSentence(t *testing.T) {
	s := NewSummarizer()
	result := s.Summarize("This is a single sentence that is long enough to pass the minimum length check.", Short)

	if !result.IsTooShort {
		t.Error("Expected IsTooShort to be true for single sentence")
	}
}
