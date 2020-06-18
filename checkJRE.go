//+build !386

package main

import (
	"bytes"
	"gitee.com/LXY1226/logging"
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
	logging.INFO("未发现JRE，准备下载...")
	initStor()
	global.Add(1)
	go func() {
		if unpackRAR(downFile("mirai-repo/shadow/jre-" + logging.RTStr + ".rar")) {
			if checkJavaBin() {
				global.Done()
				return
			}
		}
		logging.FATAL("无法获取JRE，即将退出...")
		panic("error in gathering JRE")
	}()
}

func checkJavaBin() bool {
	var stdo bytes.Buffer
	logging.DEBUG("Trying Locating JRE:", javaPath)
	cmd := exec.Command(javaPath, "-version")
	cmd.Stdout = &stdo
	cmd.Stderr = &stdo
	err := cmd.Run()
	if err != nil {
		return false
	}
	for str, err := stdo.ReadString('\n'); err == nil; {
		logging.INFO("JRE:", strings.TrimRight(str, "\r\n"))
		str, err = stdo.ReadString('\n')
	}
	return true
}
