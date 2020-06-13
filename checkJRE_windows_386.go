package main

import (
	"bytes"
	"gitee.com/LXY1226/logging"
	"os"
	"os/exec"
	"strings"
)

func checkJava() {
	defer global.Done()
	//检测本地java
	if checkJavaBin() {
		return
	}
	unpackRAR(downFile("mirai-repo/shadow/jre-" + logging.RTStr + ".rar"))
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
