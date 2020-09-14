package main

import (
	"io/ioutil"
	"os"
	"runtime"
	"time"
)

func doUpdate() {
	checkJava()
	update := true
	_, err := os.Stat(".NOUPDATE")
	if err == nil {
		update = false
	}
	if !getLibs(parseLibs(), false) {
		if !update {
			WARN("检测到NOUPDATE但是无法解析库列表，强制更新")
			update = true
		}
	}
	if !update {
		INFO("跳过mirai更新")
	} else {
		inf, err := os.Stat(".lastupdate")
		if err != nil || time.Now().Sub(inf.ModTime()) > time.Hour {
			go updateSelf()
			syncLibs()
			getLibs(parseLibs(), false)
		} else {
			INFO("删除.lastupdate来在下次启动时强制检查更新")
		}
	}
}

func updateMirai(force bool) {

}

func updateSelf() {
	rb := downFile("MiraiOK/.version")
	if rb == nil {
		ERROR("无法下载MiraiOK版本信息")
		return
	}
	data, _ := ioutil.ReadAll(rb)
	ver := string(data[:15])
	if ver != BUILDTIME {
		INFO("发现新版本", ver)
		err := os.Rename(os.Args[0], arg0+".old")
		if err != nil {
			ERROR("重命名失败", err.Error())
			return
		}
		url := "MiraiOK/miraiOK_" + RTStr
		if runtime.GOOS == "windows" {
			url += ".exe"
		}
		if !save(downFile(url), arg0) {
			ERROR("程序下载失败，取消自更新")
			os.Rename(arg0+".old", arg0)
		}
	}
}
