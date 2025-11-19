package main

import (
	"math/rand"
)

const MAX_LEVEL_CAP = 16
const PROBABILILTY float32 = 0.5

// Source - https://stackoverflow.com/a
// Posted by nmichaels, modified by community. See post 'Timeline' for change history
// Retrieved 2025-11-18, License - CC BY-SA 3.0

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

type Node struct {
	val     int
	skips   []int
	forward []*Node
}

type SkipList struct {
	head     *Node
	tail     *Node
	maxLevel int
}

func NewNode(val int, forwards int) *Node {
	return &Node{
		val:     val,
		forward: make([]*Node, forwards),
		skips:   make([]int, forwards),
	}
}

func NewSkipList() *SkipList {
	first := NewNode(MinInt, MAX_LEVEL_CAP+1)
	last := NewNode(MaxInt, 0)

	for i := 0; i <= MAX_LEVEL_CAP; i++ {
		first.forward[i] = last
		first.skips[i] = 1
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
	sl.InsertAtLevel(val, randomLevel())
}

func (sl *SkipList) InsertAtLevel(val int, lvl int) {
	heirarchy := [MAX_LEVEL_CAP + 1]*Node{}
	rank := [MAX_LEVEL_CAP + 1]int{}
	curr := sl.head
	skipped := 0

	for currLevel := sl.maxLevel; currLevel >= 0; currLevel-- {
		for curr.forward[currLevel].val <= val {
			skipped += curr.skips[currLevel]
			curr = curr.forward[currLevel]
		}

		// do nothing if the value is already added
		if curr.val == val {
			return
		}

		heirarchy[currLevel] = curr
		rank[currLevel] = skipped
	}

	if lvl > sl.maxLevel {
		for i := sl.maxLevel + 1; i <= lvl; i++ {
			heirarchy[i] = sl.head
			rank[i] = 0
		}

		sl.maxLevel = lvl
	}

	newNode := NewNode(val, lvl+1)

	for i := 0; i <= lvl; i++ {
		newNode.forward[i] = heirarchy[i].forward[i]
		heirarchy[i].forward[i] = newNode

		newNode.skips[i] = rank[i] + heirarchy[i].skips[i] - skipped
		heirarchy[i].skips[i] = skipped - rank[i] + 1
	}

	for i := lvl + 1; i <= sl.maxLevel; i++ {
		heirarchy[i].skips[i]++
	}
	for i := sl.maxLevel + 1; i <= MAX_LEVEL_CAP; i++ {
		sl.head.skips[i]++
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

func (sl *SkipList) Search(val int) (*Node, bool) {
	if val <= sl.head.val || val >= sl.tail.val {
		return nil, false
	}

	curr := sl.head

	for currLevel := sl.maxLevel; currLevel >= 0; currLevel-- {
		for curr.forward[currLevel].val <= val {
			curr = curr.forward[currLevel]
		}

		if curr.val == val {
			return curr, true
		}
	}
	return nil, false
}
