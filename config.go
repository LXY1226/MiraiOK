package main

import (
	"bufio"
	"bytes"
	"gitee.com/LXY1226/logging"
	"io/ioutil"
	"os"
)

func loadConfig() {

	f, err := os.Open("config.txt")
	if err != nil {
		logging.WARN("读取config.txt失败:", err.Error())
		logging.INFO("配置文件用于快速登录Mirai，打开config.txt，将对应部分替换即可")
		logging.INFO("Hit: 您也可以创建 .noupdate 文件来禁用更新检查")
		err := ioutil.WriteFile("config.txt", []byte("#以#号开头的行会被忽略，只有第一个有效配置会被使用\n#请不要在各种地方填空格，空格也会当作变量的一部分\n#QQ号,密码\n\n#123456789,TestMiraiOK"), 0755)
		if err != nil {
			logging.WARN("无法写入配置文件模板:", err.Error())
		}
		return
	}
	rd := bufio.NewReader(f)
	for {
		line, _, err := rd.ReadLine()
		if err != nil {
			logging.WARN("读取config失败:", err.Error())
			return
		}
		if line == nil {
			logging.INFO("config中无有效配置")
			return
		}
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		i := bytes.IndexByte(line, ',')
		if i == -1 {
			i = bytes.Index(line, []byte("，"))
			if i == -1 {
				continue
			}
			i++
		}
		args = append(args, "-Dmirai.account="+string(line[:i]))
		args = append(args, "-Dmirai.password="+string(line[i+1:]))
		argc += 2
		logging.INFO("配置已载入")
		return
	}

}
