package earthshaker

type IPool interface {
	GetPoolIndex() int
	SetPoolIndex(n int)
}

type Mempool struct {
	used     []IPool
	usedsize int
	free     []IPool
	freesize int
	max      int
}

func (m *Mempool) Ini(newfunc func() IPool, max int) {
	m.used = make([]IPool, max)
	m.usedsize = 0
	m.free = make([]IPool, max)
	m.freesize = max
	m.max = max
	for i := 0; i < max; i++ {
		m.free[i] = newfunc()
		m.free[i].SetPoolIndex(0)
	}
}

func (m *Mempool) New() (ret IPool) {
	if m.freesize > 0 {
		ret = m.free[m.freesize-1]
		m.free[m.freesize-1] = nil
		m.freesize--
		m.used[m.usedsize] = ret
		ret.SetPoolIndex(m.usedsize)
		m.usedsize++
		return
	}
	return
}

func (m *Mempool) Delete(p IPool) bool {
	if m.usedsize <= 0 {
		return false
	}
	if p == nil {
		return false
	}
	n := p.GetPoolIndex()
	if n < 0 || n >= m.max {
		return false
	}
	if m.used[n] != p {
		return false
	}

	m.free[m.freesize] = p
	m.freesize++
	p.SetPoolIndex(0)

	if n != m.usedsize-1 {
		m.used[n] = m.used[m.usedsize-1]
		m.used[m.usedsize-1] = nil
		m.usedsize--
		m.used[n].SetPoolIndex(n)
	} else {
		m.used[n] = nil
		m.usedsize--
	}

	return true
}

func (m *Mempool) Clear() {
	for i := 0; i < m.usedsize; i++ {
		p := m.used[i]
		m.used[i] = nil
		m.free[m.freesize] = p
		m.freesize++
	}
	m.usedsize = 0
}

func (m *Mempool) UsedSize() int {
	return m.usedsize
}

func (m *Mempool) FreeSize() int {
	return m.freesize
}
