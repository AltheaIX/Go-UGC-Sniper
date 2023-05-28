package main

import (
	"fmt"
	"sort"
	"testing"
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
	offsaleIDs := []int{
		13570560231, 13570534662, 13570532840, 13570531188, 13570273027, 13570191188, 13570188732, 13570154119,
		13570105324, 13570094170, 13570070214, 13570049226, 13570013624, 13569959175, 13569934343, 13569927870,
		13569926364, 13569924588, 13569827262, 13569702943, 13569613881, 13569485009, 13569250146, 13569187328,
		13569146777, 13568987638, 13568762208, 13568618919, 13568617902, 13568601396, 13568553519, 13568502194,
		13568482849, 13568222541, 13568165797, 13568123923, 13567972059, 13567913802, 13567898413, 13567880855,
		13567765138, 13567761643, 13567728168, 13567465002, 13567432037, 13567254762, 13567245588, 13567228938,
		13567220766, 13567109767, 13567004504, 13566870758, 13566835345, 13566799843, 13566796840, 13566741774,
		13565845787, 13565842073, 13565830150, 13565803844, 13565793623, 13565133870, 13565119699, 13564981015,
		13564769395, 13564765828, 13564762193, 13564758667, 13564752622, 13564434584, 13564408563, 13564400083,
		13564366253, 13564354059, 13562538536, 13560127954, 13560119143, 13560105425, 13559963707, 13559664554,
		13559462008, 13559437068, 13559320301, 13559010338, 13558960287, 13558954080, 13558941618, 13558892904,
		13558883772, 13558860994, 13558849178, 13558839193, 13558823091, 13558793911, 13558766862, 13558739693,
		13558723102, 13558698218, 13558695593, 13558650509, 13558551669, 13558397075, 13558113120, 13558070304,
		13558010756, 13557945018, 13557816827, 13557791789, 13557777442, 13557535325, 13557439606, 13557359648,
		13570759527, 13570868803, 13556983000, 13556819585, 13554363059, 13570945135, 13570931521,
	}

	sort.Sort(sort.Reverse(sort.IntSlice(offsaleIDs)))
	sort.Sort(sort.Reverse(sort.IntSlice(offsaleIDs)))
	fmt.Println(offsaleIDs)
}
