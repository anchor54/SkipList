package skiplist

import "testing"

// ------------------------------------------------------------
// Delete Span/Skip Count Test cases
// ------------------------------------------------------------

// TestDeleteSpans_SimpleDeletion tests that spans are correctly updated after deleting a single node
func TestDeleteSpans_SimpleDeletion(t *testing.T) {
	s := NewSkipList[int]()
	// Build a deterministic structure: 10 (level 1), 20 (level 2), 30 (level 1)
	s.InsertAtLevel(10, 2)
	s.InsertAtLevel(20, 2)
	s.InsertAtLevel(30, 1)

	// Verify initial spans
	checkSpans(t, s, 10, []int{1, 1, 1})
	checkSpans(t, s, 20, []int{1, 1, 2})
	checkSpans(t, s, 30, []int{1, 1})

	// Delete node 20
	s.Delete(20)

	// After deletion, node 10 should now point directly to 30
	// At level 0: 10 -> 30 (skip count should be 1)
	// At level 1: 10 -> 30 (skip count should be 1)
	// At level 2: 10 should point to tail (skip count should be 2, skipping 30)
	checkSpans(t, s, 10, []int{1, 1, 2})

	// Node 30 should remain unchanged
	checkSpans(t, s, 30, []int{1, 1})
}

// TestDeleteSpans_DeleteLevel0Node tests deleting a level 0 node
func TestDeleteSpans_DeleteLevel0Node(t *testing.T) {
	s := NewSkipList[int]()
	// Build: 10 (level 1), 20 (level 0), 30 (level 2)
	s.InsertAtLevel(10, 1)
	s.InsertAtLevel(20, 0)
	s.InsertAtLevel(30, 2)

	checkSpans(t, s, 10, []int{1, 2})

	// Delete node 20 (level 0)
	s.Delete(20)

	// Node 10 should now point directly to 30
	// At level 0: 10 -> 30 (skip count 1)
	// At level 1: 10 -> 30 (skip count 1)
	checkSpans(t, s, 10, []int{1, 1})

	// Node 30 should remain unchanged
	checkSpans(t, s, 30, []int{1, 1, 1})
}

// TestDeleteSpans_DeleteMiddleNode tests deleting a node in the middle
func TestDeleteSpans_DeleteMiddleNode(t *testing.T) {
	s := NewSkipList[int]()
	// Build: 10 (level 1), 20 (level 2), 30 (level 1), 40 (level 0)
	s.InsertAtLevel(10, 1)
	s.InsertAtLevel(20, 2)
	s.InsertAtLevel(30, 1)
	s.InsertAtLevel(40, 0)

	// Delete node 20 (middle node at level 2)
	s.Delete(20)

	// Node 10 should now point to 30
	// At level 0: 10 -> 30 (skip count 1)
	// At level 1: 10 -> 30 (skip count 1)
	checkSpans(t, s, 10, []int{1, 1})

	// Node 30 should remain unchanged
	// At level 0: 30 -> 40 (skip count 1)
	// At level 1: 30 -> Max (skip count 2)
	checkSpans(t, s, 30, []int{1, 2})
	checkSpans(t, s, 40, []int{1})
}

// TestDeleteSpans_DeleteFirstNode tests deleting the first node
func TestDeleteSpans_DeleteFirstNode(t *testing.T) {
	s := NewSkipList[int]()
	// Build: 10 (level 1), 20 (level 2), 30 (level 1)
	s.InsertAtLevel(10, 1)
	s.InsertAtLevel(20, 2)
	s.InsertAtLevel(30, 1)

	// Delete node 10 (first node)
	s.Delete(10)

	// Head should now point directly to 20
	// At level 0: head -> 20 (skip count 1)
	// At level 1: head -> 20 (skip count 1)
	// At level 2: head -> 20 (skip count 1)
	if s.head.skips[0] != 1 {
		t.Fatalf("head skip count at level 0 should be 1, got %d", s.head.skips[0])
	}
	if s.head.skips[1] != 1 {
		t.Fatalf("head skip count at level 1 should be 1, got %d", s.head.skips[1])
	}
	if s.head.skips[2] != 1 {
		t.Fatalf("head skip count at level 2 should be 1, got %d", s.head.skips[2])
	}

	// Node 20 should remain unchanged
	checkSpans(t, s, 20, []int{1, 1, 2})
}

