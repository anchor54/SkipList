package skiplist

import "testing"

// ------------------------------------------------------------
// Comprehensive InsertAtLevel Test cases
// ------------------------------------------------------------

// TestInsertAtLevel_Level0 tests inserting at level 0 (lowest level)
func TestInsertAtLevel_Level0(t *testing.T) {
	s := NewSkipList[int]()
	s.InsertAtLevel(5, 0)
	s.InsertAtLevel(10, 0)
	s.InsertAtLevel(7, 0)

	// Verify order
	assertOrder(t, s, []int{5, 7, 10})

	// Verify maxLevel is still 0
	if s.maxLevel != 0 {
		t.Fatalf("maxLevel should be 0 after inserting at level 0, got %d", s.maxLevel)
	}

	// Verify nodes only exist at level 0
	node5 := findNode(s, 5)
	if node5 == nil || len(node5.forward) != 1 {
		t.Fatalf("node 5 should have exactly 1 forward pointer (level 0)")
	}
}

// TestInsertAtLevel_MaxLevelUpdate tests that maxLevel is updated when inserting at higher levels
func TestInsertAtLevel_MaxLevelUpdate(t *testing.T) {
	s := NewSkipList[int]()

	// Insert at level 0, maxLevel should be 0
	s.InsertAtLevel(1, 0)
	if s.maxLevel != 0 {
		t.Fatalf("maxLevel should be 0, got %d", s.maxLevel)
	}

	// Insert at level 2, maxLevel should become 2
	s.InsertAtLevel(2, 2)
	if s.maxLevel != 2 {
		t.Fatalf("maxLevel should be 2, got %d", s.maxLevel)
	}

	// Insert at level 1, maxLevel should remain 2
	s.InsertAtLevel(3, 1)
	if s.maxLevel != 2 {
		t.Fatalf("maxLevel should still be 2, got %d", s.maxLevel)
	}

	// Insert at level 5, maxLevel should become 5
	s.InsertAtLevel(4, 5)
	if s.maxLevel != 5 {
		t.Fatalf("maxLevel should be 5, got %d", s.maxLevel)
	}
}

// TestInsertAtLevel_MaxLevelCap tests inserting at MAX_LEVEL_CAP
func TestInsertAtLevel_MaxLevelCap(t *testing.T) {
	s := NewSkipList[int]()
	s.InsertAtLevel(42, MaxLevelCap)

	if s.maxLevel != MaxLevelCap {
		t.Fatalf("maxLevel should be %d, got %d", MaxLevelCap, s.maxLevel)
	}

	node := findNode(s, 42)
	if node == nil {
		t.Fatalf("node 42 not found")
	}
	if len(node.forward) != MaxLevelCap+1 {
		t.Fatalf("node at MAX_LEVEL_CAP should have %d forward pointers, got %d", MaxLevelCap+1, len(node.forward))
	}
}

// TestInsertAtLevel_DuplicateValue tests that duplicate values are ignored
func TestInsertAtLevel_DuplicateValue(t *testing.T) {
	s := NewSkipList[int]()
	s.InsertAtLevel(10, 2)
	s.InsertAtLevel(20, 1)
	s.InsertAtLevel(30, 3)

	// Get initial state
	initialLevels := collectLevelValues(s)
	initialMaxLevel := s.maxLevel

	// Try to insert duplicate
	s.InsertAtLevel(20, 5) // Try to insert at higher level

	// Verify nothing changed
	afterLevels := collectLevelValues(s)
	if s.maxLevel != initialMaxLevel {
		t.Fatalf("maxLevel should not change when inserting duplicate, was %d, got %d", initialMaxLevel, s.maxLevel)
	}

	// Verify levels are unchanged
	for level := 0; level <= s.maxLevel; level++ {
		if !slicesEqual(initialLevels[level], afterLevels[level]) {
			t.Fatalf("level %d changed after duplicate insert: was %v, got %v", level, initialLevels[level], afterLevels[level])
		}
	}

	// Verify order is still correct
	assertOrder(t, s, []int{10, 20, 30})
}

