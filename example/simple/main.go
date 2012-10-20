package main

import (
	"fmt"
	"mth.com/mpi"
)


func main() {

	mpi.Init()

	fmt.Println(mpi.Comm_size(mpi.COMM_WORLD),mpi.Comm_rank(mpi.COMM_WORLD))

	mpi.Finalize()

}
