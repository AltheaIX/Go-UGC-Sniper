package main

import (
	"fmt"
	"testing"
)

func TestUnmarshalCatalog(t *testing.T) {
	responseRaw, _ := ItemRecentlyAdded()
	jsonResp := UnmarshalCatalog(responseRaw)

	fmt.Println(jsonResp.Detail[0].Id)
}

func TestGetCsrfToken(t *testing.T) {
	token := GetCsrfToken()
	fmt.Println(token)
}

func TestItemRecentlyAdded(t *testing.T) {
	responseByte, _ := ItemRecentlyAdded()
	fmt.Println(string(responseByte))
}

func TestItemDetailById(t *testing.T) {
	itemDetail, _ := ItemDetailById(13177094956)
	fmt.Println(itemDetail.Detail)
}

func TestItemRecentlyAddedAppend(t *testing.T) {
	listItems, _ := ItemRecentlyAddedAppend()
	fmt.Println(listItems)
}

func TestItemThumbnailImageById(t *testing.T) {
	thumbnailUrl, _ := ItemThumbnailImageById(13177094956)
	t.Log(thumbnailUrl)
}

func TestAddToWatcher(t *testing.T) {
	//AddToWatcher()
}

func TestNotifierWatcher(t *testing.T) {
	data, _ := ItemDetailById(13177094956)
	data.Detail[0].Image, _ = ItemThumbnailImageById(13177094956)
	err := NotifierWatcher("free", data.Detail[0])
	if err != nil {
		t.Log(err)
	}
}