// TestDeleteSpans_DeleteLastNode tests deleting the last node
func TestDeleteSpans_DeleteLastNode(t *testing.T) {
	s := NewSkipList[int]()
	// Build: 10 (level 1), 20 (level 2), 30 (level 1)
	s.InsertAtLevel(10, 1)
	s.InsertAtLevel(20, 2)
	s.InsertAtLevel(30, 1)

	// Delete node 30 (last node)
	s.Delete(30)

	// Node 20 should now point directly to tail
	// At level 0: 20 -> tail (skip count 1)
	// At level 1: 20 -> tail (skip count 1)
	// At level 2: 20 -> tail (skip count 1)
	checkSpans(t, s, 20, []int{1, 1, 1})

	// Node 10 should remain unchanged
	checkSpans(t, s, 10, []int{1, 1})
}

// TestDeleteSpans_DeleteMultiLevelNode tests deleting a node that appears at multiple levels
func TestDeleteSpans_DeleteMultiLevelNode(t *testing.T) {
	s := NewSkipList[int]()
	// Build: 10 (level 1), 20 (level 3), 30 (level 1), 40 (level 0)
	s.InsertAtLevel(10, 1)
	s.InsertAtLevel(20, 3)
	s.InsertAtLevel(30, 1)
	s.InsertAtLevel(40, 0)

	// Delete node 20 (appears at levels 0, 1, 2, 3)
	s.Delete(20)

	// Node 10 should now point to 30
	// At level 0: 10 -> 30 (skip count 1)
	// At level 1: 10 -> 30 (skip count 1)
	checkSpans(t, s, 10, []int{1, 1})

	// Node 30 should remain unchanged
	checkSpans(t, s, 30, []int{1, 2})

	// Node 30 should remain unchanged
	checkSpans(t, s, 40, []int{1})
}

// TestDeleteSpans_HeadSpansAfterDeletion tests that head spans are updated correctly
func TestDeleteSpans_HeadSpansAfterDeletion(t *testing.T) {
	s := NewSkipList[int]()
	// Build: 10 (level 0), 20 (level 1), 30 (level 2)
	s.InsertAtLevel(10, 0)
	s.InsertAtLevel(20, 1)
	s.InsertAtLevel(30, 2)

	// Verify initial head spans
	// At level 0: head -> 10 (skip count 1)
	// At level 1: head -> 20 (skip count 2)
	// At level 2: head -> 30 (skip count 3)
	if s.head.skips[0] != 1 {
		t.Fatalf("head skip count at level 0 should be 1, got %d", s.head.skips[0])
	}
	if s.head.skips[1] != 2 {
		t.Fatalf("head skip count at level 1 should be 2, got %d", s.head.skips[1])
	}
	if s.head.skips[2] != 3 {
		t.Fatalf("head skip count at level 2 should be 3, got %d", s.head.skips[2])
	}

	// Delete node 20
	s.Delete(20)

	// After deletion:
	// At level 0: head -> 10 (skip count 1)
	// At level 1: head -> 30 (skip count 1, since 10 is at level 0)
	// At level 2: head -> 30 (skip count 2, skipping 10)
	if s.head.skips[0] != 1 {
		t.Fatalf("head skip count at level 0 should be 1, got %d", s.head.skips[0])
	}
	if s.head.skips[1] != 2 {
		t.Fatalf("head skip count at level 1 should be 2, got %d", s.head.skips[1])
	}
	if s.head.skips[2] != 2 {
		t.Fatalf("head skip count at level 2 should be 2, got %d", s.head.skips[2])
	}
}

// TestDeleteSpans_MultipleDeletions tests spans after multiple deletions
func TestDeleteSpans_MultipleDeletions(t *testing.T) {
	s := NewSkipList[int]()
	// Build: 10 (level 1), 20 (level 2), 30 (level 1), 40 (level 0), 50 (level 2)
	s.InsertAtLevel(10, 1)
	s.InsertAtLevel(20, 2)
	s.InsertAtLevel(30, 1)
	s.InsertAtLevel(40, 0)
	s.InsertAtLevel(50, 2)

	// Delete node 20
	s.Delete(20)

	// Node 10 should now point to 30
	checkSpans(t, s, 10, []int{1, 1})

	// Delete node 30
	s.Delete(30)

	// Node 10 should now point to 40
	checkSpans(t, s, 10, []int{1, 2})

	// Node 40 should remain unchanged
	checkSpans(t, s, 40, []int{1})

	// Node 50 should remain unchanged
	checkSpans(t, s, 50, []int{1, 1, 1})
}

