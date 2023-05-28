package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestUnmarshalCatalog(t *testing.T) {
	responseRaw, _, _ := ItemRecentlyAdded()
	jsonResp, _ := UnmarshalCatalog(responseRaw)

	fmt.Println(jsonResp.Detail[0].Id)
}

func TestGetCsrfToken(t *testing.T) {
	// token := GetCsrfToken()
	// fmt.Println(token)
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
	accountCookie = "_|WARNING:-DO-NOT-SHARE-THIS.--Sharing-this-will-allow-someone-to-log-in-as-you-and-to-steal-your-ROBUX-and-items.|_4C6F2072DF0773910B43850465C0171A8FBB8ACDDDCBAA19E5C526C1E5CE8359512E5C9F17AA8EB20678B685395FD900AB24362ED6446F4A443A0016C4C978B014371F878FC1EF985F1E9DD18FDE3A481B1CED451C3AAE9D2219828756087DE30881F1CE07201F8E55FDADF91E97E3371F3642F0F7E9A6048FCDA2E797922051353F03A51B99E4D2108CA55E95968DE7894C346CE716590A1030EB96C883177115FCA1430A756E710612DB835725833D4C59484976F00FE2FD9C33462E2ECE9187F31C2617CF30C3E75C64859BEDEF363832272EDD0AC39B88BD6BB5D6BB04CACA220BAB6B82AE9E3EAD85F662DE2A4FD6B2FBCB60300A7E868716ACD3E80FD446A3A5E2243357F5DB4F3BEF1EA4F39435D134E514B6CDAA291094754A64A0ED05215A96E8F16E3EFA897C86117E724FF229255761A040400C353771C4E07CE3D01BB3642B92C27C47F2B6E812335F01275D67504FE4E344D0207ECBD077A1B2E99C594F1A74AF66E69680EA834892DE9E403B3B"
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
	fmt.Println("System - Offsale Tracker activated.")
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10)
	watcherId = []int{123124123, 1231234912, 12391285934, 12312482351}

	for {
		wg.Add(1)

		go func() {
			time.Sleep(5 * time.Second)
			watcherId = DeleteSlice(watcherId, 123124123)
			fmt.Println(len(watcherId))
		}()

		go func(wg *sync.WaitGroup, semaphore chan struct{}) {
			defer func() {
				<-semaphore
				wg.Done()
			}()
			semaphore <- struct{}{}

			watcherMutex.Lock()
			offsaleId := watcherId[:]
			watcherMutex.Unlock()

			for {
				fmt.Println(offsaleId)
			}
		}(&wg, semaphore)

		wg.Wait()
	}
}
