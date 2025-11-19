package core

import "golang.org/x/exp/constraints"

// https://stackoverflow.com/questions/64108933/how-to-use-math-pow-with-integers-in-go
func IntPow[T constraints.Unsigned, U constraints.Unsigned](base T, exp U) T {
    result := T(1)
    for {
        if exp & 1 == 1 {
            result *= base
        }
        exp >>= 1
        if exp == 0 {
            break
        }
        base *= base
    }

    return result
}
