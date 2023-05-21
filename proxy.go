package main

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"golang.org/x/net/proxy"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

var proxyList []string
var newProxy []string
var checkProxy bool

func ReadProxyFromFile(fileName string, forced bool) error {
	file, err := os.Open(fileName + ".txt") // Replace with the path to your file
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if checkProxy != true && fileName == "proxy_fresh" && forced != true {
			return errors.New("not empty")
		}

		line := scanner.Text()
		splitLine := strings.Split(line, ":")
		re := regexp.MustCompile("[0-9]+")
		proxyFormat := fmt.Sprintf("%s:%s", splitLine[0], re.FindAllString(splitLine[1], -1)[0])

		if fileName != "proxy_fresh" {
			newProxy = append(newProxy, proxyFormat)
			continue
		}

		proxyList = append(proxyList, proxyFormat)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	if checkProxy != true && fileName == "proxy_fresh" && forced != true {
		return errors.New("empty")
	}

	return nil
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

func CheckRequestProxy(wg *sync.WaitGroup, data string) error {
	defer wg.Done()
	fmt.Println("Checking, ", data)
	proxyURL, err := url.Parse("socks5://" + data)
	if err != nil {
		return errors.New("error on parse")
	}

	dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
	if err != nil {
		return errors.New("error on dialer")
	}

	transport := &http.Transport{
		Dial:            dialer.Dial,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Skip certificate verification
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   3 * time.Second,
	}

	// Use thumbnail because thumbnail no rate limit
	req, err := http.NewRequest("GET", "https://thumbnails.roblox.com/v1/assets?assetIds=1&returnPolicy=PlaceHolder&size=420x420&format=Png&isCircular=false", nil) // Replace with your target URL
	if err != nil {
		return errors.New("error on new request")
	}

	resp, err := client.Do(req)
	if err != nil {
		return errors.New("error on execute request")
	}
	defer resp.Body.Close()

	proxyList = append(proxyList, data)
	return nil
}

func ProxyTester() {
	runtime.GOMAXPROCS(5)

	var wg sync.WaitGroup

	err := ReadProxyFromFile("proxy_fresh", false)
	if checkProxy != true && err.Error() != "empty" {
		fmt.Println("Due to your configuration and your proxy_fresh.txt is not empty. We wont check your proxy.")
		return
	}

	_ = ReadProxyFromFile("proxy", true)
	fmt.Println("Checking proxy, we use 3s timeout for this checker to make sure proxy are fresh and fast.")

	semaphore := make(chan struct{}, 30)

	for _, data := range newProxy {
		proxyData := data

		semaphore <- struct{}{}

		wg.Add(1)
		go func() {
			defer func() {
				<-semaphore
			}()
			CheckRequestProxy(&wg, proxyData)
		}()
	}

	wg.Wait()
	WriteProxyToFile(proxyList)
	fmt.Printf("Proxy checker success, all fresh proxy are on proxy_fresh.txt and you are running this program with %d proxy.\n", len(proxyList))
	time.Sleep(time.Second * 3)

	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
