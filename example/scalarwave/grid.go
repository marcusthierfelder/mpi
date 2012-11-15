package main

import (
	"fmt"
	"github.com/marcusthierfelder/mpi"
	"log"
	"math"

	_ "os"
)

var (
	eps float64 = 1e-10
)

func compGridSize(nijk [3]int, nxyz [3]int) ([][]int, float64) {

	// compute splitting in each direction
	nx := make([]int, nijk[0])
	sum := 0
	for i := 0; i < nijk[0]-1; i++ {
		nx[i] = nxyz[0] / nijk[0]
		sum += nx[i]
	}
	nx[nijk[0]-1] = nxyz[0] - sum

	ny := make([]int, nijk[1])
	sum = 0
	for i := 0; i < nijk[1]-1; i++ {
		ny[i] = nxyz[1] / nijk[1]
		sum += ny[i]
	}
	ny[nijk[1]-1] = nxyz[1] - sum

	nz := make([]int, nijk[2])
	sum = 0
	for i := 0; i < nijk[2]-1; i++ {
		nz[i] = nxyz[2] / nijk[2]
		sum += nz[i]
	}
	nz[nijk[2]-1] = nxyz[2] - sum

	//fmt.Println(nijk,nxyz)
	//fmt.Println(nx,ny,nz)

	// now greate size of each processor
	nprocs := nijk[0] * nijk[1] * nijk[2]
	size := make([][]int, nprocs)
	//fmt.Println(nprocs,size)
	for i := range size {
		size[i] = make([]int, 3)

		size[i][0] = nx[i%nijk[0]]
		size[i][1] = ny[(i/nijk[0])%nijk[1]]
		size[i][2] = nz[(i/(nijk[0]*nijk[1]))%nijk[2]]
		//fmt.Println(i, size)
	}

	// compute volume and surface to choose best ratio
	// -> equil to most efficient parallelization
	V := float64(nxyz[0] * nxyz[1] * nxyz[2])
	A := float64(nxyz[0]*nxyz[1]*(nijk[2]+1) +
		nxyz[0]*nxyz[2]*(nijk[1]+1) +
		nxyz[1]*nxyz[2]*(nijk[0]+1))

	return size, V / A
}

func splitGrid(nxyz [3]int, nproc int) ([][]int, [3]int) {
	var procs [10000][3]int

	l := 0
	lmax := 0
	rmax := 0.
	// test all parallelization versions and test which 
	// splitting is most efficient
	for i := 1; i <= nproc; i++ {
		if nproc%i == 0 {
			for j := 1; j <= nproc; j++ {
				if nproc%(i*j) == 0 || i*j <= nproc {
					for k := 1; k <= nproc; k++ {
						if i*j*k == nproc {
							procs[l][0] = i
							procs[l][1] = j
							procs[l][2] = k

							nijk := [3]int{i, j, k}
							_, r := compGridSize(nijk, nxyz)
							//fmt.Println(g,r)
							if r > rmax {
								rmax = r
								lmax = l
							}

							l++
							if l > 10000 {
								log.Fatal("splitGrid: you use too many processors")
							}
						}
					}
				}
			}
		}
	}

	//fmt.Println(lmax)
	// take the best choice
	nijk := [3]int{procs[lmax][0], procs[lmax][1], procs[lmax][2]}
	size, _ := compGridSize(nijk, nxyz)
	//fmt.Println("")
	//fmt.Println(size,r)

	return size, nijk
}

