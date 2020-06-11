package functions

// SetToSliceInt64 transforms map int64 to slice
func SetToSliceInt64(v map[int64]struct{}) []int64 {
	result := make([]int64, 0, len(v))
	for i := range v {
		result = append(result, i)
	}

	return result
}
