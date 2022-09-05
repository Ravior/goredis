package set

import "github.com/Ravior/goredis/datastruct/dict"

type Set struct {
	dict dict.Dict
}

// Add adds member into set
func (set *Set) Add(val string) int {
	return set.dict.Put(val, nil)
}

// Remove removes member from set
func (set *Set) Remove(val string) int {
	return set.dict.Remove(val)
}

// Has returns true if the val exists in the set
func (set *Set) Has(val string) bool {
	_, exists := set.dict.Get(val)
	return exists
}

// Len returns number of members in the set
func (set *Set) Len() int {
	return set.dict.Len()
}
