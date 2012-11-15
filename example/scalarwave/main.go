package main

import (
	"fmt"
	"github.com/marcusthierfelder/mpi"
)

var (
	testMPI    bool = false
	rank, size int  = 0, 1
	proc0      bool = false
	stdout     string
)

type Grid struct {
	xyz0, dxyz [3]float64 // there are no ghost points included
	nxyz       [3]int     // ... as in here
	gh         int        // number of ghosts
	time       float64    // time ...

	box   Box     // local informations
	field []Field // data storage
}
type Box struct {
	xyz0, xyz1 [3]float64
	dxyz       [3]float64
	nxyz       [3]int
	noff       [3]int

	di, dj, dk  int
	oodx, oodx2 [3]float64

	comm Comm
	grid *Grid
}
type Field struct {
	name string
	sync bool
	data []float64
}
type Comm struct {
	neighbour  [6]int   // number of touching processor
	npts       [6]int   // number of points which have to be syncd
	send, recv [6][]int // stack of position(ijk) to sync efficiently 
}

type VarList struct {
	field []*Field
	grid  *Grid
}

func main() {

	if testMPI == false {
		mpi.Init()
		size = mpi.Comm_size(mpi.COMM_WORLD)
		rank = mpi.Comm_rank(mpi.COMM_WORLD)

		mpi.Redirect_STDOUT(mpi.COMM_WORLD)
	}
	proc0 = rank == 0
	fmt.Println(rank, proc0)

	var grid Grid
	grid.nxyz = [3]int{11, 11, 11}
	grid.dxyz = [3]float64{1, 1, 1}
	grid.xyz0 = [3]float64{-5., -5., -5.}
	grid.gh = 1
	dt := 0.1

	grid.create()
	grid.init()

	vl := grid.initialdata()
	grid.rk4_init(vl)

	grid.output()

	for ti := 0; ti < 1; ti++ {
		if proc0 {
			fmt.Printf("  %4d      %2.4f\n", ti, grid.time)
		}
		grid.rk4(vl, dt)
		grid.time += dt
		//grid.sync_all()

		//grid.output()
	}

	if testMPI == false {
		mpi.Finalize()
	}
}
