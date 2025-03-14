package common_service

import "sort"

// Function to insert a value into a sorted integer slice
func InsertSorted(arr []int, value int) []int {
	// Use binary search to find the insertion point
	index := sort.Search(len(arr), func(i int) bool {
		return arr[i] >= value
	})

	// Insert the value at the found index
	arr = append(
		arr[:index], // Slice upto but not including [index]
		append(
			[]int{value},   // Slice with the value to insert
			arr[index:]..., // Slice from [index] to the end
		)..., // Spread
	)

	return arr
}