// TestInsertAtLevel_ForwardPointers tests that forward pointers are correctly set at all levels
func TestInsertAtLevel_ForwardPointers(t *testing.T) {
	s := NewSkipList[int]()
	s.InsertAtLevel(10, 2)
	s.InsertAtLevel(20, 1)
	s.InsertAtLevel(30, 3)

	// Verify forward pointers at each level
	// Level 0: should have all nodes
	level0Vals := collectLevelValues(s)[0]
	expectedLevel0 := []int{10, 20, 30}
	if !slicesEqual(level0Vals, expectedLevel0) {
		t.Fatalf("level 0: expected %v, got %v", expectedLevel0, level0Vals)
	}

	// Level 1: should have 10, 20, and 30 (all have level >= 1)
	level1Vals := collectLevelValues(s)[1]
	expectedLevel1 := []int{10, 20, 30}
	if !slicesEqual(level1Vals, expectedLevel1) {
		t.Fatalf("level 1: expected %v, got %v", expectedLevel1, level1Vals)
	}

	// Level 2: should have 10 and 30 (both have level >= 2)
	level2Vals := collectLevelValues(s)[2]
	expectedLevel2 := []int{10, 30}
	if !slicesEqual(level2Vals, expectedLevel2) {
		t.Fatalf("level 2: expected %v, got %v", expectedLevel2, level2Vals)
	}

	// Level 3: should have only 30 (only it has level >= 3)
	level3Vals := collectLevelValues(s)[3]
	expectedLevel3 := []int{30}
	if !slicesEqual(level3Vals, expectedLevel3) {
		t.Fatalf("level 3: expected %v, got %v", expectedLevel3, level3Vals)
	}
}

// TestInsertAtLevel_ForwardPointersCorrectness tests forward pointer correctness
func TestInsertAtLevel_ForwardPointersCorrectness(t *testing.T) {
	s := NewSkipList[int]()
	s.InsertAtLevel(10, 2)
	s.InsertAtLevel(20, 1)
	s.InsertAtLevel(30, 3)

	// Verify forward pointers are non-nil and point to correct nodes
	for level := 0; level <= s.maxLevel; level++ {
		for curr := s.head; curr != s.tail && curr.forward[level] != nil; {
			next := curr.forward[level]
			if next == nil {
				t.Fatalf("nil forward pointer at level %d from node %d", level, curr.val)
			}
			if next.val < curr.val {
				t.Fatalf("forward pointer at level %d violates ordering: %d -> %d", level, curr.val, next.val)
			}
			if curr == s.head && next == s.tail && level > s.maxLevel {
				// This is fine for levels above maxLevel
				break
			}
			curr = next
			if curr == s.tail {
				break
			}
		}
	}
}

// TestInsertAtLevel_InsertAtBeginning tests inserting at the beginning of the list
func TestInsertAtLevel_InsertAtBeginning(t *testing.T) {
	s := NewSkipList[int]()
	s.InsertAtLevel(100, 1)
	s.InsertAtLevel(50, 2)
	s.InsertAtLevel(10, 0)

	assertOrder(t, s, []int{10, 50, 100})

	// Verify node 10 only appears at level 0
	node10 := findNode(s, 10)
	if node10 == nil || len(node10.forward) != 1 {
		t.Fatalf("node 10 should have exactly 1 forward pointer")
	}

	// Verify node 50 appears at levels 0, 1, 2
	node50 := findNode(s, 50)
	if node50 == nil || len(node50.forward) != 3 {
		t.Fatalf("node 50 should have 3 forward pointers")
	}
}

// TestInsertAtLevel_InsertAtEnd tests inserting at the end of the list
func TestInsertAtLevel_InsertAtEnd(t *testing.T) {
	s := NewSkipList[int]()
	s.InsertAtLevel(10, 1)
	s.InsertAtLevel(50, 2)
	s.InsertAtLevel(100, 0)

	assertOrder(t, s, []int{10, 50, 100})

	// Verify node 100 only appears at level 0
	node100 := findNode(s, 100)
	if node100 == nil || len(node100.forward) != 1 {
		t.Fatalf("node 100 should have exactly 1 forward pointer")
	}

	// Verify forward pointer of 100 points to tail
	if node100.forward[0] != s.tail {
		t.Fatalf("node 100 should point to tail at level 0")
	}
}

