package main

import (
	"fmt"
	"testing"
)

func TestUnmarshalCatalog(t *testing.T) {
	responseRaw, _, _ := ItemRecentlyAdded()
	jsonResp, _ := UnmarshalCatalog(responseRaw)

	fmt.Println(jsonResp.Detail[0].Id)
}

func TestGetCsrfToken(t *testing.T) {
	LoadConfig()

	token := GetCsrfToken()
	fmt.Println(token)
}

func TestItemRecentlyAdded(t *testing.T) {
	_ = ReadProxyFromFile("proxy_fresh", true)
	for {
		responseByte, _, err := ItemRecentlyAdded()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(responseByte))
		fmt.Println("")
	}
}

func TestDeleteSlice(t *testing.T) {
	slice1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for i := 0; i < 10; i++ {
		t.Log(slice1)
		slice1 = slice1[:len(slice1)]
		slice1 = append(slice1, 1515151)
		t.Log(slice1)
		slice1 = DeleteSlice(slice1, 1515151)
		t.Log(slice1)
		slice1 = DeleteSlice(slice1, 1515151)
		t.Log(slice1)
	}
}

func TestIsExist(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	check := IsExist(slice, 10)
	fmt.Println(check)
}

func TestItemDetailByIdProxied(t *testing.T) {
	_ = ReadProxyFromFile("proxy_fresh", true)
	for {
		responseByte, err := ItemDetailByIdProxied([]int{123123, 12412312})
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(responseByte)
	}
}

func TestItemDetailById(t *testing.T) {
	accountCookie = ""
	details, _ := ItemDetailById([]int{13558113120, 13558070304, 13558010756, 13557945018, 13557096529})
	t.Log(details.Detail)
}

func TestItemRecentlyAddedAppend(t *testing.T) {
	listItems, _, _ := ItemRecentlyAddedAppend(ItemRecentlyAdded())
	fmt.Println(listItems)
}

func TestItemThumbnailImageById(t *testing.T) {
	thumbnailUrl, _ := ItemThumbnailImageById(13177094956)
	t.Log(thumbnailUrl)
}

func TestAnything(t *testing.T) {
	/*watcherId = []int{1111, 2222, 3333, 4444, 5555}
	sort.Sort(sort.Reverse(sort.IntSlice(watcherId)))
	fmt.Println(watcherId)
	watcherId = append(watcherId, 6666)
	fmt.Println(watcherId[:])
	sort.Sort(sort.Reverse(sort.IntSlice(watcherId)))
	fmt.Println(watcherId[:])*/

	test := true

	for {
		go func() {
			externalScannerMutex.Lock()
			t.Log(test)
			test = false
			externalScannerMutex.Unlock()
		}()
	}
}
