package earthshaker

import (
	"fmt"
	"os"
	"runtime/pprof"
)

type IniParam struct {
	Name string
	Loglevel int
}

type EarthShaker struct {
	name string
	loglevel int
}

var g_EarthShaker EarthShaker

func Ini(p IniParam) bool {
	fmt.Println("start ini...")

	g_EarthShaker.name = p.Name
	g_EarthShaker.loglevel = p.Loglevel

	// log
	IniLog(g_EarthShaker.loglevel, true, g_EarthShaker.name)

	// profile
	pf, err := os.Create(g_EarthShaker.name + ".prof")
	if err != nil {
		fmt.Println(err)
		return false
	}
	pprof.StartCPUProfile(pf)

	fmt.Println("ini ok")
	return true
}

func Exit() bool {
	fmt.Println("start exit...")

	pprof.StopCPUProfile()

	fmt.Println("exit ok")

	return true
}
