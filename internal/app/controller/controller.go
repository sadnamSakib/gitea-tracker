package controller

import "runtime"

var (
	goRoutines = runtime.NumCPU() * 4
)
