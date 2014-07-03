
package  main

import (
    "fmt"
    "util"
    "os"
    "sync"
)

func test ( format string, a ...interface{} ){

   fmt.Fprintf(os.Stdout, format, a... ) // 00100000
}

func logtest(ch   chan int, done *sync.WaitGroup ) {

    for i:=0; i < 20000; i++  {
        util.LogInfo("llllllI:%d", i)
    }
    ch <- 1
    done.Done()
}
func main() {

   fmt.Println(" Let go.... ")
   fmt.Println(" Let go....%d ", 1)
   util.InitLog( util.INFO )

   util.LogInfo("just  a test:%d",  3 );
   util.LogInfo("............");
   util.Flush()
    test("sssss%d", 33)
   fmt.Printf("%08b\n", 32)             // 00100000

   chs := make( [] chan int, 1000)

   var done sync.WaitGroup
   for  i :=0; i < 1000; i++ {
       done.Add(1)
       chs[i] = make ( chan int)
       go logtest( chs[i], &done)
   }
   
   done.Wait()

   util.Flush()
}
