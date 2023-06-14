package main

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"testing"
)

func TestUnmarshalDiscord(t *testing.T) {
	_, _ = Decrypt("0rWGcWLnSc02Y2mFsG05tbl1frOBqgCQik8gkamEFeP2pDLfgst7ppOH1pyCMOEE9kUYluyv04dOzT/Q+6ZPduMBnsMBXnACbLcTEeP0YX3j1BBJswfk5Vr0HtVZzOlZ", xKey)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	response, _ := MakeRequestExternalScanner("https://discord.com/api/v9/channels/1094291863332192376/messages?limit=50", transport)
	scanner, _ := ResponseReader(response)

	pointerDiscord := UnmarshalDiscord(scanner)
	discord := *pointerDiscord

	t.Log(string(scanner))
	t.Log(discord[0].Content)
}

func TestMakeRequestExternalScannerProxied(t *testing.T) {
	err := ReadProxyFromFile("proxy_fresh", true)
	LoadConfig()
	t.Log(err)

	semaphore := make(chan struct{}, 10)

	for {
		semaphore <- struct{}{}
		go func() {
			defer func() {
				<-semaphore
			}()

			proxyURL, err := url.Parse(BuildProxyURL(proxyList[rand.Intn(len(proxyList)-1)]))
			if err != nil {
				t.Log(err)
			}

			transport := &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			}

			response, err := MakeRequestExternalScanner("https://discord.com/api/v9/channels/1094291863332192376/messages?limit=50", transport)
			if err != nil {
				return
			}

			scanner, _ := ResponseReader(response)
			t.Log(string(scanner))
		}()
	}
}

func TestMakeRequestExternalScanner(t *testing.T) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	response, err := MakeRequestExternalScanner("https://discord.com/api/v9/channels/1094291863332192376/messages?limit=50", transport)
	if err != nil {
		t.Error(err)
	}

	scanner, _ := ResponseReader(response)

	pointerDiscord := UnmarshalDiscord(scanner)
	discord := *pointerDiscord

	fmt.Println(discord)
	/*	lastExternalScannerId = 13675149661*/

	/*	for i, data := range discord {
			if len(data.Embeds) < 1 {
				continue
			}

			url := data.Embeds[0].URL
			pattern := `https:\/\/www\.roblox\.com\/catalog\/(\d+)`
			regex := regexp.MustCompile(pattern)
			matches := regex.FindStringSubmatch(url)

			if len(matches) < 1 {
				continue
			}

			itemId, err := strconv.Atoi(matches[1])
			if err != nil {
				fmt.Println(err)
				continue
			}

			if i == 0 {
				if lastExternalScannerId == itemId {
					break
				}

				lastExternalScannerId = itemId
			}
		}

		t.Log("Last Id:", lastExternalScannerId)*/
}

func TestExternalScanner(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	firstExternalScanner = false
	lastExternalScannerId = 13683406501
	LoadConfig()
	_ = ReadProxyFromFile("proxy_fresh", true)
	go ExternalScanner()
	wg.Wait()
}
