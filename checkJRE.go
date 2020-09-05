//+build !386

package main

import (
	"bytes"
	"os/exec"
	"strings"
)

func checkJava() {
	//检测本地java
	if checkJavaBin() {
		return
	}
	f, err := exec.LookPath("java")
	if err == nil {
		javaPath = f
		if checkJavaBin() {
			return
		}
	}
	INFO("未发现JRE，准备下载...")
	globalWG.Add(1)
	go func() {
		if unpackRAR(downFile("MiraiOK/jre-" + RTStr + ".rar")) {
			if checkJavaBin() {
				globalWG.Done()
				return
			}
		}
		ERROR("无法获取JRE，即将退出...")
		panic("error in gathering JRE")
	}()
}

func checkJavaBin() bool {
	var stdo bytes.Buffer
	cmd := exec.Command(javaPath, "-version")
	cmd.Stdout = &stdo
	cmd.Stderr = &stdo
	err := cmd.Run()
	if err != nil {
		return false
	}
	for str, err := stdo.ReadString('\n'); err == nil; {
		INFO(strings.TrimRight(str, "\r\n"))
		str, err = stdo.ReadString('\n')
	}
	return true
}
