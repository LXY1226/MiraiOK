package main

/*
 * Copyright 2020 LXY1226, Mamoe Technologies and contributors.
 *
 * 采用与mirai相同的LICENSE
 * 此源代码的使用受 GNU AFFERO GENERAL PUBLIC LICENSE version 3 许可证的约束, 可以在以下链接找到该许可证.
 * Use of this source code is governed by the GNU AGPLv3 license that can be found through the following link.
 *
 * https://github.com/mamoe/mirai/blob/master/LICENSE
 */

import (
	"fmt"
	"gitee.com/LXY1226/logging"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

var javaPath = "./jre/bin/java"
var wg = sync.WaitGroup{}
var args []string

var argc = 1

func main() {
	defer func() {
		if err := recover(); err != nil {
			logging.FATAL(err.(error).Error())
		}
		var str string
		println("按回车退出...")
		println("尝试访问gitee.com/LXY1226/MiraiOK试试更新一下？")
		_, _ = fmt.Scan(&str)
	}()
	if _, err := os.Stat(".DEBUG"); err == nil {
		logging.Log2Con = logging.LogDEBUG
		logging.Log2Log = logging.LogDEBUG
	} else {
		logging.Log2Log = logging.LogFATAL
		logging.Log2Log = logging.LogINFO
	}
	logging.Init("MiraiOK", BUILDTIME)
	logging.INFO("此程序以Affero GPL3.0协议发布，使用时请遵守协议")
	loadConfig()
	args = append(args, "-jar", "mirai-console-wrapper-1.1.0.jar") // Not Real jar path
	_ = ioutil.WriteFile("content/.wrapper.txt", []byte("Pure"), 0755)
	args = append(args, "--update", "keep")
	if _, err := os.Stat("content"); err != nil {
		err = os.MkdirAll("content", 0755)
		if err != nil {
			logging.ERROR("无法创建目录content:", err.Error())
			return
		}
	}
	wg.Add(1)
	go checkJava()
	_, err := os.Open(".noupdate")
	if checkWrapper(); args[argc] != "" || err != nil { //Wrapper存在且无noupdate
		inf, err := os.Stat(".lastupdate")
		if err != nil || time.Now().Sub(inf.ModTime()) > time.Hour {
			updateMirai()
		} else {
			logging.INFO("删除.lastupdate来在下次强制检查更新")
		}
	}
	wg.Wait()
	if args[argc] == "" {
		logging.ERROR("有一个或多个文件无法获取，即将退出...")
		return
	}
	logging.DEBUG(args...)
	hupChan := make(chan os.Signal)
	signal.Notify(hupChan, syscall.SIGHUP)
	go func() {
		for {
			_ = <-hupChan
			println("Mirai将被挂起但是不会停止运行")
		}
	}()
	cmd := exec.Command(javaPath, args...)

	cmd.Stdout = console
	cmd.Stderr = console
	cmd.Stdin = os.Stdin
	logging.INFO("启动Mirai...")
	//time.Sleep(time.Second) // 给用户看上面介绍的时间
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}

func checkWrapper() {
	cur, _ := os.Open(".")
	list, _ := cur.Readdirnames(-1)
	for _, name := range list {
		if strings.HasPrefix(name, "mirai-console-wrapper") {
			args[argc] = name
			return
		}
	}
	args[argc] = ""
}
