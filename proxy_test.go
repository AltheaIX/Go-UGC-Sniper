package main

import "testing"

func TestProxyTester(t *testing.T) {
	ProxyTester()
}

func TestReadProxyFromFile(t *testing.T) {
	ReadProxyFromFile()
}

func TestWriteProxyToFile(t *testing.T) {
	WriteProxyToFile(proxyList)
}
