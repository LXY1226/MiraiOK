//+build !windows

package main

import "os"

const CPSEP = ":"

func colorINFO() {
	os.Stdout.Write([]byte{27, '[', '9', '2', 'm'})
}

func colorWARN() {
	os.Stdout.Write([]byte{27, '[', '9', '1', 'm'})
}

func colorERROR() {
	os.Stdout.Write([]byte{27, '[', '3', '1', 'm'})
}

func colorReset() {
	os.Stdout.Write([]byte{27, '[', '0', 'm'})
}
