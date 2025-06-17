package utils

import (
	"strings"
	"unicode/utf8"
)

// DamerauLevenshteinDistance calculates the Damerau-Levenshtein distance between two strings
// This includes insertions, deletions, substitutions, and transpositions
func DamerauLevenshteinDistance(s1, s2 string) int {
	// Handle edge cases
	if s1 == s2 {
		return 0
	}
	if s1 == "" {
		return utf8.RuneCountInString(s2)
	}
	if s2 == "" {
		return utf8.RuneCountInString(s1)
	}

	len1 := utf8.RuneCountInString(s1)
	len2 := utf8.RuneCountInString(s2)

	// Convert strings to rune slices for proper Unicode handling
	runes1 := []rune(s1)
	runes2 := []rune(s2)

	// Create a map to store the last occurrence of each character
	da := make(map[rune]int)

	// Initialize the distance matrix
	maxLen := len1 + len2
	H := make([][]int, len1+2)
	for i := range H {
		H[i] = make([]int, len2+2)
	}

	H[0][0] = maxLen
	for i := 0; i <= len1; i++ {
		H[i+1][0] = maxLen
		H[i+1][1] = i
	}
	for j := 0; j <= len2; j++ {
		H[0][j+1] = maxLen
		H[1][j+1] = j
	}

	for i := 1; i <= len1; i++ {
		db := 0
		for j := 1; j <= len2; j++ {
			k := da[runes2[j-1]]
			l := db
			cost := 1
			if runes1[i-1] == runes2[j-1] {
				cost = 0
				db = j
			}

			H[i+1][j+1] = min(
				H[i][j]+cost,              // substitution
				H[i+1][j]+1,               // insertion
				H[i][j+1]+1,               // deletion
				H[k][l]+(i-k-1)+1+(j-l-1), // transposition
			)
		}
		da[runes1[i-1]] = i
	}

	return H[len1+1][len2+1]
}

// min returns the minimum of the given integers
func min(values ...int) int {
	if len(values) == 0 {
		return 0
	}
	minVal := values[0]
	for _, v := range values[1:] {
		if v < minVal {
			minVal = v
		}
	}
	return minVal
}

// FuzzyMatchScore calculates a similarity score between 0 and 1
// Higher scores indicate better matches
func FuzzyMatchScore(s1, s2 string) float64 {
	// Handle edge cases
	if s1 == "" && s2 == "" {
		return 1.0
	}
	if s1 == "" || s2 == "" {
		return 0.0
	}
	if s1 == s2 {
		return 1.0
	}

	maxLen := float64(max(utf8.RuneCountInString(s1), utf8.RuneCountInString(s2)))
	if maxLen == 0 {
		return 1.0
	}

	distance := float64(DamerauLevenshteinDistance(s1, s2))
	similarity := 1.0 - (distance / maxLen)

	// Ensure score is between 0 and 1
	if similarity < 0 {
		return 0.0
	}
	if similarity > 1 {
		return 1.0
	}

	return similarity
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// FuzzySearch performs fuzzy matching on a slice of strings
// Returns matches with their scores, sorted by relevance
type FuzzyMatch struct {
	Text  string
	Score float64
	Index int
}

// PerformFuzzySearch searches for fuzzy matches in a slice of strings
func PerformFuzzySearch(query string, texts []string, threshold float64) []FuzzyMatch {
	var matches []FuzzyMatch

	// Input validation
	if query == "" || len(texts) == 0 {
		return matches
	}

	query = strings.ToLower(strings.TrimSpace(query))

	for i, text := range texts {
		if text == "" {
			continue // Skip empty texts
		}

		textLower := strings.ToLower(text)
		bestScore := 0.0
		matched := false

		// Check for exact substring match first (highest priority)
		if strings.Contains(textLower, query) {
			matches = append(matches, FuzzyMatch{
				Text:  text,
				Score: 1.0,
				Index: i,
			})
			continue
		}

		// Check for fuzzy match on full text
		score := FuzzyMatchScore(query, textLower)
		if score >= threshold {
			bestScore = score
			matched = true
		}

		// Also check individual words in the text
		words := strings.Fields(textLower)
		for _, word := range words {
			if len(word) < 2 { // Skip very short words
				continue
			}
			wordScore := FuzzyMatchScore(query, word)
			if wordScore >= threshold && wordScore > bestScore {
				bestScore = wordScore
				matched = true
			}
		}

		// Add the match if we found one
		if matched {
			matches = append(matches, FuzzyMatch{
				Text:  text,
				Score: bestScore,
				Index: i,
			})
		}
	}

	// Sort matches by score (descending)
	for i := 0; i < len(matches)-1; i++ {
		for j := i + 1; j < len(matches); j++ {
			if matches[i].Score < matches[j].Score {
				matches[i], matches[j] = matches[j], matches[i]
			}
		}
	}

	return matches
}

// IsCloseMatch checks if two strings are close enough based on Damerau-Levenshtein distance
func IsCloseMatch(s1, s2 string, maxDistance int) bool {
	return DamerauLevenshteinDistance(strings.ToLower(s1), strings.ToLower(s2)) <= maxDistance
}

// SuggestCorrections suggests possible corrections for a misspelled word
func SuggestCorrections(word string, dictionary []string, maxSuggestions int, threshold float64) []string {
	matches := PerformFuzzySearch(word, dictionary, threshold)

	var suggestions []string
	count := 0
	for _, match := range matches {
		if count >= maxSuggestions {
			break
		}
		suggestions = append(suggestions, match.Text)
		count++
	}

	return suggestions
}
