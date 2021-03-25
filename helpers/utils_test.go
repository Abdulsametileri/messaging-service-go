package helpers

import (
	"github.com/lib/pq"
	"testing"
)

func TestSearchNumberInList(t *testing.T) {
	testCases := []struct {
		number int
		array  pq.Int32Array
		result bool
	}{
		{1, pq.Int32Array{1, 2, 3, 4, 5}, true},
		{3, pq.Int32Array{1, 2, 3, 4, 5}, true},
		{5, pq.Int32Array{1, 2, 3, 4, 5}, true},
		{10, pq.Int32Array{1, 2, 3, 4, 5}, false},
	}

	for _, test := range testCases {
		if result := SearchNumberInPgArray(test.number, test.array); result != test.result {
			t.Errorf("Expected: %v, Want: %v", test.result, result)
		}
	}
}

func TestSortPgArrayAscending(t *testing.T) {
	testCases := []struct {
		array       pq.Int32Array
		sortedArray pq.Int32Array
	}{
		{pq.Int32Array{4, 5, 3, 2, 1}, pq.Int32Array{1, 2, 3, 4, 5}},
		{pq.Int32Array{5, 4, 3, 2, 1}, pq.Int32Array{1, 2, 3, 4, 5}},
		{pq.Int32Array{1, 4, 3, 2, 5}, pq.Int32Array{1, 2, 3, 4, 5}},
		{pq.Int32Array{1, 2, 3, 4, 5}, pq.Int32Array{1, 2, 3, 4, 5}},
	}

	for _, test := range testCases {
		sortedArray := SortPgArrayAscending(test.array)
		for i, number := range sortedArray {
			if number != test.sortedArray[i] {
				t.Errorf("Expected: %v, Want: %v", test.sortedArray, sortedArray)
				break
			}
		}
	}
}

var arr = make(pq.Int32Array, 1000)

func init() {
	for i := 0; i < 1000; i++ {
		arr = append(arr, int32(i))
	}
}

func BenchmarkSearchNumber5InPgArray(b *testing.B) {
	benchmarkSearchNumberInPgArray(b, 5)
}

func BenchmarkSearchNumber100InPgArray(b *testing.B) {
	benchmarkSearchNumberInPgArray(b, 100)
}

func BenchmarkSearchNumber500InPgArray(b *testing.B) {
	benchmarkSearchNumberInPgArray(b, 500)
}

func BenchmarkSearchNumber1000InPgArray(b *testing.B) {
	benchmarkSearchNumberInPgArray(b, 1000)
}

func benchmarkSearchNumberInPgArray(b *testing.B, number int) {
	// run the BenchmarkSearchNumberInPgArray function b.N times
	for n := 0; n < b.N; n++ {
		SearchNumberInPgArray(number, arr)
	}
}
