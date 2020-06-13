package main

import (
	"gitee.com/LXY1226/logging"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"os"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var repos = []string{
	"http://t.imlxy.net:64724/mirai-repo/shadow/",
	"https://gitee.com/LXY1226/mirai-repo/raw/master/shadow/",
	"https://raw.githubusercontent.com/mamoe/mirai-repo/master/shadow/",
	"https://mamoe.github.io/mirai-repo/shadow/",
}

var verinfos []Verinfo

func updateMirai() {
	logging.INFO("开始检查Mirai更新... 也可以通过创建.noupdate文件来禁用更新")
	rb := downFile("mirai-repo/shadow/release.json")
	if rb == nil {
		logging.ERROR("无法下载Mirai版本信息")
		return
	}
	data, _ := ioutil.ReadAll(rb)
	err := json.Unmarshal(data, &verinfos)
	if err != nil {
		logging.ERROR("无法解析Mirai版本信息:", err.Error())
		return
	}
	global.Add(3)
	for _, info := range verinfos {
		go func(info Verinfo) {
			defer global.Done()
			fname := info.Name + "-" + info.Version + ".jar"
			if _, err := os.Stat(info.Path + fname); err == nil {
				return
			}
			logging.INFO("正在更新", info.Name, "到", info.Version, "发布于", info.Date.String())
			save(downFile("mirai-repo/shadow/"+info.Name+"/"+fname), info.Path+fname)
		}(info)
	}
	global.Wait()
	_ = ioutil.WriteFile(".lastupdate", nil, 0755)
	checkWrapper()
}

type Verinfo struct {
	Name    string    `json:"name"`
	Date    time.Time `json:"date"`
	Version string    `json:"version"`
	Info    string    `json:"verinfos"`
	Path    string    `json:"path"`
}