// TestDeleteSpans_AllLevelsUpdated tests that spans are updated at all relevant levels
func TestDeleteSpans_AllLevelsUpdated(t *testing.T) {
	s := NewSkipList[int]()
	// Build: 10 (level 2), 20 (level 3), 30 (level 1), 40 (level 0)
	s.InsertAtLevel(10, 2)
	s.InsertAtLevel(20, 3)
	s.InsertAtLevel(30, 1)
	s.InsertAtLevel(40, 0)

	// Delete node 20 (appears at levels 0, 1, 2, 3)
	s.Delete(20)

	// Node 10 should have spans updated at all levels where it pointed to 20
	// At level 0: 10 -> 30 (skip count 1)
	// At level 1: 10 -> 30 (skip count 1)
	// At level 2: 10 -> INF (skip count 3, since 30 is until level 1 and 40 is until level 0)
	checkSpans(t, s, 10, []int{1, 1, 3})

	// Head spans at levels above maxLevel should be decremented
	// maxLevel is now 2 (from node 10), so levels 3+ should have one less skip (including the right sentinel)
	for level := 3; level <= MaxLevelCap; level++ {
		if s.head.skips[level] != 4 {
			t.Fatalf("head skip count at level %d should be 4, got %d", level, s.head.skips[level])
		}
	}
}

// TestDeleteSpans_ComplexScenario tests a complex deletion scenario
func TestDeleteSpans_ComplexScenario(t *testing.T) {
	s := NewSkipList[int]()
	// Build a more complex structure
	s.InsertAtLevel(10, 1)
	s.InsertAtLevel(20, 3)
	s.InsertAtLevel(30, 1)
	s.InsertAtLevel(40, 2)
	s.InsertAtLevel(50, 0)
	s.InsertAtLevel(60, 4)

	// Delete node 40 (level 2, in the middle)
	s.Delete(40)

	// Node 30 should now point to 50
	// At level 0: 30 -> 50 (skip count 1)
	// At level 1: 30 -> 60 (skip count 2)
	checkSpans(t, s, 30, []int{1, 2})

	// Node 20 should now point to 50 at level 2
	// At level 0: 20 -> 30 (skip count 1)
	// At level 1: 20 -> 30 (skip count 1)
	// At level 2: 20 -> 60 (skip count 3, skipping 30)
	// At level 3: 20 -> 60 (skip count 3, skipping 30 and 50)
	checkSpans(t, s, 20, []int{1, 1, 3, 3})
}

// TestDeleteSpans_VerifyAllNodes tests that all remaining nodes have correct spans
func TestDeleteSpans_VerifyAllNodes(t *testing.T) {
	s := NewSkipList[int]()
	// Build: 10 (level 1), 20 (level 2), 30 (level 1)
	s.InsertAtLevel(10, 1)
	s.InsertAtLevel(20, 2)
	s.InsertAtLevel(30, 1)

	// Get initial spans for all nodes
	node10 := findNode(s, 10)
	node20 := findNode(s, 20)
	node30 := findNode(s, 30)

	initial10Spans := make([]int, len(node10.skips))
	initial20Spans := make([]int, len(node20.skips))
	initial30Spans := make([]int, len(node30.skips))
	copy(initial10Spans, node10.skips)
	copy(initial20Spans, node20.skips)
	copy(initial30Spans, node30.skips)

	// Delete node 20
	s.Delete(20)

	// Verify node 10 spans are updated
	node10After := findNode(s, 10)
	if node10After == nil {
		t.Fatalf("node 10 not found after deletion")
	}

	// Node 10 should have updated spans
	// The spans should reflect that 20 is no longer in the path
	for level := 0; level < len(node10After.skips); level++ {
		if level < len(node10After.forward) {
			// Verify the span is positive and reasonable
			if node10After.skips[level] <= 0 {
				t.Fatalf("node 10 span at level %d should be positive, got %d", level, node10After.skips[level])
			}
		}
	}

	// Verify node 30 spans remain unchanged (it wasn't affected)
	node30After := findNode(s, 30)
	if node30After == nil {
		t.Fatalf("node 30 not found after deletion")
	}
	if len(node30After.skips) != len(initial30Spans) {
		t.Fatalf("node 30 should have same number of spans, got %d, expected %d", len(node30After.skips), len(initial30Spans))
	}
	for i, span := range node30After.skips {
		if span != initial30Spans[i] {
			t.Fatalf("node 30 span at level %d should be unchanged (%d), got %d", i, initial30Spans[i], span)
		}
	}
}

