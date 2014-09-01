package main

import (
    "fmt"
    "elog"
    "es"
    "os"
    "os/signal"
    "syscall"
    proto "code.google.com/p/goprotobuf/proto"
    "netpb"
    "reflect"
)


type esserver struct {

    codec           *es.LiteMsgParse
    dispatch        *es.MsgDispatcher
    msgFactory      *es.MsgFactory
    tcpServer       *es.TcpServer 
}

func ( s *esserver ) Init() {

    s.dispatch.RegistMsgCb( 1, onMsg )
    s.msgFactory.RegistMsg( 1, reflect.TypeOf( netpb.NetMsg{} ) )
    fmt.Println(" begin listen....")
    s.tcpServer.Listen("0.0.0.0:6798")
}


func ( s *esserver ) Run() {

    s.tcpServer.Loop()

}
func NewEsServer ( ) *esserver {
    
    server := &esserver {
        
        codec      : es.NewLiteMsgParse(),
        dispatch   : es.NewMsgDispatcher(),
        msgFactory : es.NewMsgFactory(),
        tcpServer  : es.NewTcpServer(),
    }

    server.codec.SetMsgFactory( server.msgFactory )

    server.tcpServer.SetMsgParse( server.codec )
    server.tcpServer.SetMsgDispather( server.dispatch )

    return server
}


func onMsg ( es  *es.Connection, pb proto.Message  ) {
    
    msg  := pb.(*netpb.NetMsg)

    elog.LogSysln("i receieve msg ", msg.M_ID, reflect.TypeOf( pb ))
}

func main() {


    elog.InitLog( elog.INFO )
    
    server := NewEsServer()
    
    
    ch := make(chan os.Signal)
    signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
    
    server.Init()
    server.Run()

            
    fmt.Println(<-ch)

    server.tcpServer.Exit()

    elog.LogSys("Hhhhhh");

}
