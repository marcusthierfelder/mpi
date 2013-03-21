package mpi

/*
#include "mpi.h"
#cgo LDFLAGS: -lmpi
#cgo CFLAGS:  -std=gnu99 -Wall

MPI_Comm get_MPI_COMM_WORLD() {
	return (MPI_Comm)(MPI_COMM_WORLD);
}

MPI_Datatype get_MPI_Datatype(int i) {
	if (i==0) 		return (MPI_Datatype)(MPI_INT);
	else if (i==1)	return (MPI_Datatype)(MPI_LONG);
	else if (i==2)  return (MPI_Datatype)(MPI_FLOAT);
	else if (i==3)  return (MPI_Datatype)(MPI_DOUBLE);
	else 			return NULL;
}

MPI_Op get_MPI_Op(int i) {
	if (i==0) 		return (MPI_Op)(MPI_MAX);
	else if (i==1)  return (MPI_Op)(MPI_MIN);
	else if (i==2)  return (MPI_Op)(MPI_SUM);
	else if (i==3)  return (MPI_Op)(MPI_PROD);
	else 			return NULL;
}

*/
import "C"

import (
	_ "fmt"
	"log"

	_"reflect"
	"unsafe"
)

/* 
implementational details:
all #define values within mpi.h cannot be accessed directly, so
go needs c-wrappers. below there is a small subsection of all
values, which have to be extended, if needed.
(mth: I implemented currently only things I need)
*/
var (

	//communication structures
	COMM_WORLD C.MPI_Comm = C.get_MPI_COMM_WORLD()

	//datatypes
	INT     C.MPI_Datatype = C.get_MPI_Datatype(0)
	INT32   C.MPI_Datatype = C.get_MPI_Datatype(0)
	INT64   C.MPI_Datatype = C.get_MPI_Datatype(1)
	FLOAT32 C.MPI_Datatype = C.get_MPI_Datatype(2)
	FLOAT64 C.MPI_Datatype = C.get_MPI_Datatype(3)

	//operations
	MAX  C.MPI_Op = C.get_MPI_Op(0)
	MIN  C.MPI_Op = C.get_MPI_Op(1)
	SUM  C.MPI_Op = C.get_MPI_Op(2)
	PROD C.MPI_Op = C.get_MPI_Op(3)
)

// now mpi has also some types, which we directly map 
type Request C.MPI_Request
type Status C.MPI_Status

func Abort(comm C.MPI_Comm, errorcode int) {
	err := C.MPI_Abort(comm, C.int(errorcode))

	if err != 0 {
		log.Fatal(err)
	}
}

// wrapper for all types (should be used, but this is slower)
func Allreduce(sendbuf, recvbuf interface{}, op C.MPI_Op, comm C.MPI_Comm) {

	if sb, ok := sendbuf.([]int); ok {
		if rv, ok := recvbuf.([]int); ok {
			Allreduce_int(sb, rv, op, comm)
		} else {
			log.Fatal("type of sendbuf and recvbuf do not match")
		}
	} else if sb, ok := sendbuf.([]int32); ok {
		if rv, ok := recvbuf.([]int32); ok {
			Allreduce_int32(sb, rv, op, comm)
		} else {
			log.Fatal("type of sendbuf and recvbuf do not match")
		}
	} else if sb, ok := sendbuf.([]int64); ok {
		if rv, ok := recvbuf.([]int64); ok {
			Allreduce_int64(sb, rv, op, comm)
		} else {
			log.Fatal("type of sendbuf and recvbuf do not match")
		}
	} else 	if sb, ok := sendbuf.([]float32); ok {
		if rv, ok := recvbuf.([]float32); ok {
			Allreduce_float32(sb, rv, op, comm)
		} else {
			log.Fatal("type of sendbuf and recvbuf do not match")
		}
	} else 	if sb, ok := sendbuf.([]float64); ok {
		if rv, ok := recvbuf.([]float64); ok {
			Allreduce_float64(sb, rv, op, comm)
		} else {
			log.Fatal("type of sendbuf and recvbuf do not match")
		}
	} else {
		log.Fatal("type not yet supported")
	}
}

func Allreduce_int(sendbuf, recvbuf []int, op C.MPI_Op, comm C.MPI_Comm) {

	// mpi communication call
	err := C.MPI_Allreduce(
		unsafe.Pointer(&sendbuf[0]), unsafe.Pointer(&recvbuf[0]),
		C.int(len(sendbuf)), INT, op, comm)

	if err != 0 {
		log.Fatal(err)
	}
}

