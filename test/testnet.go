package main

import (
	"earthshaker"
	"fmt"
)

func main() {
	earthshaker.Ini(earthshaker.IniParam{Name: "testbuffer"})
	go start_server()
	go start_client()
	earthshaker.Exit()
}

var s earthshaker.NetServer
var port int = 8888
var n int = 8888

func start_server() {
	s.Ini(earthshaker.NetServerParam{MaxConn: 1000, SendBufSize: 10000, RecvBufSize: 10000})
	err := s.Listen(port)
	if err != nil {
		fmt.Println("Listen error")
	}
	fmt.Println("Listen OK ", port)

	for {
		err, c := s.Accept()
		if err != nil {
			fmt.Println("Accept error")
		}
		go server_processor(c)
	}
}

func server_processor(c earthshaker.Connection) {
	for c.IsConnect() {
		c.Send()
		c.Recv()
	}
}

func start_client() {

	for i := 0; i < n; i++ {

		c := earthshaker.NetClient{}
		c.Ini(earthshaker.NetClientParam{SendBufSize: 10, RecvBufSize: 10})
		err := c.Connect("127.0.0.1", port)
		if err != nil {
			fmt.Println("Connect error")
		}
	}

}
