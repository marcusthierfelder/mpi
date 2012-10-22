mpi
===
mpi-binding package for golang. It is created and tested for openmpi and
should work with all other mpi libraries.

It is tested for golang go1.0.2 and Open MPI v1.6.2
see

	http://www.open-mpi.org/doc/v1.6/



Quick Usage
===========

	go get github.com/marcusthierfelder/mpi




Detailed Usage
==============
In order to generate the binding library you need to install mpi on you system.
On ubuntu/debain use

	sudo apt-get install openmpi-dev openmpi-common

To compile the library cgo needs the location of mpi.h header and the moi library
files. You have to find it and export the path. On my system I needed

	export C_INCLUDE_PATH=/usr/include/openmpi
	export LD_LIBRARY_PATH=/usr/lib/openmpi/lib

Sometimes you don't need this since the system already "knows" these locations.



Examples
========

### simple:

simple example which writes the rank of each processor and the 
total number of mpi-jobs

### alltoall:

simple example which communicates a small number of integer between all
processors

### nonblocking:
	
simple send recv example

### scalarwave:

fancy scalarwave example in 3d with goroutines and mpi decomposition.
you can choose between several integrators and orders of finite differencing.
there are only simple boundary and output options 

note: this is not a fancy optimised code, but can be used to test clusters
for scaling etc.