func Allreduce_int32(sendbuf, recvbuf []int32, op C.MPI_Op, comm C.MPI_Comm) {

	// mpi communication call
	err := C.MPI_Allreduce(
		unsafe.Pointer(&sendbuf[0]), unsafe.Pointer(&recvbuf[0]),
		C.int(len(sendbuf)), INT32, op, comm)

	if err != 0 {
		log.Fatal(err)
	}
}

func Allreduce_int64(sendbuf, recvbuf []int64, op C.MPI_Op, comm C.MPI_Comm) {

	// mpi communication call
	err := C.MPI_Allreduce(
		unsafe.Pointer(&sendbuf[0]), unsafe.Pointer(&recvbuf[0]),
		C.int(len(sendbuf)), INT64, op, comm)

	if err != 0 {
		log.Fatal(err)
	}
}

func Allreduce_float32(sendbuf, recvbuf []float32, op C.MPI_Op, comm C.MPI_Comm) {

	// mpi communication call
	err := C.MPI_Allreduce(
		unsafe.Pointer(&sendbuf[0]), unsafe.Pointer(&recvbuf[0]),
		C.int(len(sendbuf)), FLOAT32, op, comm)

	if err != 0 {
		log.Fatal(err)
	}
}

func Allreduce_float64(sendbuf, recvbuf []float64, op C.MPI_Op, comm C.MPI_Comm) {

	// mpi communication call
	err := C.MPI_Allreduce(
		unsafe.Pointer(&sendbuf[0]), unsafe.Pointer(&recvbuf[0]),
		C.int(len(sendbuf)), FLOAT64, op, comm)

	if err != 0 {
		log.Fatal(err)
	}
}

// wrapper for all types (should be used, but this is slower)
func Alltoall(sendbuf, recvbuf interface{}, comm C.MPI_Comm) {
	
	if sb, ok := sendbuf.([]int); ok {
		if rv, ok := recvbuf.([]int); ok {
			Alltoall_int(sb, rv, comm)
		} else {
			log.Fatal("type of sendbuf and recvbuf do not match")
		}
	} else if sb, ok := sendbuf.([]int32); ok {
		if rv, ok := recvbuf.([]int32); ok {
			Alltoall_int32(sb, rv, comm)
		} else {
			log.Fatal("type of sendbuf and recvbuf do not match")
		}
	} else if sb, ok := sendbuf.([]int64); ok {
		if rv, ok := recvbuf.([]int64); ok {
			Alltoall_int64(sb, rv, comm)
		} else {
			log.Fatal("type of sendbuf and recvbuf do not match")
		}
	} else if sb, ok := sendbuf.([]float32); ok {
		if rv, ok := recvbuf.([]float32); ok {
			Alltoall_float32(sb, rv,  comm)
		} else {
			log.Fatal("type of sendbuf and recvbuf do not match")
		}
	} else if sb, ok := sendbuf.([]float64); ok {
		if rv, ok := recvbuf.([]float64); ok {
			Alltoall_float64(sb, rv,  comm)
		} else {
			log.Fatal("type of sendbuf and recvbuf do not match")
		}
	} else {
		log.Fatal("type not yet supported")
	}
}

func Alltoall_int(sendbuf, recvbuf []int, comm C.MPI_Comm) {
	lsend := len(sendbuf)
	lrecv := len(recvbuf)
	size := Comm_size(comm)
	if lsend%size != 0 || lrecv%size != 0 {
		log.Fatal("Alltoall: the bufferlength is not consistent with the number of processors")
	}

	// mpi communication call
	err := C.MPI_Alltoall(
		unsafe.Pointer(&sendbuf[0]), C.int(lsend/size), INT,
		unsafe.Pointer(&recvbuf[0]), C.int(lrecv/size), INT,
		comm)

	if err != 0 {
		log.Fatal(err)
	}
}

func Alltoall_int32(sendbuf, recvbuf []int32, comm C.MPI_Comm) {
	lsend := len(sendbuf)
	lrecv := len(recvbuf)
	size := Comm_size(comm)
	if lsend%size != 0 || lrecv%size != 0 {
		log.Fatal("Alltoall: the bufferlength is not consistent with the number of processors")
	}

	// mpi communication call
	err := C.MPI_Alltoall(
		unsafe.Pointer(&sendbuf[0]), C.int(lsend/size), INT32,
		unsafe.Pointer(&recvbuf[0]), C.int(lrecv/size), INT32,
		comm)

	if err != 0 {
		log.Fatal(err)
	}
}

