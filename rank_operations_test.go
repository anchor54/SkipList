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

// ------------------------------------------------------------
// GetRank Test cases
// ------------------------------------------------------------

// TestGetRank_Basic tests basic GetRank functionality
func TestGetRank_Basic(t *testing.T) {
	s := NewSkipList[int]()
	s.InsertAtLevel(10, 1)
	s.InsertAtLevel(20, 2)
	s.InsertAtLevel(30, 1)

	// Test rank of 10 (should be 1)
	rank, found := s.GetRank(10)
	if !found {
		t.Fatalf("GetRank(10) should find the value")
	}
	if rank != 1 {
		t.Fatalf("GetRank(10) should return rank 1, got %d", rank)
	}

	// Test rank of 20 (should be 2)
	rank, found = s.GetRank(20)
	if !found {
		t.Fatalf("GetRank(20) should find the value")
	}
	if rank != 2 {
		t.Fatalf("GetRank(20) should return rank 2, got %d", rank)
	}

	// Test rank of 30 (should be 3)
	rank, found = s.GetRank(30)
	if !found {
		t.Fatalf("GetRank(30) should find the value")
	}
	if rank != 3 {
		t.Fatalf("GetRank(30) should return rank 3, got %d", rank)
	}
}

// TestGetRank_AllValues tests GetRank for all values in a list
func TestGetRank_AllValues(t *testing.T) {
	s := NewSkipList[int]()
	values := []int{10, 20, 30, 40, 50}
	for i, val := range values {
		s.InsertAtLevel(val, i%3) // Vary levels
	}

	// Test all values
	for expectedRank, val := range values {
		rank, found := s.GetRank(val)
		if !found {
			t.Fatalf("GetRank(%d) should find the value", val)
		}
		expectedRank1Indexed := expectedRank + 1
		if rank != expectedRank1Indexed {
			t.Fatalf("GetRank(%d) should return rank %d, got %d", val, expectedRank1Indexed, rank)
		}
	}
}

// TestGetRank_NonExistentValues tests GetRank for values not in the list
func TestGetRank_NonExistentValues(t *testing.T) {
	s := NewSkipList[int]()
	s.InsertAtLevel(10, 1)
	s.InsertAtLevel(20, 1)
	s.InsertAtLevel(30, 1)

	// Test non-existent value before first
	rank, found := s.GetRank(5)
	if found {
		t.Fatalf("GetRank(5) should not find the value")
	}
	if rank != -1 {
		t.Fatalf("GetRank(5) should return rank -1, got %d", rank)
	}

	// Test non-existent value in middle
	rank, found = s.GetRank(15)
	if found {
		t.Fatalf("GetRank(15) should not find the value")
	}
	if rank != -1 {
		t.Fatalf("GetRank(15) should return rank -1, got %d", rank)
	}

	// Test non-existent value after last
	rank, found = s.GetRank(100)
	if found {
		t.Fatalf("GetRank(100) should not find the value")
	}
	if rank != -1 {
		t.Fatalf("GetRank(100) should return rank -1, got %d", rank)
	}
}

// TestGetRank_EmptyList tests GetRank on empty list
func TestGetRank_EmptyList(t *testing.T) {
	s := NewSkipList[int]()

	// Any value should fail on empty list
	testValues := []int{0, 10, -5, 100}
	for _, val := range testValues {
		rank, found := s.GetRank(val)
		if found {
			t.Fatalf("GetRank(%d) should not find the value in empty list", val)
		}
		if rank != -1 {
			t.Fatalf("GetRank(%d) should return rank -1 in empty list, got %d", val, rank)
		}
	}
}

