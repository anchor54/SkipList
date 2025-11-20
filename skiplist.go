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
	head       *Node
	tail       *Node
	maxLevel   int
	length     int
	levelCount [MAX_LEVEL_CAP + 1]int
}

func NewNode(val int, forwards int) *Node {
	return &Node{
		val:     val,
		forward: make([]*Node, forwards),
		skips:   make([]int, forwards),
	}
}

func NewSkipList() *SkipList {
	var zero int
	first := NewNode(zero, MAX_LEVEL_CAP+1)

	for i := 0; i <= MAX_LEVEL_CAP; i++ {
		first.forward[i] = nil
		first.skips[i] = 1
	}

	return &SkipList{
		head:     first,
		tail:     nil,
		maxLevel: 0,
		length:   0,
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
		for curr.forward[currLevel] != nil && curr.forward[currLevel].val <= val {
			skipped += curr.skips[currLevel]
			curr = curr.forward[currLevel]
		}

		// do nothing if the value is already added
		if curr != sl.head && curr.val == val {
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

		sl.levelCount[i]++
	}

	for i := lvl + 1; i <= sl.maxLevel; i++ {
		heirarchy[i].skips[i]++
	}
	for i := sl.maxLevel + 1; i <= MAX_LEVEL_CAP; i++ {
		sl.head.skips[i]++
	}

	sl.length++
}

func (sl *SkipList) Delete(val int) {
	heirarchy := [MAX_LEVEL_CAP + 1]*Node{}
	curr := sl.head
	var nodeToDelete *Node = nil

	for currLevel := sl.maxLevel; currLevel >= 0; currLevel-- {
		for curr.forward[currLevel] != nil && curr.forward[currLevel].val < val {
			curr = curr.forward[currLevel]
		}

		nodeToDelete = curr.forward[currLevel]
		if nodeToDelete != nil && nodeToDelete.val == val {
			curr.skips[currLevel] += nodeToDelete.skips[currLevel] - 1
			curr.forward[currLevel] = nodeToDelete.forward[currLevel]
			nodeToDelete.forward[currLevel] = nil

			sl.levelCount[currLevel]--
			if sl.levelCount[currLevel] == 0 {
				sl.maxLevel--
			}
		} else {
			heirarchy[currLevel] = curr
		}
	}

	// if the node to delete was found only then reduce the span of the remaining heirarchy
	if nodeToDelete == nil {
		return
	}

	currLevel := len(nodeToDelete.skips)
	for ; currLevel <= sl.maxLevel; currLevel++ {
		heirarchy[currLevel].skips[currLevel]--
	}
	for ; currLevel <= MAX_LEVEL_CAP; currLevel++ {
		sl.head.skips[currLevel]--
	}

	sl.length--
}

func (sl *SkipList) SearchByValue(val int) (*Node, bool) {
	curr := sl.head

	for currLevel := sl.maxLevel; currLevel >= 0; currLevel-- {
		for curr.forward[currLevel] != nil && curr.forward[currLevel].val <= val {
			curr = curr.forward[currLevel]
		}

		if curr != sl.head && curr.val == val {
			return curr, true
		}
	}
	return nil, false
}

func (sl *SkipList) SearchByRank(rank int) (*Node, bool) {
	if rank < 1 || rank > sl.length {
		return nil, false
	}

	curr := sl.head
	rankUntil := 0

	for currLevel := sl.maxLevel; currLevel >= 0; currLevel-- {
		for curr != nil && rankUntil+curr.skips[currLevel] <= rank {
			rankUntil += curr.skips[currLevel]
			curr = curr.forward[currLevel]
		}

		if rankUntil == rank {
			return curr, true
		}
	}
	return nil, false
}
