package main

import (
	"bufio"
	"gitee.com/LXY1226/logging"
	rar "github.com/nwaples/rardecode"
	"io"
	"net/http"
	"os"
	"path"
)

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

func downURL(path string) *bufio.Reader {
	for _, uri := range repos {
		resp, err := http.Get(uri + path)
		if err == nil {
			return bufio.NewReaderSize(resp.Body, 8<<20)
		}
		logging.DEBUG(path, "From", uri, "Error:", err.Error())
	}
	logging.ERROR("无法下载", path)
	return nil
}
