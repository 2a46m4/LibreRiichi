package core

import (
	"iter"
)

func All[T any](seq iter.Seq[T], f func(T) bool) (res bool) {
	for i := range seq {
		if !f(i) {
			return false
		}
	}
	return true
}


func Any[T any](seq iter.Seq[T], f func(T) bool) bool {
	for i := range seq {
		if f(i) {
			return true
		}
	}
	return false
}
