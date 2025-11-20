package skiplist

import "testing"

// ------------------------------------------------------------
// SearchByRank Test cases
// ------------------------------------------------------------

// TestSearchByRank_Basic tests basic search by rank functionality
func TestSearchByRank_Basic(t *testing.T) {
	s := NewSkipList[int]()
	s.InsertAtLevel(10, 1)
	s.InsertAtLevel(20, 2)
	s.InsertAtLevel(30, 1)

	// Rank 1 should be 10
	node, found := s.SearchByRank(1)
	if !found {
		t.Fatalf("SearchByRank(1) should find a node")
	}
	if node == nil || node.val != 10 {
		t.Fatalf("SearchByRank(1) should return node with value 10, got %v", node)
	}

	// Rank 2 should be 20
	node, found = s.SearchByRank(2)
	if !found {
		t.Fatalf("SearchByRank(2) should find a node")
	}
	if node == nil || node.val != 20 {
		t.Fatalf("SearchByRank(2) should return node with value 20, got %v", node)
	}

	// Rank 3 should be 30
	node, found = s.SearchByRank(3)
	if !found {
		t.Fatalf("SearchByRank(3) should find a node")
	}
	if node == nil || node.val != 30 {
		t.Fatalf("SearchByRank(3) should return node with value 30, got %v", node)
	}
}

// TestSearchByRank_AllRanks tests searching for all ranks in a list
func TestSearchByRank_AllRanks(t *testing.T) {
	s := NewSkipList[int]()
	values := []int{10, 20, 30, 40, 50}
	for i, val := range values {
		s.InsertAtLevel(val, i%3) // Vary levels
	}

	// Test all ranks
	for rank := 1; rank <= len(values); rank++ {
		node, found := s.SearchByRank(rank)
		if !found {
			t.Fatalf("SearchByRank(%d) should find a node", rank)
		}
		if node == nil {
			t.Fatalf("SearchByRank(%d) should not return nil", rank)
		}
		expectedVal := values[rank-1]
		if node.val != expectedVal {
			t.Fatalf("SearchByRank(%d) should return node with value %d, got %d", rank, expectedVal, node.val)
		}
	}
}

// TestSearchByRank_InvalidRanks tests invalid rank values
func TestSearchByRank_InvalidRanks(t *testing.T) {
	s := NewSkipList[int]()
	s.InsertAtLevel(10, 1)
	s.InsertAtLevel(20, 1)
	s.InsertAtLevel(30, 1)

	// Rank 0 should fail
	node, found := s.SearchByRank(0)
	if found {
		t.Fatalf("SearchByRank(0) should not find a node")
	}
	if node != nil {
		t.Fatalf("SearchByRank(0) should return nil, got %v", node)
	}

	// Negative rank should fail
	node, found = s.SearchByRank(-1)
	if found {
		t.Fatalf("SearchByRank(-1) should not find a node")
	}
	if node != nil {
		t.Fatalf("SearchByRank(-1) should return nil, got %v", node)
	}

	// Rank greater than length should fail
	node, found = s.SearchByRank(4)
	if found {
		t.Fatalf("SearchByRank(4) should not find a node")
	}
	if node != nil {
		t.Fatalf("SearchByRank(4) should return nil, got %v", node)
	}

	// Very large rank should fail
	node, found = s.SearchByRank(1000)
	if found {
		t.Fatalf("SearchByRank(1000) should not find a node")
	}
	if node != nil {
		t.Fatalf("SearchByRank(1000) should return nil, got %v", node)
	}
}

// TestSearchByRank_EmptyList tests search on empty list
func TestSearchByRank_EmptyList(t *testing.T) {
	s := NewSkipList[int]()

	// Any rank should fail on empty list
	for rank := 1; rank <= 10; rank++ {
		node, found := s.SearchByRank(rank)
		if found {
			t.Fatalf("SearchByRank(%d) should not find a node in empty list", rank)
		}
		if node != nil {
			t.Fatalf("SearchByRank(%d) should return nil in empty list, got %v", rank, node)
		}
	}
}

// TestSearchByRank_SingleElement tests search on single element list
func TestSearchByRank_SingleElement(t *testing.T) {
	s := NewSkipList[int]()
	s.InsertAtLevel(42, 1)

	// Rank 1 should work
	node, found := s.SearchByRank(1)
	if !found {
		t.Fatalf("SearchByRank(1) should find a node")
	}
	if node == nil || node.val != 42 {
		t.Fatalf("SearchByRank(1) should return node with value 42, got %v", node)
	}

	// Rank 2 should fail
	node, found = s.SearchByRank(2)
	if found {
		t.Fatalf("SearchByRank(2) should not find a node")
	}
	if node != nil {
		t.Fatalf("SearchByRank(2) should return nil, got %v", node)
	}
}

