package main

import (
	"github.com/k0kubun/go-ansi"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var console = transform.NewWriter(ansi.NewAnsiStdout(), simplifiedchinese.GBK.NewDecoder())
