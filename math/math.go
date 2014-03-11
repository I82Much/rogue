package math

/**
 * @param value The incoming value to be converted
 * @param low1  Lower bound of the value's current range
 * @param high1 Upper bound of the value's current range
 * @param low2  Lower bound of the value's target range
 * @param high2 Upper bound of the value's target range
 */
func DoMap(value, low1, high1, low2, high2 float64) float64 {
	diff := value - low1
	proportion := diff / (high1 - low1)
	return Lerp(low2, high2, proportion)
}

func IntMap(value, low1, high1, low2, high2 int) int {
	return int(DoMap(float64(value), float64(low1), float64(high1), float64(low2), float64(high2)))
}

// Linearly interpolate between two values
func Lerp(value1, value2, amt float64) float64 {
	return ((value2 - value1) * amt) + value1
}
