package es


import(
    "elog"
)

type TcpClient struct {

    
    tcpAddr             *net.TCPAddr
    conn                *Connection 
    eventDispatch       *EventDispatch
    
    isConnection        bool 
    connCb              NET_CALLBACK
    readCb              NET_CALLBACK
    closeCb             NET_CALLBACK

    msgParse            IMsgParse
    msgDispatcher       *MsgDispatcher

    //safe close 
    waitGroup           *sync.WaitGroup

    stopChan              chan bool
}

func NewTcpClient()( tcpClient *TcpClient ){

    tcpClient = new ( TcpClient )
    
    tcpCient.isConnection = false;

    tcpClient.waitGroup = &sync.WaitGroup{}
    
    tcpClient.stopChan   = make( chan bool )
    
    tcpClient.eventDispatch ,_ = NewEventDispatch()
        
    return 
}


type ( tcpClient *TcpClient ) Connect( addr string )  err {
    
    s.tcpAddr , err = net.ResolveTCPAddr("tcp", addr )
    if err != nil {
        elog.LogSysln( addr, ":resolve tcp addr fail, please usage: 0.0.0.1:2345, fail: ", err )
        return
    }
    
    conn , err := net.DialTCP("tcp", s.tcpAddr )
    if err != nil {
        elog.LogSysln( "connect server , because :", err );
        return  err
    }
    
    s.isConnection = true
    tcpClient.conn = NeWConn( conn, tcpClient.eventDispatch )
    
    //开始处理网络
    go tcpClient.conn.handleEvent().
    
    return nil
}

func ( tcpClient  *TcpClient ) onRead(  conn* Connection, id int32, msg proto.Message  ) {

    //在这里加上用户消息处理时间
    tcpClient.msgDispatcher.DispatchMsg( conn, id, msg )
    elog.LogSys(" client read msg id :%d ", id )

}

func ( tcpClient *TcpServer ) onClose(  conn* Connection ) {

    if tcpClient.closeCb != nil {
        tcpClient.closeCb( conn )
    }
    elog.LogSys("cient close")
}

func ( tcpClient *TcpClient  )  SetConnCb( cb NET_CALLBACK )  {
     tcpClient.connCb = cb
}

func ( tcpClient *TcpClient )  SetReadCb( cb NET_CALLBACK )  {
    tcpClient.readCb = cb
}

func ( tcpClient *TcpClient )  SetCloseCb( cb NET_CALLBACK ) {
    tcpClient.closeCb = cb
}

func ( s *TcpClient ) SetMsgParse( parse IMsgParse ) {
    s.msgParse = parse
}

func ( tcpCient* TcpClient ) SetMsgDispather( dispatcher *MsgDispatcher ) {
    tcpClient.msgDispatcher = dispatcher
}

func ( tcpClient *TcpClient  ) Exit(  ) {

   tcpClient.eventDispatch.exit();
   elog.LogSys(" close new conn ch ")
   s.waitGroup.Wait()
   //关闭连接
   tcpClient.conn.Close()
   elog.Flush()
    
}


