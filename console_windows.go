package main

import (
	"os"
	"syscall"
)

var colorFunc *syscall.Proc
var fd = os.Stdout.Fd()

func init() {
	dll, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		return
	}
	colorFunc, err = dll.FindProc("SetConsoleTextAttribute")
	if err != nil {
		return
	}
}

func colorINFO() {
	colorFunc.Call(os.Stdout.Fd(), uintptr(0xa)) // GREEN
}

func colorWARN() {
	colorFunc.Call(os.Stdout.Fd(), uintptr(0x4)) // RED
}

func colorERROR() {
	colorFunc.Call(os.Stdout.Fd(), uintptr(0x4)) //?
}
