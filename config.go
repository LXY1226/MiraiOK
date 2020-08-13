package main

import (
	"bytes"
	"gitee.com/LXY1226/logging"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const configTemplate = "#DEBUG\n#NOUPDATE\n" +
	"#在----------下面可以添加需要在每次启动时输入得指令\n" +
	"#请注意，指令部分中#并不起效，miraiOK会原样输入到console\n" +
	"例如:\n" +
	"login 123456789 TestMiraiOK\n" +
	"say 655057127 MiraiOK_published!\n" +
	"----------\n"

//type loginCommands struct {
//	buf   []byte
//	drain bool
//}

type loginCommands []byte

var loginCommand loginCommands

func (loginCommands) Read(p []byte) (n int, err error) {
	if loginCommand != nil {
		n = len(loginCommand)
		copy(p, loginCommand)
		loginCommand = nil
		return n, nil
	} else {
		return os.Stdin.Read(p)
	}
}

func readConfig() {
	logging.Log2Con = logging.LogINFO
	logging.Log2Log = logging.LogERROR
	data, err := ioutil.ReadFile("config.txt")
	var isOldConf = true // 会在一个月后[2020/07/28]移除
	if err == nil {
		pos := 0
		next := 0
	preConf:
		for {
			next = bytes.IndexByte(data[pos:], '\n')
			if next == -1 {
				break
			}
			next += pos
			if data[pos] == '#' {
				goto end
			}
			switch string(data[pos:next]) {
			case "DEBUG":
				logging.Log2Con = logging.LogDEBUG
				logging.Log2Log = logging.LogDEBUG
			case "NOUPDATE":
				doUpdate = false
			case "----------":
				isOldConf = false
				loginCommand = data[next+1:]
				break preConf
			}
		end:
			pos = next + 1

		}
	}
	logging.Init("MiraiOK", BUILDTIME)
	logging.INFO("此程序以Affero GPL3.0协议发布，使用时请遵守协议")
	logging.INFO("代码库: github.com/LXY1226/MiraiOK gitee.com/LXY1226/MiraiOK")
	logging.DEBUG("如你所愿开启MiraiOK调试输出")
	if !doUpdate {
		logging.INFO("如你所愿不会更新Mirai本体")
	}
	if err == nil {
		if isOldConf {
			loginCommand = make([]byte, 256)[:0]
			logging.INFO("检测到旧版配置，备份至config.bak.txt")
			logging.INFO("账号会迁移至新config.txt，建议查看文件内的相关改动")
			for _, line := range strings.Split(string(data), "\n") {
				if len(line) == 0 || line[0] == '#' {
					continue
				}
				i := strings.Index(line, ",")
				if i == -1 {
					i = strings.Index(line, "，")
					i++
				}
				if i > 0 {
					logging.INFO("用户", line[:i], "已迁移")
					loginCommand = append(loginCommand, "login "...)
					loginCommand = append(loginCommand, line[:i]...)
					loginCommand = append(loginCommand, ' ')
					loginCommand = append(loginCommand, line[i+1:]...)
					loginCommand = append(loginCommand, '\n')
				}
			}
			if err := os.Rename("config.txt", "config.bak.txt"); err != nil {
				logging.WARN("无法移动config.txt至config.bak.txt", err.Error())
			} else if err := ioutil.WriteFile("config.txt", append([]byte(configTemplate), loginCommand...), 0755); err != nil {
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
