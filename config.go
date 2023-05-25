package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var username string

func SetAccountDetails(user User) {
	accountId = user.Id
	username = user.Username
}

func GetAccountDetails(accountCookie string) *User {
	client := &http.Client{Timeout: 10 * time.Second}

	cookie := &http.Cookie{
		Name:    ".ROBLOSECURITY",
		Value:   accountCookie,
		Expires: time.Now().Add(time.Hour * 1000),
	}

	req, err := http.NewRequest("GET", "https://www.roblox.com/mobileapi/userinfo", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie)

	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		err = errors.New("status code is not 200")
	}

	scanner, _ := ResponseReader(response)

	userPtr := UnmarshalAccount(scanner)
	user := *userPtr
	SetAccountDetails(user)
	return userPtr
}

func LoadConfig() (Config, error) {
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
	accountCookie = config.Cookie
	checkProxy = config.Options.AlwaysCheckProxy
	threads = config.Options.Threads
	lastItemId = config.LastId
	watcherId = config.OffsaleId
	return config, err
}

func SaveConfig(filename string, config Config) error {
	file, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, file, 0644)
	if err != nil {
		return err
	}

	fmt.Println("System - Config saved successfully.")

	return nil
}
