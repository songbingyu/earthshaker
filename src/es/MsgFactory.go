/*
 *  CopyRight  2014 , bingyu.song   All Right Reserved
    I believe  Spring brother
 */

package  es

import (
    "reflect"
    
    "es/elog"
    proto "code.google.com/p/goprotobuf/proto"
)

type  IMsgFactory interface {
    
   createMsgById( id int32 ) proto.Message  

}


type MsgFactory   struct {
    msgMap  [ MSG_MAX]reflect.Type
}

func NewMsgFactory() *MsgFactory {
    
    factory := &MsgFactory{
        
    }

    return factory
}

func ( factory *MsgFactory ) createMsgById( id int32 ) proto.Message {
   if id < 0 || id > MSG_MAX  {
        return nil
   }

   msg :=  reflect.New( factory.msgMap[id ] ).Interface().( proto.Message )
    
   elog.LogSysln("create msg :", reflect.TypeOf( msg ))
   return msg 
}

func ( factory *MsgFactory  ) RegistMsg( id int32, msgType  reflect.Type   ) {
    
    factory.msgMap[id] = msgType
}

func ( factory *MsgFactory ) UnregistMsg( id int32  ) {
    
    factory.msgMap[id] = nil
}


















