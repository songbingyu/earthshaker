/*
 *  CopyRight  2014 , bingyu.song   All Right Reserved
    I believe  Spring brother
 */


package es


import (
    "elog"
    "sync"
)

const  (

    NEW = iota
    READ
    CLOSE
)

type netEvent struct {

    eventType   int
    conn        *Connection
}


type EventDispatch  struct {

    mu       sync.Mutex
     
    eventCh  chan *netEvent
    
}

func  NewEventDispatch()( ed *EventDispatch, err error ){

    ed = &EventDispatch{
            eventCh : make( chan *netEvent ),
    }

    err =  nil

    return
}

func ( eq *EventDispatch ) pushEvent( ne *netEvent ){

    // eq.mu.Lock()
    // defer eq.mu.Unlock()
    elog.LogSys("push event to event channel ")     
    eq.eventCh <- ne
}

func ( eq *EventDispatch ) dispatch() (ne *netEvent ) { 

    //eq.mu.Lock()
    //defer eq.mu.Unlock()
    
    // select ?

    /*select {
        case ne =  <- eq.eventCh:
             return ne
    }*/


    return nil
}

func ( eq *EventDispatch ) getEvents()( ch  chan *netEvent ){
    
    ch = eq.eventCh
    return
}


func ( eq *EventDispatch ) exit() {

    elog.LogSys("close event channel")
    close( eq.eventCh )
}





