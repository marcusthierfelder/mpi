package main

import (
	"fmt"
	"github.com/marcusthierfelder/mpi"
)

func main() {

	mpi.Init()

	fmt.Println(mpi.Comm_size(mpi.COMM_WORLD), mpi.Comm_rank(mpi.COMM_WORLD))

	mpi.Finalize()

}
