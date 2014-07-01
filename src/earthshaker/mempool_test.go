package earthshaker

import (
	"testing"
	"math/rand"
	"sort"
)

type teststruct struct {
	poolindex int
	Value int
}

type bytest []*teststruct

func (a bytest) Len() int           { return len(a) }
func (a bytest) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a bytest) Less(i, j int) bool { return a[i].Value < a[j].Value }

func (t *teststruct) GetPoolIndex() int {
	return t.poolindex
}

func (t *teststruct) SetPoolIndex(n int) {
	t.poolindex = n
}

func TestMempool(t *testing.T) {

	m := Mempool{}
	m.Ini(func() IPool { return new(teststruct) }, 10)

	c := m.New()

	if c == nil {
		t.Errorf("New Error")
	}

	b := m.Delete(c)
	if !b {
		t.Errorf("Delete Error")
	}

	b = m.Delete(c)
	if b {
		t.Errorf("Delete Double Error")
	}

	for i := 0; i < 10; i++ {
		c := m.New()
		if c == nil {
			t.Errorf("New Error")
		}
	}

	c = m.New()
	if c != nil {
		t.Errorf("New Error")
	}

	m.Clear()

	for i := 0; i < 10; i++ {
		c := m.New()
		if c == nil {
			t.Errorf("New Error")
		}
	}

	c = m.New()
	if c != nil {
		t.Errorf("New Error")
	}
}


func TestMempoolData(t *testing.T) {

	m := Mempool{}
	m.Ini(func() IPool { return new(teststruct) }, 10)

	buf := make([]*teststruct, 10)
	for i := 0; i < 10; i++ {
		c := m.New().(*teststruct)
		if c == nil {
			t.Errorf("New Error")
		} else {
			c.Value = i
			buf[i] = c
		}
	}

	for m.UsedSize() > 0 {
		n := rand.Intn(10)
		m.Delete(buf[n])
	}

	newbuf := make([]*teststruct, 10)
	for i := 0; i < 10; i++ {
		c := m.New().(*teststruct)
		if c == nil {
			t.Errorf("New Error")
		} else {
			newbuf[i] = c
		}
	}

	sort.Sort(bytest(newbuf))

	for i := 0; i < 10; i++ {
		if newbuf[i].Value != buf[i].Value {
			t.Error("New Error ", newbuf[i], " ", buf[i])
		}
	}
}
