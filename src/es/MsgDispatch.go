/*
 *  CopyRight  2014 , bingyu.song   All Right Reserved
    I believe  Spring brother
 */

package  es

import (
    //"reflect"
    
    "elog"
    proto "code.google.com/p/goprotobuf/proto"
)

type CALLBACK     func( *Connection,  proto.Message  ) 


type  MsgDispatcher struct {
    
   cbMap   [MSG_MAX] CALLBACK 
}


func  NewMsgDispatcher() *MsgDispatcher {

    md := &MsgDispatcher {

    }
    return md
}

func ( md *MsgDispatcher ) RegistMsgCb( id int, cb CALLBACK ) {
    
    md.cbMap[id] = cb
}

func ( md *MsgDispatcher ) UnregistMsgCb( id int  ) {
    
    md.cbMap[id] = nil
}

func ( md* MsgDispatcher ) DispatchMsg( conn *Connection, id int32,msg proto.Message ) {

    if id >=  MSG_MAX {
        return 
    }

    elog.LogSys(" dispatcher msg :%d ", id )
    
    cb  := md.cbMap[id]
    
    if cb == nil {
        
        elog.LogSys(" dispatcher msg :%d  not register cb ", id )
        return
    }

    cb( conn, msg  )
}

func (md *MsgDispatcher ) GetMsgCb( id int32 )  CALLBACK  {

    cb := md.cbMap[id]
    return cb 
}





























