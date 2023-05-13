package main

import (
	jsoniter "github.com/json-iterator/go"
	"strings"
	"testing"
)

func TestAddToWatcher(t *testing.T) {
	AddToWatcher()
}

func TestNotifierWatcher(t *testing.T) {
	data, _ := ItemDetailById(13186590783)
	data.Detail[0].Image, _ = ItemThumbnailImageById(13186590783)
	_name := strings.Replace(string(data.Detail[0].Name), `"`, "", 2)
	t.Log(_name)
	data.Detail[0].Name = jsoniter.RawMessage(_name)

	err := NotifierWatcher("paid", data.Detail[0])
	if err != nil {
		t.Log(err)
	}
}