// TestInsertAtLevel_InsertInMiddle tests inserting in the middle
func TestInsertAtLevel_InsertInMiddle(t *testing.T) {
	s := NewSkipList[int]()
	s.InsertAtLevel(10, 2)
	s.InsertAtLevel(50, 2)
	s.InsertAtLevel(30, 1) // Insert in middle

	assertOrder(t, s, []int{10, 30, 50})

	// Verify node 30 appears at levels 0 and 1
	node30 := findNode(s, 30)
	if node30 == nil || len(node30.forward) != 2 {
		t.Fatalf("node 30 should have 2 forward pointers")
	}

	// Verify forward pointers
	if node30.forward[0].val != 50 {
		t.Fatalf("node 30 at level 0 should point to 50, got %d", node30.forward[0].val)
	}
}

// TestInsertAtLevel_SkipCountsComprehensive tests skip counts comprehensively
func TestInsertAtLevel_SkipCountsComprehensive(t *testing.T) {
	s := NewSkipList[int]()
	// Insert: 10 (level 1), 20 (level 2), 30 (level 1), 40 (level 0)
	s.InsertAtLevel(10, 1)
	s.InsertAtLevel(20, 2)
	s.InsertAtLevel(30, 1)
	s.InsertAtLevel(40, 0)

	// Check spans for each node
	// Note: Skip counts depend on the structure - verifying they're set correctly
	node10 := findNode(s, 10)
	node20 := findNode(s, 20)
	node30 := findNode(s, 30)
	node40 := findNode(s, 40)

	if node10 == nil || len(node10.skips) != 2 {
		t.Fatalf("node 10 should have 2 skip counts")
	}
	if node20 == nil || len(node20.skips) != 3 {
		t.Fatalf("node 20 should have 3 skip counts")
	}
	if node30 == nil || len(node30.skips) != 2 {
		t.Fatalf("node 30 should have 2 skip counts")
	}
	if node40 == nil || len(node40.skips) != 1 {
		t.Fatalf("node 40 should have 1 skip count")
	}

	// Verify skip counts are positive
	for i, span := range node10.skips {
		if span <= 0 {
			t.Fatalf("node 10 skip count at level %d should be positive, got %d", i, span)
		}
	}
	for i, span := range node20.skips {
		if span <= 0 {
			t.Fatalf("node 20 skip count at level %d should be positive, got %d", i, span)
		}
	}
}

// TestInsertAtLevel_IncreasingLevels tests inserting nodes with increasing levels
func TestInsertAtLevel_IncreasingLevels(t *testing.T) {
	s := NewSkipList[int]()
	for i := 0; i <= 5; i++ {
		s.InsertAtLevel((i+1)*10, i)
	}

	if s.maxLevel != 5 {
		t.Fatalf("maxLevel should be 5, got %d", s.maxLevel)
	}

	assertOrder(t, s, []int{10, 20, 30, 40, 50, 60})

	// Verify each node has correct number of forward pointers
	for i := 0; i <= 5; i++ {
		val := (i + 1) * 10
		node := findNode(s, val)
		if node == nil {
			t.Fatalf("node %d not found", val)
		}
		expectedLevel := i + 1 // +1 because level i means i+1 forward pointers
		if len(node.forward) != expectedLevel {
			t.Fatalf("node %d should have %d forward pointers, got %d", val, expectedLevel, len(node.forward))
		}
	}
}

// TestInsertAtLevel_MultipleSameLevel tests inserting multiple nodes at the same level
func TestInsertAtLevel_MultipleSameLevel(t *testing.T) {
	s := NewSkipList[int]()
	s.InsertAtLevel(10, 1)
	s.InsertAtLevel(20, 1)
	s.InsertAtLevel(30, 1)
	s.InsertAtLevel(40, 1)

	assertOrder(t, s, []int{10, 20, 30, 40})

	// All nodes should appear at level 0 and 1
	level0Vals := collectLevelValues(s)[0]
	level1Vals := collectLevelValues(s)[1]
	expected := []int{10, 20, 30, 40}

	if !slicesEqual(level0Vals, expected) {
		t.Fatalf("level 0: expected %v, got %v", expected, level0Vals)
	}
	if !slicesEqual(level1Vals, expected) {
		t.Fatalf("level 1: expected %v, got %v", expected, level1Vals)
	}
}

