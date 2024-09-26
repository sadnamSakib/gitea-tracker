package service

import "runtime"

var (
	goRoutines = runtime.NumCPU() * 2
)