func (grid *Grid) create() {
	defer un(trace("createGrid"))
	//fmt.Println(rank, size, proc0)

	g, nijk := splitGrid(grid.nxyz, size)
	//fmt.Println(g, nijk)

	// set box sizes for each processor
	grid.box.nxyz[0] = g[rank][0] + 2*grid.gh
	grid.box.dxyz[0] = grid.dxyz[0]
	grid.box.xyz0[0] = grid.xyz0[0]
	for i := 0; i < rank%nijk[0]; i++ {
		grid.box.xyz0[0] += grid.dxyz[0] * float64(g[i][0])
	}
	grid.box.xyz0[0] -= grid.dxyz[0] * float64(grid.gh)
	grid.box.xyz1[0] = grid.box.xyz0[0] + float64(grid.box.nxyz[0]-1)*grid.dxyz[0]

	grid.box.nxyz[1] = g[rank][1] + 2*grid.gh
	grid.box.dxyz[1] = grid.dxyz[1]
	grid.box.xyz0[1] = grid.xyz0[1]
	for i := 0; i < (rank/nijk[0])%nijk[1]; i++ {
		grid.box.xyz0[1] += grid.dxyz[1] * float64(g[i][1])
	}
	grid.box.xyz0[1] -= grid.dxyz[1] * float64(grid.gh)
	grid.box.xyz1[1] = grid.box.xyz0[1] + float64(grid.box.nxyz[1]-1)*grid.dxyz[1]

	grid.box.nxyz[2] = g[rank][2] + 2*grid.gh
	grid.box.dxyz[2] = grid.dxyz[2]
	grid.box.xyz0[2] = grid.xyz0[2]
	for i := 0; i < (rank/(nijk[0]*nijk[1]))%nijk[2]; i++ {
		grid.box.xyz0[2] += grid.dxyz[2] * float64(g[i][2])
	}
	grid.box.xyz0[2] -= grid.dxyz[2] * float64(grid.gh)
	grid.box.xyz1[2] = grid.box.xyz0[2] + float64(grid.box.nxyz[2]-1)*grid.dxyz[2]

	//fmt.Println(grid.box)

	// helpers
	for i := 0; i < 3; i++ {
		grid.box.oodx[i] = 1. / grid.box.dxyz[i]
		grid.box.oodx2[i] = 1. / (grid.box.dxyz[i] * grid.box.dxyz[i])
	}
	grid.box.di = 1
	grid.box.dj = grid.box.nxyz[0]
	grid.box.dk = grid.box.nxyz[0] * grid.box.nxyz[1]

	// find neighbours
	i := rank % nijk[0]
	j := (rank / nijk[0]) % nijk[1]
	k := (rank / (nijk[0] * nijk[1])) % nijk[2]
	if i > 0 {
		grid.box.comm.neighbour[0] = (i - 1) + j*nijk[0] + k*nijk[0]*nijk[1]
	} else {
		grid.box.comm.neighbour[0] = -1
	}
	if i < nijk[0]-1 {
		grid.box.comm.neighbour[1] = (i + 1) + j*nijk[0] + k*nijk[0]*nijk[1]
	} else {
		grid.box.comm.neighbour[1] = -1
	}
	if j > 0 {
		grid.box.comm.neighbour[2] = i + (j-1)*nijk[0] + k*nijk[0]*nijk[1]
	} else {
		grid.box.comm.neighbour[2] = -1
	}
	if j < nijk[1]-1 {
		grid.box.comm.neighbour[3] = i + (j+1)*nijk[0] + k*nijk[0]*nijk[1]
	} else {
		grid.box.comm.neighbour[3] = -1
	}
	if k > 0 {
		grid.box.comm.neighbour[4] = i + j*nijk[0] + (k-1)*nijk[0]*nijk[1]
	} else {
		grid.box.comm.neighbour[4] = -1
	}
	if k < nijk[2]-1 {
		grid.box.comm.neighbour[5] = i + j*nijk[0] + (k+1)*nijk[0]*nijk[1]
	} else {
		grid.box.comm.neighbour[5] = -1
	}

	// set communication buffers
	// for each of the 6 sides, go through gh points and store ijk in a 
	// buffer and send it to the neighbour processor - if there are one
	grid.box.comm.npts[0] = grid.gh * grid.box.nxyz[1] * grid.box.nxyz[2]
	grid.box.comm.npts[1] = grid.gh * grid.box.nxyz[1] * grid.box.nxyz[2]
	grid.box.comm.npts[2] = grid.gh * grid.box.nxyz[0] * grid.box.nxyz[2]
	grid.box.comm.npts[3] = grid.gh * grid.box.nxyz[0] * grid.box.nxyz[2]
	grid.box.comm.npts[4] = grid.gh * grid.box.nxyz[0] * grid.box.nxyz[1]
	grid.box.comm.npts[5] = grid.gh * grid.box.nxyz[0] * grid.box.nxyz[1]
	for i := 0; i < 6; i++ {
		grid.box.comm.send[i] = make([]int, grid.box.comm.npts[i])
		grid.box.comm.recv[i] = make([]int, grid.box.comm.npts[i])
	}

	ijk0 := [12]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	ijk := 0
	for k := 0; k < grid.box.nxyz[2]; k++ {
		for j := 0; j < grid.box.nxyz[1]; j++ {
			for i := 0; i < grid.box.nxyz[0]; i++ {

				// x dir
				if i < grid.gh {
					grid.box.comm.send[0][ijk0[0]] = ijk
					ijk0[0]++
				}
				if i >= grid.gh && i < 2*grid.gh {
					grid.box.comm.recv[0][ijk0[6]] = ijk
					ijk0[6]++
				}
				if i >= grid.box.nxyz[0]-grid.gh {
					grid.box.comm.send[1][ijk0[1]] = ijk
					ijk0[1]++
				}
				if i < grid.box.nxyz[0]-grid.gh && i >= grid.box.nxyz[0]-2*grid.gh {
					grid.box.comm.recv[1][ijk0[7]] = ijk
					ijk0[7]++
				}

				// y dir
				if j < grid.gh {
					grid.box.comm.send[2][ijk0[2]] = ijk
					ijk0[2]++
				}
				if j >= grid.gh && j < 2*grid.gh {
					grid.box.comm.recv[2][ijk0[8]] = ijk
					ijk0[8]++
				}
				if j >= grid.box.nxyz[1]-grid.gh {
					grid.box.comm.send[3][ijk0[3]] = ijk
					ijk0[3]++
				}
				if j < grid.box.nxyz[1]-grid.gh && j >= grid.box.nxyz[1]-2*grid.gh {
					grid.box.comm.recv[3][ijk0[9]] = ijk
					ijk0[9]++
				}

				// z dir
				if k < grid.gh {
					grid.box.comm.send[4][ijk0[4]] = ijk
					ijk0[4]++
				}
				if k >= grid.gh && k < 2*grid.gh {
					grid.box.comm.recv[4][ijk0[10]] = ijk
					ijk0[10]++
				}
				if k >= grid.box.nxyz[2]-grid.gh {
					grid.box.comm.send[5][ijk0[5]] = ijk
					ijk0[5]++
				}
				if k < grid.box.nxyz[2]-grid.gh && k >= grid.box.nxyz[2]-2*grid.gh {
					grid.box.comm.recv[5][ijk0[11]] = ijk
					ijk0[11]++
				}

				ijk++
			}
		}
	}

	if false {
		fmt.Println(grid.box)
	}
}

