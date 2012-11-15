package main

import (
	_ "fmt"
	_ "math"
)

func (box *Box) innerpoint(i, j, k int) bool {
	n := 1
	if i > n && j > n && k > n && i < box.nxyz[0]-n && j < box.nxyz[1]-n && k < box.nxyz[2]-n {
		return true
	}

	return false
}
