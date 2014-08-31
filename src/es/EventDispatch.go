/*
 *  CopyRight  2014 , bingyu.song   All Right Reserved
    I believe  Spring brother
 */


package es


import (
    "net"
    "sync"
    "elog"
    proto "code.google.com/p/goprotobuf/proto"
)

type  EventType   int16

const  (

    CONN_NEW EventType = iota
    // read write
    CONN_READ
    CONN_WRITE
    CONN_CLOSE
    EVENT_MAX
)


//Fixme : need test
const defEventNum = 1000;

type IEvent   interface {
    
    GetEventType() EventType 
    Next() IEvent
    Clear()
}

type ConnNewEvent struct {
    conn        *net.TCPConn
    next        *ConnNewEvent 
}

func ( ev *ConnNewEvent ) GetEventType() EventType {
    return CONN_NEW
}


func ( ev *ConnNewEvent ) Clear() {
    ev.conn = nil
    ev.next = nil
}

func ( ev *ConnNewEvent ) Next() IEvent {
   return ev.next
}

type ConnReadEvent struct {
    conn        *Connection
    id          int32
    msg         proto.Message
    next        *ConnReadEvent 
}

func ( ev *ConnReadEvent ) Clear() {
    
    ev.conn = nil
    ev.id = 0
    ev.msg = nil
    ev.next = nil
}

func ( ev *ConnReadEvent ) Next() IEvent {
   return ev.next
}

func ( ev *ConnReadEvent ) GetEventType() EventType {
    return CONN_READ
}

type ConnWriteEvent struct {
    id          int32
    msg         proto.Message
    next        *ConnWriteEvent 
}

func ( ev *ConnWriteEvent ) GetEventType() EventType {
    return CONN_WRITE
}

func ( ev *ConnWriteEvent ) Clear() {
    
    ev.id = 0
    ev.msg = nil
    ev.next = nil
}

func ( ev *ConnWriteEvent ) Next() IEvent {
   return ev.next
}

type ConnCloseEvent struct {
    conn        *Connection
    next        *ConnCloseEvent 
}

func ( ev *ConnCloseEvent ) GetEventType() EventType {
    return CONN_CLOSE
}

func ( ev *ConnCloseEvent ) Clear() {
    ev.conn = nil
    ev.next = nil
}

func ( ev *ConnCloseEvent ) Next() IEvent {
    return ev.next
}


type EventDispatch  struct {

    eventMu           [EVENT_MAX]sync.Mutex
    eventList         [EVENT_MAX] IEvent
    eventCh           [EVENT_MAX]chan IEvent
    
}


func ( e *EventDispatch ) NewEventByType( et EventType )( ne IEvent ) {
   
   switch et {
        case CONN_NEW:
            return new(ConnNewEvent )
        case CONN_READ:
            return new(ConnReadEvent)
        case CONN_WRITE:
            return new(ConnWriteEvent)
        case CONN_CLOSE:
            return new(ConnCloseEvent)
   }

   return nil
}

func ( e *EventDispatch ) AllocEvent( et EventType ) ( ev IEvent, err error ){
   
   e.eventMu[ et ].Lock()
    ev  = e.eventList[ et ]
    if  ev != nil {
        e.eventList[et] = ev.Next()
    }
   e.eventMu[et].Unlock()
   if ev == nil {
         ev = e.NewEventByType( et )
   } else{
        ev.Clear()
   }
   err = nil
   return 
}

func ( e *EventDispatch) FreeEvent( ev IEvent ) {

    et := ev.GetEventType()
    e.eventMu[et].Lock()
    ne := ev.Next()
    // make compliler happy 
    if ne != nil {
        return
    }
    ne = e.eventList[et]
    e.eventList[et] = ev
    e.eventMu[et].Unlock()

    elog.LogSys(" free event type : %d", ev.GetEventType())
}

func ( e *EventDispatch ) initAllocEvent( ) {

    // init alloc some event
    for j := 0; j < int(EVENT_MAX); j++ {
        e.eventMu[j].Lock()
        for i := 0; i < defEventNum ; i++  {
            ev := e.NewEventByType( EventType(j))
            next := ev.Next()
            //make compiler happy 
            if next != nil {
                return
            }
            next = e.eventList[j]
            e.eventList[j] = ev
        }
        e.eventMu[j].Unlock()

    }
}


func  NewEventDispatch()( ed *EventDispatch, err error ){

    ed = &EventDispatch{
            //eventCh : make( [EVENT_MAX]chan IEvent),
    }
    
    for  i :=  range ed.eventCh {
        ed.eventCh[i] = make( chan IEvent )
    }
    err =  nil
    
    ed.initAllocEvent()
    return
}

func ( ed *EventDispatch ) pushEvent( ev IEvent ){

    elog.LogSys("push %d event to event channel ", ev.GetEventType() )     
    ed.eventCh[ ev.GetEventType() ] <- ev
}

func ( ed *EventDispatch ) AddNewConnEvent( c *net.TCPConn ) { 
    
    e, _:= ed.AllocEvent( CONN_NEW)
    ev := e.( *ConnNewEvent )
    ev.conn = c
    ed.pushEvent( ev )
    return  
}

func ( ed *EventDispatch ) AddConnReadEvent( conn *Connection, id int32, msg proto.Message ) {
    e, _:= ed.AllocEvent( CONN_READ )
    ev := e.( *ConnReadEvent )
    ev.conn = conn
    ev.id = id
    ev.msg = msg 
    ed.pushEvent( ev )
    return  
}

func ( ed *EventDispatch ) AddConnWriteEvent( conn *Connection, id int32, msg proto.Message ) {
    e, _:= ed.AllocEvent( CONN_READ )
    ev := e.( *ConnReadEvent )
    ev.conn = conn
    ev.id = id
    ev.msg = msg 
    ed.pushEvent( ev )
    return  
}
func ( ed *EventDispatch ) AddConnCloseEvent( conn *Connection  ) {
    
    e, _:= ed.AllocEvent( CONN_CLOSE )
    ev := e.( *ConnCloseEvent )
    ev.conn = conn
    ed.pushEvent( ev )
    return  
}

func ( ed *EventDispatch  ) DelEvent( ev IEvent ) {
    ed.FreeEvent( ev )    
}

/*func ( ed *EventDispatch ) getEvents()( ch  chan *NetEvent ){
    
    ch = ed.eventCh
    return
}*/



func ( ed *EventDispatch ) exit() {

    elog.LogSys("close event channel")
    
    for i := range ed.eventCh {
        close( ed.eventCh[i] )
    }
}





