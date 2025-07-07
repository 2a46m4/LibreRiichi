package core

import "math/rand"

// Creates a random permutation of the array
// Modifies the existing array
func PermuteArray[T any](array []T) []T {
	for i := len(array); i > 0; i -= 1 {
		rand := rand.Intn(i)
		temp := array[i]
		array[i] = array[rand]
		array[rand] = temp
	}
	return array
}

// Rotates the array to the left, e.g. [1, 2, 3] becomes [2, 3, 1]
func RotateArrayLeft[T any](array []T, by int) []T {
	temp := make([]T, by)
	copy(temp, array[:by])
	copy(array, array[by:])
	copy(array[len(array)-by:], temp)
	return array
}

func Last[T any](array []T) T {
	return array[len(array)-1]
}

func LastPtr[T any](array []T) *T {
	return &array[len(array)-1]
}

// Always modifies the array
func Pop[T any](array *[]T) T {
	*array = (*array)[:len(*array)-1]
	return Last(*array)
}

func Swap[T any](array []T, first, second uint) []T {
	temp := array[first]
	array[first] = array[second]
	array[second] = temp
	return array
}

// Does not preserve order
func Remove[T any, U int | uint](array *[]T, idx U) *[]T {
	Swap(*array, uint(idx), uint(len(*array)-1))
	Pop(array)
	return array
}

func Count[T comparable](array []T, find T) (count int) {
	for _, i := range array {
		if i == find {
			count += 1
		}
	}
	return count
}
