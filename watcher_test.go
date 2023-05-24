package main

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"strings"
	"testing"
)

func TestAddToWatcher(t *testing.T) {
	// AddToWatcher()
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
	data, _ := ItemDetailById(13186590783)
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
	newItemId := []int{13502640961, 13502643452}
	NotifierWatcherHandler(newItemId)
}
