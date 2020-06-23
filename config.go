package main

import (
	"gitee.com/LXY1226/logging"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const configTemplate = "#DEBUG\n#NOUPDATE\n#以#号开头的行会被忽略，所有有效配置会被使用\n" +
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

func readConfig() {
	logging.Log2Con = logging.LogINFO
	logging.Log2Log = logging.LogERROR
	data, err := ioutil.ReadFile("config.txt")
	var lines []string
	var isOldConf = true // 会在一个月后[2020/07/23]移除
	if err == nil {
		lines = strings.Split(string(data), "\n")
		for _, line := range lines {
			if line == "DEBUG" {
				logging.Log2Con = logging.LogDEBUG
				logging.Log2Log = logging.LogDEBUG
				isOldConf = false
				break
			} else if line == "#DEBUG" {
				isOldConf = false
				break
			}
		}
	}
	logging.Init("MiraiOK", BUILDTIME)
	logging.INFO("此程序以Affero GPL3.0协议发布，使用时请遵守协议")
	logging.DEBUG("如你所愿开启MiraiOK调试输出")
	if err == nil {
		lines = strings.Split(string(data), "\n")
		for _, line := range lines {
			switch true {
			case len(line) == 0:
				continue
			case line[0] == '#':
				continue
			case line == "NOUPDATE":
				logging.INFO("如你所愿不会更新Mirai本体")
				doUpdate = false
			}
			i := strings.Index(line, ",")
			if i == -1 {
				i = strings.Index(line, "，")
				i++
			}
			if i > 0 {
				loginCommand.WriteAccount(line[:i], line[i+1:])
			}
		}
		if isOldConf {
			logging.INFO("检测到旧版配置，将会写入新配置项，建议查看文件内的相关改动")
			err := ioutil.WriteFile("config.txt", append([]byte("#DEBUG\n#NOUPDATE\n"), data...), 0755)
			if err != nil {
				logging.WARN("无法写入新配置文件模板:", err.Error())
			}
			logging.INFO("五秒后继续...")
			time.Sleep(5 * time.Second)
		}
	} else {
		logging.WARN("读取config.txt失败:", err.Error())
		err = ioutil.WriteFile("config.txt", []byte(configTemplate), 0755)
		if err != nil {
			logging.WARN("无法写入配置文件模板:", err.Error())
		}
	}
}

//func loadConfig() {
//	f, err := os.Open("config.txt")
//	if err != nil {
//		logging.WARN("读取config.txt失败:", err.Error())
//		logging.INFO("配置文件用于快速登录Mirai，打开config.txt，将对应部分替换即可")
//		err := ioutil.WriteFile("config.txt", []byte(configTemplate), 0755)
//		if err != nil {
//			logging.WARN("无法写入配置文件模板:", err.Error())
//		}
//		return
//	}
//	rd := bufio.NewReader(f)
//	for {
//		line, _, err := rd.ReadLine()
//		if err != nil {
//			if err == io.EOF {
//				logging.INFO("查看config.txt来添加自动登录")
//				return
//			}
//			logging.WARN("读取config失败:", err.Error())
//			return
//		}
//		if len(line) == 0 || line[0] == '#' {
//			continue
//		}
//		i := bytes.IndexByte(line, ',')
//		if i == -1 {
//			i = bytes.Index(line, []byte("，"))
//			if i == -1 {
//				continue
//			}
//			i++
//		}
//		loginCommand.WriteAccount(string(line[:i]), string(line[i+1:]))
//	}
//}
