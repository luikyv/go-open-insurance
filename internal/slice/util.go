package slice

import "slices"

func ContainsAll[T comparable](superSet []T, subSet ...T) bool {
	for _, t := range subSet {
		if !slices.Contains(superSet, t) {
			return false
		}
	}

	return true
}

func ContainsAny[T comparable](superSet []T, subSet ...T) bool {
	for _, t := range subSet {
		if slices.Contains(superSet, t) {
			return true
		}
	}

	return false
}

func FindFirst[T any](elements []T, condition func(t T) bool) (T, bool) {
	for _, e := range elements {
		if condition(e) {
			return e, true
		}
	}

	return *new(T), false
}
