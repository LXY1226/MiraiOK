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
	"github.com/kardianos/osext"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

var javaPath = "./jre/bin/java"
var global = sync.WaitGroup{}
var args = []string{"-jar", "", "--update", "keep"}
var doUpdate = true

var arg0 string

func main() {
	defer func() {
		if err := recover(); err != nil {
			logging.FATAL(err.(error).Error())
		}
		var str string
		println("按回车退出...")
		println("请尝试清空文件，重新下载此程序")
		_, _ = fmt.Scan(&str)
	}()
	readConfig()
	_ = ioutil.WriteFile("content/.wrapper.txt", []byte("Pure"), 0755)
	if _, err := os.Stat("content"); err != nil {
		err = os.MkdirAll("content", 0755)
		if err != nil {
			logging.ERROR("无法创建content目录", err.Error())
			return
		}
	}
	arg0, _ = osext.Executable()
	checkJava()
	_, err := os.Open(".noupdate")
	if checkWrapper(); doUpdate {
		inf, err := os.Stat(".lastupdate")
		if err != nil || time.Now().Sub(inf.ModTime()) > time.Hour {
			go updateSelf()
			updateMirai()
		} else {
			logging.INFO("删除.lastupdate来在下次强制检查更新")
		}
	}
	global.Wait()
	if args[1] == "" {
		logging.ERROR("Mirai本体下载失败，准备退出...")
		return
	}
	logging.DEBUG(args...)
	go noStop()
	cmd := exec.Command(javaPath, args...)
	//cmd.Env = append(cmd.Env, "JAVAPATH=xxxxx")
	cmd.Stdout = console
	cmd.Stderr = console
	cmd.Stdin = loginCommand
	logging.INFO("启动Mirai...")
	_ = os.Remove(arg0 + ".old")
	//time.Sleep(time.Second) // 给用户看上面介绍的时间
	err = cmd.Run()
	if err != nil {
		logging.ERROR("运行失败，尝试更新mirai三件套", err.Error())
		updateMirai()
		err = cmd.Run()
		if err != nil {
			logging.ERROR("无法启动", err.Error())
		}
	}
}

func checkWrapper() {
	cur, _ := os.Open(".")
	list, _ := cur.Readdirnames(-1)
	for _, name := range list {
		if strings.HasPrefix(name, "mirai-console-wrapper") {
			args[1] = name
			return
		}
	}
	doUpdate = true
}

func updateSelf() {
	rb := downFile("mirai/MiraiOK/.version")
	if rb == nil {
		logging.ERROR("无法下载MiraiOK版本信息")
		return
	}
	data, _ := ioutil.ReadAll(rb)

	ver := string(data[:15])
	if ver != BUILDTIME {
		logging.INFO("发现新版本", ver)
		err := os.Rename(arg0, arg0+".old")
		if err != nil {
			logging.ERROR("重命名失败", err.Error())
			return
		}
		url := "mirai/MiraiOK/miraiOK_" + runtime.GOOS + "_" + runtime.GOARCH
		if runtime.GOOS == "windows" {
			url += ".exe"
		}
		if !save(downFile(url), arg0) {
			logging.ERROR("程序下载失败，取消自更新")
			os.Rename(arg0+".old", arg0)
		}
	}
}