// TestDeleteSpans_DeleteFromEmptyLevel tests deleting when a node doesn't exist at a level
func TestDeleteSpans_DeleteFromEmptyLevel(t *testing.T) {
	s := NewSkipList[int]()
	// Build: 10 (level 1), 20 (level 0), 30 (level 2)
	s.InsertAtLevel(10, 1)
	s.InsertAtLevel(20, 0)
	s.InsertAtLevel(30, 2)

	// Delete node 20 (only at level 0)
	s.Delete(20)

	// Node 10 should now point directly to 30
	// At level 0: 10 -> 30 (skip count 1)
	// At level 1: 10 -> 30 (skip count 1)
	checkSpans(t, s, 10, []int{1, 1})

	// Node 30 should remain unchanged
	checkSpans(t, s, 30, []int{1, 1, 1})
}

// TestDeleteSpans_BasicIntegrity tests that spans remain positive and nodes are still findable after deletion
// This test should pass even with the current buggy Delete() implementation
func TestDeleteSpans_BasicIntegrity(t *testing.T) {
	s := NewSkipList[int]()
	// Build: 10 (level 1), 20 (level 2), 30 (level 1)
	s.InsertAtLevel(10, 1)
	s.InsertAtLevel(20, 2)
	s.InsertAtLevel(30, 1)

	// Verify nodes exist before deletion
	node10 := findNode(s, 10)
	node20 := findNode(s, 20)
	node30 := findNode(s, 30)

	if node10 == nil || node20 == nil || node30 == nil {
		t.Fatalf("nodes should exist before deletion")
	}

	// Delete node 20
	s.Delete(20)

	// Verify node 20 is gone
	if findNode(s, 20) != nil {
		t.Fatalf("node 20 should be deleted")
	}

	// Verify remaining nodes still exist
	node10After := findNode(s, 10)
	node30After := findNode(s, 30)

	if node10After == nil {
		t.Fatalf("node 10 should still exist after deletion")
	}
	if node30After == nil {
		t.Fatalf("node 30 should still exist after deletion")
	}

	// Verify spans arrays still exist and have positive values
	// (Even if they're not correctly updated, they should at least be valid)
	for i, span := range node10After.skips {
		if span <= 0 {
			t.Fatalf("node 10 span at level %d should be positive, got %d", i, span)
		}
	}
	for i, span := range node30After.skips {
		if span <= 0 {
			t.Fatalf("node 30 span at level %d should be positive, got %d", i, span)
		}
	}

	// Verify forward pointers are still valid
	for i := 0; i < len(node10After.forward); i++ {
		if node10After.forward[i] == nil {
			t.Fatalf("node 10 forward pointer at level %d should not be nil", i)
		}
	}
	for i := 0; i < len(node30After.forward); i++ {
		if node30After.forward[i] != nil {
			t.Fatalf("node 30 forward pointer at level %d should be nil", i)
		}
	}
}

