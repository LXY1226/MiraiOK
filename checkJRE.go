//+build !386

package main

import (
	"bytes"
	"gitee.com/LXY1226/logging"
	"os"
	"os/exec"
	"strings"
)

func checkJava() {
	defer wg.Done()
	//检测本地java
	if checkJavaBin() {
		return
	}
	f, err := exec.LookPath("java")
	if err != nil {
		logging.INFO("未发现JRE，准备下载...")
		if unpackRAR(downURL("jre-" + logging.RTStr + ".rar")) {
			return
		}
	} else {
		javaPath = f
	}
	if checkJavaBin() {
		return
	}
	logging.FATAL("无法获取JRE，即将退出...")
	os.Exit(0)
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
