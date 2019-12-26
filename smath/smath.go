package smath

func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func LCM(a, b int, values ...int) int {
	result := (a / GCD(a, b)) * b
	for _, x := range values {
		result = LCM(result, x)
	}
	return result
}

func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}