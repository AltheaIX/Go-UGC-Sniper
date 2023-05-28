package main

import (
	"io"
	"net/http"
	"os"
)

func DownloadFile(filePath string) error {
	db, err := ReadFirebase()
	if err != nil {
		panic(err)
	}

	response, err := http.Get(db.Version.Newlink)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
