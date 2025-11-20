package main

import "testing"

// ------------------------------------------------------------
// Helper functions
// ------------------------------------------------------------

func assertOrder(t *testing.T, sl *SkipList[int], expected []int) {
	t.Helper()	
	
	i := 0
	for curr := sl.head.forward[0]; curr != nil && curr.forward[0] != nil; curr = curr.forward[0] {
		if i >= len(expected) {
			t.Fatalf("skiplist has more elements than expected; extra element %d", curr.val)
		}

		if curr.val != expected[i] {
			t.Fatalf("value mismatch at index %d: got %d, want %d", i, curr.val, expected[i])
		}
		i++
	}

	if i != len(expected)-1 {
		t.Fatalf("skiplist has fewer elements than expected; checked %d, expected %d", i+1, len(expected))
	}
}

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}

// Collect all internal values for each level (excluding head/tail).
func collectLevelValues(sl *SkipList[int]) [][]int {
	levels := make([][]int, sl.maxLevel+1)
	for level := 0; level <= sl.maxLevel; level++ {
		for curr := sl.head.forward[level]; curr != nil && curr != sl.tail; curr = curr.forward[level] {
			levels[level] = append(levels[level], curr.val)
		}
	}
	return levels
}

func assertStrictlyIncreasing(t *testing.T, vals []int, level int) {
	t.Helper()
	for i := 1; i < len(vals); i++ {
		if vals[i] <= vals[i-1] {
			t.Fatalf("values not strictly increasing at level %d index %d: %d followed by %d",
				level, i-1, vals[i-1], vals[i])
		}
	}
}

func slicesEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// removeOne removes a single occurrence of target from vals (if present).
func removeOne(vals []int, target int) ([]int, bool) {
	out := make([]int, 0, len(vals))
	removed := false
	for _, v := range vals {
		if v == target && !removed {
			removed = true
			continue
		}
		out = append(out, v)
	}
	return out, removed
}

// Find a node that appears at more than one level (height > 0).
func findMultiLevelNode(sl *SkipList[int]) *Node[int] {
	for curr := sl.head.forward[0]; curr != nil && curr != sl.tail; curr = curr.forward[0] {
		if len(curr.forward) > 1 {
			return curr
		}
	}
	return nil
}

// findNode returns the node with value v (or nil).
func findNode(s *SkipList[int], v int) *Node[int] {
	for n := s.head; n != nil; {
		// traverse level 0 list
		if len(n.forward) == 0 {
			return nil
		}
		n = n.forward[0]
		if n == nil {
			return nil
		}
		if n.val == v {
			return n
		}
		if n.val > v {
			return nil
		}
	}
	return nil
}

// checkSpans asserts that the found node for value `v` has expected spans across all levels.
// expectedSpans must be of length == levelCount for the node (or test will fail).
func checkSpans(t *testing.T, s *SkipList[int], v int, expectedSpans []int) {
	t.Helper()
	n := findNode(s, v)
	if n == nil {
		t.Fatalf("value %d not found in skiplist", v)
	}
	if len(n.skips) != len(expectedSpans) {
		t.Fatalf("value %d: span length mismatch: got %d expected %d", v, len(n.skips), len(expectedSpans))
	}
	for i, want := range expectedSpans {
		got := n.skips[i]
		if got != want {
			t.Fatalf("value %d: span level %d: got %d, want %d", v, i, got, want)
		}
	}
}

// ------------------------------------------------------------
// Add Test cases
// ------------------------------------------------------------

func TestAddInit(t *testing.T) {
	skipList := NewSkipList[int]()

	if skipList.head == nil {
		t.Error("Skiplist's head cannot be nil")
	}

	if skipList.head.forward[0] != skipList.tail {
		t.Error("Skiplist's first element is not connected to the last element")
	}
}