// TestDeleteSpans_AddDeleteSequence tests spans after a sequence of add and delete operations
func TestDeleteSpans_AddDeleteSequence(t *testing.T) {
	s := NewSkipList[int]()

	// Step 1: Add initial nodes
	s.InsertAtLevel(10, 1)
	s.InsertAtLevel(20, 2)
	s.InsertAtLevel(30, 1)
	s.InsertAtLevel(40, 0)

	// Verify initial spans
	checkSpans(t, s, 10, []int{1, 1})
	checkSpans(t, s, 20, []int{1, 1, 3})
	checkSpans(t, s, 30, []int{1, 2})
	checkSpans(t, s, 40, []int{1})

	// Step 2: Delete a node
	s.Delete(20)

	// Verify spans after first deletion
	checkSpans(t, s, 10, []int{1, 1})
	checkSpans(t, s, 30, []int{1, 2})
	checkSpans(t, s, 40, []int{1})

	// Step 3: Add more nodes
	s.InsertAtLevel(25, 1)
	s.InsertAtLevel(35, 2)

	// Verify spans after additions
	checkSpans(t, s, 10, []int{1, 1})
	checkSpans(t, s, 25, []int{1, 1})
	checkSpans(t, s, 30, []int{1, 1})
	checkSpans(t, s, 35, []int{1, 2, 2})
	checkSpans(t, s, 40, []int{1})

	// Step 4: Delete another node
	s.Delete(30)

	// Verify spans after second deletion
	checkSpans(t, s, 10, []int{1, 1})
	checkSpans(t, s, 25, []int{1, 1})
	checkSpans(t, s, 35, []int{1, 2, 2})
	checkSpans(t, s, 40, []int{1})

	// Step 5: Add one more node
	s.InsertAtLevel(15, 0)

	// Verify final spans
	checkSpans(t, s, 10, []int{1, 2})
	checkSpans(t, s, 15, []int{1})
	checkSpans(t, s, 25, []int{1, 1})
	checkSpans(t, s, 35, []int{1, 2, 2})
	checkSpans(t, s, 40, []int{1})
}

// TestDeleteSpans_AddDeleteAddPattern tests a pattern of add-delete-add operations
func TestDeleteSpans_AddDeleteAddPattern(t *testing.T) {
	s := NewSkipList[int]()

	// Add nodes: 10, 20, 30
	s.InsertAtLevel(10, 1)
	s.InsertAtLevel(20, 2)
	s.InsertAtLevel(30, 1)

	// Delete 20
	s.Delete(20)

	// Add 20 back at a different level
	s.InsertAtLevel(20, 1)

	// Verify spans - 20 should now be at level 1 instead of level 2
	checkSpans(t, s, 10, []int{1, 1})
	checkSpans(t, s, 20, []int{1, 1})
	checkSpans(t, s, 30, []int{1, 1})

	// Delete 10
	s.Delete(10)

	// Verify spans after deleting 10
	checkSpans(t, s, 20, []int{1, 1})
	checkSpans(t, s, 30, []int{1, 1})

	// Add 10 back
	s.InsertAtLevel(10, 0)

	// Verify final spans
	checkSpans(t, s, 10, []int{1})
	checkSpans(t, s, 20, []int{1, 1})
	checkSpans(t, s, 30, []int{1, 1})
}

// TestDeleteSpans_MultipleAddDeleteCycles tests multiple cycles of add and delete
func TestDeleteSpans_MultipleAddDeleteCycles(t *testing.T) {
	s := NewSkipList[int]()

	// Cycle 1: Add nodes
	s.InsertAtLevel(10, 1)
	s.InsertAtLevel(20, 2)
	s.InsertAtLevel(30, 1)
	checkSpans(t, s, 10, []int{1, 1})
	checkSpans(t, s, 20, []int{1, 1, 2})
	checkSpans(t, s, 30, []int{1, 1})

	// Cycle 1: Delete middle node
	s.Delete(20)
	checkSpans(t, s, 10, []int{1, 1})
	checkSpans(t, s, 30, []int{1, 1})

	// Cycle 2: Add more nodes
	s.InsertAtLevel(15, 0)
	s.InsertAtLevel(25, 2)
	s.InsertAtLevel(35, 1)
	checkSpans(t, s, 10, []int{1, 2})
	checkSpans(t, s, 15, []int{1})
	checkSpans(t, s, 25, []int{1, 1, 3})
	checkSpans(t, s, 30, []int{1, 1})
	checkSpans(t, s, 35, []int{1, 1})

	// Cycle 2: Delete first and last
	s.Delete(10)
	s.Delete(35)
	checkSpans(t, s, 15, []int{1})
	checkSpans(t, s, 25, []int{1, 1, 2})
	checkSpans(t, s, 30, []int{1, 1})

	// Cycle 3: Add nodes at various levels
	s.InsertAtLevel(12, 1)
	s.InsertAtLevel(18, 0)
	s.InsertAtLevel(22, 2)
	checkSpans(t, s, 12, []int{1, 3})
	checkSpans(t, s, 15, []int{1})
	checkSpans(t, s, 18, []int{1})
	checkSpans(t, s, 22, []int{1, 1, 1})
	checkSpans(t, s, 25, []int{1, 1, 2})
	checkSpans(t, s, 30, []int{1, 1})

	// Cycle 3: Delete multiple nodes
	s.Delete(15)
	s.Delete(22)
	checkSpans(t, s, 12, []int{1, 2})
	checkSpans(t, s, 18, []int{1})
	checkSpans(t, s, 25, []int{1, 1, 2})
	checkSpans(t, s, 30, []int{1, 1})
}