// TestSearchByRank_AfterInsertions tests search after multiple insertions
func TestSearchByRank_AfterInsertions(t *testing.T) {
	s := NewSkipList[int]()

	// Insert in non-sequential order
	s.InsertAtLevel(30, 2)
	s.InsertAtLevel(10, 1)
	s.InsertAtLevel(50, 1)
	s.InsertAtLevel(20, 0)
	s.InsertAtLevel(40, 1)

	// Expected order: 10, 20, 30, 40, 50
	expected := []int{10, 20, 30, 40, 50}

	for rank := 1; rank <= len(expected); rank++ {
		node, found := s.SearchByRank(rank)
		if !found {
			t.Fatalf("SearchByRank(%d) should find a node", rank)
		}
		if node == nil || node.val != expected[rank-1] {
			t.Fatalf("SearchByRank(%d) should return node with value %d, got %v", rank, expected[rank-1], node)
		}
	}
}

// TestSearchByRank_AfterDeletions tests search after deletions
func TestSearchByRank_AfterDeletions(t *testing.T) {
	s := NewSkipList[int]()
	s.InsertAtLevel(10, 1)
	s.InsertAtLevel(20, 2)
	s.InsertAtLevel(30, 1)
	s.InsertAtLevel(40, 0)
	s.InsertAtLevel(50, 1)

	// Delete 20
	s.Delete(20)

	// Expected order after deletion: 10, 30, 40, 50
	expected := []int{10, 30, 40, 50}

	for rank := 1; rank <= len(expected); rank++ {
		node, found := s.SearchByRank(rank)
		if !found {
			t.Fatalf("SearchByRank(%d) should find a node after deletion", rank)
		}
		if node == nil || node.val != expected[rank-1] {
			t.Fatalf("SearchByRank(%d) should return node with value %d, got %v", rank, expected[rank-1], node)
		}
	}

	// Rank 5 should fail (only 4 elements remain)
	node, found := s.SearchByRank(5)
	if found {
		t.Fatalf("SearchByRank(5) should not find a node")
	}
	if node != nil {
		t.Fatalf("SearchByRank(5) should return nil, got %v", node)
	}
}

// TestSearchByRank_MixedOperations tests search after mixed add/delete operations
func TestSearchByRank_MixedOperations(t *testing.T) {
	s := NewSkipList[int]()

	// Add nodes
	s.InsertAtLevel(10, 1)
	s.InsertAtLevel(20, 2)
	s.InsertAtLevel(30, 1)

	// Verify ranks
	node, _ := s.SearchByRank(1)
	if node.val != 10 {
		t.Fatalf("After initial insert, rank 1 should be 10")
	}

	// Delete 20
	s.Delete(20)

	// Verify ranks after deletion
	node, _ = s.SearchByRank(1)
	if node.val != 10 {
		t.Fatalf("After deleting 20, rank 1 should be 10")
	}
	node, _ = s.SearchByRank(2)
	if node.val != 30 {
		t.Fatalf("After deleting 20, rank 2 should be 30")
	}

	// Add 25
	s.InsertAtLevel(25, 1)

	// Verify ranks after addition
	node, _ = s.SearchByRank(1)
	if node.val != 10 {
		t.Fatalf("After adding 25, rank 1 should be 10")
	}
	node, _ = s.SearchByRank(2)
	if node.val != 25 {
		t.Fatalf("After adding 25, rank 2 should be 25")
	}
	node, _ = s.SearchByRank(3)
	if node.val != 30 {
		t.Fatalf("After adding 25, rank 3 should be 30")
	}

	// Delete 10
	s.Delete(10)

	// Verify ranks after deleting first element
	node, _ = s.SearchByRank(1)
	if node.val != 25 {
		t.Fatalf("After deleting 10, rank 1 should be 25")
	}
	node, _ = s.SearchByRank(2)
	if node.val != 30 {
		t.Fatalf("After deleting 10, rank 2 should be 30")
	}
}

