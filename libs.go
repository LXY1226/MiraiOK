package main

import (
	"bytes"
	"os"
)

type lib struct {
	name, version string
}

func checkLibs() bool {
	for _, lib := range libs {
		_, err := os.Stat(lib.LibPath())
		if err != nil {
			WARN("检查lib：", lib.LibPath(), "出错 ", err)
			return false
		}
	}
	return true
}

func parseLibs(data []byte) []lib {
	libs := make([]lib, 4)[:0]
	pos := 0
	for {
		pos = bytes.IndexByte(data, '\n')
		if pos == -1 {
			return libs
		}
		line := data[:pos]
		data = data[pos+1:]
		pos = bytes.IndexByte(line, ':')
		if pos == -1 {
			WARN("jar列表解析出错")
			return nil
		}
		var l lib
		l.name = string(line[:pos])
		if line[pos+1] == ' ' {
			pos++
		}
		l.version = string(line[pos+1:])
		libs = append(libs, l)
	}
}

func (l lib) Path() string {
	return l.name + "/" + l.fName()
}

func (l lib) LibPath() string {
	return libDIR + l.fName()
}

func (l lib) fName() string {
	return l.name + "-" + l.version + ".jar"
}
