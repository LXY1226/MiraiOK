//+build !windows

package main

import "os"

func colorINFO() {
	os.Stdout.Write([]byte{0x1b, '[', '9', '2', 'm'})
}

func colorWARN() {
	os.Stdout.Write([]byte{0x1b, '[', '9', '1', 'm'})
}

func colorERROR() {
	os.Stdout.Write([]byte{0x1b, '[', '3', '1', 'm'})
}
