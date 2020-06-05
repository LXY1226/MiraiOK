package main

import (
	"bytes"
	"gitee.com/LXY1226/logging"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var repos = []string{
	"http://t.imlxy.net:64724/mirai-repo/shadow/",
	"https://raw.githubusercontent.com/mamoe/mirai-repo/master/shadow/",
	"https://gitee.com/LXY1226/mirai-repo/raw/master/shadow/",
	"https://mamoe.github.io/mirai-repo/shadow/",
}

const (
	WrapperName = "mirai-console-wrapper"
	ConsoleName = "mirai-console"
	CoreName    = "mirai-core"
	CoreJarName = "mirai-core-qqandroid"
)

func updateWrapper() {
	defer wg.Done()
	info := getInfoFromGithub(WrapperName)
	if info.TagName == "" {
		return
	}
	fname := WrapperName + "-" + info.TagName + ".jar"
	if _, err := os.Stat(fname); err == nil {
		return
	}
	logging.INFO("更新", WrapperName, "到", info.TagName, "发布于", info.PublishedAt.String())
	for _, repo := range repos {
		if save(downURL(repo+WrapperName+"/"+fname), fname) {
			args[argc] = fname
			return
		}
	}
	logging.ERROR("下载", fname, "失败")
}

func updateConsole() {
	defer wg.Done()
	info := getInfoFromMaven(ConsoleName)
	if info.TagName == "" {
		return
	}
	fname := ConsoleName + "-" + info.TagName + ".jar"
	if _, err := os.Stat("content/" + fname); err == nil {
		return
	}
	logging.INFO("更新", ConsoleName, "到", info.TagName, "发布于", info.PublishedAt.String())
	for _, repo := range repos {
		if save(downURL(repo+ConsoleName+"/"+fname), "content/"+fname) {
			return
		}
	}
	logging.ERROR("下载", fname, "失败")
}

func updateCore() {
	defer wg.Done()
	info := getInfoFromMaven(CoreName)
	if info.TagName == "" {
		return
	}
	fname := CoreJarName + "-" + info.TagName + ".jar"
	if _, err := os.Stat("content/" + fname); err == nil {
		return
	}
	logging.INFO("更新", CoreName, "到", info.TagName, "发布于", info.PublishedAt.String())
	for _, repo := range repos {
		if save(downURL(repo+CoreJarName+"/"+fname), "content/"+fname) {
			return
		}
	}
	logging.ERROR("下载", fname, "失败")
}

func getInfoFromGithub(projectName string) (info GHinfo) {
	resp, err := http.Get("https://api.github.com/repos/mamoe/" + projectName + "/releases/latest")
	if err != nil {
		logging.ERROR("从Github获取", projectName, "版本信息失败：", err.Error())
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logging.ERROR("从Github读取", projectName, "版本信息失败：", err.Error())
		return
	}
	err = json.Unmarshal(body, &info)
	if err != nil {
		logging.ERROR("从Github解析", projectName, "版本信息失败：", err.Error())
		return
	}
	return
}

func getInfoFromMaven(projectName string) (info GHinfo) {
	resp, err := http.Get("https://mirrors.huaweicloud.com/repository/maven/net/mamoe/" + projectName + "/maven-metadata.xml")
	if err != nil {
		logging.ERROR("从Maven获取", projectName, "版本信息失败：", err.Error())
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logging.ERROR("从Maven读取", projectName, "版本信息失败：", err.Error())
		return
	}
	a := bytes.Index(data, []byte("<latest>"))
	b := bytes.Index(data, []byte("</latest>"))
	c := bytes.Index(data, []byte("<lastUpdated>"))
	d := bytes.Index(data, []byte("</lastUpdated>"))
	if (a != -1 || b != -1 || a < b) && (c != -1 || d != -1 || c < d) {
		info.TagName = string(data[a+8 : b])
		info.PublishedAt, _ = time.Parse("20060102150405", string(data[c+13:d]))
	} else {
		logging.ERROR("从Maven解析", projectName, "版本信息失败")
	}
	return
}

type GHinfo struct {
	TagName     string    `json:"tag_name"`
	PublishedAt time.Time `json:"published_at"`
	Body        string    `json:"body"`
	Assets      []struct {
		Name string `json:"name"`
		URL  string `json:"browser_download_url"`
	} `json:"assets"`
}
