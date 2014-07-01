package earthshaker

import (
	"errors"
	"fmt"
	"net"
	"time"
	"netpb"
	proto "code.google.com/p/goprotobuf/proto"
)

type Connection struct {
	conn       net.Conn
	isConnect  bool
	SendPacket int
	RecvPacket int
	SendByte   int
	RecvByte   int
	sendBuffer CircleBuffer
	recvBuffer CircleBuffer
}

type IEventProcessor interface {
	OnRecvMsg(m * netpb.NetMsg) bool
}

type NetClient struct {
	Connection
	addr string
	port int
	parsePaket Packet
	parseMsg netpb.NetMsg
	iep IEventProcessor
}

type NetServer struct {
	lis         net.Listener
	isValid     bool
	maxConn     int
	sendBufSize int
	recvBufSize int
	conns       []Connection
}

type ConnectionParam struct {
	sendBufSize int
	recvBufSize int
}

type NetClientParam struct {
	SendBufSize int
	RecvBufSize int
	Addr string
	Port int
	Iep IEventProcessor
}

type NetServerParam struct {
	MaxConn     int
	SendBufSize int
	RecvBufSize int
}

func (c *Connection) Ini(p ConnectionParam) bool {
	c.isConnect = false
	c.SendPacket = 0
	c.RecvPacket = 0
	c.SendByte = 0
	c.RecvByte = 0
	if !c.sendBuffer.Ini(p.sendBufSize) {
		return false
	}
	if !c.recvBuffer.Ini(p.recvBufSize) {
		return false
	}
	return true
}

func (nc *NetClient) Ini(p NetClientParam) bool {
	nc.addr = p.Addr
	nc.port = p.Port
	return nc.Connection.Ini(ConnectionParam{sendBufSize: p.SendBufSize, recvBufSize: p.RecvBufSize})
}

func (ns *NetServer) Ini(p NetServerParam) bool {
	ns.isValid = false
	ns.maxConn = p.MaxConn
	ns.sendBufSize = p.SendBufSize
	ns.recvBufSize = p.RecvBufSize
	return true
}

func (c *Connection) IsConnect() bool {
	return c.isConnect
}

func (c *Connection) send() {
	c.conn.SetWriteDeadline(time.Now().Add(1 * time.Millisecond))
	n, err := c.sendBuffer.Output(c.conn.Write)
	c.SendByte += n
	if err != nil {
		c.isConnect = false
	}
}

func (c *Connection) recv() {
	c.conn.SetReadDeadline(time.Now().Add(1 * time.Millisecond))
	n, err := c.recvBuffer.Output(c.conn.Read)
	c.RecvByte += n
	if err != nil {
		c.isConnect = false
	}
}

func (nc *NetClient) connect(addr string, port int) error {
	err := error(nil)
	addrex := fmt.Sprint(addr, ":", port)
	nc.conn, err = net.Dial("tcp", addrex)
	if err != nil {
		nc.isConnect = false
	} else {
		nc.isConnect = true
	}
	return err
}

func (nc *NetClient) Heartbeat() {

	if nc.IsConnect() {
		nc.send()
		nc.recv()
		nc.process_msg()
	} else {
		nc.reconnect()
	}

}

func (nc *NetClient) process_msg() {
	var err error
	for nc.parsePaket.Deserialize(&nc.Connection.recvBuffer) {
		err = proto.Unmarshal(nc.parsePaket.data, &nc.parseMsg)
		if err != nil {
			continue
		}
		if !nc.iep.OnRecvMsg(&nc.parseMsg) {
			break
		}
	}
}

func (nc *NetClient) SendMsg(m * netpb.NetMsg) bool {
	nc.parsePaket.head = 2014
	nc.parsePaket.flg = 0
	var err error
	nc.parsePaket.data, err = proto.Marshal(m)
	if err != nil {
		return false
	}
	nc.parsePaket.size = int16(len(nc.parsePaket.data))
	if !nc.parsePaket.Serialize(&nc.Connection.sendBuffer) {
		return false
	}
	return true
}

func (nc *NetClient) reconnect() {
	nc.connect(nc.addr, nc.port)
}

func (ns *NetServer) Listen(port int) error {
	if ns.isValid {
		return errors.New("Listen Fail:Already Valid")
	}
	err := error(nil)
	portex := fmt.Sprint(":", port)
	ns.lis, err = net.Listen("tcp", portex)
	if err != nil {
		ns.isValid = false
	} else {
		ns.isValid = true
	}
	return err
}

func (ns *NetServer) Accept() (err error, c Connection) {
	if !ns.isValid {
		err = errors.New("Accept Fail:Closed Listener")
		return
	}
	ns.lis.(*net.TCPListener).SetDeadline(time.Now().Add(1 * time.Millisecond))
	c.conn, err = ns.lis.Accept()
	if err != nil {
		c.isConnect = false
	} else {
		c.isConnect = true
		c.Ini(ConnectionParam{sendBufSize: ns.sendBufSize, recvBufSize: ns.recvBufSize})
	}
	return
}

func (ns *NetServer) Close() {
	if ns.isValid {
		ns.lis.Close()
		ns.isValid = false
		// TODO : 连接的处理
	}
}