func (grid *Grid) init() {
	defer un(trace("initGrid"))

	grid.AddVar("x", false)
	grid.AddVar("y", false)
	grid.AddVar("z", false)

	xp := grid.GetVar("x")
	yp := grid.GetVar("y")
	zp := grid.GetVar("z")

	ijk := 0
	for k := 0; k < grid.box.nxyz[2]; k++ {
		for j := 0; j < grid.box.nxyz[1]; j++ {
			for i := 0; i < grid.box.nxyz[0]; i++ {
				xp[ijk] = grid.box.xyz0[0] + float64(i)*grid.box.dxyz[0]
				yp[ijk] = grid.box.xyz0[1] + float64(j)*grid.box.dxyz[1]
				zp[ijk] = grid.box.xyz0[2] + float64(k)*grid.box.dxyz[2]

				ijk++
			}
		}
	}
}

/* pure mpi synchronization */
func (grid *Grid) sync_all() {
	grid.sync_one(grid.GetVar("f"))
}

func (grid *Grid) sync_vl(vl VarList) {
	for _, v := range vl.field {
		grid.sync_one(v.data)
	}
}

func (grid *Grid) sync_one(data []float64) {
	c := grid.box.comm
	mpi.Barrier(mpi.COMM_WORLD)

	/* go through all 6 sides seperatly, first x-dir both sides,
	then the other sides */
	for d := 0; d < 6; d++ {
		e := (d/2)*2 + 1 - d%2

		// fist one direction
		sendbuf := make([]float64, c.npts[d])
		recvbuf := make([]float64, c.npts[e])
		if c.neighbour[d] != -1 {
			for i := 0; i < c.npts[d]; i++ {
				sendbuf[i] = data[c.send[d][i]]
			}
		}

		if c.neighbour[e] != -1 {
			mpi.Recv_float64(recvbuf, c.neighbour[e], 123, mpi.COMM_WORLD)
			//mpi.Wait(&request1, &status)
			fmt.Println("recv---", rank, c.neighbour[e], recvbuf)
		}
		if c.neighbour[d] != -1 {
			fmt.Println("send---", rank, c.neighbour[d], sendbuf)
			mpi.Send_float64(sendbuf, c.neighbour[d], 123, mpi.COMM_WORLD)
			//mpi.Wait(&request2, &status)
		}

		if c.neighbour[e] != -1 {
			for i := 0; i < c.npts[e]; i++ {
				data[c.recv[e][i]] = recvbuf[i]
			}
		}

		mpi.Barrier(mpi.COMM_WORLD)

	}
}