func Alltoall_int64(sendbuf, recvbuf []int64, comm C.MPI_Comm) {
	lsend := len(sendbuf)
	lrecv := len(recvbuf)
	size := Comm_size(comm)
	if lsend%size != 0 || lrecv%size != 0 {
		log.Fatal("Alltoall: the bufferlength is not consistent with the number of processors")
	}

	// mpi communication call
	err := C.MPI_Alltoall(
		unsafe.Pointer(&sendbuf[0]), C.int(lsend/size), INT64,
		unsafe.Pointer(&recvbuf[0]), C.int(lrecv/size), INT64,
		comm)

	if err != 0 {
		log.Fatal(err)
	}
}

func Alltoall_float32(sendbuf, recvbuf []float32, comm C.MPI_Comm) {
	lsend := len(sendbuf)
	lrecv := len(recvbuf)
	size := Comm_size(comm)
	if lsend%size != 0 || lrecv%size != 0 {
		log.Fatal("Alltoall: the bufferlength is not consistent with the number of processors")
	}

	// mpi communication call
	err := C.MPI_Alltoall(
		unsafe.Pointer(&sendbuf[0]), C.int(lsend/size), FLOAT32,
		unsafe.Pointer(&recvbuf[0]), C.int(lrecv/size), FLOAT32,
		comm)

	if err != 0 {
		log.Fatal(err)
	}
}

func Alltoall_float64(sendbuf, recvbuf []float64, comm C.MPI_Comm) {
	lsend := len(sendbuf)
	lrecv := len(recvbuf)
	size := Comm_size(comm)
	if lsend%size != 0 || lrecv%size != 0 {
		log.Fatal("Alltoall: the bufferlength is not consistent with the number of processors")
	}

	// mpi communication call
	err := C.MPI_Alltoall(
		unsafe.Pointer(&sendbuf[0]), C.int(lsend/size), FLOAT64,
		unsafe.Pointer(&recvbuf[0]), C.int(lrecv/size), FLOAT64,
		comm)

	if err != 0 {
		log.Fatal(err)
	}
}

func Barrier(comm C.MPI_Comm) {
	//
	C.MPI_Barrier(comm)
}

func Comm_size(comm C.MPI_Comm) int {
	n := C.int(-1)
	C.MPI_Comm_size(comm, &n)
	return int(n)
}

func Comm_rank(comm C.MPI_Comm) int {
	n := C.int(-1)
	C.MPI_Comm_rank(comm, &n)
	return int(n)
}

func Init() {
	// initialize the actual mpi stuff, but with two nil
	err := C.MPI_Init(nil, nil)

	if err != 0 {
		log.Fatal(err)
	}
}

func Finalize() {
	err := C.MPI_Finalize()

	if err != 0 {
		log.Fatal(err)
	}
}

// wrapper for all types (should be used, but this is slower)
func Irecv(recvbuf interface{}, source, tag int, comm C.MPI_Comm, request *Request) {
	
	if rb, ok := recvbuf.([]int); ok {		
		Irecv_int(rb, source, tag, comm, request)
	} else if rb, ok := recvbuf.([]int32); ok {		
		Irecv_int32(rb, source, tag, comm, request)
	} else if rb, ok := recvbuf.([]int64); ok {		
		Irecv_int64(rb, source, tag, comm, request)
	} else if rb, ok := recvbuf.([]float32); ok {
		Irecv_float32(rb, source, tag, comm, request)
	} else if rb, ok := recvbuf.([]float64); ok {
		Irecv_float64(rb, source, tag, comm, request)
	} else {
		log.Fatal("type not yet supported")
	}
}

func Irecv_int(recvbuf []int, source, tag int, comm C.MPI_Comm, request *Request) {

	err := C.MPI_Irecv(
		unsafe.Pointer(&recvbuf[0]), C.int(len(recvbuf)), INT,
		C.int(source), C.int(tag), comm, (*C.MPI_Request)(request))

	if err != 0 {
		log.Fatal(err)
	}
}

func Irecv_int32(recvbuf []int32, source, tag int, comm C.MPI_Comm, request *Request) {

	err := C.MPI_Irecv(
		unsafe.Pointer(&recvbuf[0]), C.int(len(recvbuf)), INT32,
		C.int(source), C.int(tag), comm, (*C.MPI_Request)(request))

	if err != 0 {
		log.Fatal(err)
	}
}

func Irecv_int64(recvbuf []int64, source, tag int, comm C.MPI_Comm, request *Request) {

	err := C.MPI_Irecv(
		unsafe.Pointer(&recvbuf[0]), C.int(len(recvbuf)), INT64,
		C.int(source), C.int(tag), comm, (*C.MPI_Request)(request))

	if err != 0 {
		log.Fatal(err)
	}
}

