package common

import "math"

const floatThreshold = 1e-3

// compareFloat compares 2 floats and return true if they are the "same"
func CompareFloat(x, y float64) bool {
	diff := math.Abs(x - y)
	mean := math.Abs(x+y) / 2.0
	if math.IsNaN(diff / mean) {
		return false
	}
	return (diff / mean) > floatThreshold
}
