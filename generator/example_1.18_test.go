//go:build go1.18
// +build go1.18

package generator

// SumIntsOrFloats sums the values of map m. It supports both int64 and float64
// as types for map values.
func SumIntsOrFloats[K comparable, V int64 | float64](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

// ChangeType is a type of change detected.
/* ENUM(
	Create
	Update
	Delete
) */
type ChangeType int
