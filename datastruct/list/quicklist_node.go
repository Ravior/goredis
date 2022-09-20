package list

import "container/list"

type quickListNode struct {
	node   *list.Element
	offset int
	ql     *QuickList
}

func (qln *quickListNode) page() []interface{} {
	return qln.node.Value.([]interface{})
}

func (qln *quickListNode) get() interface{} {
	return qln.page()[qln.offset]
}

func (qln *quickListNode) next() bool {
	page := qln.page()
	if qln.offset < len(page)-1 {
		qln.offset++
		return true
	}
	if qln.node == qln.ql.data.Back() {
		qln.offset = len(page)
		return false
	}
	qln.offset = 0
	qln.node = qln.node.Next()
	return true
}

func (qln *quickListNode) prev() bool {
	if qln.offset > 0 {
		qln.offset--
		return true
	}
	if qln.node == qln.ql.data.Front() {
		qln.offset = -1
		return false
	}
	qln.node = qln.node.Prev()
	prevPage := qln.node.Value.([]interface{})
	qln.offset = len(prevPage) - 1
	return true
}

func (qln *quickListNode) atEnd() bool {
	if qln.ql.data.Len() == 0 {
		return true
	}
	if qln.node != qln.ql.data.Back() {
		return false
	}
	page := qln.page()
	return qln.offset == len(page)
}

func (qln *quickListNode) atBegin() bool {
	if qln.ql.data.Len() == 0 {
		return true
	}
	if qln.node != qln.ql.data.Front() {
		return false
	}
	return qln.offset == -1
}

func (qln *quickListNode) set(val interface{}) {
	page := qln.page()
	page[qln.offset] = val
}

func (qln *quickListNode) remove() interface{} {
	page := qln.page()
	val := page[qln.offset]
	page = append(page[:qln.offset], page[qln.offset+1:]...)
	if len(page) > 0 {
		qln.node.Value = page
		// page is not empty, update iter.offset only
		qln.node.Value = page
		if qln.offset == len(page) {
			// removed page[-1], doubleLinkedListNode should move to next page
			if qln.node != qln.ql.data.Back() {
				qln.node = qln.node.Next()
				qln.offset = 0
			}
			// else: assert iter.atEnd() == true
		}
	} else {
		// page is empty, update iter.doubleLinkedListNode and iter.offset
		if qln.node == qln.ql.data.Back() {
			// removed last element, ql is empty now
			qln.ql.data.Remove(qln.node)
			qln.node = nil
			qln.offset = 0
		} else {
			nextNode := qln.node.Next()
			qln.ql.data.Remove(qln.node)
			qln.node = nextNode
			qln.offset = 0
		}
	}
	qln.ql.size--
	return val
}
