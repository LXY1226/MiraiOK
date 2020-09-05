package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

var lDate = tTime{
	buf: []byte("xxxx-xx-xx xx:xx:xx"),
}

var lLock sync.Mutex

func INFO(v ...interface{}) {
	lLock.Lock()
	colorINFO()
	os.Stdout.Write(lDate.str())
	os.Stdout.WriteString(" I/MiraiOK: ")
	fmt.Println(v...)
	lLock.Unlock()
}

func WARN(v ...interface{}) {
	lLock.Lock()
	colorWARN()
	os.Stdout.Write(lDate.str())
	os.Stdout.WriteString(" W/MiraiOK: ")
	fmt.Println(v...)
	lLock.Unlock()
}

func ERROR(v ...interface{}) {
	lLock.Lock()
	colorERROR()
	os.Stdout.Write(lDate.str())
	os.Stdout.WriteString(" E/MiraiOK: ")
	fmt.Println(v...)
	lLock.Unlock()
}

type tTime struct {
	Year, Day, Hour, Minute, Second int
	Month                           time.Month
	buf                             []byte
}

func (t *tTime) str() []byte {
	ct := time.Now()
	year, month, day := ct.Date()
	hour, minute, second := ct.Clock()
	if year != t.Year {
		t.buf[0] = byte('0' + year/1000%10)
		t.buf[1] = byte('0' + year/100%10)
		t.buf[2] = byte('0' + year/10%10)
		t.buf[3] = byte('0' + year%10)
		t.Year = year
	}
	if month != t.Month {
		t.buf[5] = byte('0' + month/10%10)
		t.buf[6] = byte('0' + month%10)
		t.Month = month
	}
	if day != t.Day {
		t.buf[8] = byte('0' + day/10%10)
		t.buf[9] = byte('0' + day%10)
		t.Day = day
	}
	if hour != t.Hour {
		t.buf[11] = byte('0' + hour/10%10)
		t.buf[12] = byte('0' + hour%10)
		t.Hour = hour
	}
	if minute != t.Minute {
		t.buf[14] = byte('0' + minute/10%10)
		t.buf[15] = byte('0' + minute%10)
		t.Minute = minute
	}
	if second != t.Second {
		t.buf[17] = byte('0' + second/10%10)
		t.buf[18] = byte('0' + second%10)
		t.Second = second
	}
	return t.buf
}