// TestGetRank_SingleElement tests GetRank on single element list
func TestGetRank_SingleElement(t *testing.T) {
	s := NewSkipList[int]()
	s.InsertAtLevel(42, 1)

	// Existing value should work
	rank, found := s.GetRank(42)
	if !found {
		t.Fatalf("GetRank(42) should find the value")
	}
	if rank != 1 {
		t.Fatalf("GetRank(42) should return rank 1, got %d", rank)
	}

	// Non-existent value should fail
	rank, found = s.GetRank(10)
	if found {
		t.Fatalf("GetRank(10) should not find the value")
	}
	if rank != -1 {
		t.Fatalf("GetRank(10) should return rank -1, got %d", rank)
	}
}

// TestGetRank_AfterInsertions tests GetRank after multiple insertions
func TestGetRank_AfterInsertions(t *testing.T) {
	s := NewSkipList[int]()

	// Insert in non-sequential order
	s.InsertAtLevel(30, 2)
	s.InsertAtLevel(10, 1)
	s.InsertAtLevel(50, 1)
	s.InsertAtLevel(20, 0)
	s.InsertAtLevel(40, 1)

	// Expected order: 10, 20, 30, 40, 50
	testCases := []struct {
		value        int
		expectedRank int
	}{
		{10, 1},
		{20, 2},
		{30, 3},
		{40, 4},
		{50, 5},
	}

	for _, tc := range testCases {
		rank, found := s.GetRank(tc.value)
		if !found {
			t.Fatalf("GetRank(%d) should find the value", tc.value)
		}
		if rank != tc.expectedRank {
			t.Fatalf("GetRank(%d) should return rank %d, got %d", tc.value, tc.expectedRank, rank)
		}
	}
}

// TestGetRank_AfterDeletions tests GetRank after deletions
func TestGetRank_AfterDeletions(t *testing.T) {
	s := NewSkipList[int]()
	s.InsertAtLevel(10, 1)
	s.InsertAtLevel(20, 2)
	s.InsertAtLevel(30, 1)
	s.InsertAtLevel(40, 0)
	s.InsertAtLevel(50, 1)

	// Delete 20
	s.Delete(20)

	// Test remaining values
	testCases := []struct {
		value        int
		expectedRank int
		shouldFind   bool
	}{
		{10, 1, true},
		{20, -1, false}, // Deleted
		{30, 2, true},
		{40, 3, true},
		{50, 4, true},
	}

	for _, tc := range testCases {
		rank, found := s.GetRank(tc.value)
		if found != tc.shouldFind {
			t.Fatalf("GetRank(%d) found=%v, expected found=%v", tc.value, found, tc.shouldFind)
		}
		if found && rank != tc.expectedRank {
			t.Fatalf("GetRank(%d) should return rank %d, got %d", tc.value, tc.expectedRank, rank)
		}
		if !found && rank != -1 {
			t.Fatalf("GetRank(%d) should return rank -1 when not found, got %d", tc.value, rank)
		}
	}
}

// TestGetRank_FirstAndLast tests GetRank for first and last elements
func TestGetRank_FirstAndLast(t *testing.T) {
	s := NewSkipList[int]()
	s.InsertAtLevel(10, 1)
	s.InsertAtLevel(20, 2)
	s.InsertAtLevel(30, 1)
	s.InsertAtLevel(40, 0)
	s.InsertAtLevel(50, 1)

	// Test first element
	rank, found := s.GetRank(10)
	if !found {
		t.Fatalf("GetRank(10) should find first element")
	}
	if rank != 1 {
		t.Fatalf("GetRank(10) should return rank 1, got %d", rank)
	}

	// Test last element
	rank, found = s.GetRank(50)
	if !found {
		t.Fatalf("GetRank(50) should find last element")
	}
	if rank != 5 {
		t.Fatalf("GetRank(50) should return rank 5, got %d", rank)
	}
}

