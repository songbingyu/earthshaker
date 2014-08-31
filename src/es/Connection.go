/*
 *  CopyRight  2014 , bingyu.song   All Right Reserved
    I believe  Spring brother
 */


package es

import (
    "net"
    //"bytes"
    "sync"
    "sync/atomic"
    proto "code.google.com/p/goprotobuf/proto"
    "elog"
)

const defSendBufSize  = 1024*1024
const defRecvBufSize = 1024*1024

const (
    NOCONNECT uint32 = iota
    CONNECTED
    //TODO: add pool manger
)

type Connection struct {
    tcpConn         *net.TCPConn
    //Fixme  should atomic 
    connectState    uint32    
    eventDispatch   *EventDispatch
    sendBufSize     int
    recvBufSize     int
    sendBuf         CircleBuffer
    recvBuf         CircleBuffer
    //sendBuf         bytes.Buffer
    //recvBuf         bytes.Buffer
    //Fixme :   
    sendMutex       sync.Mutex
    sendChan        chan IEvent
    //msg parse 

    msgParse        IMsgParse
    
    connId          uint64
    connCb          NET_CALLBACK
    readCb          NET_CALLBACK
    closeCb         NET_CALLBACK
    
    closeCh         chan bool
    waitGroup       *sync.WaitGroup
     
    //monitor
    sendBytes       int64
    recvBytes       int64

    sendPacket     int
    recvPacket     int

    lastRecvTime    int
    
    
    
}

func NewConn( conn *net.TCPConn ,dispatch *EventDispatch, wait *sync.WaitGroup )( connection *Connection, err error ) {
    
    // TODO: connection pool manger

    connection = &Connection{ 
            tcpConn     : conn,
            eventDispatch: dispatch,
            connectState : NOCONNECT,
            sendBufSize : defSendBufSize,
            recvBufSize : defRecvBufSize,
            sendBytes : 0,
            recvBytes : 0,
            sendPacket : 0,
            recvPacket : 0,
            lastRecvTime : 0,
            waitGroup     : wait,
            connId        : 0,
            sendChan      : make( chan IEvent, 100 ),
            closeCh       : make( chan bool ),
            //msgParse    : &LiteMsgParse{},
        }

    connection.sendBuf.Ini( connection.sendBufSize)
    connection.recvBuf.Ini( connection.recvBufSize)
    elog.LogSys("create New conn ")
    
   return
}

func ( conn  *Connection ) IsConnected( ) bool {

    return conn.connectState ==  CONNECTED 
}

func ( conn *Connection )Close() {

    if  atomic.CompareAndSwapUint32(&conn.connectState, CONNECTED, NOCONNECT ) {
        
        conn.eventDispatch.AddConnCloseEvent( conn )
        conn.tcpConn.Close()
        close(  conn.closeCh )
        close( conn.sendChan )
    }
}

func (conn *Connection ) handleEvent( ) {
 
    conn.waitGroup.Add(1)
    defer conn.waitGroup.Done()
    for {
        
        select {
            case <- conn.closeCh :
            break
        }

        elog.LogSys(" recvieve msg ........... ")
        // let reading....
        err := conn.Read()
        if err != nil {
            elog.LogSysln(" recvieve msg fail ", err )
            conn.Close()
            break
        }

        elog.LogSys(" recvieve msg .... ,begin encode ")
        
        //parse
        conn.msgParse.Encode( conn )

    }
    
}

func ( conn *Connection ) handleWrite() {

    // write 
    go conn.handleWrite()

    conn.waitGroup.Add(1)
    defer conn.waitGroup.Done()
    for {

        select {
            
            case  e :=  <- conn.sendChan :
                    defer conn.eventDispatch.DelEvent( e )
                    if e.GetEventType() == CONN_WRITE {
                        ev := e.(*ConnWriteEvent )
                        //这里的错误该如何处理？
                        conn.msgParse.Decode( conn, ev.id, ev.msg )
                        err := conn.Write( )
                        if err != nil {
                            conn.Close()
                            return
                        }
                    }
            case  <- conn.closeCh:
                    return            
        }
            
    }
    return
}
//Fixme ? should do something?
func ( conn *Connection ) Read( ) error {
     n, err  := conn.recvBuf.Input( conn.tcpConn.Read) 
     conn.recvBytes += int64(n)
     return err
}

func ( conn *Connection ) Write() error {

     n, err := conn.sendBuf.Output( conn.tcpConn.Read) 
     conn.sendBytes += int64(n)
     return err
}

func ( conn *Connection ) Send ( id int32, msg proto.Message ){
    e, _:= conn.eventDispatch.AllocEvent( CONN_WRITE )
    ev := e.( *ConnWriteEvent )
    ev.id =  id
    ev.msg = msg

    conn.sendChan <- ev 

}

func ( conn *Connection ) setMsgParse( parse IMsgParse ) {
    
    conn.msgParse = parse
}

func ( conn *Connection ) SetConnId( id uint64) {

    conn.connId = id
}

func ( conn *Connection ) SetConnStatus( state uint32 ) {
   
   atomic.StoreUint32( &conn.connectState, state)
}


