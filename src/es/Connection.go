/*
 *  CopyRight  2014 , bingyu.song   All Right Reserved
    I believe  Spring brother
 */


package es

import (
    "net"
    "bytes"

    "elog"
)

const defReadBufSize  = 1024*1024
const defWriteBufSize = 1024*1024

type Connection struct {
    tcpConn         *net.TCPConn
    readBuf         *bytes.Buffer
    writeBuf        *bytes.Buffer
    eventDispatch   *EventDispatch
    readBufSize     int
    writeBufSize    int
}

func NewConn( conn *net.TCPConn ,dispatch *EventDispatch )( connection *Connection, err error ) {
    
    connection = &Connection { 
            tcpConn     : conn,
            eventDispatch: dispatch,
        }


    err = nil

    elog.LogSys("create New conn ")
    
    ne := new ( netEvent)
    ne.eventType = NEW;
    ne.conn = connection;
    connection.eventDispatch.pushEvent( ne )
    return
}

func (c *Connection ) handleEvent( ) {
 
    for {

        data := make([] byte, defReadBufSize )
        ne := new ( netEvent)
        ne.conn = c
        _, err := c.tcpConn.Read( data )
        if err != nil {
            ne.eventType = CLOSE;
            c.eventDispatch.pushEvent( ne )
            break;
        }
        
        ne.eventType = READ;
        c.eventDispatch.pushEvent( ne )
    }
    
}

















