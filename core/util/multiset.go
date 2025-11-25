package core

type MultiSet[T comparable] map[T]int

func NewMultiSet[T comparable](objs ...T) MultiSet[T] {
	set := MultiSet[T]{}
	for _, obj := range objs {
		set[obj] += 1
	}
	return set
}

func (set MultiSet[T]) Add(obj ...T) {
	for _, o := range obj {
		set[o] += 1
	}
}

func (set MultiSet[T]) In(obj T) bool {
	_, ok := set[obj]
	return ok
}

func (set MultiSet[T]) Count(obj T) int {
	count, ok := set[obj]
	if !ok {
		return 0
	} else {
		return count
	}

}

func (set MultiSet[T]) Remove(objs ...T) {
	for _, obj := range objs {
		count := set[obj]
		if count == 0 {
			delete(set, obj)
		} else {
			set[obj] = count - 1
		}
	}
}

func (set *MultiSet[T]) Clear() {
	*set = MultiSet[T]{}
}

func (set MultiSet[T]) Iter(yield func(T) bool) {
	for t, count := range set {
		for range count {
			if !yield(t) {
				return
			}
		}
	}
}

func (set MultiSet[T]) ToSlice() (ret []T) {
	for t, count := range set {
		for range count {
			ret = append(ret, t)
		}
	}
	return ret
}

func (set MultiSet[T]) Size() (size int) {
	for _, count := range set {
		size += count
	}
	return size
}

func (set MultiSet[T]) UniqueElementLen() int {
	return len(set)
}

func (set MultiSet[T]) UniqueElements() (ret []T) {
	for t := range set {
		ret = append(ret, t)
	}
	return ret
}