func TestAdd(t *testing.T) {
	// 1. setup
	skipList := NewSkipList[int]()
	itemsToAdd := []int{10, 5, 53, 32}
	itemsToVerify := []int{5, 10, 32, 53}
	// 2. action
	for _, item := range itemsToAdd {
		skipList.Add(item)
	}

	// 3. verification
	assertOrder(t, skipList, itemsToVerify)
}

// Adding a duplicate item should not increase the count of elements
func TestAddDuplicate(t *testing.T) {
	// 1. setup
	skipList := NewSkipList[int]()
	itemsToAdd := []int{10, 5, 53, 32}
	itemsToVerify := []int{5, 10, 32, 53}
	// 2. action
	for _, item := range itemsToAdd {
		skipList.Add(item)
	}

	// adding a duplicate item
	skipList.Add(10)

	// 3. verification
	assertOrder(t, skipList, itemsToVerify)
}

func TestAddWithNegativeAndZero(t *testing.T) {
	sl := NewSkipList[int]()
	toAdd := []int{0, -10, 15, -1, 7}
	expected := []int{-10, -1, 0, 7, 15}

	for _, v := range toAdd {
		sl.Add(v)
	}

	assertOrder(t, sl, expected)
}

func TestAddAscendingInput(t *testing.T) {
	sl := NewSkipList[int]()
	toAdd := []int{1, 2, 3, 4, 5}
	expected := []int{1, 2, 3, 4, 5}

	for _, v := range toAdd {
		sl.Add(v)
	}

	assertOrder(t, sl, expected)
}

func TestAddDescendingInput(t *testing.T) {
	sl := NewSkipList[int]()
	toAdd := []int{5, 4, 3, 2, 1}
	expected := []int{1, 2, 3, 4, 5}

	for _, v := range toAdd {
		sl.Add(v)
	}

	assertOrder(t, sl, expected)
}

func TestAddMultipleDuplicates(t *testing.T) {
	sl := NewSkipList[int]()
	toAdd := []int{10, 10, 10, 10}
	expected := []int{10}

	for _, v := range toAdd {
		sl.Add(v)
	}

	assertOrder(t, sl, expected)
}

func TestMaxLevelAndNodeLevels(t *testing.T) {
	sl := NewSkipList[int]()

	// Add a bunch of elements to exercise randomLevel
	for i := range 1000 {
		sl.Add(i)
	}

	if sl.maxLevel > MAX_LEVEL_CAP {
		t.Fatalf("maxLevel exceeded MAX_LEVEL_CAP: got %d, cap %d", sl.maxLevel, MAX_LEVEL_CAP)
	}

	// Check each node's level is within bounds and consistent
	for curr := sl.head; curr != nil; curr = curr.forward[0] {
		nodeLevel := len(curr.forward) - 1
		if curr == sl.tail && nodeLevel != -1 {
			t.Fatalf("last node should have 0 forward pointers")
		}
		if curr != sl.tail && nodeLevel < 0 {
			t.Fatalf("node %d has invalid level %d", curr.val, nodeLevel)
		}
		if curr == sl.head && nodeLevel != MAX_LEVEL_CAP {
			t.Fatalf("last node should have the highest possible level")
		}
		if curr != sl.head && nodeLevel > sl.maxLevel {
			t.Fatalf("node %d has level %d greater than list maxLevel %d", curr.val, nodeLevel, sl.maxLevel)
		}

		if curr == sl.tail {
			break
		}
	}
}

func TestForwardPointersMonotonic(t *testing.T) {
	sl := NewSkipList[int]()

	vals := []int{10, 5, 53, 32, 7, 1, 100}
	for _, v := range vals {
		sl.Add(v)
	}

	for curr := sl.head; curr != sl.tail; curr = curr.forward[0] {
		for level := 0; level < len(curr.forward) && curr.forward[level] != nil; level++ {
			next := curr.forward[level]
			if next.val < curr.val {
				t.Fatalf("forward pointer at level %d is not monotonic: %d -> %d", level, curr.val, next.val)
			}
		}
	}
}

