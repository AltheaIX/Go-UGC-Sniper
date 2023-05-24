package main

import (
	"fmt"
	"log"
	"net/http"
)

func UnmarshalDatabase(responseRaw []byte) *Database {
	trial := &Database{}

	err := json.Unmarshal(responseRaw, &trial)
	if err != nil {
		fmt.Println(err)
	}
	return trial
}

func ReadFirebase() (*Database, error) {
	dataBase := &Database{}
	// Set up the Firebase Realtime Database URL
	databaseURL := "https://ugc-authentication-default-rtdb.asia-southeast1.firebasedatabase.app/"

	requestURL := databaseURL + ".json"

	response, err := http.Get(requestURL)
	if err != nil {
		panic("Failed to read database.")
		return dataBase, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		panic("Failed to read database.")
		return dataBase, err
	}

	scanner, err := ResponseReader(response)
	if err != nil {
		log.Println(err)
		return dataBase, err
	}

	dataBase = UnmarshalDatabase(scanner)
	return dataBase, nil
}
