package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"gitee.com/LXY1226/logging"
	rar "github.com/nwaples/rardecode"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

const UA = "MiraiOK|" + BUILDTIME + "|" + logging.RTStr

var accessToken string

func save(br *bufio.Reader, fname string) bool {
	if br == nil {
		return false
	}
	f, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}
	_, err = br.WriteTo(f)
	if err != nil {
		panic(err)
	}
	_ = f.Close()
	return true
}

func unpackRAR(br *bufio.Reader) bool {
	if br == nil {
		return false
	}
	r, err := rar.NewReader(br, "")
	if err != nil {
		logging.ERROR("解压出错:", err.Error())
		return false
	}

	for f, err := r.Next(); err == nil; f, err = r.Next() {
		if f.IsDir {
			continue
		}
		dir, fname := path.Split(f.Name)
		print("解压：" + fname + "        \r")
		if _, err := os.Stat(dir); err != nil {
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				panic(err)
			}
		}
		f, err := os.OpenFile(f.Name, os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(f, r)
		if err != nil {
			panic(err)
		}
	}
	return true
}

func initStor() {
	if accessToken != "" {
		return
	}
	http.DefaultClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: caPool(),
		},
	}
	req, err := http.NewRequest("POST", torURL, strings.NewReader(tor))
	if err != nil {
		logging.WARN("初始化远程存储失败")
		return
	}
	req.Header.Set("User-Agent", ua)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logging.WARN("初始化远程存储失败")
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logging.WARN("初始化远程存储失败")
		return
	}
	accessToken = "Bearer " + dumpASToken(data)
	logging.INFO("初始化远程存储成功")
}

func downFile(path string) *bufio.Reader {
	if accessToken == "" {
		return nil
	}
	req, err := http.NewRequest("GET", dowURL+path+":/content", nil)
	if err != nil {
		logging.WARN("URL初始化失败", path)
		return nil
	}
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Authorization", accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logging.WARN("访问远程存储失败")
		return nil
	}
	if resp.StatusCode != 200 {
		logging.WARN("下载失败", path)
		return nil
	}
	return bufio.NewReaderSize(resp.Body, 1<<20)
}

func dumpASToken(data []byte) string {
	i := bytes.Index(data, []byte(`"access_token":"`))
	if i == -1 {
		return ""
	}
	data = data[i+16:]
	j := bytes.Index(data, []byte(`"`))
	if j == -1 {
		return ""
	}
	return string(data[:j])
}
