# Skip List - A Generic Probabilistic Data Structure in Go

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.25-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

A high-performance, generic skip list implementation in Go with support for custom comparators, rank-based queries, and comprehensive testing.

## 📚 Table of Contents

- [Features](#-features)
- [Installation](#-installation)
- [Quick Start](#-quick-start)
- [API Reference](#-api-reference)
- [Examples](#-examples)
- [Performance](#-performance)
- [Implementation Details](#-implementation-details)
- [Testing](#-testing)
- [Contributing](#-contributing)

## ✨ Features

- **Generic Implementation**: Works with any type that satisfies Go's `cmp.Ordered` constraint or implements a custom `Comparable` interface
- **Set Semantics**: No duplicate values allowed
- **Rank-Based Queries**: Search by position (1-indexed) with O(log n) complexity
- **Custom Comparators**: Support for complex custom comparison logic
- **Comprehensive API**: Methods for add, delete, search, contains, range iteration, and more
- **Well-Tested**: 2000+ lines of comprehensive tests covering edge cases and correctness
- **Deterministic Testing**: Support for inserting at specific levels for reproducible tests

## 📦 Installation

```bash
go get github.com/anchor54/SkipList
```

## 🚀 Quick Start

```go
package main

import (
    "fmt"
    "github.com/anchor54/SkipList"
)

func main() {
    // Create a new skip list for integers
    sl := skiplist.NewSkipList[int]()

    // Add elements
    sl.Add(10)
    sl.Add(5)
    sl.Add(20)
    sl.Add(15)

    // Search by value
    if node, found := sl.SearchByValue(15); found {
        fmt.Println("Found:", node)
    }

    // Search by rank (position)
    if node, found := sl.SearchByRank(3); found {
        fmt.Println("3rd element:", node)
    }

    // Iterate over all elements
    sl.Range(func(val int) bool {
        fmt.Printf("%d ", val)
        return true // continue iteration
    })
    // Output: 5 10 15 20

    // Check if value exists
    exists := sl.Contains(15) // true

    // Get length
    length := sl.Len() // 4

    // Delete an element
    sl.Delete(15)

    // Clear all elements
    sl.Clear()
}
```

## 📖 API Reference

### Creating Skip Lists

#### `NewSkipList[T cmp.Ordered]() *SkipList[T]`
Creates a new skip list for types that satisfy the `cmp.Ordered` constraint (int, float, string, etc.).

```go
sl := skiplist.NewSkipList[int]()
sl := skiplist.NewSkipList[string]()
sl := skiplist.NewSkipList[float64]()
```

#### `NewComparableSkipList[T Comparable[T]]() *SkipList[T]`
Creates a new skip list for custom types that implement the `Comparable` interface.

```go
type Person struct {
    Name string
    Age  int
}

func (p Person) Compare(other Person) int {
    if p.Age < other.Age {
        return -1
    } else if p.Age > other.Age {
        return 1
    }
    return 0
}

sl := skiplist.NewComparableSkipList[Person]()
```

### Core Operations

#### `Add(val T)`
Adds a value to the skip list. Duplicates are ignored.

**Time Complexity**: O(log n) average

```go
sl.Add(10)
sl.Add(5)
sl.Add(10) // ignored - no duplicates
```

#### `Delete(val T)`
Removes a value from the skip list. No-op if value doesn't exist.

**Time Complexity**: O(log n) average

```go
sl.Delete(10)
```

#### `SearchByValue(val T) (*Node[T], bool)`
Searches for a value and returns the node and a boolean indicating if found.

**Time Complexity**: O(log n) average

```go
if node, found := sl.SearchByValue(10); found {
    fmt.Println("Found:", node.val)
}
```

#### `SearchByRank(rank int) (*Node[T], bool)`
Searches for the element at a given position (1-indexed). Returns nil if rank is out of bounds.

**Time Complexity**: O(log n)

```go
// Get the 3rd smallest element
if node, found := sl.SearchByRank(3); found {
    fmt.Println("3rd element:", node.val)
}
```

### Utility Methods

#### `Len() int`
Returns the number of elements in the skip list.

```go
length := sl.Len()
```

#### `IsEmpty() bool`
Returns true if the skip list is empty.

```go
if sl.IsEmpty() {
    fmt.Println("Skip list is empty")
}
```

#### `Contains(val T) bool`
Checks if a value exists in the skip list.

```go
exists := sl.Contains(10)
```

#### `Range(fn func(val T) bool)`
Iterates over all elements in ascending order. Iteration stops if the function returns false.

```go
// Print all elements
sl.Range(func(val int) bool {
    fmt.Println(val)
    return true // continue
})

// Print first 5 elements
count := 0
sl.Range(func(val int) bool {
    fmt.Println(val)
    count++
    return count < 5 // stop after 5
})
```

#### `Clear()`
Removes all elements from the skip list.

```go
sl.Clear()
```

### Testing Methods

#### `InsertAtLevel(val T, level int)`
Inserts a value at a specific level. Useful for deterministic testing.

```go
sl.InsertAtLevel(10, 2) // Insert at level 2
```

## 💡 Examples

### Basic Usage

See [examples/basic_usage.go](examples/basic_usage.go) for a complete example.

### Custom Comparator

```go
type Person struct {
    Name string
    Age  int
}

func (p Person) Compare(other Person) int {
    // Compare by age first, then by name
    if p.Age != other.Age {
        return cmp.Compare(p.Age, other.Age)
    }
    return cmp.Compare(p.Name, other.Name)
}

sl := skiplist.NewComparableSkipList[Person]()
sl.Add(Person{"Alice", 30})
sl.Add(Person{"Bob", 25})
```

See [examples/custom_comparator.go](examples/custom_comparator.go) for a complete example.

### String Skip List

See [examples/string_skiplist.go](examples/string_skiplist.go) for a complete example.

## ⚡ Performance

### Time Complexity

| Operation | Average | Worst Case |
|-----------|---------|------------|
| Add       | O(log n) | O(n) |
| Delete    | O(log n) | O(n) |
| SearchByValue | O(log n) | O(n) |
| SearchByRank | O(log n) | O(n) |
| Contains  | O(log n) | O(n) |

### Space Complexity

- **Average**: O(n)
- **Worst Case**: O(n × MAX_LEVEL_CAP) where MAX_LEVEL_CAP = 16

### Benchmark Results

Run benchmarks with:

```bash
# Run all benchmarks
go test -bench=. -benchmem

# Run specific benchmark
go test -bench=BenchmarkAdd -benchmem

# Run benchmarks with more iterations (for better accuracy)
go test -bench=. -benchmem -benchtime=5s

# Run benchmarks excluding tests
go test -bench=. -benchmem -run=^$
```

**Actual benchmark results** (AMD Ryzen 7 5800H, Go 1.25):

```
BenchmarkAdd-16                 3537198    449.8 ns/op    96 B/op    3 allocs/op
BenchmarkSearch-16              6044464    196.5 ns/op     0 B/op    0 allocs/op
BenchmarkSearch_Found-16        9906860    116.5 ns/op     0 B/op    0 allocs/op
BenchmarkSearch_NotFound-16     22851115     75.50 ns/op    0 B/op    0 allocs/op
BenchmarkDelete-16            281995209      4.162 ns/op   0 B/op    0 allocs/op
BenchmarkRankSearch-16          9261488    133.5 ns/op     0 B/op    0 allocs/op
BenchmarkContains-16            6286678    201.3 ns/op     0 B/op    0 allocs/op
BenchmarkRange-16                506406   2285 ns/op       0 B/op    0 allocs/op
BenchmarkLen-16             1000000000      0.2600 ns/op   0 B/op    0 allocs/op
```

**Size-based benchmarks** (showing how performance scales):

```
BenchmarkSearch_Size10-16       57187696    25.05 ns/op    0 B/op    0 allocs/op
BenchmarkSearch_Size100-16      27296677    74.88 ns/op    0 B/op    0 allocs/op
BenchmarkSearch_Size1000-16     14100397    89.81 ns/op    0 B/op    0 allocs/op
BenchmarkSearch_Size10000-16    10060653   120.8 ns/op     0 B/op    0 allocs/op
```

**Key observations**:
- ✅ `Len()` is O(1) - extremely fast (0.26 ns/op)
- ✅ Search operations scale well with size (logarithmic)
- ✅ `Delete` is very efficient when elements exist
- ✅ `Range` iteration is O(n) as expected
- ✅ Memory allocations are minimal (most operations have 0 allocs/op)

## 🔧 Implementation Details

### Skip List Structure

A skip list is a layered data structure where:
- **Level 0** contains all elements in sorted order
- **Higher levels** contain subsets of elements, forming "express lanes"
- Each element has a randomly determined height
- Search starts from the highest level and drops down to lower levels

### Probability and Levels

- **Probability**: 0.5 (50% chance of promoting to next level)
- **Max Level**: 16 (configurable via `MAX_LEVEL_CAP`)
- Expected height for n elements: log₂(n)

### Node Structure

Each node contains:
- `val`: The stored value
- `forward`: Array of forward pointers (one per level)
- `skips`: Array of skip distances (for rank-based queries)

### Rank-Based Search

The skip list maintains skip distances at each level, enabling O(log n) rank-based queries:
- Each forward pointer stores how many elements it skips
- Accumulated skips provide the rank during traversal

## 🧪 Testing

The project includes comprehensive tests organized into modules:

- **helpers_test.go**: Test helper functions
- **basic_operations_test.go**: Core add, delete, search tests
- **insert_level_test.go**: Deterministic insertion tests
- **delete_spans_test.go**: Skip distance correctness tests
- **rank_operations_test.go**: Rank-based query tests

### Running Tests

```bash
# Run all tests
go test -v

# Run specific test file
go test -v -run TestAdd

# Run with coverage
go test -cover

# Generate coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Running Benchmarks

```bash
# Run all benchmarks
go test -bench=. -benchmem

# Run specific benchmark
go test -bench=BenchmarkAdd -benchmem

# Run benchmarks with longer duration (more accurate)
go test -bench=. -benchmem -benchtime=5s

# Run benchmarks excluding unit tests
go test -bench=. -benchmem -run=^$

# Compare benchmarks (useful for performance regression testing)
go test -bench=. -benchmem -count=5 > old.txt
# ... make changes ...
go test -bench=. -benchmem -count=5 > new.txt
benchcmp old.txt new.txt
```

**Note**: Install `benchcmp` for comparing benchmark results:
```bash
go install golang.org/x/tools/cmd/benchcmp@latest
```

### Test Coverage

The test suite includes:
- ✅ 2000+ lines of tests
- ✅ Edge cases (empty lists, single elements, duplicates)
- ✅ Correctness verification (ordering, skip distances)
- ✅ Mixed operations (add/delete sequences)
- ✅ Rank-based query validation

## 📂 Project Structure

```
skiplist/
├── skiplist.go                  # Core implementation
├── basic_operations_test.go     # Basic tests
├── insert_level_test.go         # Deterministic insertion tests
├── delete_spans_test.go         # Skip distance tests
├── rank_operations_test.go      # Rank query tests
├── helpers_test.go              # Test utilities
├── benchmark_test.go            # Performance benchmarks
├── examples/
│   ├── basic_usage.go          # Basic usage example
│   ├── custom_comparator.go    # Custom type example
│   └── string_skiplist.go      # String skiplist example
├── go.mod
└── README.md
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit issues, fork the repository, and create pull requests.

### Development Setup

```bash
# Clone the repository
git clone https://github.com/anchor54/SkipList.git
cd SkipList

# Run tests
go test -v

# Run with race detector
go test -race -v

# Format code
go fmt ./...

# Run linter (if golangci-lint is installed)
golangci-lint run ./...
```

### Pre-commit Hook

This repository includes a git pre-commit hook that automatically runs:
1. **Code formatting** (`go fmt`)
2. **Static analysis** (`go vet`)
3. **Linting** (`golangci-lint` if installed)
4. **Tests** (`go test`)

The hook is located at `.git/hooks/pre-commit` and will run automatically when you commit.

**To manually test the hook:**
```bash
.git/hooks/pre-commit
```

**To bypass the hook (not recommended):**
```bash
git commit --no-verify
```

**Installing golangci-lint (optional but recommended):**
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

The project includes a `.golangci.yml` configuration file for consistent linting rules.

## 📝 License

This project is licensed under the MIT License. See the LICENSE file for details.

## 🙏 Acknowledgments

- Based on the original skip list paper by William Pugh (1990)
- Inspired by Redis sorted sets implementation

## 📚 References

- [Skip Lists: A Probabilistic Alternative to Balanced Trees](https://15721.courses.cs.cmu.edu/spring2018/papers/08-oltpindexes1/pugh-skiplists-cacm1990.pdf) - William Pugh
- [Skip List on Wikipedia](https://en.wikipedia.org/wiki/Skip_list)

---

**Made with ❤️ in Go**

