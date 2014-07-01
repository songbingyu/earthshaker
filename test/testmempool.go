package main

import (
	"earthshaker"
)

func main() {
	earthshaker.Ini(earthshaker.IniParam{Name: "testbuffer"})
	testmempool()
	earthshaker.Exit()
}

type teststruct struct {
	poolindex int
}

func (t *teststruct) GetPoolIndex() int {
	return t.poolindex
}

func (t *teststruct) SetPoolIndex(n int) {
	t.poolindex = n
}

func testmempool() {

	m := earthshaker.Mempool{}
	m.Ini(func() earthshaker.IPool { return new(teststruct) }, 10000)

	n := 10000

	buf := make([]earthshaker.IPool, n)

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			c := m.New()
			buf[j] = c
		}
		for j := 0; j < n; j++ {
			m.Delete(buf[j])
		}
	}
}
