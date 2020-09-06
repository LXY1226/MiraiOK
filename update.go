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
	if err != nil && os.IsNotExist(err) {
		update = false
	}
	data, err := ioutil.ReadFile(libDIR + "version.txt")
	if err != nil {
		INFO("读取库列表出现错误", err)
		if !update {
			update = true
		}
	} else {
		libs = parseLibs(data)
	}
	if !checkLibs() {
		update = true
	}
	if libs != nil && len(libs) != 0 && !update {
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
	rb := downFile("shadow/latest.txt")
	if rb == nil {
		ERROR("无法下载Mirai版本信息")
		return
	}
	data, _ := ioutil.ReadAll(rb)
	libs = parseLibs(data)
	if libs == nil {
		ERROR("mirai-repo出错...")

	}
	os.MkdirAll(libDIR[:len(libDIR)-1], 0666)
	err := ioutil.WriteFile(libDIR+"version.txt", data, 0777)
	if err != nil {
		ERROR("无法写入库列表", err)
	}
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
	_ = ioutil.WriteFile(".lastupdate", nil, 0777)
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
