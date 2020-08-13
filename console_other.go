//+build !windows

package main

import (
	"os"
	"os/signal"
	"syscall"
)

var console = os.Stdout

func noStop() {
	hupChan := make(chan os.Signal)
	signal.Notify(hupChan, syscall.SIGHUP, syscall.SIGTSTP)
	for {
		_ = <-hupChan
		println("Mirai将被挂起但是不会停止运行")
	}
}
