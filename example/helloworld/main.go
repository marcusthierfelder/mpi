package main

import (
        "fmt"
        "github.com/marcusthierfelder/mpi"
)

func main() {

        mpi.Init()
        var rank int = mpi.Comm_rank(mpi.COMM_WORLD)
        var size int = mpi.Comm_rank(mpi.COMM_WORLD)
        fmt.Printf("Hello world from rank %d of %d\n",rank,size)
        mpi.Finalize()
}

