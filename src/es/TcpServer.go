/*
 *  CopyRight  2014 , bingyu.song   All Right Reserved
    I believe  Spring brother
 */


package es


import (
    "net"
    "sync"
    "sync/atomic"
    "es/elog"
)


const defMaxConn = 50000

type  NET_CALLBACK  func( conn *Connection )  error  



type TcpServer struct {

    tcpAddr             *net.TCPAddr

    isListen            bool

    acceptor            *Acceptor

    newConnCh           chan *Connection

    maxConn             int 

    eventDispatch       *EventDispatch

    connCb              NET_CALLBACK
    readCb              NET_CALLBACK
    closeCb             NET_CALLBACK

    msgParse            IMsgParse
    msgDispatcher       *MsgDispatcher
    
    // slice 
    connMap             map[uint64]*Connection
    connMutex           sync.Mutex
    maxConnId           uint64
    
    //safe close 
    waitGroup           *sync.WaitGroup

    stopChan              chan bool
}


func NewTcpServer()( s *TcpServer ){

    s = new ( TcpServer )
    
    s.isListen = false;

    s.waitGroup = &sync.WaitGroup{}
    
    s.newConnCh  = make( chan *Connection)
    
    s.stopChan   = make( chan bool )
    
    s.eventDispatch ,_ = NewEventDispatch()
        
    s.maxConn = defMaxConn

    s.acceptor  = newAcceptor( s.waitGroup, s.eventDispatch, s.stopChan )

    s.connMap  = make( map[uint64]*Connection , s.maxConn )
    
    s.maxConnId     = 0;

    s.eventDispatch.RegistEvCb( CONN_NEW, s.onConn )
    
    s.eventDispatch.RegistEvCb( CONN_READ, s.onRead )
   
    s.eventDispatch.RegistEvCb( CONN_CLOSE, s.onClose )
   
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
    // main socket process
    go s.acceptor.loop( s.newConnCh )
   
    return  
}


func ( s* TcpServer ) Loop( ){

    if s.isListen == false {
    
        elog.LogSys(" not listen, please first listen ")

        return 
    }
    s.eventDispatch.Loop()
}


func ( s *TcpServer ) onConn( e IEvent  ) {
    
    ev := e.(*ConnNewEvent)
    //控制数量
    if (len( s.connMap ) >= s.maxConn ){
        ev.conn.Close()
        return
    }
    
    conn,_ := NewConn( ev.conn, s.eventDispatch, s.waitGroup ) 
    //设置消息解析器
    conn.setMsgParse( s.msgParse )
    //回调用户
    if s.connCb != nil {
        s.connCb( conn )
    }
    //hi, go 
    conn.SetConnStatus(  CONNECTED )    
    conn.SetConnId(atomic.AddUint64(&s.maxConnId, 1) )
    s.AddConn( conn )
    go conn.handleEvent(); 
    elog.LogSys("receive conn id: %s , total : %d", conn.connId, len(s.connMap ))
}

func ( s *TcpServer ) onRead( e IEvent ) {

    ev := e.(*ConnReadEvent)
    //在这里加上用户消息处理时间
    s.msgDispatcher.DispatchMsg( ev.conn, ev.id, ev.msg )
    elog.LogSys(" server  read msg id :%d ", ev.id )

}

func ( s *TcpServer ) onClose(  e IEvent ) {

    ev := e.(*ConnCloseEvent)
    if s.closeCb != nil {
        s.closeCb( ev.conn )
    }
    s.DelConn( ev.conn )
    elog.LogSys("cient close")
}

func ( s *TcpServer ) AddConn( conn *Connection ){
    s.connMutex.Lock()
    defer s.connMutex.Unlock()
    s.connMap[ conn.connId ] = conn
}


func ( s *TcpServer ) DelConn( conn *Connection ){
    s.connMutex.Lock()
    defer s.connMutex.Unlock()
    delete( s.connMap ,conn.connId )
}


func (s *TcpServer )  SetConnCb( cb NET_CALLBACK )  {
    s.connCb = cb
}

func (s *TcpServer )  SetReadCb( cb NET_CALLBACK )  {
    s.readCb = cb
}

func (s *TcpServer )  SetCloseCb( cb NET_CALLBACK ) {
    s.closeCb = cb
}

func ( s *TcpServer ) SetMsgParse( parse IMsgParse ) {
    s.msgParse = parse
}

func ( s* TcpServer ) SetMsgDispather( dispatcher *MsgDispatcher ) {
    s.msgDispatcher = dispatcher
}

func ( s* TcpServer ) Exit(  ) {

   s.eventDispatch.Close();
   s.stopChan <- false
   close( s.stopChan  )
   elog.LogSys(" close new conn ch ")
   close( s.newConnCh )
   elog.LogSys(" end  new conn ch ")
   elog.LogSys("wait...")

   s.waitGroup.Wait()
   elog.Flush()
    
}


