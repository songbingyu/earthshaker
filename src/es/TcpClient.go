package es


import(
    "net"
   "sync"
   "es/elog"
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
    
    tcpClient.isConnection = false;

    tcpClient.waitGroup = &sync.WaitGroup{}
    
    tcpClient.stopChan   = make( chan bool )
    
    tcpClient.eventDispatch ,_ = NewEventDispatch()
    
    tcpClient.eventDispatch.RegistEvCb( CONN_READ, tcpClient.onRead )
   
    tcpClient.eventDispatch.RegistEvCb( CONN_CLOSE, tcpClient.onClose )
 
    return 
}


func  ( tcpClient *TcpClient ) Connect( addr string )  error {
    
    var err error
    tcpClient.tcpAddr , err = net.ResolveTCPAddr("tcp", addr )
    if err != nil {
        elog.LogSysln( addr, ":resolve tcp addr fail, please usage: 0.0.0.1:2345, fail: ", err )
        return err
    }
    
    var localaddr *net.TCPAddr
    localaddr , err = net.ResolveTCPAddr("tcp", "0.0.0.0")
    conn , err := net.DialTCP("tcp", localaddr,tcpClient.tcpAddr )
    if err != nil {
        elog.LogSysln( "connect server , because :", err );
        return  err
    }
    
    tcpClient.isConnection = true
    tcpClient.conn, err = NewConn( conn, tcpClient.eventDispatch, tcpClient.waitGroup )
    

    if tcpClient.connCb != nil {
        tcpClient.connCb( tcpClient.conn )
    }
    //开始处理网络
    go tcpClient.conn.handleEvent()
    
    return err
}


func ( tcpClient *TcpClient ) Loop( ) {
    if !tcpClient.isConnection {
        elog.LogSys(" client not connevt server ")
    }

    tcpClient.eventDispatch.Loop()

}
func ( tcpClient  *TcpClient ) onRead( e IEvent  ) {

    ev := e.(*ConnReadEvent)
    //在这里加上用户消息处理时间
    tcpClient.msgDispatcher.DispatchMsg( ev.conn, ev.id, ev.msg )
    elog.LogSys(" client read msg id :%d ", ev.id )

}

func ( tcpClient *TcpClient ) onClose(e IEvent ) {

    ev := e.(*ConnCloseEvent)
    if tcpClient.closeCb != nil {
        tcpClient.closeCb( ev.conn )
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

func ( tcpClient  *TcpClient ) SetMsgParse( parse IMsgParse ) {
    tcpClient.msgParse = parse
}

func ( tcpClient* TcpClient ) SetMsgDispather( dispatcher *MsgDispatcher ) {
    tcpClient.msgDispatcher = dispatcher
}

func ( tcpClient *TcpClient  ) Close(  ) {

   tcpClient.eventDispatch.Close();
   elog.LogSys(" close  conn  ")
   tcpClient.conn.Close()
   tcpClient.waitGroup.Wait()
   //关闭连接
   elog.Flush()
    
}


