mpi
===
mpi-binding package for golang. It is created and tested for openmpi and
should work with all other mpi libraries. Only a small fraction of  methods
are already implemented... more are comming if needed. The names are quite
similar to the mpi-library and they should do the same.

It is tested for golang go1.0.2,go1.0.3 and Open MPI v1.6.2
see

	http://www.open-mpi.org/doc/v1.6/



Quick Usage
===========

	go get github.com/marcusthierfelder/mpi



Detailed Usage (linux)
======================
In order to generate the binding library you need to install mpi on you system.
On ubuntu/debain use:

	sudo apt-get install openmpi-dev openmpi-common

To install this library, cgo needs the location of mpi-header (mpi.h) and the 
mpi-library (mpi.a). Sometimes the system already "knows" these locations. 
If not, you have to find them and export the path. On my system I needed:

	export C_INCLUDE_PATH=/usr/include/openmpi
	export LD_LIBRARY_PATH=/usr/lib/openmpi/lib

On some machines the compiler does not use LD_LIBRARY_PATH, then try:

	export LIBRARY_PATH=/usr/lib/openmpi/lib

To start a parallel job on 4 cores do:
    mpirun -np 4 my_prog



Detailed Usage (mac osx)
========================
You need gcc in order to compile mpi and bind it to golang. Easiest way to get gcc is by
installing xcode and and the command line tools (xcode -> preferences -> downloads). 
Afterwards gcc should work.

To install mpi use this page:

	https://sites.google.com/site/dwhipp/tutorials/installing-open-mpi-on-mac-os-x

or in short, download the newest version of openmpi (currently it is 1.6.5) and 
untar it somewhere. Use the terminal and go into the folder

	./configure
	make
	sudo make install



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
there are only simple boundary and output options. 
does not work properly yet.

note: this is not optimised, but can be used to test clusters for scaling etc.



Errors/Problems
===============

If something is wrong or not working or missing, feel free to contact me 
or post an issue on github.com. 





