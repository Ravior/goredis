package sortedset

import "math/rand"

const (
	maxLevel = 16
)

// Element is a key-score pair
type Element struct {
	Member string
	Score  float64
}

// Level aspect of a node
type Level struct {
	forward *node // forward node has greater score
	span    int64 // forward node的节点跨度，即到达一个节点中间有几个节点(包括目标下一个节点)
}

type node struct {
	Element
	backward *node
	level    []*Level // level[0] is base level
}

func createNode(level int16, score float64, member string) *node {
	n := &node{
		Element: Element{
			Member: member,
			Score:  score,
		},
		level: make([]*Level, level),
	}
	for i := range n.level {
		n.level[i] = &Level{}
	}
	return n
}

type skipList struct {
	header *node
	tail   *node
	length int64
	level  int16
}

func CreateSkipList() *skipList {
	return &skipList{
		header: createNode(maxLevel, 0, ""),
		level:  1,
	}
}

func randomLevel() int16 {
	level := int16(1)
	for float32(rand.Int31()&0xFFFF) < (0.25 * 0xFFFF) {
		level++
	}
	if level < maxLevel {
		return level
	}
	return maxLevel
}

func (skiplist *skipList) insert(score float64, member string) *node {
	// 存储搜索路径
	update := make([]*node, maxLevel)
	// 存储经过的节点跨度
	rank := make([]int64, maxLevel)

	x := skiplist.header
	// 找到插入的位置
	for i := skiplist.level; i >= 0; i-- {
		// store rank that is crossed to reach the insert position
		if i == skiplist.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}

		// 如果 score 相等，还需要比较 value 值
		if x.level[i] != nil {
			for x.level[i].forward != nil &&
				(x.level[i].forward.Score < score ||
					(x.level[i].forward.Score == score && x.level[i].forward.Member < member)) {
				rank[i] += x.level[i].span
				x = x.level[i].forward
			}
		}
		// 记录 "搜索路径"
		update[i] = x
	}

	level := randomLevel()
	// 如果随机生成的 level 超过了当前最大 level 需要更新跳跃表的信息
	if level > skiplist.level {
		for i := skiplist.level; i < level; i++ {
			update[i] = skiplist.header
			update[i].level[i].span = skiplist.length
		}
		skiplist.level = level
	}

	// 创建新节点
	x = createNode(level, score, member)
	// 插入节点
	for i := int16(0); i < level; i++ {
		x.level[i].forward = update[i].level[i].forward
		update[i].level[i].forward = x

		// update span covered by update[i] as x is inserted here
		x.level[i].span = update[i].level[i].span - (rank[0] - rank[i])
		update[i].level[i].span = (rank[0] - rank[i]) + 1
	}

	// increment span for untouched levels
	for i := level; i < skiplist.level; i++ {
		update[i].level[i].span++
	}

	// set backward node
	if update[0] == skiplist.header {
		x.backward = nil
	} else {
		x.backward = update[0]
	}
	if x.level[0].forward != nil {
		x.level[0].forward.backward = x
	} else {
		skiplist.tail = x
	}

	// 节点数量+1
	skiplist.length++

	return x
}

func (skiplist *skipList) removeNode(node *node, update []*node) {
	for i := int16(0); i < skiplist.level; i++ {
		if update[i].level[i].forward == node {
			update[i].level[i].span += node.level[i].span - 1
			update[i].level[i].forward = node.level[i].forward
		} else {
			update[i].level[i].span--
		}
	}
	if node.level[0].forward != nil {
		node.level[0].forward.backward = node.backward
	} else {
		skiplist.tail = node.backward
	}
	for skiplist.level > 1 && skiplist.header.level[skiplist.level-1].forward == nil {
		skiplist.level--
	}
	skiplist.length--
}

func (skiplist *skipList) remove(member string, score float64) bool {
	update := make([]*node, maxLevel)
	x := skiplist.header
	for i := skiplist.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil &&
			(x.level[i].forward.Score < score ||
				(x.level[i].forward.Score == score && x.level[i].forward.Member < member)) {
			x = x.level[i].forward
		}
		update[i] = x
	}
	/* We may have multiple elements with the same score, what we need
	 * is to find the element with both the right score and object. */
	x = x.level[0].forward
	if x != nil && x.Score == score && x.Member == member {
		skiplist.removeNode(x, update)
		return true
	}
	return false
}