// TestDeleteSpans_DeleteThenAddSameValue tests deleting and then adding the same value
func TestDeleteSpans_DeleteThenAddSameValue(t *testing.T) {
	s := NewSkipList[int]()

	// Add nodes
	s.InsertAtLevel(10, 2)
	s.InsertAtLevel(20, 1)
	s.InsertAtLevel(30, 2)

	// l2: -INF --> 10 ---------> 30 --> INF
	// l1: -INF --> 10 --> 20 --> 30 --> INF
	// l0: -INF --> 10 --> 20 --> 30 --> INF

	// Verify initial spans
	checkSpans(t, s, 10, []int{1, 1, 2})
	checkSpans(t, s, 20, []int{1, 1})
	checkSpans(t, s, 30, []int{1, 1, 1})

	// Delete 20
	s.Delete(20)

	// l2: -INF --> 10 --> 30 --> INF
	// l1: -INF --> 10 --> 30 --> INF
	// l0: -INF --> 10 --> 30 --> INF

	// Verify spans after deletion
	checkSpans(t, s, 10, []int{1, 1, 1})
	checkSpans(t, s, 30, []int{1, 1, 1})

	// Add 20 back at a different level
	s.InsertAtLevel(20, 0)

	// l2: -INF --> 10 ---------> 30 --> INF
	// l1: -INF --> 10 ---------> 30 --> INF
	// l0: -INF --> 10 --> 20 --> 30 --> INF

	// Verify spans - 20 should now be at level 0
	checkSpans(t, s, 10, []int{1, 2, 2})
	checkSpans(t, s, 20, []int{1})
	checkSpans(t, s, 30, []int{1, 1, 1})

	// Delete 20 again
	s.Delete(20)

	// l2: -INF --> 10 --> 30 --> INF
	// l1: -INF --> 10 --> 30 --> INF
	// l0: -INF --> 10 --> 30 --> INF

	// Verify spans
	checkSpans(t, s, 10, []int{1, 1, 1})
	checkSpans(t, s, 30, []int{1, 1, 1})

	// Add 40 back at yet another level
	s.InsertAtLevel(40, 0)

	// l2: -INF --> 10 --> 30 ---------> INF
	// l1: -INF --> 10 --> 30 ---------> INF
	// l0: -INF --> 10 --> 30 --> 40 --> INF

	// Verify spans - 40 should now be at level 0
	checkSpans(t, s, 10, []int{1, 1, 1})
	checkSpans(t, s, 30, []int{1, 2, 2})
	checkSpans(t, s, 40, []int{1})

	// Add 20 back at yet another level
	s.InsertAtLevel(20, 3)

	// l3: -INF ---------> 20 ----------------> INF
	// l2: -INF --> 10 --> 20 --> 30 ---------> INF
	// l1: -INF --> 10 --> 20 --> 30 ---------> INF
	// l0: -INF --> 10 --> 20 --> 30 --> 40 --> INF

	// Verify spans - 40 should now be at level 0
	checkSpans(t, s, 10, []int{1, 1, 1})
	checkSpans(t, s, 20, []int{1, 1, 1, 3})
	checkSpans(t, s, 30, []int{1, 2, 2})
	checkSpans(t, s, 40, []int{1})

	if s.head.skips[3] != 2 {
		t.Fatalf("head skip count at level 3 should be 2, got %d", s.head.skips[3])
	}
}

