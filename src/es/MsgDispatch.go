/*
 *  CopyRight  2014 , bingyu.song   All Right Reserved
    I believe  Spring brother
 */

package  es

import (
    "reflect"
    
    //"elog"
    proto "code.google.com/p/goprotobuf/proto"
)

type CALLBACK     func( *Connection, *proto.Message ) 

type  MsgInfo  struct {
    cb         CALLBACK
    msgType    reflect.Type
}

type  MsgDispatcher struct {
    
   msgMap   [MSG_MAX]*MsgInfo 
}

var  instance  *MsgDispatcher

func  GetMsgDisptcher() *MsgDispatcher {

    if instance == nil {
        instance = &MsgDispatcher {
        }
    }
    return instance
}

func ( md *MsgDispatcher ) RegistMsg( id int, cb CALLBACK, msg reflect.Type   ) {
    
    md.msgMap[id] = &MsgInfo { cb, msg }
}

func ( md *MsgDispatcher ) UnregistMsg( id int  ) {
    
    md.msgMap[id] = nil
}

func (md *MsgDispatcher ) GetMsgInfo( id int  )  ( info *MsgInfo ) {

    info = md.msgMap[id]
    return 
}































