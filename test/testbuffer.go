package main

import (
	"earthshaker"
	"math/rand"
)

func main() {
	earthshaker.Ini(earthshaker.IniParam{Name: "testbuffer"})
	testbuffer()
	earthshaker.Exit()
}

func testbuffer() {
	n := 100000000
	srcdata := make([]byte, n)
	desdata := make([]byte, n)
	for i := 0; i < n; i++ {
		srcdata[i] = byte(i % 255)
	}

	b := earthshaker.CircleBuffer{}
	b.Ini(100)

	srcstep := 0
	desstep := 0
	for srcstep < n || desstep < n {
		b.Input(func(data []byte) (int, error) {
			size := len(data)

			srcnum := rand.Intn(50)
			if srcstep+srcnum > n {
				srcnum = n - srcstep
			}

			if srcnum > size {
				srcnum = size
			}

			copy(data, srcdata[srcstep:srcstep+srcnum])

			srcstep += srcnum

			return srcnum, nil
		})

		b.Output(func(data []byte) (int, error) {
			size := len(data)

			desnum := rand.Intn(50)
			if desstep+desnum > n {
				desnum = n - desstep
			}

			if desnum > size {
				desnum = size
			}

			copy(desdata[desstep:desstep+desnum], data)

			desstep += desnum

			return desnum, nil
		})
	}
}
