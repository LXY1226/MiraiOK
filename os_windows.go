package main

import (
	"os"
	"syscall"
)

const CPSEP = ";"

var colorFunc *syscall.Proc
var fd = os.Stdout.Fd()
var defaultColor uintptr

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

func colorReset() {
	colorFunc.Call(fd, uintptr(0x7))
}

func colorWARN() {
	colorFunc.Call(fd, uintptr(0xc)) // Light Red
}

func colorERROR() {
	colorFunc.Call(fd, uintptr(0x4)) // Red
}
