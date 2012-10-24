package main

import (
	"fmt"
	_ "os"
)

var pr = true

/* some helpers to see which function is called */
func trace(s string) string {
	if pr == true {
		fmt.Println("entering:", s)
	} else {
		fmt.Print("")
	}
	return s
}

func un(s string) {
	if pr == true {
		fmt.Println("leaving: ", s)
	} else {
		fmt.Print("")
	}
}

/* lagrange interpolation functions */
func coefficients_lagrange_N(x, xmin, h float64, c []float64) {
	n := len(c)

	for i := 0; i < n; i++ {
		d := 1.
		for j := 0; j < n; j++ {
			if j != i {
				d *= (x - xmin - float64(j)*h) / (float64(i-j) * h)
			}
		}
		c[i] = d
	}
}

func interpolate_lagrange_N(x, xmin, h float64, c, u []float64) float64 {
	n := len(c)
	sum := 0.

	for i := 0; i < n; i++ {
		sum += c[i] * u[i]
	}

	return sum
}

func interpolate_TriN(x, xmin, dx [3]float64, u [][][]float64) float64 {

	//fmt.Println("-----")
	//fmt.Println(x, xmin, dx)

	n := len(u)
	sum := 3.141
	c := make([]float64, n)
	v := make([][]float64, n)
	for i := 0; i < n; i++ {
		v[i] = make([]float64, n)
	}
	w := make([]float64, n)

	//fmt.Println(u)

	coefficients_lagrange_N(x[2], xmin[2], dx[2], c)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			v[i][j] = interpolate_lagrange_N(x[2], xmin[2], dx[2], c, u[i][j])
		}
	}
	//fmt.Println(c,v)

	coefficients_lagrange_N(x[1], xmin[1], dx[1], c)
	for i := 0; i < n; i++ {
		w[i] = interpolate_lagrange_N(x[1], xmin[1], dx[1], c, v[i])
	}
	//fmt.Println(c,w)

	coefficients_lagrange_N(x[0], xmin[0], dx[0], c)
	sum = interpolate_lagrange_N(x[0], xmin[0], dx[0], c, w)

	//fmt.Println(c,sum)
	//os.Exit(1)

	return sum
}

/* some other stuff */
func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}
