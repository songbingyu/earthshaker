package main

import (
        "fmt"
        "net"
        //"time"
        "bufio"
        "os"
        "bytes"
       "encoding/binary"
        "netpb"
        //"unsafe"
        proto"code.google.com/p/goprotobuf/proto"
       )

func main() {
          conn, err := net.Dial("tcp", "127.0.0.1:6798")
          checkError(err)
          buf := new ( bytes.Buffer )
          var  id  int32 = 1
          
          msg := &netpb.NetMsg{ 
           
          }

          msg.M_ID = new( netpb.MsgID  )
          *msg.M_ID = netpb.MsgID_EM_CS_LOGIN 
          
          buffer, _  := proto.Marshal(msg)

          var  l int32 = int32(4 +len(buffer))
          
          binary.Write( buf,binary.LittleEndian, l )
          binary.Write( buf,binary.LittleEndian, id )
          binary.Write( buf,binary.LittleEndian, buffer)
          
          fmt.Println(buf.Bytes(), l, buf.Len() )
          reader := bufio.NewReader(os.Stdin)
          _, err = conn.Write( buf.Bytes() )
          for {
                  
                  _, err := reader.ReadString('\n')
                  fmt.Println(err)
                  if err != nil {
                      conn.Close()
                          break

                  }
          }
}
func checkError(err error) {
    if err != nil {
        fmt.Println("Fatal error ", err.Error())
    }

}
