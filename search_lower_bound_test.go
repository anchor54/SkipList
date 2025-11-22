package skiplist

import "testing"

func TestGetLowerBound_EmptyList(t *testing.T) {
	sl := NewSkipList[int]()

	if node, found := sl.GetLowerBound(10); node != nil || found {
		t.Errorf("GetLowerBound on empty list should return (nil, false), got (%v, %v)", node, found)
	}
}

func TestGetLowerBound_SingleElementList(t *testing.T) {
	sl := NewSkipList[int]()
	sl.Add(10)

	// Search for the exact element
	if node, found := sl.GetLowerBound(10); !found || node.val != 10 {
		t.Errorf("expected to find 10, got node %v, found %v", node, found)
	}

	// Search for a smaller value
	if node, found := sl.GetLowerBound(5); !found || node.val != 10 {
		t.Errorf("expected to find 10 for lower bound of 5, got node %v, found %v", node, found)
	}

	// Search for a larger value
	if _, found := sl.GetLowerBound(15); found {
		t.Errorf("expected not to find any element for lower bound of 15, but found one")
	}
}

func TestGetLowerBound_MultiElementList(t *testing.T) {
	sl := NewSkipList[int]()
	items := []int{10, 20, 30, 40, 50}
	for _, item := range items {
		sl.Add(item)
	}

	testCases := []struct {
		name          string
		searchValue   int
		expectedValue int
		shouldBeFound bool
	}{
		{"Exact match (middle)", 30, 30, true},
		{"Exact match (first)", 10, 10, true},
		{"Exact match (last)", 50, 50, true},
		{"Between two elements", 25, 30, true},
		{"Smaller than all", 5, 10, true},
		{"Larger than all", 55, 0, false},
		{"Just below an element", 29, 30, true},
		{"Just above an element", 31, 40, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			node, found := sl.GetLowerBound(tc.searchValue)

			if found != tc.shouldBeFound {
				t.Fatalf("expected found to be %v, but got %v", tc.shouldBeFound, found)
			}

			if tc.shouldBeFound && (node == nil || node.val != tc.expectedValue) {
				var nodeVal int
				if node != nil {
					nodeVal = node.val
				}
				t.Errorf("expected to find value %d, but got %d", tc.expectedValue, nodeVal)
			}

			if !tc.shouldBeFound && node != nil {
				t.Errorf("expected not to find any node, but got %v", node)
			}
		})
	}
}

func TestGetLowerBound_WithDuplicatesInInput(t *testing.T) {
	sl := NewSkipList[int]()
	items := []int{10, 20, 20, 30}
	for _, item := range items {
		sl.Add(item)
	}

	// Search for the duplicate value
	if node, found := sl.GetLowerBound(20); !found || node.val != 20 {
		t.Errorf("expected to find 20, got node %v, found %v", node, found)
	}

	// Search for a value between the duplicate
	if node, found := sl.GetLowerBound(21); !found || node.val != 30 {
		t.Errorf("expected to find 30 for lower bound of 21, got node %v, found %v", node, found)
	}
}
