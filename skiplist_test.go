package main

import "testing"

// ------------------------------------------------------------
// Helper functions
// ------------------------------------------------------------

func createSkipList() *SkipList {
	first := &Node{val: MinInt, forward: make([]*Node, MAX_LEVEL_CAP + 1)}
	last := &Node{val: MaxInt, forward: make([]*Node, 0)}

	for i := 0; i <= MAX_LEVEL_CAP; i++ {
		first.forward[i] = last
	}

	return &SkipList{
		head: first,
		tail: last,
		maxLevel: 0,
	}
}

func assertOrder(t *testing.T, sl *SkipList, expected []int) {
    t.Helper()

    curr := sl.head
    i := 0

    for {
        if i >= len(expected) {
            t.Fatalf("skiplist has more elements than expected; extra element %d", curr.val)
        }

        if curr.val != expected[i] {
            t.Fatalf("value mismatch at index %d: got %d, want %d", i, curr.val, expected[i])
        }

        if curr == sl.tail {
            break
        }

        curr = curr.forward[0]
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
func collectLevelValues(sl *SkipList) [][]int {
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
func findMultiLevelNode(sl *SkipList) *Node {
	for curr := sl.head.forward[0]; curr != nil && curr != sl.tail; curr = curr.forward[0] {
		if len(curr.forward) > 1 {
			return curr
		}
	}
	return nil
}

// ------------------------------------------------------------
// Add Test cases
// ------------------------------------------------------------

func TestAddInit(t *testing.T) {
	skipList := createSkipList()
	
	if skipList.head == nil {
		t.Error("Skiplist's head cannot be nil")
	}

	if skipList.tail == nil {
		t.Error("Skiplist's tail cannot be nil")
	}

	if skipList.head.val != MinInt {
		t.Error("Skiplist's first element has to be", MinInt)
	}

	if skipList.tail.val != MaxInt {
		t.Error("Skiplist's last element has to be", MaxInt)
	}

	if skipList.head.forward[0] != skipList.tail {
		t.Error("Skiplist's first element is not connected to the last element")
	}
}

func TestAdd(t *testing.T) {
	// 1. setup
	skipList := createSkipList()
	itemsToAdd := []int{10, 5, 53, 32}
	itemsToVerify := []int{MinInt, 5, 10, 32, 53, MaxInt}
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
	skipList := createSkipList()
	itemsToAdd := []int{10, 5, 53, 32}
	itemsToVerify := []int{MinInt, 5, 10, 32, 53, MaxInt}
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
    sl := createSkipList()
    toAdd := []int{0, -10, 15, -1, 7}
    expected := []int{MinInt, -10, -1, 0, 7, 15, MaxInt}

    for _, v := range toAdd {
        sl.Add(v)
    }

    assertOrder(t, sl, expected)
}

func TestAddAscendingInput(t *testing.T) {
    sl := createSkipList()
    toAdd := []int{1, 2, 3, 4, 5}
    expected := []int{MinInt, 1, 2, 3, 4, 5, MaxInt}

    for _, v := range toAdd {
        sl.Add(v)
    }

    assertOrder(t, sl, expected)
}

func TestAddDescendingInput(t *testing.T) {
    sl := createSkipList()
    toAdd := []int{5, 4, 3, 2, 1}
    expected := []int{MinInt, 1, 2, 3, 4, 5, MaxInt}

    for _, v := range toAdd {
        sl.Add(v)
    }

    assertOrder(t, sl, expected)
}

func TestAddMultipleDuplicates(t *testing.T) {
    sl := createSkipList()
    toAdd := []int{10, 10, 10, 10}
    expected := []int{MinInt, 10, MaxInt}

    for _, v := range toAdd {
        sl.Add(v)
    }

    assertOrder(t, sl, expected)
}

func TestMaxLevelAndNodeLevels(t *testing.T) {
    sl := createSkipList()

    // Add a bunch of elements to exercise randomLevel
    for i := range 1000 {
        sl.Add(i)
    }

    if sl.maxLevel > MAX_LEVEL_CAP {
        t.Fatalf("maxLevel exceeded MAX_LEVEL_CAP: got %d, cap %d", sl.maxLevel, MAX_LEVEL_CAP)
    }

    // Check each node's level is within bounds and consistent
    for curr := sl.head; ; {
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
        curr = curr.forward[0]
    }
}

func TestForwardPointersMonotonic(t *testing.T) {
    sl := createSkipList()

    vals := []int{10, 5, 53, 32, 7, 1, 100}
    for _, v := range vals {
        sl.Add(v)
    }

    for curr := sl.head; curr != sl.tail; curr = curr.forward[0] {
        for level := 0; level < len(curr.forward); level++ {
            next := curr.forward[level]
            if next == nil {
                t.Fatalf("node %d has nil forward pointer at level %d", curr.val, level)
            }
            if next.val < curr.val {
                t.Fatalf("forward pointer at level %d is not monotonic: %d -> %d", level, curr.val, next.val)
            }
        }
    }
}

// ------------------------------------------------------------
// Delete Test cases
// ------------------------------------------------------------

func TestDelete(t *testing.T) {
	skipList := createSkipList()
	itemsToAdd := []int{10, 5, 53, 32}
	itemsToVerify := []int{MinInt, 5, 32, 53, MaxInt}

	for _, item := range(itemsToAdd) {
		skipList.Add(item)
	}

	skipList.Delete(10)

	assertOrder(t, skipList, itemsToVerify)
}

func TestDeleteAbsentElement(t *testing.T) {
	skipList := createSkipList()
	itemsToAdd := []int{10, 5, 53, 32}
	itemsToVerify := []int{MinInt, 5, 10, 32, 53, MaxInt}

	for _, item := range(itemsToAdd) {
		skipList.Add(item)
	}

	skipList.Delete(100)

	assertOrder(t, skipList, itemsToVerify)
}

func TestDeleteFirstOrLastElement(t *testing.T) {
	skipList := createSkipList()
	itemsToAdd := []int{10, 5, 53, 32}

	for _, item := range(itemsToAdd) {
		skipList.Add(item)
	}

	assertPanic(t, func() {
		skipList.Delete(MinInt)
	})

	assertPanic(t, func() {
		skipList.Delete(MaxInt)
	})
}

func TestDeleteSingleElementList(t *testing.T) {
	skipList := createSkipList()
	skipList.Add(42)

	skipList.Delete(42)

	expected := []int{MinInt, MaxInt}
	assertOrder(t, skipList, expected)
}

func TestDeleteMultipleElements(t *testing.T) {
	skipList := createSkipList()
	itemsToAdd := []int{10, 5, 53, 32, 7}
	for _, item := range itemsToAdd {
		skipList.Add(item)
	}

	// delete a couple of elements in the middle
	skipList.Delete(5)
	skipList.Delete(32)

	expected := []int{MinInt, 7, 10, 53, MaxInt}
	assertOrder(t, skipList, expected)
}

func TestDeleteSameElementTwice(t *testing.T) {
	skipList := createSkipList()
	itemsToAdd := []int{10, 5, 53, 32}
	for _, item := range itemsToAdd {
		skipList.Add(item)
	}

	skipList.Delete(10)
	// second delete should be a no-op
	skipList.Delete(10)

	expected := []int{MinInt, 5, 32, 53, MaxInt}
	assertOrder(t, skipList, expected)
}

func TestDeleteAllElements(t *testing.T) {
	skipList := createSkipList()
	itemsToAdd := []int{10, 5, 53, 32}
	for _, item := range itemsToAdd {
		skipList.Add(item)
	}

	for _, item := range itemsToAdd {
		skipList.Delete(item)
	}

	expected := []int{MinInt, MaxInt}
	assertOrder(t, skipList, expected)
}

func TestDeleteMultiLevelNodePreservesLevels(t *testing.T) {
	skipList := createSkipList()

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
// Search Test cases
// ------------------------------------------------------------

func TestSearch(t *testing.T) {
	skipList := createSkipList()
	itemsToAdd := []int{10, 5, 53, 32}
	for _, item := range itemsToAdd {
		skipList.Add(item)
	}

	for _, item := range itemsToAdd {	
		if !skipList.Search(item) {
			t.Fatalf("item %d not found in skip list", item)
		}
	}
}

func TestSearchAbsentElement(t *testing.T) {
	skipList := createSkipList()
	itemsToAdd := []int{10, 5, 53, 32}
	for _, item := range itemsToAdd {
		skipList.Add(item)
	}

	if skipList.Search(100) {
		t.Fatalf("item 100 found in skip list but should not be present")
	}

	if skipList.Search(0) {
		t.Fatalf("item 0 found in skip list but should not be present")
	}

	if skipList.Search(3) {
		t.Fatalf("item 0 found in skip list but should not be present")
	}

	if skipList.Search(33) {
		t.Fatalf("item 0 found in skip list but should not be present")
	}
}

func TestSearchSentinelElements(t *testing.T) {
	skipList := createSkipList()

	if skipList.Search(MinInt) {
		t.Fatalf("item %d found in skip list but should not be present", MinInt)
	}

	if skipList.Search(MaxInt) {
		t.Fatalf("item %d found in skip list but should not be present", MaxInt)
	}
}

func TestDeleteAndSearch(t *testing.T) {
	skipList := createSkipList()
	itemsToAdd := []int{10, 5, 53, 32}
	for _, item := range itemsToAdd {
		skipList.Add(item)
	}

	skipList.Delete(53)
	if skipList.Search(53) {
		t.Fatalf("item %d found in skip list but should not be present", MaxInt)
	}
}