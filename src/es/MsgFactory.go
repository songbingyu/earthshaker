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

type  IMsgFactory interface {
    
   createMsgById( id int ) proto.Message  

}


type MsgFactory   struct {
    msgMap  [ MSG_MAX]reflect.Type
}


func ( factory *MsgFactory ) createMsgById( id int ) proto.Message {
   if id < 0 || id > MSG_MAX  {
        return nil
   }

   msg :=  reflect.New( factory.msgMap[id ] ).Interface().( proto.Message )

   return msg 
}

func ( factory *MsgFactory  ) RegistMsg( id int, cb CALLBACK, msgType  reflect.Type   ) {
    
    factory.msgMap[id] = msgType
}

func ( factory *MsgFactory ) UnregistMsg( id int  ) {
    
    factory.msgMap[id] = nil
}


















