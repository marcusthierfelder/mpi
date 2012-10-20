package main

import (
	"fmt"
	"mth.com/mpi"
)

func main() {

	mpi.Init()

	chunk := 10
	rank := mpi.Comm_rank(mpi.COMM_WORLD)
	size := mpi.Comm_size(mpi.COMM_WORLD)

	sb := make([]int, chunk)
	rb := make([]int, chunk)

	for i := 0; i < chunk; i++ {
		sb[i] = rank + 1
		rb[i] = 0
	}

    var request,request2 mpi.Request
    var status mpi.Status

    right := (rank + 1) % size
    left := rank - 1
    if left < 0 {
        left = size - 1
    }
    mpi.Irecv_int(&rb, left, 123, mpi.COMM_WORLD, &request);
    mpi.Isend_int(&sb, right, 123, mpi.COMM_WORLD, &request2);
    mpi.Wait(&request, &status);
    mpi.Wait(&request2, &status);

	fmt.Println(sb, rb)

	mpi.Finalize()

}