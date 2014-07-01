package earthshaker

import (
	"math/rand"
	"testing"
)

func TestIni(t *testing.T) {
	b := CircleBuffer{}
	b.Ini(10)
	if b.Capacity() != 10 {
		t.Errorf("Capacity Error")
	}
}

func TestSize(t *testing.T) {
	b := CircleBuffer{}
	b.Ini(10)

	if !b.Empty() {
		t.Errorf("Empty Error")
	}

	b.Write([]byte{1, 2, 3})
	if b.Size() != 3 {
		t.Errorf("Size Error")
	}

	b.Write([]byte{4, 5, 6})
	if b.Size() != 6 {
		t.Errorf("Size Error")
	}

	out := make([]byte, 10)
	b.Read(out[0:4])
	if b.Size() != 2 {
		t.Errorf("Size Error")
	}

	b.Write([]byte{7, 8, 9, 10, 11, 12, 13, 14})
	if b.Size() != 10 {
		t.Errorf("Size Error")
	}

	if !b.Full() {
		t.Errorf("Full Error")
	}

	b.Read(out[0:9])
	if b.Size() != 1 {
		t.Errorf("Size Error")
	}
}

func TestData(t *testing.T) {
	n := 1000000
	srcdata := make([]byte, n)
	desdata := make([]byte, n)
	for i := 0; i < n; i++ {
		srcdata[i] = byte(rand.Intn(255))
	}

	b := CircleBuffer{}
	b.Ini(100)
	srcstep := 0
	desstep := 0
	for srcstep < n || desstep < n {
		srcnum := rand.Intn(50)

		if srcstep+srcnum > n {
			srcnum = n - srcstep
		}

		if b.Write(srcdata[srcstep : srcstep+srcnum]) {
			srcstep += srcnum
		}

		desnum := rand.Intn(30)
		if desnum > b.Size() {
			desnum = b.Size()
		}

		if b.Read(desdata[desstep : desstep+desnum]) {
			desstep += desnum
		}
	}

	if srcstep != n || desstep != n {
		t.Errorf("Step Error, %v, %v, %v", srcstep, desstep, n)
	}

	for i := 0; i < n; i++ {
		if srcdata[i] != desdata[i] {
			t.Errorf("Data Error, %v, %v, %v", i, srcdata[i], desdata[i])
			break
		}
	}
}

func TestInputOutputData(t *testing.T) {
	n := 1000000
	srcdata := make([]byte, n)
	desdata := make([]byte, n)
	for i := 0; i < n; i++ {
		srcdata[i] = byte(i % 255)
	}

	b := CircleBuffer{}
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

			desnum := rand.Intn(30)
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

	for i := 0; i < n; i++ {
		if srcdata[i] != desdata[i] {
			t.Errorf("Data Error, %v, %v, %v", i, srcdata[i], desdata[i])
			break
		}
	}
}

func TestFreeSpace(t *testing.T) {

	b := CircleBuffer{}
	b.Ini(10)

	b.Write([]byte{1, 2, 3})

	if b.FreeSpace() != 7 {
		t.Errorf("FreeSpace Error")
	}

	if b.CanWrite(8) || !b.CanWrite(7) {
		t.Errorf("FreeSpace Error")
	}

	if b.CanRead(4) || !b.CanRead(3) {
		t.Errorf("FreeSpace Error")
	}
}

func TestStore(t *testing.T) {

	b := CircleBuffer{}
	b.Ini(100)

	b.Write([]byte{1, 2, 3})
	oldsize := b.Size()
	b.Store()

	b.Write([]byte{1, 2, 3})
	if oldsize == b.Size() {
		t.Errorf("Store Error")
	}

	b.Restore()
	if oldsize != b.Size() {
		t.Errorf("Store Error")
	}
}
