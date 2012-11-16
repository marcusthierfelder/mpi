package mpi

/*
#include <stdio.h>
#include <stdlib.h>
#include "mpi.h"
#cgo LDFLAGS: -lmpi
#cgo CFLAGS:  -std=gnu99 -Wall

void call_freopen(char *file) {
	freopen(file, "w", stdout);
	freopen(file, "w", stderr);
}

*/
import "C"

import (
	"fmt"
)


// if you want to redirect each stdout per process in a file 
// in order to avoid a parallel stdout mess (no standard mpi function)
func Redirect_STDOUT(comm C.MPI_Comm) {
	rank := Comm_rank(comm)
	file := fmt.Sprintf("stdout.%04d", rank)

	if rank != 0 {
		C.call_freopen(C.CString(file))
	}

}
