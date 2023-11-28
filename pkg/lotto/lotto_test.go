package lotto

import (
	"fmt"
	"testing"
)

func TestDoubleColorBallCount(t *testing.T) {
	x := DoubleColorBallCount(8, 3)
	fmt.Println("x: ", x)
}
