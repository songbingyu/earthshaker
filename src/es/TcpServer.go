/*
 *  CopyRight  2014 , bingyu.song   All Right Reserved
    I believe  Spring brother
 */


package es


import (
    "net"
    "sync"

    "elog"
)


const defMaxNetGoroutines = 1;



type TcpServer struct {

    tcpAddr             *net.TCPAddr

    isListen            bool

    acceptor            *Acceptor

    newConnCh           chan *Connection

    maxNetGoroutines    int 

    eventLoop           *EventLoop

    eventDispatch       *EventDispatch

    //safe close 
    waitGroup           *sync.WaitGroup

    stopChan              chan bool
}


func NewTcpServer()( s *TcpServer, err error ){

    s = new ( TcpServer )
    
    s.isListen = false;

    s.waitGroup = &sync.WaitGroup{}
    
    s.newConnCh  = make( chan *Connection)
    
    s.stopChan   = make( chan bool )
    
    s.eventDispatch ,_ = NewEventDispatch()
        
    s.maxNetGoroutines = defMaxNetGoroutines

    s.eventLoop = newEventLoop( s.waitGroup )

    s.acceptor  = newAcceptor( s.waitGroup, s.eventDispatch, s.stopChan )

    err = nil 
    return 
}


func ( s* TcpServer ) Listen( addr string ) ( err error ) {

    s.tcpAddr , err = net.ResolveTCPAddr("tcp", addr )
    if err != nil {
        elog.LogSysln( addr, ":resolve tcp addr fail, please usage: 0.0.0.1:2345, fail: ", err )
        return
    }
    
    s.acceptor.listener , err = net.ListenTCP("tcp", s.tcpAddr )
    if err != nil {
        elog.LogSysln( "listen tcp fail, because :", err );
        return 
    }
    
    s.isListen = true
    
    // IO process
    for i:=0; i < s.maxNetGoroutines; i++  {
            go s.eventLoop.loop( s.newConnCh )
    }
    // main socket process
    go s.acceptor.loop( s.newConnCh )
   
    return  
}


func ( s* TcpServer ) Run( ){

    if s.isListen == false {
    
        elog.LogSys(" not listen, please first listen ")

        return 
    }
    
    elog.LogSys("begon recv event ") 

   go func (){
    elog.LogSys("begin wait events ")
    s.waitGroup.Add(1)
    defer s.waitGroup.Done()
    
    for  ne := range s.eventDispatch.getEvents() { 
        switch ne.eventType {
            case NEW:
                s.OnConn( ne.conn );
            case READ:
                s.OnRead( ne.conn )
            case CLOSE:
                s.OnClose( ne.conn )
            default:
                elog.LogSysln(" unknow event type :", ne.eventType)
        }
    }
    elog.LogSys("end wait events ")
  }()
}


func ( s *TcpServer ) OnConn(  conn* Connection ) {

    elog.LogSys("receive conn")
}

func ( s *TcpServer ) OnRead(  conn* Connection ) {

    elog.LogSys(" client read")

}

func ( s *TcpServer ) OnClose(  conn* Connection ) {

    elog.LogSys("cient close")
}


func ( s* TcpServer ) Exit(  ) {

   s.eventDispatch.exit();
   s.stopChan <- false
   close( s.stopChan  )
   elog.LogSys(" close new conn ch ")
   close( s.newConnCh )
   elog.LogSys(" end  new conn ch ")
   elog.LogSys("wait...")

   s.waitGroup.Wait()
   elog.Flush()
    
}


