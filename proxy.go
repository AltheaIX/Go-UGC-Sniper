package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"golang.org/x/net/proxy"
	"net/http"
	"net/url"
	"os"
	"time"
)

var proxyList []string
var newProxy []string

func ReadProxyFromFile(fileName string) {
	file, err := os.Open(fileName + ".txt") // Replace with the path to your file
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if fileName != "proxy_fresh" {
			newProxy = append(newProxy, line)
		}

		proxyList = append(proxyList, line)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func WriteProxyToFile(proxyList []string) {
	file, err := os.Create("proxy_fresh.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Loop through the array and write each element to the file
	for _, element := range proxyList {
		_, err := fmt.Fprintln(file, element)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Data written to file successfully!")
}

func ProxyTester() {
	ReadProxyFromFile("proxy")
	fmt.Println("Checking proxy, we use 3s timeout for this checker to make sure proxy are fresh and fast.")
	for _, data := range newProxy {
		proxyURL, err := url.Parse("socks5://" + data)
		fmt.Println(proxyURL)
		if err != nil {
			fmt.Println("error on parser.")
			panic(err)
		}

		dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
		if err != nil {
			fmt.Println("error on dialer.")
			panic(err)
		}

		transport := &http.Transport{
			Dial:            dialer.Dial,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Skip certificate verification
		}

		client := &http.Client{
			Transport: transport,
			Timeout:   3 * time.Second,
		}

		req, err := http.NewRequest("GET", "https://thumbnails.roblox.com/v1/assets?assetIds=1&returnPolicy=PlaceHolder&size=420x420&format=Png&isCircular=false", nil) // Replace with your target URL
		if err != nil {
			fmt.Println("error on new request.")
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("error on execute request.")
			continue
		}
		defer resp.Body.Close()

		proxyList = append(proxyList, data)
	}
	WriteProxyToFile(proxyList)
	fmt.Println("Proxy checker success, all fresh proxy are on proxy_fresh.txt")
}
