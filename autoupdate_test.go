package main

import (
	"os"
	"testing"
)

func TestDownloadFile(t *testing.T) {
	executePath, err := os.Executable()
	if err != nil {
		panic(executePath)
	}
	t.Log(executePath)

	err = DownloadFile(executePath)
	t.Log(err)
}
