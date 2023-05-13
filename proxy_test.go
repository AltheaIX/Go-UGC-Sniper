package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestProxyTester(t *testing.T) {
	ProxyTester()
}

func TestCheckRequestProxy(t *testing.T) {
	var wg sync.WaitGroup
	ReadProxyFromFile("proxy")
	for _, data := range newProxy {
		proxyData := data
		wg.Add(1)
		err := CheckRequestProxy(&wg, proxyData)
		if err != nil {
			fmt.Println(err)
		}
	}
	wg.Wait()
}

func TestReadProxyFromFile(t *testing.T) {
	ReadProxyFromFile("proxy_fresh")
}

func TestWriteProxyToFile(t *testing.T) {
	WriteProxyToFile(proxyList)
}