// TestDeleteSpans_ComplexAddDeleteMix tests a complex mix of add and delete operations
func TestDeleteSpans_ComplexAddDeleteMix(t *testing.T) {
	s := NewSkipList[int]()

	// Build initial structure: 10, 20, 30, 40, 50
	s.InsertAtLevel(10, 1)
	s.InsertAtLevel(20, 2)
	s.InsertAtLevel(30, 1)
	s.InsertAtLevel(40, 0)
	s.InsertAtLevel(50, 3)

	// l3: -INF ------------------------------> 50 --> INF
	// l2: -INF ---------> 20 ----------------> 50 --> INF
	// l1: -INF --> 10 --> 20 --> 30 ---------> 50 --> INF
	// l0: -INF --> 10 --> 20 --> 30 --> 40 --> 50 --> INF

	// Verify initial spans
	checkSpans(t, s, 10, []int{1, 1})
	checkSpans(t, s, 20, []int{1, 1, 3})
	checkSpans(t, s, 30, []int{1, 2})
	checkSpans(t, s, 40, []int{1})
	checkSpans(t, s, 50, []int{1, 1, 1, 1})

	// Delete 20 and 40
	s.Delete(20)
	s.Delete(40)

	// l3: -INF ------------------------------> 50 --> INF
	// l2: -INF ------------------------------> 50 --> INF
	// l1: -INF --> 10 ---------> 30 ---------> 50 --> INF
	// l0: -INF --> 10 ---------> 30 ---------> 50 --> INF

	// Verify spans after deletions
	checkSpans(t, s, 10, []int{1, 1})
	checkSpans(t, s, 30, []int{1, 1})
	checkSpans(t, s, 50, []int{1, 1, 1, 1})

	// Add new nodes: 15, 25, 35, 45
	s.InsertAtLevel(15, 0)
	s.InsertAtLevel(25, 1)
	s.InsertAtLevel(35, 2)
	s.InsertAtLevel(45, 0)

	// l3: -INF --------------------------------------------> 50 --> INF
	// l2: -INF ------------------------------> 35 ---------> 50 --> INF
	// l1: -INF --> 10 ---------> 25 --> 30 --> 35 ---------> 50 --> INF
	// l0: -INF --> 10 --> 15 --> 25 --> 30 --> 35 --> 45 --> 50 --> INF

	// Verify spans after additions
	checkSpans(t, s, 10, []int{1, 2})
	checkSpans(t, s, 15, []int{1})
	checkSpans(t, s, 25, []int{1, 1})
	checkSpans(t, s, 30, []int{1, 1})
	checkSpans(t, s, 35, []int{1, 2, 2})
	checkSpans(t, s, 45, []int{1})
	checkSpans(t, s, 50, []int{1, 1, 1, 1})

	// Delete 10, 30, and 50
	s.Delete(10)
	s.Delete(30)
	s.Delete(50)

	// l3: -INF ----------------------------------------------> INF
	// l2: -INF ------------------------------> 35 -----------> INF
	// l1: -INF ----------------> 25 ---------> 35 -----------> INF
	// l0: -INF ---------> 15 --> 25 ---------> 35 --> 45 ----> INF

	// Verify spans after more deletions
	checkSpans(t, s, 15, []int{1})
	checkSpans(t, s, 25, []int{1, 1})
	checkSpans(t, s, 35, []int{1, 2, 2})
	checkSpans(t, s, 45, []int{1})

	// Add final nodes: 5, 55
	s.InsertAtLevel(5, 1)
	s.InsertAtLevel(55, 2)

	// l3: -INF -------------------------------------------> INF
	// l2: -INF ----------------------> 35 ---------> 55 --> INF
	// l1: -INF --> 5 ---------> 25 --> 35 ---------> 55 --> INF
	// l0: -INF --> 5 --> 15 --> 25 --> 35 --> 45 --> 55 --> INF

	// Verify final spans
	checkSpans(t, s, 5, []int{1, 2})
	checkSpans(t, s, 15, []int{1})
	checkSpans(t, s, 25, []int{1, 1})
	checkSpans(t, s, 35, []int{1, 2, 2})
	checkSpans(t, s, 45, []int{1})
	checkSpans(t, s, 55, []int{1, 1, 1})
}

