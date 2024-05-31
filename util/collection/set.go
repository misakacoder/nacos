package collection

type HashSet[E comparable] struct {
	data map[E]struct{}
}

func NewHashSet[E comparable](elements []E) *HashSet[E] {
	data := map[E]struct{}{}
	for _, element := range elements {
		data[element] = struct{}{}
	}
	return &HashSet[E]{data: data}
}

func (hashSet *HashSet[E]) Add(value E) {
	hashSet.data[value] = struct{}{}
}

func (hashSet *HashSet[E]) Contains(value E) bool {
	_, ok := hashSet.data[value]
	return ok
}
