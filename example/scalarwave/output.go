package main

import (
	"fmt"
	"github.com/marcusthierfelder/mpi"
	"io"
	"log"
	"os"
	_ "reflect"
)

func (grid *Grid) output() {
	defer un(trace("output"))

	grid.output_1d("f", "f.x", 0)
	grid.output_1d("g", "g.x", 0)
}

/* this a a very slow and lazy implementation for x-direction */
func (grid *Grid) output_1d(data string, file string, d int) {

	/* init buffer */
	buffer := make([]float64, 2*grid.nxyz[d])
	x := make([]float64, grid.nxyz[d])
	ptr := grid.GetVar(data)
	pos := [3]float64{0., 0., 0.}

	for i := 0; i < grid.nxyz[d]; i++ {
		x[i] = grid.xyz0[d] + float64(i)*grid.dxyz[d]
		pos[d] = x[i]
		v, b := grid.box.interpolate(pos, ptr)
		buffer[i], buffer[grid.nxyz[d]+i] = v, float64(btoi(b))
	}
	//fmt.Println(rank, x,buffer[:grid.nxyz[d]])

	/* find data using interpolation */
	recvbuf := make([]float64, 2*grid.nxyz[d])
	mpi.Allreduce_float64(buffer, recvbuf, mpi.SUM, mpi.COMM_WORLD)
	for i := 0; i < grid.nxyz[d]; i++ {
		buffer[i] = recvbuf[i] / recvbuf[i+grid.nxyz[d]]
	}
	//fmt.Println(recvbuf)

	/* create or append? */
	if proc0 {
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

		fmt.Println("write " + file)
		io.WriteString(f, fmt.Sprintf("#Time = %e\n", grid.time))
		for i := 0; i < grid.nxyz[d]; i++ {
			io.WriteString(f, fmt.Sprintf("%e %e\n", x[i], buffer[i]))
		}
		io.WriteString(f, "\n")

		f.Close()
	}
}
