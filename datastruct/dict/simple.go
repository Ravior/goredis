package dict

type SimpleDict struct {
	m map[string]interface{}
}

func NewSimpleDict() *SimpleDict {
	return &SimpleDict{
		m: make(map[string]interface{}),
	}
}

func (dict *SimpleDict) Get(key string) (val interface{}, exists bool) {
	val, ok := dict.m[key]
	return val, ok
}

// Len returns the number of dict
func (dict *SimpleDict) Len() int {
	if dict.m == nil {
		panic("m is nil")
	}
	return len(dict.m)
}

// Put puts key value into dict and returns the number of new inserted key-value
func (dict *SimpleDict) Put(key string, val interface{}) (result int) {
	_, existed := dict.m[key]
	dict.m[key] = val
	if existed {
		return 0
	}
	return 1
}

// Remove removes the key and return the number of deleted key-value
func (dict *SimpleDict) Remove(key string) (result int) {
	_, existed := dict.m[key]
	delete(dict.m, key)
	if existed {
		return 1
	}
	return 0
}

// ForEach traversal the dict
func (dict *SimpleDict) ForEach(iterator Iterator) {
	for k, v := range dict.m {
		if !iterator(k, v) {
			break
		}
	}
}

// Keys returns all keys in dict
func (dict *SimpleDict) Keys() []string {
	result := make([]string, len(dict.m))
	i := 0
	for k := range dict.m {
		result[i] = k
		i++
	}
	return result
}

// RandomKeys randomly returns keys of the given number, may contain duplicated key
func (dict *SimpleDict) RandomKeys(limit int) []string {
	result := make([]string, limit)
	for i := 0; i < limit; i++ {
		for k := range dict.m {
			result[i] = k
			break
		}
	}
	return result
}

// RandomDistinctKeys randomly returns keys of the given number, won't contain duplicated key
func (dict *SimpleDict) RandomDistinctKeys(limit int) []string {
	size := limit
	if size > len(dict.m) {
		size = len(dict.m)
	}
	result := make([]string, size)
	i := 0
	for k := range dict.m {
		if i == size {
			break
		}
		result[i] = k
		i++
	}
	return result
}

func (dict *SimpleDict) Clear() {
	*dict = *NewSimpleDict()
}