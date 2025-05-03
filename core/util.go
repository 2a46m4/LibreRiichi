package core

import "math/rand"

type UnitType struct{}

var Unit = UnitType{}

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
