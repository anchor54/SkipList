package main

import "math/rand"

var MAX_LEVEL_CAP = 16
var PROBABILILTY float32 = 0.5

type Node struct {
	val int
	forward []*Node
}

type SkipList struct {
	head *Node
	maxLevel int
}

func randomLevel() int {
	lvl := 0
	for rand.Float32() >= PROBABILILTY && lvl < MAX_LEVEL_CAP {
		lvl++
	}
	return lvl
}

func (sl *SkipList) Add(val int) {
	heirarchy := make([]*Node, MAX_LEVEL_CAP + 1)
	curr := sl.head

	for currLevel := sl.maxLevel; currLevel >= 0; currLevel-- {
		for curr.forward[currLevel].val < val {
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
		// get last node of the level
		last := sl.head
		for len(last.forward) > 0 {
			last = last.forward[sl.maxLevel]
		}
		
		for i := sl.maxLevel + 1; i <= lvl; i++ {
			sl.head.forward[i] = last
			heirarchy[i] = sl.head
		}

		sl.maxLevel = lvl
	}

	newNode := Node{
		val: val,
		forward: make([]*Node, MAX_LEVEL_CAP + 1),
	}

	for i := 0; i <= lvl; i++ {
		newNode.forward[i] = heirarchy[i]
		heirarchy[i].forward[i] = &newNode
	}
}