package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestUnmarshalDiscord(t *testing.T) {
	url, _ := Decrypt("0rWGcWLnSc02Y2mFsG05tbl1frOBqgCQik8gkamEFeP2pDLfgst7ppOH1pyCMOEE9kUYluyv04dOzT/Q+6ZPduMBnsMBXnACbLcTEeP0YX3j1BBJswfk5Vr0HtVZzOlZ", xKey)

	response, _ := MakeRequestExternalScanner(url)
	scanner, _ := ResponseReader(response)

	pointerDiscord := UnmarshalDiscord(scanner)
	discord := *pointerDiscord

	t.Log(string(scanner))
	t.Log(discord[0].Content)
}

func TestMakeRequestExternalScanner(t *testing.T) {
	urlLink, _ := Decrypt("ePZKQrzSNR8O58R+Badoos2o4qh3bi4Y1YzWs0tepQke4y+XZ3pj+15qma92TJabziivo7H0CP9z2OtuTpvSsw==", xKey)

	response, err := MakeRequestExternalScanner(urlLink)
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
	LoadConfig()
	_ = ReadProxyFromFile("proxy_fresh", true)
	go ExternalScanner()
	wg.Wait()
}
