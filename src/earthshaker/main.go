package main

import (
	"earthshaker"
	"fmt"
)

func main() {
	fmt.Println("enter main")
	earthshaker.Ini(earthshaker.IniParam{Name:"main", Loglevel:earthshaker.DEBUG})
	
	earthshaker.Exit()
}
