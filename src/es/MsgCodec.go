/*
 *  CopyRight  2014 , bingyu.song   All Right Reserved
    I believe  Spring brother
 */


package es

import (
    //"fmt"
    "errors"
    //"bytes"
    "unsafe"
    "reflect"
    "encoding/binary"
    "es/elog"
    proto "code.google.com/p/goprotobuf/proto"
)

type  Msg struct  {
 len        int32 
 id         int32
 //name       []byte
 body       []byte
 //check      int  check is must ?
}

const  c_HeaderLen  int32  = int32( unsafe.Sizeof( int32(1) ) )
const  c_MinMsgLen  int32  = 2*c_HeaderLen     // 
const  c_MaxMsgLen  int32  = 64*1024*1024      //

// 通用接口
type  IMsgParse  interface  {

    Encode(conn *Connection ) error 
    Decode(conn *Connection ,id int32 , msg proto.Message ) error 

} 

type   MSG_CALLBACK         func( conn *Connection, id int32 , pb proto.Message  ) 
type   ERROR_CALLBACK       func( conn *Connection ) 

const (
    c_NoError = iota
    c_NeedContinue
    c_ErrLen
    c_ErrMsgId
    c_ParseErr
)


type   LiteMsgParse struct {
    
   msgCb        MSG_CALLBACK
   errCb        ERROR_CALLBACK
   msgFactory   IMsgFactory
   //添加一些统计信息
}



func defErrorCallback( conn *Connection ) {

   if conn.IsConnected() {
        elog.LogSys(" read msg error, close socket ")     
        conn.Close()
   }
    
}

func defMsgCallback( conn *Connection, id int32,  msg proto.Message ) {

    conn.eventDispatch.AddConnReadEvent( conn, id, msg )
}

func ( lmp *LiteMsgParse ) parse( conn *Connection ) ( id int32, errNo int, msg proto.Message ) {
    
    var len int32
    //Fixme  add little big endian cfg 
    err := binary.Read( &conn.recvBuf, binary.LittleEndian, &len ) 
    if( err != nil ) {
        elog.LogSysln( "read  msg len  err : " , err )
        errNo = c_ErrLen
        return  
    }
    
    elog.LogSysln( " parse  msg len : " , len )
    
    if  len > c_MaxMsgLen || len < c_HeaderLen {
        elog.LogSysln( " msg len is  err : " )
        errNo = c_ErrLen
        return 
    }

    if conn.recvBuf.Size()  >= int(len - c_HeaderLen ) {
        err = binary.Read( &conn.recvBuf, binary.LittleEndian, &id ) 
        if err != nil || id >= MSG_MAX   {
            elog.LogSysln( "read  msg id err : " , err )
            errNo = c_ErrMsgId
            return 
        }
        elog.LogSysln( " parsr  msg id : " , id )

        /*info := GetMsgDispatcher().GetMsgInfo( int(id) )
          if info == nil {
          elog.LogSys("not register msg is : %d", id )
          return c_ErrMsgId
        }*/
        
        // msg may be nil
        if ( len - c_HeaderLen ) > 0 {
            b := make( []byte, len- c_HeaderLen )
            err = binary.Read( &conn.recvBuf, binary.LittleEndian, &b )
            if err != nil {
                elog.LogSys("read  msg body fail ", id )
                errNo  =c_ParseErr
                return
            }

            msg = lmp.msgFactory.createMsgById( id )
            elog.LogSysln(" create msg type ", reflect.TypeOf( msg ))
            err = proto.Unmarshal( b, msg) 
            if err != nil {
                elog.LogSys(" protobuf parse  msg body fail ", err )
                errNo = c_ParseErr
                return 
            }
            errNo = c_NoError
            return 
        }

    } 
    elog.LogSysln( " recv msg len is big than buffer  " )
    errNo = c_NeedContinue
    id = 0
    return 
}


// create new LiteMsgParse

func NewLiteMsgParse( ) *LiteMsgParse  {

    lmp := &LiteMsgParse {
        msgCb : defMsgCallback,     
        errCb : defErrorCallback,
    }
    return lmp
}



/*******************************************************
    *           *       *                   * 
    *           *       *                   *
         len       id           body
                   ^             ^
                   |             |

                         -
                        len
********************************************************/

func ( lmp * LiteMsgParse )  Encode ( conn *Connection) error {
    
    // let go ...
    var err error
    elog.LogSysln( "recv msg :",conn.recvBuf.Size() )
    for conn.recvBuf.Size() >= int(c_MinMsgLen) {
        
         id, ret ,msg := lmp.parse( conn )
         elog.LogSysln(" parse msg type ", reflect.TypeOf( msg ))
         if ret == c_NoError {
            //call back 
            lmp.msgCb( conn, id, msg )                 
         } else if( ret == c_NeedContinue ) {
            return  err
         } else {
            lmp.errCb(conn);
            err = errors.New( "read err " )
            return err
         }
    }
    return err
}

func ( lmp *LiteMsgParse ) Decode( conn *Connection, id int32, msg proto.Message  ) error   {
    
    
    var err error
    err = binary.Write( &conn.sendBuf, binary.LittleEndian, &id ) 
    if err != nil {
        elog.LogSysln(" write id  error", err)
        return err
    }
    
    var buf []byte
    buf ,err = proto.Marshal(msg)
    if err != nil {
        elog.LogSysln(" parse msg body  error", err)
        return err
    }
    err = binary.Write( &conn.sendBuf, binary.LittleEndian, int32(len(buf)) ) 
    if err != nil {
        elog.LogSysln(" write msg len error", err)
        return err
    }
    err = binary.Write( &conn.sendBuf, binary.LittleEndian, buf ) 
    if err != nil {
        elog.LogSysln(" write msg body error", err)
        return err
    }
    return err  
}

func ( lmp *LiteMsgParse ) SetMsgCb( cb MSG_CALLBACK ) {

    lmp.msgCb = cb
}

func ( lmp *LiteMsgParse ) SetErrCb( cb ERROR_CALLBACK ) {
    lmp.errCb = cb
}

func ( lmp *LiteMsgParse ) SetMsgFactory( factory IMsgFactory ) {
    lmp.msgFactory = factory;
}