// TestSearchByRank_FirstAndLast tests searching for first and last ranks
func TestSearchByRank_FirstAndLast(t *testing.T) {
	s := NewSkipList[int]()
	s.InsertAtLevel(10, 1)
	s.InsertAtLevel(20, 2)
	s.InsertAtLevel(30, 1)
	s.InsertAtLevel(40, 0)
	s.InsertAtLevel(50, 1)

	// Test first rank
	node, found := s.SearchByRank(1)
	if !found {
		t.Fatalf("SearchByRank(1) should find first node")
	}
	if node == nil || node.val != 10 {
		t.Fatalf("SearchByRank(1) should return first node (10), got %v", node)
	}

	// Test last rank
	node, found = s.SearchByRank(5)
	if !found {
		t.Fatalf("SearchByRank(5) should find last node")
	}
	if node == nil || node.val != 50 {
		t.Fatalf("SearchByRank(5) should return last node (50), got %v", node)
	}
}

// TestSearchByRank_LargeList tests search on a larger list
func TestSearchByRank_LargeList(t *testing.T) {
	s := NewSkipList[int]()

	// Insert 20 nodes
	for i := 1; i <= 20; i++ {
		s.InsertAtLevel(i*10, i%4) // Vary levels
	}

	// Test various ranks
	testCases := []struct {
		rank        int
		expectedVal int
	}{
		{1, 10},
		{5, 50},
		{10, 100},
		{15, 150},
		{20, 200},
	}

	for _, tc := range testCases {
		node, found := s.SearchByRank(tc.rank)
		if !found {
			t.Fatalf("SearchByRank(%d) should find a node", tc.rank)
		}
		if node == nil || node.val != tc.expectedVal {
			t.Fatalf("SearchByRank(%d) should return node with value %d, got %v", tc.rank, tc.expectedVal, node)
		}
	}
}

// TestSearchByRank_AllRanksSequential tests all ranks sequentially
func TestSearchByRank_AllRanksSequential(t *testing.T) {
	s := NewSkipList[int]()

	// Insert nodes
	values := []int{5, 10, 15, 20, 25, 30, 35, 40}
	for _, val := range values {
		s.InsertAtLevel(val, val%3) // Vary levels
	}

	// Test all ranks from 1 to length
	for rank := 1; rank <= len(values); rank++ {
		node, found := s.SearchByRank(rank)
		if !found {
			t.Fatalf("SearchByRank(%d) should find a node", rank)
		}
		if node == nil {
			t.Fatalf("SearchByRank(%d) should not return nil", rank)
		}
		expectedVal := values[rank-1]
		if node.val != expectedVal {
			t.Fatalf("SearchByRank(%d) should return node with value %d, got %d", rank, expectedVal, node.val)
		}
	}
}

// TestSearchByRank_AfterMultipleDeletions tests search after multiple deletions
func TestSearchByRank_AfterMultipleDeletions(t *testing.T) {
	s := NewSkipList[int]()

	// Insert nodes
	for i := 1; i <= 10; i++ {
		s.InsertAtLevel(i*10, i%3)
	}

	// Delete several nodes
	s.Delete(20)
	s.Delete(50)
	s.Delete(80)

	// Expected remaining: 10, 30, 40, 60, 70, 90, 100
	expected := []int{10, 30, 40, 60, 70, 90, 100}

	for rank := 1; rank <= len(expected); rank++ {
		node, found := s.SearchByRank(rank)
		if !found {
			t.Fatalf("SearchByRank(%d) should find a node after deletions", rank)
		}
		if node == nil || node.val != expected[rank-1] {
			t.Fatalf("SearchByRank(%d) should return node with value %d, got %v", rank, expected[rank-1], node)
		}
	}
}

// TestSearchByRank_ConsistencyWithOrder tests that ranks match sequential order
func TestSearchByRank_ConsistencyWithOrder(t *testing.T) {
	s := NewSkipList[int]()

	// Insert nodes in non-sequential order
	insertOrder := []int{50, 10, 30, 20, 40}
	for _, val := range insertOrder {
		s.InsertAtLevel(val, val%3)
	}

	// Get all nodes in order by traversing level 0
	var orderedValues []int
	for curr := s.head.forward[0]; curr != s.tail; curr = curr.forward[0] {
		orderedValues = append(orderedValues, curr.val)
	}

	// Verify SearchByRank returns nodes in the same order
	for rank := 1; rank <= len(orderedValues); rank++ {
		node, found := s.SearchByRank(rank)
		if !found {
			t.Fatalf("SearchByRank(%d) should find a node", rank)
		}
		if node == nil || node.val != orderedValues[rank-1] {
			t.Fatalf("SearchByRank(%d) should return node with value %d (from order), got %v", rank, orderedValues[rank-1], node)
		}
	}
}
