package main

import "math/rand"

var MAX_LEVEL_CAP = 16
var PROBABILILTY float32 = 0.5

// Source - https://stackoverflow.com/a
// Posted by nmichaels, modified by community. See post 'Timeline' for change history
// Retrieved 2025-11-18, License - CC BY-SA 3.0

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

type Node struct {
	val     int
	forward []*Node
}

type SkipList struct {
	head     *Node
	tail     *Node
	maxLevel int
}

func NewSkipList() *SkipList {
	first := &Node{val: MinInt, forward: make([]*Node, MAX_LEVEL_CAP+1)}
	last := &Node{val: MaxInt, forward: make([]*Node, 0)}

	for i := 0; i <= MAX_LEVEL_CAP; i++ {
		first.forward[i] = last
	}

	return &SkipList{
		head:     first,
		tail:     last,
		maxLevel: 0,
	}
}

func randomLevel() int {
	lvl := 0
	for rand.Float32() >= PROBABILILTY && lvl < MAX_LEVEL_CAP {
		lvl++
	}
	return lvl
}

func (sl *SkipList) Add(val int) {
	heirarchy := make([]*Node, MAX_LEVEL_CAP+1)
	curr := sl.head

	for currLevel := sl.maxLevel; currLevel >= 0; currLevel-- {
		for curr.forward[currLevel].val <= val {
			curr = curr.forward[currLevel]
		}

		// do nothing if the value is already added
		if curr.val == val {
			return
		}

		heirarchy[currLevel] = curr
	}

	lvl := randomLevel()

	if lvl > sl.maxLevel {
		for i := sl.maxLevel + 1; i <= lvl; i++ {
			heirarchy[i] = sl.head
		}

		sl.maxLevel = lvl
	}

	newNode := Node{
		val:     val,
		forward: make([]*Node, lvl+1),
	}

	for i := 0; i <= lvl; i++ {
		newNode.forward[i] = heirarchy[i].forward[i]
		heirarchy[i].forward[i] = &newNode
	}
}

func (sl *SkipList) Delete(val int) {
	if val == sl.head.val {
		panic("Invalid operation: Cannot delete head of list")
	}
	if val == sl.tail.val {
		panic("Invalid operation: Cannot delete tail of list")
	}
	curr := sl.head

	for currLevel := sl.maxLevel; currLevel >= 0; currLevel-- {
		for curr.forward[currLevel].val < val {
			curr = curr.forward[currLevel]
		}

		nodeToDelete := curr.forward[currLevel]
		if nodeToDelete.val == val {
			curr.forward[currLevel] = nodeToDelete.forward[currLevel]
			nodeToDelete.forward[currLevel] = nil
		}
	}
}

func (sl *SkipList) Search(val int) bool {
	if val <= sl.head.val || val >= sl.tail.val {
		return false
	}

	curr := sl.head

	for currLevel := sl.maxLevel; currLevel >= 0; currLevel-- {
		for curr.forward[currLevel].val <= val {
			curr = curr.forward[currLevel]
		}

		if curr.val == val {
			return true
		}
	}
	return false
}
