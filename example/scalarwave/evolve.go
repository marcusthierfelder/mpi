package main

import (
	_ "fmt"
	"math"
)

var (
	tmp_vl []VarList
)

/* intial data */
func (grid *Grid) initialdata() VarList {
	defer un(trace("initaldata"))
	grid.AddVar("f", true)
	grid.AddVar("g", true)

	vl := grid.vlalloc()
	vl.AddVar("f")
	vl.AddVar("g")

	box := grid.box
	f := grid.GetVar("f")
	g := grid.GetVar("g")
	x := grid.GetVar("x")
	y := grid.GetVar("y")
	z := grid.GetVar("z")

	ijk := 0
	for k := 0; k < box.nxyz[2]; k++ {
		for j := 0; j < box.nxyz[1]; j++ {
			for i := 0; i < box.nxyz[0]; i++ {
				r := math.Sqrt(x[ijk]*x[ijk] + y[ijk]*y[ijk] + z[ijk]*z[ijk])

				f[ijk] = math.Exp(-r * r)
				g[ijk] = 0.

				ijk++
			}
		}
	}

	return vl
}

/* rhs computation */
func (grid *Grid) rhs(r VarList, uc VarList) {

	box := grid.box
	f := uc.GetVar(0)
	g := uc.GetVar(1)
	rf := r.GetVar(0)
	rg := r.GetVar(1)

	di := box.di
	dj := box.dj
	dk := box.dk

	ijk := 0
	for k := 0; k < box.nxyz[2]; k++ {
		for j := 0; j < box.nxyz[1]; j++ {
			for i := 0; i < box.nxyz[0]; i++ {

				if box.innerpoint(i, j, k) {
					// second order laplacian
					laplace := box.oodx2[0] * (-6.*f[ijk] +
						f[ijk-di] + f[ijk+di] + f[ijk-dj] + f[ijk+dj] + f[ijk-dk] + f[ijk+dk])

					rf[ijk] = g[ijk]
					rg[ijk] = laplace

				} else {
					// simple boundary condition
					rf[ijk] = 0.
					rg[ijk] = 0.

				}
				ijk++
			}
		}
	}

	grid.sync_vl(r)
}

/* helper functions for the time integrator */
func (grid *Grid) cpy(vl1 VarList, vl2 VarList) {
	ntot := grid.box.nxyz[0] * grid.box.nxyz[1] * grid.box.nxyz[2]

	for i, v1 := range vl1.field {
		v2 := vl2.field[i]

		for ijk := 0; ijk < ntot; ijk++ {
			v1.data[ijk] = v2.data[ijk]
		}
	}
}

func (grid *Grid) addto(vl1 VarList, c float64, vl2 VarList) {
	ntot := grid.box.nxyz[0] * grid.box.nxyz[1] * grid.box.nxyz[2]

	for i, v1 := range vl1.field {
		v2 := vl2.field[i]

		for ijk := 0; ijk < ntot; ijk++ {
			v1.data[ijk] += c * v2.data[ijk]
		}
	}
}

func (grid *Grid) add(vl1 VarList, c2 float64, vl2 VarList, c3 float64, vl3 VarList) {
	ntot := grid.box.nxyz[0] * grid.box.nxyz[1] * grid.box.nxyz[2]

	for i, v1 := range vl1.field {
		v2 := vl2.field[i]
		v3 := vl3.field[i]

		for ijk := 0; ijk < ntot; ijk++ {
			v1.data[ijk] = c2*v2.data[ijk] + c3*v3.data[ijk]

		}
	}
}

func (grid *Grid) rk4_init(uc VarList) {
	defer un(trace("rk4_init"))

	tmp_vl = []VarList{grid.vlalloc(), grid.vlalloc(), grid.vlalloc()}

	for _, v := range uc.field {
		grid.AddVar(v.name+"_p", v.sync)
		tmp_vl[0].AddVar(v.name + "_p")
		grid.AddVar(v.name+"_r", v.sync)
		tmp_vl[1].AddVar(v.name + "_r")
		grid.AddVar(v.name+"_v", v.sync)
		tmp_vl[2].AddVar(v.name + "_v")
	}
}

func (grid *Grid) rk4(uc VarList, dt float64) {

	p := tmp_vl[0]
	r := tmp_vl[1]
	v := tmp_vl[2]

	grid.cpy(p, uc)

	grid.rhs(r, uc)
	grid.addto(uc, dt/6., r)

	grid.add(v, 1.0, p, dt/2., r)
	grid.rhs(r, v)
	grid.addto(uc, dt/3., r)

	grid.add(v, 1.0, p, dt/2., r)
	grid.rhs(r, v)
	grid.addto(uc, dt/3., r)

	grid.add(v, 1.0, p, dt, r)
	grid.rhs(r, v)
	grid.addto(uc, dt/6., r)

}