// TestDeleteSpans_HeadSpansThroughOperations tests head spans through multiple operations
func TestDeleteSpans_HeadSpansThroughOperations(t *testing.T) {
	s := NewSkipList[int]()

	// Add nodes at different levels
	s.InsertAtLevel(10, 0)
	s.InsertAtLevel(20, 1)
	s.InsertAtLevel(30, 2)

	// l2: -INF ----------------> 30 --> INF
	// l1: -INF ---------> 20 --> 30 --> INF
	// l0: -INF --> 10 --> 20 --> 30 --> INF

	// Verify initial head spans
	if s.head.skips[0] != 1 {
		t.Fatalf("head skip count at level 0 should be 1, got %d", s.head.skips[0])
	}
	if s.head.skips[1] != 2 {
		t.Fatalf("head skip count at level 1 should be 2, got %d", s.head.skips[1])
	}
	if s.head.skips[2] != 3 {
		t.Fatalf("head skip count at level 2 should be 3, got %d", s.head.skips[2])
	}

	// Delete 20
	s.Delete(20)

	// l2: -INF ---------> 30 --> INF
	// l1: -INF ---------> 30 --> INF
	// l0: -INF --> 10 --> 30 --> INF

	// Verify head spans after deletion
	if s.head.skips[0] != 1 {
		t.Fatalf("head skip count at level 0 should be 1, got %d", s.head.skips[0])
	}
	if s.head.skips[1] != 2 {
		t.Fatalf("head skip count at level 1 should be 2, got %d", s.head.skips[1])
	}
	if s.head.skips[2] != 2 {
		t.Fatalf("head skip count at level 2 should be 2, got %d", s.head.skips[2])
	}

	// Add 25 at level 1
	s.InsertAtLevel(25, 1)

	// l2: -INF ----------------> 30 --> INF
	// l1: -INF ---------> 25 --> 30 --> INF
	// l0: -INF --> 10 --> 25 --> 30 --> INF

	// Verify head spans after addition
	if s.head.skips[0] != 1 {
		t.Fatalf("head skip count at level 0 should be 1, got %d", s.head.skips[0])
	}
	if s.head.skips[1] != 2 {
		t.Fatalf("head skip count at level 1 should be 2, got %d", s.head.skips[1])
	}
	if s.head.skips[2] != 3 {
		t.Fatalf("head skip count at level 2 should be 3, got %d", s.head.skips[2])
	}

	// Delete 10 (first node)
	s.Delete(10)

	// l2: -INF ---------> 30 --> INF
	// l1: -INF --> 25 --> 30 --> INF
	// l0: -INF --> 25 --> 30 --> INF

	// Verify head spans after deleting first node
	if s.head.skips[0] != 1 {
		t.Fatalf("head skip count at level 0 should be 1, got %d", s.head.skips[0])
	}
	if s.head.skips[1] != 1 {
		t.Fatalf("head skip count at level 1 should be 1, got %d", s.head.skips[1])
	}
	if s.head.skips[2] != 2 {
		t.Fatalf("head skip count at level 2 should be 2, got %d", s.head.skips[2])
	}
}

// TestDeleteSpans_AllNodesSpansAfterOperations tests all nodes have correct spans after mixed operations
func TestDeleteSpans_AllNodesSpansAfterOperations(t *testing.T) {
	s := NewSkipList[int]()

	// Initial setup
	s.InsertAtLevel(10, 1)
	s.InsertAtLevel(20, 2)
	s.InsertAtLevel(30, 1)
	s.InsertAtLevel(40, 0)

	// l2: -INF ---------> 20 ----------------> INF
	// l1: -INF --> 10 --> 20 --> 30 ---------> INF
	// l0: -INF --> 10 --> 20 --> 30 --> 40 --> INF

	// Delete 20
	s.Delete(20)

	// Add 25
	s.InsertAtLevel(25, 1)

	// Delete 10
	s.Delete(10)

	// Add 15
	s.InsertAtLevel(15, 0)

	// Delete 30
	s.Delete(30)

	// Add 35
	s.InsertAtLevel(35, 2)

	// l2: -INF ----------------> 35 ----------> INF
	// l1: -INF ---------> 25 --> 35 ---------> INF
	// l0: -INF --> 15 --> 25 --> 35 --> 40 --> INF

	// Final verification of all remaining nodes
	checkSpans(t, s, 15, []int{1})
	checkSpans(t, s, 25, []int{1, 1})
	checkSpans(t, s, 35, []int{1, 2, 2})
	checkSpans(t, s, 40, []int{1})

	// Verify order is correct
	assertOrder(t, s, []int{15, 25, 35, 40})
}