/* box stuff */
func (box *Box) inside(pos [3]float64) bool {
	if pos[0] < box.xyz0[0]-eps ||
		pos[1] < box.xyz0[1]-eps ||
		pos[2] < box.xyz0[2]-eps ||
		pos[0] > box.xyz0[0]+float64(box.nxyz[0]-1)*box.dxyz[0]+eps ||
		pos[1] > box.xyz0[1]+float64(box.nxyz[1]-1)*box.dxyz[1]+eps ||
		pos[2] > box.xyz0[2]+float64(box.nxyz[2]-1)*box.dxyz[2]+eps {
		return false
	}
	return true
}

func (box *Box) interpolate(pos [3]float64, data []float64) (float64, bool) {

	fmt.Println(box.xyz0, box.xyz1, pos)

	// test if point is inside box
	if box.inside(pos) == false {
		return 0., false
	}

	// find interpolation offset for 2nd order interpolation
	var n [3]int
	var x0 [3]float64
	for i := 0; i < 3; i++ {
		n[i] = int(math.Trunc((pos[i] - box.xyz0[i]) / box.dxyz[i]))
		if n[i] == box.nxyz[i] {
			n[i]--
			//return 0., false
		}
		x0[i] = box.xyz0[i] + float64(n[i])*box.dxyz[i]
	}
	//fmt.Println("1 ",x0)

	// find values of the cube around the point
	v := [][][]float64{{{0, 0}, {0, 0}}, {{0, 0}, {0, 0}}}
	for k := 0; k < 2; k++ {
		for j := 0; j < 2; j++ {
			for i := 0; i < 2; i++ {
				ijk := (n[0] + i) + (n[1]+j)*box.nxyz[0] + (n[2]+k)*box.nxyz[0]*box.nxyz[1]
				v[i][j][k] = data[ijk]
			}
		}
	}
	//fmt.Println(v)

	// now interpolte tri-polynomial 2nd order
	interp := interpolate_TriN(pos, x0, box.dxyz, v)

	return interp, true
}
