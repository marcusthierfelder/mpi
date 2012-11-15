package main

import (
	"fmt"
	"log"
	_ "reflect"
)

/* variable storage stuff */
func (grid *Grid) AddVar(name string, sync bool) {
	fmt.Println("AddVar: ", name)

	if false {
		fmt.Println(grid.field, len(grid.field))
	}

	l := len(grid.field)
	tmp := make([]Field, l+1)
	copy(tmp, grid.field)
	grid.field = tmp

	f := Field{name: name, sync: sync,
		data: make([]float64, grid.box.nxyz[0]*grid.box.nxyz[1]*grid.box.nxyz[2])}

	grid.field[l] = f
}

func (grid *Grid) GetVar(name string) []float64 {
	ptr := grid.GetField(name)
	return ptr.data
}

func (grid *Grid) GetField(name string) *Field {
	var ptr *Field

	ptr = nil
	for i, f := range grid.field {
		if f.name == name {
			ptr = &(grid.field[i])
			// this DOES NOT WORK, why?
			// ptr =&f
		}
	}

	if ptr == nil {
		log.Fatal("var \"" + name + "\" does not exist")
	}
	return ptr
}

func (grid *Grid) PrintVars() {
	fmt.Println("Variables:")
	for i, v := range grid.field {
		fmt.Println("  ", i, v.name)
	}
}

/* varlist stuff */
func (grid *Grid) vlalloc() VarList {
	var vl VarList
	vl.grid = grid

	return vl
}

func (vl *VarList) AddVar(name string) {
	l := len(vl.field)
	tmp := make([]*Field, l+1)
	copy(tmp, vl.field)
	vl.field = tmp

	vl.field[l] = vl.grid.GetField(name)
}

func (vl *VarList) GetVar(i int) []float64 {
	if i >= len(vl.field) {
		log.Fatal("GetVar, number out of list")
	}
	return vl.field[i].data
}

func (vl *VarList) PrintVars() {
	fmt.Println("VarList:")
	for i, v := range vl.field {
		fmt.Println("  ", i, v.name)
	}
}
