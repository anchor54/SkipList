package main

import (
	"fmt"

	skiplist "github.com/anchor54/SkipList"
)

func main() {
	// Create a new skip list for integers
	sl := skiplist.NewSkipList[int]()

	// Add elements
	fmt.Println("Adding elements: 10, 5, 20, 15, 30")
	sl.Add(10)
	sl.Add(5)
	sl.Add(20)
	sl.Add(15)
	sl.Add(30)

	// Check length
	fmt.Printf("Skip list length: %d\n", sl.Len())

	// Search for a value
	if node, found := sl.SearchByValue(15); found {
		fmt.Printf("Found value 15 in the skip list\n")
		_ = node
	}

	// Check if a value exists
	if sl.Contains(25) {
		fmt.Println("25 exists in the skip list")
	} else {
		fmt.Println("25 does not exist in the skip list")
	}

	// Iterate over all elements
	fmt.Print("Elements in order: ")
	sl.Range(func(val int) bool {
		fmt.Printf("%d ", val)
		return true // continue iteration
	})
	fmt.Println()

	// Search by rank (position)
	if node, found := sl.SearchByRank(3); found {
		fmt.Printf("The 3rd element is: %d\n", node.Value())
	}

	// Delete an element
	sl.Delete(15)
	fmt.Printf("After deleting 15, length: %d\n", sl.Len())

	// Iterate again
	fmt.Print("Elements after deletion: ")
	sl.Range(func(val int) bool {
		fmt.Printf("%d ", val)
		return true
	})
	fmt.Println()

	// Clear the skip list
	sl.Clear()
	fmt.Printf("After clearing, length: %d, isEmpty: %v\n", sl.Len(), sl.IsEmpty())
}
