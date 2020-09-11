package main

import (
	"bytes"
	"io/ioutil"
	"os"
)

type lib struct {
	name, version string
}

func getLibs(libs []lib, force bool) bool {
	_ = ioutil.WriteFile(".lastupdate", nil, 0755)
	if libs == nil || len(libs) == 0 {
		return false
	}
	for _, l := range libs {
		_, err := os.Stat(l.LibPath())
		if force || err != nil {
			WARN("获取lib：", l.LibPath())
			globalWG.Add(1)
			go func(l lib) {
				if save(downFile("shadow/"+l.Path()), l.LibPath()) {
					globalWG.Done()
				}
			}(l)
		}
	}
	globalWG.Wait()
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

func syncLibs() {
	os.MkdirAll(libDIR[:len(libDIR)-1], 0755)
	INFO("同步最新库列表...")
	os.MkdirAll(libDIR, 0755)
	if !save(downFile("shadow/latest.txt"), libDIR+"version.txt") {
		ERROR("无法下载库列表")
		return
	}
	libs := parseLibs()
	if libs == nil {
		ERROR("mirai-repo出错... (libs == nil) 请截图联系miraiOK")
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
