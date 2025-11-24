package core

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](obj ...T) Set[T] {
	set := Set[T]{}
	set.Add(obj...)
	return set
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

func Reduce[T comparable, E any](set Set[T], init E, fn func(acc E, itr T) E) E {
	for i := range set {
		init = fn(init, i)
	}
	return init
}
