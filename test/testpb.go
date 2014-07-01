package main

import (
	proto "code.google.com/p/goprotobuf/proto"
	"earthshaker"
	"netpb"
	"fmt"
)

func main() {
	earthshaker.Ini(earthshaker.IniParam{Name: "testbuffer"})
	testnetpb()
	earthshaker.Exit()
}

func testnetpb() {
	msg := &netpb.NetMsg {
		M_ID : netpb.ToMsgID(netpb.MsgID_EM_CS_LOGIN),
		M_CSNetMsg : &netpb.CSNetMsg {
			M_CSLoginReq : &netpb.CSLoginReq {
				M_Name : proto.String("aaa"),
				M_Pwd : proto.String("bbb"),
			},
		},
	}

	newmsg := new(netpb.NetMsg)

	for i := 0; i < 1000000; i++ {
		data, err := proto.Marshal(msg) //SerializeToOstream
		if err != nil {
			fmt.Print("marshaling error: ", err)
		}
		err = proto.Unmarshal(data, newmsg)
		if err != nil {
			fmt.Print("unmarshaling error: ", err)
		}
	}
}