// TestGetRank_LargeList tests GetRank on a larger list
func TestGetRank_LargeList(t *testing.T) {
	s := NewSkipList[int]()

	// Insert 20 nodes
	for i := 1; i <= 20; i++ {
		s.InsertAtLevel(i*10, i%4) // Vary levels
	}

	// Test various values
	testCases := []struct {
		value        int
		expectedRank int
	}{
		{10, 1},
		{50, 5},
		{100, 10},
		{150, 15},
		{200, 20},
	}

	for _, tc := range testCases {
		rank, found := s.GetRank(tc.value)
		if !found {
			t.Fatalf("GetRank(%d) should find the value", tc.value)
		}
		if rank != tc.expectedRank {
			t.Fatalf("GetRank(%d) should return rank %d, got %d", tc.value, tc.expectedRank, rank)
		}
	}
}

// TestGetRank_ConsistencyWithSearchByRank tests that GetRank and SearchByRank are inverse operations
func TestGetRank_ConsistencyWithSearchByRank(t *testing.T) {
	s := NewSkipList[int]()

	// Insert nodes in non-sequential order
	insertOrder := []int{50, 10, 30, 20, 40}
	for _, val := range insertOrder {
		s.InsertAtLevel(val, val%3)
	}

	// For each value, GetRank should return the rank that SearchByRank can use to find it
	for _, val := range []int{10, 20, 30, 40, 50} {
		rank, found := s.GetRank(val)
		if !found {
			t.Fatalf("GetRank(%d) should find the value", val)
		}

		// SearchByRank should find the same value at that rank
		node, foundByRank := s.SearchByRank(rank)
		if !foundByRank {
			t.Fatalf("SearchByRank(%d) should find a node", rank)
		}
		if node == nil || node.val != val {
			t.Fatalf("GetRank(%d) returned rank %d, but SearchByRank(%d) returned value %d instead of %d",
				val, rank, rank, node.val, val)
		}
	}
}

// TestGetRank_AllValuesSequential tests GetRank for all values sequentially
func TestGetRank_AllValuesSequential(t *testing.T) {
	s := NewSkipList[int]()

	// Insert nodes
	values := []int{5, 10, 15, 20, 25, 30, 35, 40}
	for _, val := range values {
		s.InsertAtLevel(val, val%3) // Vary levels
	}

	// Test all values
	for expectedRank, val := range values {
		rank, found := s.GetRank(val)
		if !found {
			t.Fatalf("GetRank(%d) should find the value", val)
		}
		expectedRank1Indexed := expectedRank + 1
		if rank != expectedRank1Indexed {
			t.Fatalf("GetRank(%d) should return rank %d, got %d", val, expectedRank1Indexed, rank)
		}
	}
}

// TestGetRank_AfterMultipleDeletions tests GetRank after multiple deletions
func TestGetRank_AfterMultipleDeletions(t *testing.T) {
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
	testCases := []struct {
		value        int
		expectedRank int
		shouldFind   bool
	}{
		{10, 1, true},
		{20, -1, false}, // Deleted
		{30, 2, true},
		{40, 3, true},
		{50, -1, false}, // Deleted
		{60, 4, true},
		{70, 5, true},
		{80, -1, false}, // Deleted
		{90, 6, true},
		{100, 7, true},
	}

	for _, tc := range testCases {
		rank, found := s.GetRank(tc.value)
		if found != tc.shouldFind {
			t.Fatalf("GetRank(%d) found=%v, expected found=%v", tc.value, found, tc.shouldFind)
		}
		if found && rank != tc.expectedRank {
			t.Fatalf("GetRank(%d) should return rank %d, got %d", tc.value, tc.expectedRank, rank)
		}
		if !found && rank != -1 {
			t.Fatalf("GetRank(%d) should return rank -1 when not found, got %d", tc.value, rank)
		}
	}
}

