package main

import "testing"

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