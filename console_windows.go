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
	colorFunc.Call(fd, uintptr(0xa)) // Green
}

func colorWARN() {
	colorFunc.Call(fd, uintptr(0xc)) // Light Red
}

func colorERROR() {
	colorFunc.Call(fd, uintptr(0x4)) // Red
}