// ---------- Tests: deterministic InsertWithLevel (recommended) ----------

func TestSpans_AfterDeterministicInsert_Simple(t *testing.T) {
	s := NewSkipList[int]()
	s.InsertAtLevel(1, 1) // level 1
	s.InsertAtLevel(2, 2) // level 2
	s.InsertAtLevel(3, 1) // level 1
	s.InsertAtLevel(4, 3) // level 3

	checkSpans(t, s, 1, []int{1, 1})
	checkSpans(t, s, 2, []int{1, 1, 2})
	checkSpans(t, s, 3, []int{1, 1})
	checkSpans(t, s, 4, []int{1, 1, 1, 1})
}

func TestSpans_MultiLevelInsert_Positioning(t *testing.T) {
	// Build a slightly bigger deterministic list and assert spans for a middle node.
	s := NewSkipList[int]()

	levels := map[int]int{
		1: 1,
		2: 3,
		3: 1,
		4: 2,
		5: 1,
		6: 4,
		7: 1,
	}
	for v := 1; v <= 7; v++ {
		s.InsertAtLevel(v, levels[v])
	}

	checkSpans(t, s, 4, []int{1, 1, 2})
}

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
	s.InsertAtLevel(42, MAX_LEVEL_CAP)

	if s.maxLevel != MAX_LEVEL_CAP {
		t.Fatalf("maxLevel should be %d, got %d", MAX_LEVEL_CAP, s.maxLevel)
	}

	node := findNode(s, 42)
	if node == nil {
		t.Fatalf("node 42 not found")
	}
	if len(node.forward) != MAX_LEVEL_CAP+1 {
		t.Fatalf("node at MAX_LEVEL_CAP should have %d forward pointers, got %d", MAX_LEVEL_CAP+1, len(node.forward))
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
	for level := 3; level <= MAX_LEVEL_CAP; level++ {
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

// Helper function to check if a slice contains a value
func contains(slice []int, val int) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
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

// ------------------------------------------------------------
// Delete Test cases
// ------------------------------------------------------------

func TestDelete(t *testing.T) {
	skipList := NewSkipList[int]()
	itemsToAdd := []int{10, 5, 53, 32}
	itemsToVerify := []int{5, 32, 53}

	for _, item := range itemsToAdd {
		skipList.Add(item)
	}

	skipList.Delete(10)

	assertOrder(t, skipList, itemsToVerify)
}

func TestDeleteAbsentElement(t *testing.T) {
	skipList := NewSkipList[int]()
	itemsToAdd := []int{10, 5, 53, 32}
	itemsToVerify := []int{5, 10, 32, 53}

	for _, item := range itemsToAdd {
		skipList.Add(item)
	}

	skipList.Delete(100)

	assertOrder(t, skipList, itemsToVerify)
}

func TestDeleteSingleElementList(t *testing.T) {
	skipList := NewSkipList[int]()
	skipList.Add(42)

	skipList.Delete(42)

	if skipList.length != 0 {
		t.Fatalf("expected length to be 0 got %d", skipList.length)
	}
}

func TestDeleteMultipleElements(t *testing.T) {
	skipList := NewSkipList[int]()
	itemsToAdd := []int{10, 5, 53, 32, 7}
	for _, item := range itemsToAdd {
		skipList.Add(item)
	}

	// delete a couple of elements in the middle
	skipList.Delete(5)
	skipList.Delete(32)

	expected := []int{7, 10, 53}
	assertOrder(t, skipList, expected)
}

func TestDeleteSameElementTwice(t *testing.T) {
	skipList := NewSkipList[int]()
	itemsToAdd := []int{10, 5, 53, 32}
	for _, item := range itemsToAdd {
		skipList.Add(item)
	}

	skipList.Delete(10)
	// second delete should be a no-op
	skipList.Delete(10)

	expected := []int{5, 32, 53}
	assertOrder(t, skipList, expected)
}

func TestDeleteAllElements(t *testing.T) {
	skipList := NewSkipList[int]()
	itemsToAdd := []int{10, 5, 53, 32}
	for _, item := range itemsToAdd {
		skipList.Add(item)
	}

	for _, item := range itemsToAdd {
		skipList.Delete(item)
	}

	if skipList.length != 0 {
		t.Fatalf("expected length to be 0, but got %d", skipList.length)
	}
}

func TestDeleteMultiLevelNodePreservesLevels(t *testing.T) {
	skipList := NewSkipList[int]()

	// Insert enough elements to very likely get multi-level nodes.
	for i := 1; i <= 200; i++ {
		skipList.Add(i)
	}

	// Choose a node that we know has > 1 level.
	node := findMultiLevelNode(skipList)
	if node == nil {
		t.Fatalf("no node with level > 0 found; randomLevel may have produced only level 0 nodes, please rerun the test")
	}
	deleteVal := node.val

	// Snapshot of all levels before deletion.
	beforeLevels := collectLevelValues(skipList)

	// Delete the chosen multi-level node.
	skipList.Delete(deleteVal)

	// Snapshot of all levels after deletion.
	afterLevels := collectLevelValues(skipList)

	// For each level: ordering valid, and only deleteVal removed.
	for level := 0; level <= skipList.maxLevel; level++ {
		before := beforeLevels[level]
		after := afterLevels[level]

		assertStrictlyIncreasing(t, before, level)
		assertStrictlyIncreasing(t, after, level)

		expectedAfter, removed := removeOne(before, deleteVal)
		if !removed {
			// Node was not present at this level; the level should be unchanged.
			if !slicesEqual(before, after) {
				t.Fatalf("level %d changed even though node %d was not present at this level", level, deleteVal)
			}
		} else {
			// Node was present; level should be identical except for that single removal.
			if !slicesEqual(expectedAfter, after) {
				t.Fatalf("level %d sequence incorrect after deleting %d; expected %v, got %v",
					level, deleteVal, expectedAfter, after)
			}
		}
	}
}

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
	for level := 3; level <= MAX_LEVEL_CAP; level++ {
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

// ------------------------------------------------------------
// Search Test cases
// ------------------------------------------------------------

func TestSearch(t *testing.T) {
	skipList := NewSkipList[int]()
	itemsToAdd := []int{10, 5, 53, 32}
	itemsToVerify := []int{5, 10, 32, 53}
	for _, item := range itemsToAdd {
		skipList.Add(item)
	}

	for len(itemsToVerify) > 0 {
		item := itemsToVerify[0]
		if node, present := skipList.SearchByValue(item); present {
			for i := 0; i < len(itemsToVerify); i++ {
				if node == nil {
					t.Fatalf("list ended at %d", itemsToVerify[i])
				}
				if node.val != itemsToVerify[i] {
					t.Fatalf("item %d not present in skip list at correct position", itemsToVerify[i])
				}
				if len(node.forward) == 0 {
					t.Fatal("node reached end!")
				}
				node = node.forward[0]
			}
		} else {
			t.Fatalf("item %d not found in skip list", item)
		}
		itemsToVerify = itemsToVerify[1:]
	}
}

func TestSearchAbsentElement(t *testing.T) {
	skipList := NewSkipList[int]()
	itemsToAdd := []int{10, 5, 53, 32}
	for _, item := range itemsToAdd {
		skipList.Add(item)
	}

	if _, present := skipList.SearchByValue(100); present {
		t.Fatalf("item 100 found in skip list but should not be present")
	}

	if _, present := skipList.SearchByValue(0); present {
		t.Fatalf("item 0 found in skip list but should not be present")
	}

	if _, present := skipList.SearchByValue(3); present {
		t.Fatalf("item 0 found in skip list but should not be present")
	}

	if _, present := skipList.SearchByValue(33); present {
		t.Fatalf("item 0 found in skip list but should not be present")
	}
}

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
