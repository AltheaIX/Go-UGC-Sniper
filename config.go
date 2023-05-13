package main

import (
	"io/ioutil"
	"os"
)

func LoadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var config Config

	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	freeWebhookUrl = config.FreeWebhook
	paidWebhookUrl = config.PaidWebhook
	lastItemId = config.LastId
}
