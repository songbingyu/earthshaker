/*
 *  CopyRight  2014 , bingyu.song   All Right Reserved
    I believe  Spring brother
 */


package es

import (
    "sync"

    "elog"
)

type EventLoop struct {

    waitGroup   *sync.WaitGroup
}

func newEventLoop( w *sync.WaitGroup )( el *EventLoop ){

    el = &EventLoop{ waitGroup : w ,
    
                   }

    return el
    
}

func ( el *EventLoop ) loop( conns <-chan *Connection ) {

 
    el.waitGroup.Add(1)
    defer el.waitGroup.Done()
    
    elog.LogSys(" io goroutine begin proc")
    for conn := range conns {
        el.handleEvent( conn )
    }
    
    elog.LogSys(" io goroutine end proc")

}

func ( el *EventLoop ) handleEvent( conn *Connection ) {

    elog.LogSys(" handle event ")
    conn.handleEvent()

}


