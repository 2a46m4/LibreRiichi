package core

type Set[T comparable] map[T]struct{}

func NewSet[T comparable]() Set[T] {
	return Set[T]{}
}

func (set Set[T]) Add(obj ...T) {
	for _, o := range obj {
		set[o] = struct{}{}
	}
}

func (set Set[T]) In(obj T) bool {
	_, ok := set[obj]
	return ok
}

func (set Set[T]) Remove(obj T) {
	delete(set, obj)
}

func (set Set[T]) Clear(obj T) {
	set = Set[T]{}
}