// TestGetRank_MixedOperations tests GetRank after mixed add/delete operations
func TestGetRank_MixedOperations(t *testing.T) {
	s := NewSkipList[int]()

	// Add nodes
	s.InsertAtLevel(10, 1)
	s.InsertAtLevel(20, 2)
	s.InsertAtLevel(30, 1)

	// Verify ranks
	rank, found := s.GetRank(10)
	if !found || rank != 1 {
		t.Fatalf("After initial insert, GetRank(10) should return rank 1")
	}

	// Delete 20
	s.Delete(20)

	// Verify ranks after deletion
	rank, found = s.GetRank(10)
	if !found || rank != 1 {
		t.Fatalf("After deleting 20, GetRank(10) should return rank 1")
	}
	_, found = s.GetRank(20)
	if found {
		t.Fatalf("After deleting 20, GetRank(20) should not find the value")
	}
	rank, found = s.GetRank(30)
	if !found || rank != 2 {
		t.Fatalf("After deleting 20, GetRank(30) should return rank 2")
	}

	// Add 25
	s.InsertAtLevel(25, 1)

	// Verify ranks after addition
	rank, found = s.GetRank(10)
	if !found || rank != 1 {
		t.Fatalf("After adding 25, GetRank(10) should return rank 1")
	}
	rank, found = s.GetRank(25)
	if !found || rank != 2 {
		t.Fatalf("After adding 25, GetRank(25) should return rank 2")
	}
	rank, found = s.GetRank(30)
	if !found || rank != 3 {
		t.Fatalf("After adding 25, GetRank(30) should return rank 3")
	}

	// Delete 10
	s.Delete(10)

	// Verify ranks after deleting first element
	_, found = s.GetRank(10)
	if found {
		t.Fatalf("After deleting 10, GetRank(10) should not find the value")
	}
	rank, found = s.GetRank(25)
	if !found || rank != 1 {
		t.Fatalf("After deleting 10, GetRank(25) should return rank 1")
	}
	rank, found = s.GetRank(30)
	if !found || rank != 2 {
		t.Fatalf("After deleting 10, GetRank(30) should return rank 2")
	}
}

// TestGetRank_NegativeAndZeroValues tests GetRank with negative and zero values
func TestGetRank_NegativeAndZeroValues(t *testing.T) {
	s := NewSkipList[int]()
	s.InsertAtLevel(-10, 1)
	s.InsertAtLevel(-1, 1)
	s.InsertAtLevel(0, 0)
	s.InsertAtLevel(7, 1)
	s.InsertAtLevel(15, 2)

	// Expected order: -10, -1, 0, 7, 15
	testCases := []struct {
		value        int
		expectedRank int
		shouldFind   bool
	}{
		{-10, 1, true},
		{-1, 2, true},
		{0, 3, true},
		{7, 4, true},
		{15, 5, true},
		{-5, -1, false}, // Not in list
		{5, -1, false},  // Not in list
	}

	for _, tc := range testCases {
		rank, found := s.GetRank(tc.value)
		if found != tc.shouldFind {
			t.Fatalf("GetRank(%d) found=%v, expected found=%v", tc.value, found, tc.shouldFind)
		}
		if found && rank != tc.expectedRank {
			t.Fatalf("GetRank(%d) should return rank %d, got %d", tc.value, tc.expectedRank, rank)
		}
		if !found && rank != -1 {
			t.Fatalf("GetRank(%d) should return rank -1 when not found, got %d", tc.value, rank)
		}
	}
}

// TestGetRank_DuplicateInsertion tests that GetRank handles duplicate insertions correctly
func TestGetRank_DuplicateInsertion(t *testing.T) {
	s := NewSkipList[int]()
	s.InsertAtLevel(10, 1)
	s.InsertAtLevel(20, 1)
	s.InsertAtLevel(30, 1)

	// Try to insert duplicate
	s.InsertAtLevel(20, 2) // Should be ignored

	// GetRank should still work correctly
	rank, found := s.GetRank(20)
	if !found {
		t.Fatalf("GetRank(20) should find the value")
	}
	if rank != 2 {
		t.Fatalf("GetRank(20) should return rank 2, got %d", rank)
	}

	// Verify length is still 3 (duplicate not added)
	if s.Len() != 3 {
		t.Fatalf("Expected length 3 after duplicate insertion, got %d", s.Len())
	}
}
