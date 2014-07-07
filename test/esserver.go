package main

import (
    "fmt"
    "elog"
    "es"
    "os"
    "os/signal"
    "syscall"
)


func main() {


    elog.InitLog( elog.INFO )
    tcpServer ,err := es.NewTcpServer()
    if err != nil {
        elog.LogInfo("crate tcp socket fail ")
    }
    
    ch := make(chan os.Signal)
    signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
    
    fmt.Println(" begin listen....")
    tcpServer.Listen("0.0.0.0:6798")

    tcpServer.Run()
    fmt.Println(" begin listen....")

            
    fmt.Println(<-ch)

    tcpServer.Exit()

    elog.LogSys("Hhhhhh");


}