func Irecv_float32(recvbuf []float32, source, tag int, comm C.MPI_Comm, request *Request) {

	err := C.MPI_Irecv(
		unsafe.Pointer(&recvbuf[0]), C.int(len(recvbuf)), FLOAT32,
		C.int(source), C.int(tag), comm, (*C.MPI_Request)(request))

	if err != 0 {
		log.Fatal(err)
	}
}

func Irecv_float64(recvbuf []float64, source, tag int, comm C.MPI_Comm, request *Request) {

	err := C.MPI_Irecv(
		unsafe.Pointer(&recvbuf[0]), C.int(len(recvbuf)), FLOAT64,
		C.int(source), C.int(tag), comm, (*C.MPI_Request)(request))

	if err != 0 {
		log.Fatal(err)
	}
}

// wrapper for all types (should be used, but this is slower)
func Isend(sendbuf interface{}, dest, tag int, comm C.MPI_Comm, request *Request) {
	
	if sb, ok := sendbuf.([]int); ok {		
		Isend_int(sb, dest, tag, comm, request)
	} else if sb, ok := sendbuf.([]int32); ok {		
		Isend_int32(sb, dest, tag, comm, request)
	} else if sb, ok := sendbuf.([]int64); ok {		
		Isend_int64(sb, dest, tag, comm, request)
	} else if sb, ok := sendbuf.([]float32); ok {
		Isend_float32(sb, dest, tag, comm, request)
	} else if sb, ok := sendbuf.([]float64); ok {
		Isend_float64(sb, dest, tag, comm, request)
	} else {
		log.Fatal("type not yet supported")
	}
}

func Isend_int(sendbuf []int, dest, tag int, comm C.MPI_Comm, request *Request) {

	err := C.MPI_Isend(
		unsafe.Pointer(&sendbuf[0]), C.int(len(sendbuf)), INT,
		C.int(dest), C.int(tag), comm, (*C.MPI_Request)(request))

	if err != 0 {
		log.Fatal(err)
	}
}

func Isend_int32(sendbuf []int32, dest, tag int, comm C.MPI_Comm, request *Request) {

	err := C.MPI_Isend(
		unsafe.Pointer(&sendbuf[0]), C.int(len(sendbuf)), INT32,
		C.int(dest), C.int(tag), comm, (*C.MPI_Request)(request))

	if err != 0 {
		log.Fatal(err)
	}
}

func Isend_int64(sendbuf []int64, dest, tag int, comm C.MPI_Comm, request *Request) {

	err := C.MPI_Isend(
		unsafe.Pointer(&sendbuf[0]), C.int(len(sendbuf)), INT64,
		C.int(dest), C.int(tag), comm, (*C.MPI_Request)(request))

	if err != 0 {
		log.Fatal(err)
	}
}

func Isend_float32(sendbuf []float32, dest, tag int, comm C.MPI_Comm, request *Request) {

	err := C.MPI_Isend(
		unsafe.Pointer(&sendbuf[0]), C.int(len(sendbuf)), FLOAT32,
		C.int(dest), C.int(tag), comm, (*C.MPI_Request)(request))

	if err != 0 {
		log.Fatal(err)
	}
}

func Isend_float64(sendbuf []float64, dest, tag int, comm C.MPI_Comm, request *Request) {

	err := C.MPI_Isend(
		unsafe.Pointer(&sendbuf[0]), C.int(len(sendbuf)), FLOAT64,
		C.int(dest), C.int(tag), comm, (*C.MPI_Request)(request))

	if err != 0 {
		log.Fatal(err)
	}
}

// wrapper for all types (should be used, but this is slower)
func Recv(recvbuf interface{}, source, tag int, comm C.MPI_Comm) {
	
	if rb, ok := recvbuf.([]int); ok {		
		Recv_int(rb, source, tag, comm)
	} else if rb, ok := recvbuf.([]int32); ok {		
		Recv_int32(rb, source, tag, comm)
	} else if rb, ok := recvbuf.([]int64); ok {		
		Recv_int64(rb, source, tag, comm)
	} else if rb, ok := recvbuf.([]float32); ok {
		Recv_float32(rb, source, tag, comm)
	} else if rb, ok := recvbuf.([]float64); ok {
		Recv_float64(rb, source, tag, comm)
	} else {
		log.Fatal("type not yet supported")
	}
}

