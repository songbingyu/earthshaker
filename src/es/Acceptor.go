/*
 *  CopyRight  2014 , bingyu.song   All Right Reserved
    I believe  Spring brother
 */


 //Accept wraper

package es

import (
    "net"
    "sync"
    "time"
    "elog"
)


type Acceptor struct {
    listener       *net.TCPListener
    eventDispatch  *EventDispatch
    waitGroup      *sync.WaitGroup
    stopChan       chan bool
}

func newAcceptor(   w *sync.WaitGroup, ed *EventDispatch, ch chan bool  )( acc *Acceptor ){
    
    acc = &Acceptor{
        waitGroup: w,
        eventDispatch : ed,
        stopChan : ch,
    }
    return
}

func( a *Acceptor )loop(  ch chan *Connection ) {

    elog.LogSys(" acceptor is begin ")
    a.waitGroup.Add(1)
    defer a.waitGroup.Done()

    for {

        select {
            case <-a.stopChan:
                 elog.LogSys(" recv stop chan ,stop Accept ")
                 a.listener.Close()
                 return
            default:

        }
        
        a.listener.SetDeadline( time.Now().Add( 1e9 ) )
        conn, err := a.listener.AcceptTCP()
        if err != nil {
            if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
                continue
            }
            elog.LogSysln(" accetp fail :", err )
            break
        }
        
        elog.LogSys("receive new conn")
        
        //newConn ,err := NewConn( conn, a.eventDispatch, a.waitGroup )
        
        a.eventDispatch.AddNewConnEvent( conn )

    }
    
    elog.LogSys(" acceptor is end ")
}



