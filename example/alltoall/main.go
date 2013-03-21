package main

import (
	"fmt"
	"github.com/marcusthierfelder/mpi"
)

func main() {

	mpi.Init()

	chunk := 4 * mpi.Comm_size(mpi.COMM_WORLD)
	rank := mpi.Comm_rank(mpi.COMM_WORLD)
	//size := mpi.Comm_size(mpi.COMM_WORLD)

	sb := make([]int, chunk)
	rb := make([]int, chunk)

	for i := 0; i < chunk; i++ {
		sb[i] = rank + 1
		rb[i] = 0
	}

	//mpi.Alltoall_int(sb, rb, mpi.COMM_WORLD)
	mpi.Alltoall(sb, rb, mpi.COMM_WORLD)

	fmt.Println(sb, rb)

	mpi.Finalize()

}
