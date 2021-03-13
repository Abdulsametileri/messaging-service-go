package helpers

import (
	"github.com/lib/pq"
	"sort"
)

func SearchNumberInPgArray(number int, array pq.Int32Array) bool {
	numberAsInt32 := int32(number)
	n := len(array)

	i := sort.Search(n, func(i int) bool { return array[i] >= numberAsInt32 })

	if i < n && array[i] == numberAsInt32 {
		return true
	}
	return false
}

func SortPgArrayAscending(numbers pq.Int32Array) pq.Int32Array {
	sort.Slice(numbers, func(i, j int) bool { return numbers[i] < numbers[j] })
	return numbers
}
