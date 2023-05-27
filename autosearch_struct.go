package main

type ListItemData struct {
	ID       int    `json:"id"`
	ItemType string `json:"itemType"`
}

type ListItem struct {
	PreviousPageCursor any            `json:"previousPageCursor"`
	NextPageCursor     string         `json:"nextPageCursor"`
	Data               []ListItemData `json:"data"`
}
