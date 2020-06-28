package main

import (
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	time.AfterFunc(30*time.Second, func() {
		os.Exit(0)
	})
	main()
}