func (skiplist *skipList) getRank(member string, score float64) int64 {
	var rank int64 = 0
	x := skiplist.header
	for i := skiplist.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil &&
			(x.level[i].forward.Score < score || x.level[i].forward.Score == score &&
				x.level[i].forward.Member <= member) {
			rank += x.level[i].span
			x = x.level[i].forward
		}

		if x.Member == member {
			return rank
		}
	}
	return 0
}

func (skiplist *skipList) getByRank(rank int64) *node {
	var i int64 = 0
	x := skiplist.header
	for level := skiplist.level - 1; level >= 0; level-- {
		for x.level[level].forward != nil && (i+x.level[level].span) <= rank {
			i += x.level[level].span
			x = x.level[level].forward
		}
		if i == rank {
			return x
		}
	}
	return nil
}

func (skiplist *skipList) hasInRange(min *ScoreBorder, max *ScoreBorder) bool {
	if min.Value > max.Value || (min.Value == max.Value && (min.Exclude || max.Exclude)) {
		return false
	}
	x := skiplist.tail
	if x == nil || !min.less(x.Score) {
		return false
	}
	x = skiplist.header.level[0].forward
	if x == nil || !max.greater(x.Score) {
		return false
	}
	return true
}

func (skiplist *skipList) getFirstInScoreRange(min *ScoreBorder, max *ScoreBorder) *node {
	if !skiplist.hasInRange(min, max) {
		return nil
	}
	n := skiplist.header
	// scan from top level
	for level := skiplist.level - 1; level >= 0; level-- {
		// if forward is not in range than move forward
		for n.level[level].forward != nil && !min.less(n.level[level].forward.Score) {
			n = n.level[level].forward
		}
	}
	/* This is an inner range, so the next node cannot be NULL. */
	n = n.level[0].forward
	if !max.greater(n.Score) {
		return nil
	}
	return n
}

func (skiplist *skipList) getLastInScoreRange(min *ScoreBorder, max *ScoreBorder) *node {
	if !skiplist.hasInRange(min, max) {
		return nil
	}
	n := skiplist.header
	// scan from top level
	for level := skiplist.level - 1; level >= 0; level-- {
		for n.level[level].forward != nil && max.greater(n.level[level].forward.Score) {
			n = n.level[level].forward
		}
	}
	if !min.less(n.Score) {
		return nil
	}
	return n
}

/*
 * return removed elements
 */
func (skiplist *skipList) RemoveRangeByScore(min *ScoreBorder, max *ScoreBorder, limit int) (removed []*Element) {
	update := make([]*node, maxLevel)
	removed = make([]*Element, 0)
	// find backward nodes (of target range) or last x of each level
	x := skiplist.header
	for i := skiplist.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil {
			if min.less(x.level[i].forward.Score) { // already in range
				break
			}
			x = x.level[i].forward
		}
		update[i] = x
	}

	// x is the first one within range
	x = x.level[0].forward

	// remove nodes in range
	for x != nil {
		if !max.greater(x.Score) { // already out of range
			break
		}
		next := x.level[0].forward
		removedElement := x.Element
		removed = append(removed, &removedElement)
		skiplist.removeNode(x, update)
		if limit > 0 && len(removed) == limit {
			break
		}
		x = next
	}
	return removed
}

// 1-based rank, including start, exclude stop
func (skiplist *skipList) RemoveRangeByRank(start int64, stop int64) (removed []*Element) {
	var i int64 = 0 // rank of iterator
	update := make([]*node, maxLevel)
	removed = make([]*Element, 0)

	// scan from top level
	x := skiplist.header
	for level := skiplist.level - 1; level >= 0; level-- {
		for x.level[level].forward != nil && (i+x.level[level].span) < start {
			i += x.level[level].span
			x = x.level[level].forward
		}
		update[level] = x
	}

	i++
	x = x.level[0].forward // first x in range

	// remove nodes in range
	for x != nil && i < stop {
		next := x.level[0].forward
		removedElement := x.Element
		removed = append(removed, &removedElement)
		skiplist.removeNode(x, update)
		x = next
		i++
	}
	return removed
}