func Recv_int(recvbuf []int, source, tag int, comm C.MPI_Comm) {

	err := C.MPI_Recv(
		unsafe.Pointer(&recvbuf[0]), C.int(len(recvbuf)), INT,
		C.int(source), C.int(tag), comm, nil)

	if err != 0 {
		log.Fatal(err)
	}
}

func Recv_int32(recvbuf []int32, source, tag int, comm C.MPI_Comm) {

	err := C.MPI_Recv(
		unsafe.Pointer(&recvbuf[0]), C.int(len(recvbuf)), INT32,
		C.int(source), C.int(tag), comm, nil)

	if err != 0 {
		log.Fatal(err)
	}
}

func Recv_int64(recvbuf []int64, source, tag int, comm C.MPI_Comm) {

	err := C.MPI_Recv(
		unsafe.Pointer(&recvbuf[0]), C.int(len(recvbuf)), INT64,
		C.int(source), C.int(tag), comm, nil)

	if err != 0 {
		log.Fatal(err)
	}
}

func Recv_float32(recvbuf []float32, source, tag int, comm C.MPI_Comm) {

	err := C.MPI_Recv(
		unsafe.Pointer(&recvbuf[0]), C.int(len(recvbuf)), FLOAT32,
		C.int(source), C.int(tag), comm, nil)

	if err != 0 {
		log.Fatal(err)
	}
}

func Recv_float64(recvbuf []float64, source, tag int, comm C.MPI_Comm) {

	err := C.MPI_Recv(
		unsafe.Pointer(&recvbuf[0]), C.int(len(recvbuf)), FLOAT64,
		C.int(source), C.int(tag), comm, nil)

	if err != 0 {
		log.Fatal(err)
	}
}

// wrapper for all types (should be used, but this is slower)
func Send(sendbuf interface{}, dest, tag int, comm C.MPI_Comm) {
	
	if sb, ok := sendbuf.([]int); ok {		
		Send_int(sb, dest, tag, comm)
	} else if sb, ok := sendbuf.([]int32); ok {		
		Send_int32(sb, dest, tag, comm)
	} else if sb, ok := sendbuf.([]int64); ok {		
		Send_int64(sb, dest, tag, comm)
	} else if sb, ok := sendbuf.([]float32); ok {
		Send_float32(sb, dest, tag, comm)
	} else if sb, ok := sendbuf.([]float64); ok {
		Send_float64(sb, dest, tag, comm)
	} else {
		log.Fatal("type not yet supported")
	}
}

func Send_int(sendbuf []int, dest, tag int, comm C.MPI_Comm) {

	err := C.MPI_Send(
		unsafe.Pointer(&sendbuf[0]), C.int(len(sendbuf)), INT,
		C.int(dest), C.int(tag), comm)

	if err != 0 {
		log.Fatal(err)
	}
}

func Send_int32(sendbuf []int32, dest, tag int, comm C.MPI_Comm) {

	err := C.MPI_Send(
		unsafe.Pointer(&sendbuf[0]), C.int(len(sendbuf)), INT32,
		C.int(dest), C.int(tag), comm)

	if err != 0 {
		log.Fatal(err)
	}
}

func Send_int64(sendbuf []int64, dest, tag int, comm C.MPI_Comm) {

	err := C.MPI_Send(
		unsafe.Pointer(&sendbuf[0]), C.int(len(sendbuf)), INT64,
		C.int(dest), C.int(tag), comm)

	if err != 0 {
		log.Fatal(err)
	}
}

func Send_float32(sendbuf []float32, dest, tag int, comm C.MPI_Comm) {

	err := C.MPI_Send(
		unsafe.Pointer(&sendbuf[0]), C.int(len(sendbuf)), FLOAT32,
		C.int(dest), C.int(tag), comm)

	if err != 0 {
		log.Fatal(err)
	}
}

func Send_float64(sendbuf []float64, dest, tag int, comm C.MPI_Comm) {

	err := C.MPI_Send(
		unsafe.Pointer(&sendbuf[0]), C.int(len(sendbuf)), FLOAT64,
		C.int(dest), C.int(tag), comm)

	if err != 0 {
		log.Fatal(err)
	}
}

func Wait(request *Request, status *Status) {
	err := C.MPI_Wait((*C.MPI_Request)(request), (*C.MPI_Status)(status))

	if err != 0 {
		log.Fatal(err)
	}
}

func Waitall() {

	/*
		int MPI_Waitall(int count, MPI_Request *array_of_requests,
		    MPI_Status *array_of_statuses)
	*/
}
