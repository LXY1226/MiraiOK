package main

import (
	"bufio"
	"bytes"
	"gitee.com/LXY1226/logging"
	"io"
	"io/ioutil"
	"os"
)

const configTemplate = "#以#号开头的行会被忽略，所有有效配置会被使用\n" +
	"#请不要在各种地方填空格，空格也会当作变量的一部分\n" +
	"#QQ号,密码\n\n" +
	"#123456789,TestMiraiOK"

type loginCommands struct {
	buf   []byte
	drain bool
}

var loginCommand = loginCommands{
	buf:   make([]byte, 256)[:0],
	drain: false,
}

func (loginCommands) Read(p []byte) (n int, err error) {
	if !loginCommand.drain {
		p = append(p[:0], loginCommand.buf...)
		//p = loginCommand.buf
		loginCommand.drain = true
		return len(p), nil
	} else {
		return os.Stdin.Read(p)
	}
}

func (loginCommands) WriteAccount(QQ, Password string) {
	logging.INFO("用户", QQ, "已载入")
	loginCommand.buf = append(loginCommand.buf, "login "...)
	loginCommand.buf = append(loginCommand.buf, QQ...)
	loginCommand.buf = append(loginCommand.buf, ' ')
	loginCommand.buf = append(loginCommand.buf, Password...)
	loginCommand.buf = append(loginCommand.buf, '\n')
}

func loadConfig() {
	f, err := os.Open("config.txt")
	if err != nil {
		logging.WARN("读取config.txt失败:", err.Error())
		logging.INFO("配置文件用于快速登录Mirai，打开config.txt，将对应部分替换即可")
		logging.INFO("Hit: 您也可以创建 .noupdate 文件来禁用更新检查")
		err := ioutil.WriteFile("config.txt", []byte(configTemplate), 0755)
		if err != nil {
			logging.WARN("无法写入配置文件模板:", err.Error())
		}
		return
	}
	rd := bufio.NewReader(f)
	for {
		line, _, err := rd.ReadLine()
		if err != nil {
			if err == io.EOF {
				logging.INFO("查看config.txt来添加自动登录")
				return
			}
			logging.WARN("读取config失败:", err.Error())
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
		loginCommand.WriteAccount(string(line[:i]), string(line[i+1:]))
	}
}
