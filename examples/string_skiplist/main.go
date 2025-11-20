package main

import (
	"fmt"
	"github.com/anchor54/SkipList"
)

func main() {
	// Create a skip list for strings
	sl := skiplist.NewSkipList[string]()

	// Add programming languages
	languages := []string{
		"Go", "Python", "JavaScript", "Rust", "Java",
		"C++", "TypeScript", "Ruby", "Swift", "Kotlin",
	}

	fmt.Println("Adding programming languages:")
	for _, lang := range languages {
		sl.Add(lang)
	}

	fmt.Printf("Total languages: %d\n\n", sl.Len())

	// Display all languages in alphabetical order
	fmt.Println("Languages in alphabetical order:")
	sl.Range(func(lang string) bool {
		fmt.Printf("  - %s\n", lang)
		return true
	})

	// Find languages starting with specific letters
	fmt.Println("\nSearching for languages:")
	searchTerms := []string{"Go", "Python", "Haskell", "Rust"}
	for _, term := range searchTerms {
		if sl.Contains(term) {
			fmt.Printf("  ✓ Found: %s\n", term)
		} else {
			fmt.Printf("  ✗ Not found: %s\n", term)
		}
	}

	// Get languages by position
	fmt.Println("\nLanguages by rank:")
	for rank := 1; rank <= 3; rank++ {
		if node, found := sl.SearchByRank(rank); found {
			fmt.Printf("  Rank %d: %s\n", rank, node.Value())
		}
	}

	// Remove some languages
	fmt.Println("\nRemoving JavaScript and Java...")
	sl.Delete("JavaScript")
	sl.Delete("Java")
	fmt.Printf("Remaining languages: %d\n", sl.Len())

	// Early termination example
	fmt.Println("\nFirst 5 languages after removal:")
	count := 0
	sl.Range(func(lang string) bool {
		count++
		fmt.Printf("  %d. %s\n", count, lang)
		return count < 5 // stop after 5 items
	})
}