// TestInsertAtLevel_HeadSkipCounts tests that head skip counts are updated correctly
func TestInsertAtLevel_HeadSkipCounts(t *testing.T) {
	s := NewSkipList[int]()
	s.InsertAtLevel(10, 0)

	// Head should have skip count of 1 at level 0 (pointing to 10)
	if s.head.skips[0] != 1 {
		t.Fatalf("head skip count at level 0 should be 1, got %d", s.head.skips[0])
	}

	s.InsertAtLevel(20, 2)

	// Head should have skip count of 1 at level 0 (pointing to 10)
	if s.head.skips[0] != 1 {
		t.Fatalf("head skip count at level 0 should be 1, got %d", s.head.skips[0])
	}

	// Head should have skip count of 2 at level 1 (pointing to 20, skipping 10)
	if s.head.skips[1] != 2 {
		t.Fatalf("head skip count at level 1 should be 2, got %d", s.head.skips[1])
	}

	// Head should have skip count of 2 at level 2 (pointing to 20, skipping 10)
	if s.head.skips[2] != 2 {
		t.Fatalf("head skip count at level 2 should be 2, got %d", s.head.skips[2])
	}

	// Head should have skip count of 3 at levels above maxLevel
	// (initial 1 + 1 for node 10 + 1 for node 20 = 3)
	for level := 3; level <= MaxLevelCap; level++ {
		if s.head.skips[level] != 3 {
			t.Fatalf("head skip count at level %d should be 3, got %d", level, s.head.skips[level])
		}
	}
}

// TestInsertAtLevel_LevelConsistency tests that nodes only appear at their designated levels
func TestInsertAtLevel_LevelConsistency(t *testing.T) {
	s := NewSkipList[int]()
	s.InsertAtLevel(10, 1) // Should appear at levels 0, 1
	s.InsertAtLevel(20, 3) // Should appear at levels 0, 1, 2, 3
	s.InsertAtLevel(30, 0) // Should appear only at level 0

	// Verify node 10 appears at levels 0 and 1, but not 2
	level0Vals := collectLevelValues(s)[0]
	level1Vals := collectLevelValues(s)[1]
	level2Vals := collectLevelValues(s)[2]

	if !contains(level0Vals, 10) {
		t.Fatalf("node 10 should appear at level 0")
	}
	if !contains(level1Vals, 10) {
		t.Fatalf("node 10 should appear at level 1")
	}
	if contains(level2Vals, 10) {
		t.Fatalf("node 10 should NOT appear at level 2")
	}

	// Verify node 30 appears only at level 0
	if !contains(level0Vals, 30) {
		t.Fatalf("node 30 should appear at level 0")
	}
	if contains(level1Vals, 30) {
		t.Fatalf("node 30 should NOT appear at level 1")
	}
}

// TestInsertAtLevel_ComplexScenario tests a complex insertion scenario
func TestInsertAtLevel_ComplexScenario(t *testing.T) {
	s := NewSkipList[int]()
	// Insert in non-sequential order with varying levels
	insertions := []struct {
		val   int
		level int
	}{
		{50, 3},
		{10, 1},
		{80, 2},
		{30, 0},
		{60, 4},
		{20, 1},
		{70, 0},
	}

	for _, ins := range insertions {
		s.InsertAtLevel(ins.val, ins.level)
	}

	// Verify order
	assertOrder(t, s, []int{10, 20, 30, 50, 60, 70, 80})

	// Verify maxLevel is 4
	if s.maxLevel != 4 {
		t.Fatalf("maxLevel should be 4, got %d", s.maxLevel)
	}

	// Verify all nodes can be found
	for _, ins := range insertions {
		node := findNode(s, ins.val)
		if node == nil {
			t.Fatalf("node %d not found after insertion", ins.val)
		}
		expectedLevel := ins.level + 1
		if len(node.forward) != expectedLevel {
			t.Fatalf("node %d should have %d forward pointers, got %d", ins.val, expectedLevel, len(node.forward))
		}
	}
}
