package fibonacci

import "testing"

func TestFibonacciNumberMatchesDemoPromise(t *testing.T) {
	if got := fibonacciNumber(10); got != 55 {
		t.Fatalf("fibonacciNumber(10) = %d", got)
	}
}
