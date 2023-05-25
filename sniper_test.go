package main

import "testing"

func TestUnmarshalMarketplaceDetail(t *testing.T) {
	// responseRaw, err := MarketplaceDetailByCollectibleItemId("50987837-dc54-48cf-a1f0-f96ad1a26a32")
	responseRaw := []byte(`[
    {
        "collectibleItemId": "50987837-dc54-48cf-a1f0-f96ad1a26a32",
        "name": "Black Luxury Purse",
        "description": "Black Luxury Purse\nC FOR coded clothing\n\n👛 Shop more coded clothing:\nhttps://www.roblox.com/catalog?Category=1&CreatorName=coded%20clothing&CreatorType=Group&salesTypeFilter=1\n\ncreated by heartician",
        "collectibleProductId": "3f2092f9-9125-46c5-abd8-2ac0c69cbee1",
        "itemRestrictions": null,
        "creatorHasVerifiedBadge": false,
        "creatorType": "User",
        "itemTargetId": 13073745492,
        "creatorId": 1424338327,
        "creatorName": "codedcosmetics",
        "price": 0,
        "lowestPrice": 20,
        "hasResellers": true,
        "unitsAvailableForConsumption": 0,
        "offSaleDeadline": "0001-01-01T00:00:00",
        "assetStock": 250000,
        "errorCode": null,
        "saleLocationType": "ShopAndAllExperiences",
        "universeIds": [],
        "sales": 250000,
        "lowestResalePrice": 20
    }
]`)
	marketplaceDetail := UnmarshalMarketplaceDetail(responseRaw)
	t.Log(marketplaceDetail.Data[0].CreatorId)
}

func TestDeleteIntSlice(t *testing.T) {
	slice := []int{12312415123, 123121951234, 15283417528, 1231195182}
	updatedSlice := DeleteSlice(slice, 123212415123)
	t.Log(updatedSlice)
}

func TestMarketplaceDetailByCollectibleItemId(t *testing.T) {
	LoadConfig()
	detail, err := MarketplaceDetailByCollectibleItemId("7a88889e-70b0-4bc8-afa2-ceacb67d7b84")
	if err != nil {
		t.Log(err)
	}
	t.Log(detail)
}

func TestSniper(t *testing.T) {
	LoadConfig()
	GetAccountDetails(accountCookie)
	listFreeItem = append(listFreeItem, "7a88889e-70b0-4bc8-afa2-ceacb67d7b84")
	listFreeItem = append(listFreeItem, "302910d4-f926-491f-835b-6369e1ce60ef")
	listFreeItem = append(listFreeItem, "9db32ebc-b4b5-4b50-8735-01a337567d37")
	SniperHandler()
}
