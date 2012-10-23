package main

import (
	"fmt"
	"testing"
)


func TestInterpolate(t *testing.T) {

	u := [][][]float64{{{0,0},{0,0}},{{1,1},{1,1}}}
x := [3]float64{1,1,1}
dx := [3]float64{2,2,2}
x0 := [3]float64{0,0,0}

	fmt.Println(u)

	fmt.Println(interpolate_TriN(x,x0,dx, u) )

	t.Errorf("---")

}