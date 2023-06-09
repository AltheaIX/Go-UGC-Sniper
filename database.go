package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func UnmarshalDatabase(responseRaw []byte) *Database {
	database := &Database{}

	err := json.Unmarshal(responseRaw, &database)
	if err != nil {
		fmt.Println(err)
	}
	return database
}

func ReadFirebase() (*Database, error) {
	dataBase := &Database{}
	// Set up the Firebase Realtime Database URL
	encryptedText := "8ZayCxzy4sXzrPsv3tov6UXY4XO0jTXCcBVJdQjS/L6MC/BqDJz9zUdZNxWI5oMUu3dbKMo2zEHEo6eYwCStv8A0es8KwKuMaGDo1nfB5WWesLboaaMQwGbjZPEard4o"
	databaseURL, err := Decrypt(encryptedText, xKey)
	if err != nil {
		fmt.Println("error on encryption.")
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}

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
