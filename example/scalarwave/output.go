package main

import (
	"io"
	"fmt"
	"log"
	"os"
)


func (grid *Grid) output() {
	defer un(trace("output"))

	grid.output_1d("x", "f.x", 0)
}



/* this a a very slow and lazy implementation for x-direction */
func (grid *Grid) output_1d(data string, file string, d int) {

	/* init buffer */
	var buffer [2][]float64 
	buffer[0] = make([]float64,grid.nxyz[d])
	buffer[1] = make([]float64,grid.nxyz[d])

	for i:=0; i<grid.nxyz[d]; i++ {
		buffer[0][i] = grid.xyz0[d] + float64(i)*grid.dxyz[d]
	}

	/* find data using interpolation */


	/* create or append? */
	var flag int
	if grid.time == 0. {
		flag = os.O_CREATE | os.O_TRUNC | os.O_RDWR
	} else {
		flag = os.O_APPEND | os.O_RDWR
	}

	f, err := os.OpenFile(file, flag, 0666)
	if err != nil {
		log.Fatal(err)
		return
	}

	io.WriteString(f, fmt.Sprintf("#Time = %e\n", grid.time))
	for i := 0; i < grid.nxyz[d]; i++ {
		io.WriteString(f, fmt.Sprintf("%e %e\n", buffer[0][i], buffer[1][i]))
	}
	io.WriteString(f, "\n")

	f.Close()
}