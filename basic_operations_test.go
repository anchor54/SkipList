package skiplist

import "testing"

// ------------------------------------------------------------
// Basic Add Test cases
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

	if sl.maxLevel > MaxLevelCap {
		t.Fatalf("maxLevel exceeded MAX_LEVEL_CAP: got %d, cap %d", sl.maxLevel, MaxLevelCap)
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
		if curr == sl.head && nodeLevel != MaxLevelCap {
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
// Basic Delete Test cases
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
// Basic Search Test cases
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
