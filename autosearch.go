package main

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
)

var allItems []int

func UnmarshalListItem(responseRaw []byte) (*ListItem, error) {
	listItem := &ListItem{}

	err := json.Unmarshal(responseRaw, &listItem)
	if err != nil {
		return listItem, err
	}
	return listItem, nil
}

func AllItems() error {
	defer handlePanic()

	fmt.Println("Autosearch - Please wait, we will get an old offsale items.")
	listItem := &ListItem{}

	var response *http.Response
	var err error
	var itemDetail *ItemDetail

	for {
		response, err = MakeRequest("https://catalog.roblox.com/v1/search/items?category=Accessories&includeNotForSale=true&limit=120&salesTypeFilter=1&sortType=3&subcategory=Accessories")
		if err != nil {
			fmt.Println(err)
			continue
		}

		defer response.Body.Close()
		if response.StatusCode != 200 {
			continue
		}

		break
	}

	scanner, _ := ResponseReader(response)

	if string(scanner) == "" {
		return errors.New("empty body")
	}

	listItem, err = UnmarshalListItem(scanner)

	for _, data := range listItem.Data {
		allItems = append(allItems, data.ID)
	}

	for {
		itemDetail, err = ItemDetailById(allItems)
		if err != nil {
			continue
		}
		break
	}

	for _, data := range itemDetail.Detail {
		if data.PriceStatus == "Off Sale" && data.Quantity == 0 {
			if IsExist(watcherId, data.Id) {
				continue
			}

			/*if len(watcherId) == 120 {
				watcherId = watcherId[:118]
			}*/

			watcherId = append(watcherId, data.Id)
			continue
		}

		if data.Quantity == 0 {
			continue
		}
	}

	for i := 0; i <= 10; i++ {
		allItems = []int{}

		if listItem.NextPageCursor != "" {
			for {
				response, err = MakeRequest(fmt.Sprintf("https://catalog.roblox.com/v1/search/items?category=Accessories&includeNotForSale=true&limit=120&salesTypeFilter=1&sortType=3&subcategory=Accessories&cursor=%v", listItem.NextPageCursor))
				if err != nil {
					fmt.Println(err)
					continue
				}

				defer response.Body.Close()
				if response.StatusCode != 200 {
					continue
				}

				break
			}

			scanner, _ := ResponseReader(response)

			if string(scanner) == "" {
				return errors.New("empty body")
			}

			listItem, err = UnmarshalListItem(scanner)

			for _, data := range listItem.Data {
				allItems = append(allItems, data.ID)
			}

			for {
				itemDetail, err = ItemDetailById(allItems)
				if err != nil {
					continue
				}
				break
			}

			for _, data := range itemDetail.Detail {
				if data.PriceStatus == "Off Sale" && data.Quantity == 0 {
					if IsExist(watcherId, data.Id) {
						continue
					}

					sort.Sort(sort.Reverse(sort.IntSlice(watcherId)))
					/*if len(watcherId) == 120 {
						watcherId = watcherId[:118]
						break
					}*/

					watcherId = append(watcherId, data.Id)
					continue
				}

				if data.Quantity == 0 {
					continue
				}
			}
		}
		break
	}

	fmt.Printf("Autosearch - Watching %d of offsale items now.\n", len(watcherId))
	fmt.Println("Autosearch - You are ready to go.")

	return nil
}
