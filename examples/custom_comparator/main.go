package main

import (
	"fmt"

	skiplist "github.com/anchor54/SkipList"
)

// Person represents a person with name and age
type Person struct {
	Name string
	Age  int
}

// Implement the Comparable interface
func (p Person) Compare(other Person) int {
	// Compare by age first, then by name
	if p.Age < other.Age {
		return -1
	} else if p.Age > other.Age {
		return 1
	}

	// If ages are equal, compare by name
	if p.Name < other.Name {
		return -1
	} else if p.Name > other.Name {
		return 1
	}
	return 0
}

func main() {
	// Create a skip list for Person objects
	sl := skiplist.NewComparableSkipList[Person]()

	// Add people
	people := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 35},
		{"David", 25},
		{"Eve", 30},
	}

	fmt.Println("Adding people to skip list:")
	for _, person := range people {
		sl.Add(person)
		fmt.Printf("  Added: %s (age %d)\n", person.Name, person.Age)
	}

	// Display all people in sorted order
	fmt.Println("\nPeople in sorted order (by age, then name):")
	sl.Range(func(p Person) bool {
		fmt.Printf("  %s - Age: %d\n", p.Name, p.Age)
		return true
	})

	// Search for a specific person
	searchPerson := Person{"Bob", 25}
	if sl.Contains(searchPerson) {
		fmt.Printf("\nFound %s (age %d) in the skip list\n", searchPerson.Name, searchPerson.Age)
	}

	// Get the youngest person (rank 1)
	if node, found := sl.SearchByRank(1); found {
		person := node.Value()
		fmt.Printf("\nYoungest person: %s (age %d)\n", person.Name, person.Age)
	}

	// Get the oldest person (last rank)
	if node, found := sl.SearchByRank(sl.Len()); found {
		person := node.Value()
		fmt.Printf("Oldest person: %s (age %d)\n", person.Name, person.Age)
	}

	// Remove a person
	sl.Delete(Person{"Charlie", 35})
	fmt.Printf("\nAfter removing Charlie, skip list length: %d\n", sl.Len())
}
