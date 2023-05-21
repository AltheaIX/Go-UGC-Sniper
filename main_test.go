package main

import (
	"fmt"
	"testing"
)

func TestUnmarshalCatalog(t *testing.T) {
	responseRaw, _, _ := ItemRecentlyAdded()
	jsonResp := UnmarshalCatalog(responseRaw)

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

func TestItemDetailById(t *testing.T) {
	accountCookie = "_|WARNING:-DO-NOT-SHARE-THIS.--Sharing-this-will-allow-someone-to-log-in-as-you-and-to-steal-your-ROBUX-and-items.|_B991E48E5442983989581996A9EC0E7416FFBB08CC15E95500FE7CF54B91344E164FA66BD8B248A4936CA482AA27ED79E3F64205F427EE155243BF2F8D61648A6D9ADFEA48227F087FE55902D95291655894E8EA01F0346E83369F435978A1CB5B3641663125679F24ADF3F754BFF93BE91FE27FA65A08EA0F074EDE345FB06678A22C930D50B8FDFD5B05B32734396A410E04FF67E0847F6AB9D6A31C199F0D0D8F3B4D6C7D27F32D9BAAE1268EA54A51C343981E037B6322EC973CDC11E01D713094F71CB952B726C4900DC54C2D454757DCBC7BAB61BE43415F0111B607D36658B069FD1C05D96928FF79E2E0CB6CDED3CA179C9422A214DAC35DE277832B671CAA8FDB6B35871CF8482EE450858E27732D507C88F22D089BAB14B5BE6BD2117E6B3061C8F0370C7E492780298C1DD016DF8B4AEBAAF9BD18DB83F5E5E208649BF123D6EA240C71E1526FFC66303242F7A9551B784D1D4EFE6D6A771FA03CFFF6F49831527D66EF135E54A6C3192876EEA810"
	for i := 0; i <= 250; i++ {
		fmt.Print(i)
		_, _ = ItemDetailById(13177094956)
	}
}

func TestItemRecentlyAddedAppend(t *testing.T) {
	listItems, _, _ := ItemRecentlyAddedAppend(ItemRecentlyAdded())
	fmt.Println(listItems)
}

func TestItemThumbnailImageById(t *testing.T) {
	thumbnailUrl, _ := ItemThumbnailImageById(13177094956)
	t.Log(thumbnailUrl)
}
