package main

import (
	"fmt"
	"log"
	_ "reflect"
)

/* storage stuff */
func (grid *Grid) AddVar(name string, sync bool) {

	if false {
		fmt.Println(grid.field, len(grid.field))
	}

	l := len(grid.field)
	tmp := make([](*Field), l+1)
	copy(tmp, grid.field)
	grid.field = tmp

	f := Field{name: name, sync: sync,
		data: make([]float64, grid.nxyz[0]*grid.nxyz[1]*grid.nxyz[2])}

	grid.field[l] = &f
}

func (grid *Grid) GetVar(name string) *[]float64 {

	var ptr *[]float64
	for _, f := range grid.field {
		if f.name == name {
			ptr = &(f.data)
		}
	}

	if ptr == nil {
		log.Fatal("var \"" + name + "\" does not exist")
	}
	return ptr
}
