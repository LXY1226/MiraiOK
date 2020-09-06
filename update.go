package main

import (
	"io/ioutil"
	"os"
	"runtime"
	"time"
)

func doUpdate() {
	checkJava()
	go updateSelf()
	update := true
	_, err := os.Stat(".NOUPDATE")
	if err == nil {
		update = false
	}
	if !checkLibs(parseLibs()) {
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
			updateMirai(false)
		} else {
			INFO("删除.lastupdate来在下次启动时强制检查更新")
		}
	}
}

func updateMirai(force bool) {
	INFO("开始检查Mirai更新... 也可以通过创建.noupdate文件来禁用更新")
	os.MkdirAll(libDIR, 0755)
	if !save(downFile("shadow/latest.txt"), libDIR+"version.txt") {
		ERROR("无法下载Mirai版本信息")
		return
	}
	libs := parseLibs()
	if libs == nil {
		ERROR("mirai-repo出错... [libs == nil] 请截图联系miraiOK")
	}
	os.MkdirAll(libDIR[:len(libDIR)-1], 0755)
	globalWG.Add(len(libs))
	for _, li := range libs {
		go func(li lib) {
			if _, err := os.Stat(li.LibPath()); !force && err == nil {
				goto done
			}
			INFO("下载", li.name, "版本", li.version)
			save(downFile("shadow/"+li.Path()), li.LibPath())
		done:
			globalWG.Done()
		}(li)
	}
	globalWG.Wait()
	_ = ioutil.WriteFile(".lastupdate", nil, 0755)
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
