package common

// KeysFromItems returns the int64 keys from a response Items field.
func KeysFromItems[T any](items *[]T, getKey func(T) int64) []int64 {
	if items == nil {
		return nil
	}
	out := make([]int64, 0, len(*items))
	for _, it := range *items {
		out = append(out, getKey(it))
	}
	return out
}
