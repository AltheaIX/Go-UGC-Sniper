package main

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"strings"
	"sync"
	"testing"
)

func TestAddToWatcher(t *testing.T) {
	// AddToWatcher()
}

func TestReleaseSemaphore(t *testing.T) {
	semaphore := make(chan struct{}, 3)
	var wg sync.WaitGroup

	for {
		select {
		case <-pauseChan:
			fmt.Println("Paused")
			resumeGoroutines()
			continue
		default:
			go func() {
				defer func() {
					wg.Done()
					ReleaseSemaphore(semaphore)
				}()

				semaphore <- struct{}{}
				wg.Add(1)
				fmt.Println("Worker started")

				pauseGoroutines()
			}()
		}
	}

	wg.Wait()
}

func TestOffsaleTracker(t *testing.T) {
	_ = ReadProxyFromFile("proxy_fresh", true)
	watcherId = []int{13562538536, 13570957619, 13569927870, 13571520410, 13570759527, 13570660992}
	var wg sync.WaitGroup

	wg.Add(1)
	go OffsaleTrackerHandler()
	wg.Wait()
}

func TestIsFieldSet(t *testing.T) {
	jsonRaw := `{
    "data": [
        {
            "id": 13452810701,
            "itemType": "Asset",
            "assetType": 43,
            "name": "Neck Kawaii Strawberry Camera",
            "description": "\n🍓🎀🧸made by tany1360 ☁️🎀🍓\n\n🍋join my group for more kawaii items and to chat:)🍋\nhttps://www.roblox.com/groups/11855791/The-Lemon-Land#!/store\n\n💗take a look at my catalog for more items💗\nhttps://www.roblox.com/catalog/?Category=1&CreatorName=SimplyALemon&SortType=3",
            "productId": 1541545429,
            "genres": [
                "All"
            ],
            "itemStatus": [],
            "itemRestrictions": [],
            "creatorHasVerifiedBadge": true,
            "creatorType": "Group",
            "creatorTargetId": 11855791,
            "creatorName": "The Lemon Land",
            "price": 20,
            "favoriteCount": 0,
            "offSaleDeadline": null,
            "saleLocationType": "NotApplicable"
        }
    ]
}`
	itemDetail, _ := UnmarshalCatalog([]byte(jsonRaw))
	fmt.Println(IsFieldSet(itemDetail.Detail[0], "Price"))
}

func TestNotifierWatcher(t *testing.T) {
	data, _ := ItemDetailById([]int{13186590783})
	data.Detail[0].Image, _ = ItemThumbnailImageById(13186590783)
	_name := strings.Replace(string(data.Detail[0].Name), `"`, "", 2)
	t.Log(_name)
	data.Detail[0].Name = jsoniter.RawMessage(_name)

	err := NotifierWatcher("paid", data.Detail[0])
	if err != nil {
		t.Log(err)
	}
}

func TestNotifierWatcherHandler(t *testing.T) {
	newItemId := []int{13502640961, 13502643452, 13557096529}
	NotifierWatcherHandler(newItemId)
}

func TestBoughtNotifier(t *testing.T) {
	LoadConfig()
	GetAccountDetails(accountCookie)

	detail, err := MarketplaceDetailByCollectibleItemId("71053164-baa8-449a-b131-f7fd96d68278")
	if err != nil {
		fmt.Println(err)
	}

	_name := strings.Replace(string(detail.Data[0].Name), `"`, "", 2)
	detail.Data[0].Name = jsoniter.RawMessage(_name)
	t.Log(fmt.Sprintf("%s", detail.Data[0].Name))
}
