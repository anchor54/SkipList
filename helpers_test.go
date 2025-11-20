package skiplist

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

// Helper function to check if a slice contains a value
func contains(slice []int, val int) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}
