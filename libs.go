package main

import (
	"bytes"
	"io/ioutil"
	"os"
)

type lib struct {
	name, version string
}

func checkLibs(libs []lib) bool {
	if libs == nil || len(libs) == 0 {
		return false
	}
	for _, lib := range libs {
		_, err := os.Stat(lib.LibPath())
		if err != nil {
			WARN("检查lib：", lib.LibPath(), "出错 ", err)
			return false
		}
	}
	return true
}

func parseLibs() []lib {
	data, err := ioutil.ReadFile(libDIR + "version.txt")
	if err != nil {
		WARN("读取库列表失败", err)
		return nil
	}
	libs := make([]lib, 4)[:0]
	pos := 0
	classpath = "CLASSPATH="
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
		classpath += l.LibPath() + CPSEP
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
