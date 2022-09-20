package list

// Expected check whether given item is equals to expected value
type Expected func(a interface{}) bool

// Iterator traverses list.
// It receives index and value as params, returns true to continue traversal, while returns false to break
type Iterator func(i int, v interface{}) bool

type List interface {
	Add(val interface{})
	Get(index int) (val interface{})
	Set(index int, val interface{})
	Insert(index int, val interface{})
	Remove(index int) (val interface{})
	RemoveLast() (val interface{})
	RemoveAllByVal(expected Expected) int
	RemoveByVal(expected Expected, count int) int
	ReverseRemoveByVal(expected Expected, count int) int
	Len() int
	ForEach(iterator Iterator)
	Contains(expected Expected) bool
	Range(start int, stop int) []interface{}
}
