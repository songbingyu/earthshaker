/*
 *  CopyRight  2014 , bingyu.song   All Right Reserved
    I believe  Spring brother
 */


package util

import (
        "os"
        "os/signal"
        "syscall"
        "elog"
)

type signalHandler func ( s os.Signal )

type signalHelper struct {

    signMap map[os.Signal]signalHandler
    ch      chan os.Signal
}


var sh signalHelper

type ( s *signalHelper ) handler( s os.Signal )( err error ){

    if handler, ok :=  s.signMap[s]; ok != nil {
       
        handler( sig )
        return nil
    }else {
        return fmt.Errorf(" not add sign handler %d ",sig )
    }
}

func  AddSign( s os.Signal ,signalHandler ) {

    //TODO: check
    sh.signMap[ s ] = singleHandler;
}

func SignProc() {

    go func() {
        for{

            sh.ch = make( chan os.Signal )
            var sigs [] os.Signal
            for sig := range sh.signMap {
                sigs = append( sigs, sig )
            }

            signal.Notify( sh.ch , sigs)

            sig := <- sh.ch
            err := sh.handler( sig )
            if err  != nil {
                elog->LogSysln( err )
            }
            return 

        }

    }()

    
}




 









