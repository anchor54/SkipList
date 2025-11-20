package skiplist

import (
	"math/rand"
	"testing"
)

// ------------------------------------------------------------
// Benchmark Functions
// ------------------------------------------------------------

// BenchmarkAdd benchmarks adding elements to an empty skip list
func BenchmarkAdd(b *testing.B) {
	sl := NewSkipList[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sl.Add(i)
	}
}

// BenchmarkAdd_Existing benchmarks adding elements to a skip list with existing elements
func BenchmarkAdd_Existing(b *testing.B) {
	sl := NewSkipList[int]()
	// Pre-populate with 1000 elements
	for i := 0; i < 1000; i++ {
		sl.Add(i * 2) // Even numbers
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sl.Add(i*2 + 1) // Odd numbers (new elements)
	}
}

// BenchmarkSearch benchmarks searching for elements in a skip list
func BenchmarkSearch(b *testing.B) {
	sl := NewSkipList[int]()
	// Pre-populate with elements
	size := 10000
	for i := 0; i < size; i++ {
		sl.Add(i)
	}

	// Create random search values
	searchValues := make([]int, b.N)
	for i := range searchValues {
		searchValues[i] = rand.Intn(size)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = sl.SearchByValue(searchValues[i])
	}
}

// BenchmarkSearch_Found benchmarks searching for elements that exist
func BenchmarkSearch_Found(b *testing.B) {
	sl := NewSkipList[int]()
	size := 10000
	for i := 0; i < size; i++ {
		sl.Add(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = sl.SearchByValue(i % size)
	}
}

// BenchmarkSearch_NotFound benchmarks searching for elements that don't exist
func BenchmarkSearch_NotFound(b *testing.B) {
	sl := NewSkipList[int]()
	size := 10000
	// Add even numbers only
	for i := 0; i < size; i++ {
		sl.Add(i * 2)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = sl.SearchByValue(i*2 + 1) // Search for odd numbers (not present)
	}
}

// BenchmarkDelete benchmarks deleting elements from a skip list
func BenchmarkDelete(b *testing.B) {
	sl := NewSkipList[int]()
	size := 10000
	// Pre-populate
	for i := 0; i < size; i++ {
		sl.Add(i)
	}

	// Create random delete values
	deleteValues := make([]int, b.N)
	for i := range deleteValues {
		deleteValues[i] = rand.Intn(size)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sl.Delete(deleteValues[i])
	}
}

// BenchmarkRankSearch benchmarks searching by rank
func BenchmarkRankSearch(b *testing.B) {
	sl := NewSkipList[int]()
	size := 10000
	// Pre-populate
	for i := 0; i < size; i++ {
		sl.Add(i)
	}

	// Create random ranks
	ranks := make([]int, b.N)
	for i := range ranks {
		ranks[i] = rand.Intn(size) + 1 // Rank is 1-indexed
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = sl.SearchByRank(ranks[i])
	}
}

// BenchmarkContains benchmarks the Contains method
func BenchmarkContains(b *testing.B) {
	sl := NewSkipList[int]()
	size := 10000
	// Pre-populate
	for i := 0; i < size; i++ {
		sl.Add(i)
	}

	// Create random search values
	searchValues := make([]int, b.N)
	for i := range searchValues {
		searchValues[i] = rand.Intn(size)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sl.Contains(searchValues[i])
	}
}

// BenchmarkRange benchmarks iterating over all elements
func BenchmarkRange(b *testing.B) {
	sl := NewSkipList[int]()
	size := 1000
	// Pre-populate
	for i := 0; i < size; i++ {
		sl.Add(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sl.Range(func(val int) bool {
			_ = val
			return true
		})
	}
}

// BenchmarkLen benchmarks the Len method
func BenchmarkLen(b *testing.B) {
	sl := NewSkipList[int]()
	size := 10000
	for i := 0; i < size; i++ {
		sl.Add(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sl.Len()
	}
}

// ------------------------------------------------------------
// Size-based Benchmarks
// ------------------------------------------------------------

// BenchmarkAdd_Size10 benchmarks Add with 10 elements
func BenchmarkAdd_Size10(b *testing.B) {
	benchmarkAddWithSize(b, 10)
}

// BenchmarkAdd_Size100 benchmarks Add with 100 elements
func BenchmarkAdd_Size100(b *testing.B) {
	benchmarkAddWithSize(b, 100)
}

// BenchmarkAdd_Size1000 benchmarks Add with 1000 elements
func BenchmarkAdd_Size1000(b *testing.B) {
	benchmarkAddWithSize(b, 1000)
}

// BenchmarkAdd_Size10000 benchmarks Add with 10000 elements
func BenchmarkAdd_Size10000(b *testing.B) {
	benchmarkAddWithSize(b, 10000)
}

func benchmarkAddWithSize(b *testing.B, size int) {
	sl := NewSkipList[int]()
	// Pre-populate
	for i := 0; i < size; i++ {
		sl.Add(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sl.Add(size + i)
	}
}

// BenchmarkSearch_Size10 benchmarks Search with 10 elements
func BenchmarkSearch_Size10(b *testing.B) {
	benchmarkSearchWithSize(b, 10)
}

// BenchmarkSearch_Size100 benchmarks Search with 100 elements
func BenchmarkSearch_Size100(b *testing.B) {
	benchmarkSearchWithSize(b, 100)
}

// BenchmarkSearch_Size1000 benchmarks Search with 1000 elements
func BenchmarkSearch_Size1000(b *testing.B) {
	benchmarkSearchWithSize(b, 1000)
}

// BenchmarkSearch_Size10000 benchmarks Search with 10000 elements
func BenchmarkSearch_Size10000(b *testing.B) {
	benchmarkSearchWithSize(b, 10000)
}

func benchmarkSearchWithSize(b *testing.B, size int) {
	sl := NewSkipList[int]()
	// Pre-populate
	for i := 0; i < size; i++ {
		sl.Add(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = sl.SearchByValue(i % size)
	}
}

// ------------------------------------------------------------
// Mixed Operation Benchmarks
// ------------------------------------------------------------

// BenchmarkMixedOperations benchmarks a mix of add, search, and delete operations
func BenchmarkMixedOperations(b *testing.B) {
	sl := NewSkipList[int]()
	size := 10000

	// Pre-populate
	for i := 0; i < size; i++ {
		sl.Add(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Mix of operations
		sl.Add(size + i)
		_, _ = sl.SearchByValue(i % size)
		sl.Delete(i % size)
	}
}

// BenchmarkSequentialAdd benchmarks adding elements in sequential order
func BenchmarkSequentialAdd(b *testing.B) {
	sl := NewSkipList[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sl.Add(i)
	}
}

// BenchmarkRandomAdd benchmarks adding elements in random order
func BenchmarkRandomAdd(b *testing.B) {
	sl := NewSkipList[int]()
	values := make([]int, b.N)
	for i := range values {
		values[i] = rand.Intn(b.N * 2)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sl.Add(values[i])
	}
}

// BenchmarkStringOperations benchmarks string operations
func BenchmarkStringOperations(b *testing.B) {
	sl := NewSkipList[string]()

	// Pre-populate with strings
	strings := []string{"apple", "banana", "cherry", "date", "elderberry"}
	for i := 0; i < 1000; i++ {
		sl.Add(strings[i%len(strings)] + string(rune(i)))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		searchVal := strings[i%len(strings)] + string(rune(i%1000))
		_, _ = sl.SearchByValue(searchVal)
	}
}
