package main

import (
	"bufio"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

const VERSION = "v1.2.7"

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var listId []int

var lastItemId int

func ResponseReader(response *http.Response) ([]byte, error) {
	var body []byte
	var err error

	reader := bufio.NewReader(response.Body)
	for {
		chunk, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				body = append(body, chunk...)
				break
			}
			return nil, err
		}
		body = append(body, chunk...)
	}
	return body, err
}

func DeleteSlice[T comparable](list []T, elementToRemove T) []T {
	var newSlice []T
	elementFound := false

	for _, element := range list {
		if element == elementToRemove {
			elementFound = true
			continue
		}

		newSlice = append(newSlice, element)
	}

	if elementFound {
		return newSlice
	}

	return list
}

func IsExist[T comparable](list []T, elementToCheck T) bool {
	encountered := make(map[T]bool)

	for _, element := range list {
		encountered[element] = true
	}

	if encountered[elementToCheck] {
		return true
	}

	return false
}

func UnmarshalCatalog(responseRaw []byte) (*ItemDetail, error) {
	itemDetail := &ItemDetail{}

	err := json.Unmarshal(responseRaw, &itemDetail)
	if err != nil {
		return itemDetail, err
	}
	return itemDetail, nil
}

func UnmarshalAccount(responseRaw []byte) *User {
	user := &User{}

	err := json.Unmarshal(responseRaw, &user)
	if err != nil {
		return user
	}

	return user
}

func ItemRecentlyAddedAppend(scanner []byte, proxy *url.URL, err error) (int, *url.URL, error) {
	if err != nil {
		return 0, proxy, err
	}

	listId = nil

	listItems, err := UnmarshalCatalog(scanner)
	if err != nil {
		return lastItemId, proxy, err
	}

	for _, data := range listItems.Detail {
		if data.Id == lastItemId {
			break
		}

		listId = append(listId, data.Id)
	}
	return listItems.Detail[0].Id, proxy, nil
}

func AutoUpdate(executePath string) {
	fmt.Println("New update found, auto update on process...")

	var updateFilePath string
	updateDir := filepath.Dir(executePath)
	updateFilePath = filepath.Join(updateDir, "Go-UGC-Sniper.exe")

	if !strings.Contains(executePath, "Go-UGC-Sniper-update.exe") {
		updateFilePath = filepath.Join(updateDir, "Go-UGC-Sniper-update.exe")
	}

	err := DownloadFile(updateFilePath)
	if err != nil {
		fmt.Println("Failed to auto update, but you can still update it from Wagoogus Discord.")
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}

	// Start the command to delete the current executable in the background
	err = os.Remove(executePath)
	if err != nil {
		fmt.Println("Failed to replace program, but it is still saved.")
		return
	}

	err = os.Rename(updateFilePath, executePath)
	if err != nil {
		fmt.Println("Failed to replace executable, ", err)
		return
	}

	fmt.Println("Update success, program will closed and you can just reopen it.")
	time.Sleep(5 * time.Second)
	os.Exit(1)
}

// Comment this for linux build
func setConsoleTitle(title string) error {
	handle, err := syscall.LoadLibrary("kernel32.dll")
	if err != nil {
		return err
	}
	defer syscall.FreeLibrary(handle)

	proc, err := syscall.GetProcAddress(handle, "SetConsoleTitleW")
	if err != nil {
		return err
	}

	_, _, callErr := syscall.Syscall(proc, 1, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))), 0, 0)
	if callErr != 0 {
		return callErr
	}

	return nil
}

func isDebuggerPresent() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	procCheckRemoteDebuggerPresent := kernel32.NewProc("CheckRemoteDebuggerPresent")
	var isDebuggerPresent int32
	handle, err := syscall.GetCurrentProcess()
	if err != nil {
		panic(err)
	}
	_, _, _ = procCheckRemoteDebuggerPresent.Call(uintptr(unsafe.Pointer(uintptr(handle))), uintptr(unsafe.Pointer(&isDebuggerPresent)))
	if isDebuggerPresent != 0 {
		fmt.Println("Debugger detected, closing in 10 seconds.")
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}
}

func handlePanic() {
	if r := recover(); r != nil {
		// Open or create the crash_log.txt file
		file, err := os.OpenFile("crash_log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// Set the log output to the crash_log.txt file
		log.SetOutput(file)

		// Log the panic details to crash_log.txt
		log.Printf("Panic occurred: %v", r)

		// Re-panic to ensure the program terminates
		panic(r)
	}
}

func main() {
	log.SetFlags(log.LstdFlags)
	defer handlePanic()

	// Comment this for linux build
	isDebuggerPresent()

	executePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	if strings.Contains(executePath, "Go-UGC-Sniper-update.exe") {
		os.Remove(filepath.Join(filepath.Dir(executePath), "Go-UGC-Sniper.exe"))
	}

	config, err := LoadConfig()

	// Comment this for linux build
	err = setConsoleTitle(fmt.Sprintf("Go UGC Sniper - Beta Version - %v - Threads %d", VERSION, threads))
	if err != nil {
		panic(err)
	}

	database, err := ReadFirebase()
	if err != nil {
		panic(err)
	}

	if database.Trial.Status != "active" {
		fmt.Println("Trial ended, please purchase it on Wagoogus discord.")
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}

	if database.Version.Version != VERSION {
		AutoUpdate(executePath)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	defer wg.Done()
	wg.Add(1)

	go func() {
		<-sig

		config.LastId = lastItemId
		err = SaveConfig("config.json", config)
		if err != nil {
			fmt.Println("System - Config unsaved.")
		}

		wg.Done()
	}()

	if err != nil {
		panic(err)
	}

	ProxyTester()
	_ = ReadProxyFromFile("proxy_fresh", true)

	if len(proxyList) == 0 {
		fmt.Println("Exiting program, no proxy are fresh.")
		time.Sleep(3 * time.Second)
		os.Exit(1)
	}

	userDetail := GetAccountDetails(accountCookie)

	if userDetail.Id == 0 {
		fmt.Println("Invalid Cookie, program exiting within 5 seconds.")
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}

	fmt.Printf("Logging in as %v and id %d\n\n", userDetail.Username, userDetail.Id)
	time.Sleep(time.Second * 5)

	err = AllItems()
	if err != nil {
		panic(err)
		time.Sleep(5 * time.Second)
	}

	AddToWatcher(sig)

	fmt.Println("Program exited.")
	wg.Wait()
}
