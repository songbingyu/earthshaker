package earthshaker

import (
	"unsafe"
)

type IPacket interface {
	Serialize(b *CircleBuffer) bool
	Deserialize(b *CircleBuffer) bool
}

type Packet struct {
	head int16
	flg  byte
	size int16
	data []byte
}

const PacketHeaderSize int = 2 + 1 + 2

func (p *Packet) Serialize(b *CircleBuffer) bool {

	sesize := PacketHeaderSize + int(p.size)

	if !b.CanWrite(sesize) {
		return false
	}

	b.Write((*(*[2]byte)(unsafe.Pointer(&p.head)))[0:])
	b.Write((*(*[1]byte)(unsafe.Pointer(&p.flg)))[0:])
	b.Write((*(*[2]byte)(unsafe.Pointer(&p.size)))[0:])
	b.Write(p.data[0:p.size])

	return true
}

func (p *Packet) Deserialize(b *CircleBuffer) bool {

	b.Store()

	if !b.CanRead(PacketHeaderSize) {
		b.Restore()
		return false
	}

	b.Read((*(*[2]byte)(unsafe.Pointer(&p.head)))[0:])
	b.Read((*(*[1]byte)(unsafe.Pointer(&p.flg)))[0:])
	b.Read((*(*[2]byte)(unsafe.Pointer(&p.size)))[0:])

	if !b.CanRead(int(p.size)) {
		b.Restore()
		return false
	}

	p.data = make([]byte, p.size)
	b.Read(p.data[0:p.size])

	return true
}
