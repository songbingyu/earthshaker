/*
 *  CopyRight  2014 , bingyu.song   All Right Reserved
    I believe  Spring brother
 */


package es

import (
    "fmt"
    //"errors"
    //"bytes"
    //"unsafe"
    "reflect"
    "encoding/binary"
    "elog"
    proto "code.google.com/p/goprotobuf/proto"
)

type  Msg struct  {
 len        int32 
 id         int32
 //name       []byte
 body       []byte
 //check      int  check is must ?
}

const  c_HeaderLen  int32  = 4 //unsafe.Sizeof( int16(0) )
const  c_MinMsgLen  int32  = 2*c_HeaderLen     // 
const  c_MaxMsgLen  int32  = 64*1024*1024      //

// 通用接口
type  IMsgParse  interface  {

    Encode(conn *Connection ) error 
    Decode(conn *Connection ) error 

} 


type   LiteMsgParse struct {
    
    //添加一些统计信息

}


func ( lmp * LiteMsgParse )  Encode ( conn *Connection) error {
    
    // let go ...
    var len int32
    var id  int32
    var err error
    elog.LogSysln( "recv msg :",conn.recvBuf.Size() )
    for conn.recvBuf.Size() >= int(c_MinMsgLen) {
        //Fixme  add little big endian cfg 
        err = binary.Read( &conn.recvBuf, binary.LittleEndian, &len ) 
        if( err != nil ) {
            elog.LogSysln( " recv msg len  err : " , err )
            return err 
        }

        elog.LogSysln( " recv msg len : " , len )
        
        if  len > c_MaxMsgLen || len < c_MinMsgLen {
            err = fmt.Errorf("msg len is error")
            break
        }
        
        if conn.recvBuf.Size()  >= int(len - c_HeaderLen ) {
            err = binary.Read( &conn.recvBuf, binary.LittleEndian, &id ) 
            if err != nil || id >= MSG_MAX   {
                elog.LogSysln( " recv msg id err : " , err )
                break
            }
            elog.LogSysln( " recv msg id : " , id )
            b := make( []byte, len- c_HeaderLen )
            err = binary.Read( &conn.recvBuf, binary.LittleEndian, &b )
            if err != nil {
                break
            }
            info := GetMsgDisptcher().GetMsgInfo( int(id) )
            if info == nil {
               elog.LogSys("not register msg is : %d", id )
               continue
            }
            var msg proto.Message = reflect.New( info.msgType ).Elem().Interface().( proto.Message )
            err = proto.Unmarshal( b, msg) 
            if err != nil {
                continue 
            }
            
        } else {
            elog.LogSysln( " recv msg len is big than buffer  " )
            return err
        }
    }
    return err
}

func ( lmp *LiteMsgParse ) Decode( conn *Connection ) error  {
    
    return  nil
}













