package earthshaker

import (
	"testing"
	"unsafe"
)

func TestPacket(t *testing.T) {

	i := int32(255 * 2)
	is := *(*[4]byte)(unsafe.Pointer(&i))
	if len(is) != 4 {
		t.Errorf("Pointer Error %v", len(is))
	}

	b := CircleBuffer{}
	b.Ini(20)

	p := Packet{}
	p.head = 12
	p.flg = 1
	p.size = 10
	p.data = make([]byte, p.size)
	p.Serialize(&b)

	if b.Size() != 15 {
		t.Errorf("Serialize Error %v", b.Size())
	}

	p2 := Packet{}
	if !p2.Deserialize(&b) {
		t.Errorf("Deserialize Error")
	}

	if b.Size() != 0 {
		t.Errorf("Deserialize Error %v", b.Size())
	}

	if p.head != p2.head || p.flg != p2.flg || p.size != p2.size {
		t.Errorf("Packet Error %v %v", p, p2)
	}
}
