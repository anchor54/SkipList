package main

import (
	"cmp"
	"math/rand"
)

const (
	MAX_LEVEL_CAP = 16
	PROBABILILTY float32 = 0.5
)

type Comparator[T any] func(a, b T) int

type Comparable[T any] interface {
	Compare(other T) int
}

type Node[T any] struct {
	val     T
	skips   []int
	forward []*Node[T]
}

type SkipList[T any] struct {
	head       *Node[T]
	tail       *Node[T]
	maxLevel   int
	length     int
	levelCount [MAX_LEVEL_CAP + 1]int
	comparator Comparator[T]
}

func NewNode[T any](val T, forwards int) *Node[T] {
	return &Node[T]{
		val:     val,
		forward: make([]*Node[T], forwards),
		skips:   make([]int, forwards),
	}
}

func NewSkipList[T cmp.Ordered]() *SkipList[T] {
	var zero T
	first := NewNode(zero, MAX_LEVEL_CAP + 1)

	for i := 0; i <= MAX_LEVEL_CAP; i++ {
		first.forward[i] = nil
		first.skips[i] = 1
	}

	return &SkipList[T]{
		head:     first,
		tail:     nil,
		maxLevel: 0,
		length:   0,
		comparator: cmp.Compare[T],
	}
}

func NewComparableSkipList[T Comparable[T]]() *SkipList[T] {
	var zero T
	first := NewNode(zero, MAX_LEVEL_CAP + 1)

	for i := 0; i <= MAX_LEVEL_CAP; i++ {
		first.forward[i] = nil
		first.skips[i] = 1
	}

	return &SkipList[T]{
		head:     first,
		tail:     nil,
		maxLevel: 0,
		length:   0,
		comparator: func(a, b T) int {
			return a.Compare(b)
		},
	}
}

func randomLevel() int {
	lvl := 0
	for rand.Float32() >= PROBABILILTY && lvl < MAX_LEVEL_CAP {
		lvl++
	}
	return lvl
}

func (sl *SkipList[T]) Add(val T) {
	sl.InsertAtLevel(val, randomLevel())
}

func (sl *SkipList[T]) InsertAtLevel(val T, lvl int) {
	heirarchy := [MAX_LEVEL_CAP + 1]*Node[T]{}
	rank := [MAX_LEVEL_CAP + 1]int{}
	curr := sl.head
	skipped := 0

	for currLevel := sl.maxLevel; currLevel >= 0; currLevel-- {
		for curr.forward[currLevel] != nil && sl.comparator(curr.forward[currLevel].val, val) <= 0 {
			skipped += curr.skips[currLevel]
			curr = curr.forward[currLevel]
		}

		// do nothing if the value is already added
		if curr != sl.head && sl.comparator(curr.val, val) == 0 {
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

func (sl *SkipList[T]) Delete(val T) {
	heirarchy := [MAX_LEVEL_CAP + 1]*Node[T]{}
	curr := sl.head
	var nodeToDelete *Node[T] = nil

	for currLevel := sl.maxLevel; currLevel >= 0; currLevel-- {
		for curr.forward[currLevel] != nil && sl.comparator(curr.forward[currLevel].val, val) < 0 {
			curr = curr.forward[currLevel]
		}

		nodeToDelete = curr.forward[currLevel]
		if nodeToDelete != nil && sl.comparator(nodeToDelete.val, val) == 0 {
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

func (sl *SkipList[T]) SearchByValue(val T) (*Node[T], bool) {
	curr := sl.head

	for currLevel := sl.maxLevel; currLevel >= 0; currLevel-- {
		for curr.forward[currLevel] != nil && sl.comparator(curr.forward[currLevel].val, val) <= 0 {
			curr = curr.forward[currLevel]
		}

		if curr != sl.head && sl.comparator(curr.val, val) == 0 {
			return curr, true
		}
	}
	return nil, false
}

func (sl *SkipList[T]) SearchByRank(rank int) (*Node[T], bool) {
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
