package main

type Discord []struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Author  struct {
		ID               string `json:"id"`
		Username         string `json:"username"`
		GlobalName       any    `json:"global_name"`
		Avatar           string `json:"avatar"`
		Discriminator    string `json:"discriminator"`
		PublicFlags      int    `json:"public_flags"`
		Bot              bool   `json:"bot"`
		AvatarDecoration any    `json:"avatar_decoration"`
	} `json:"author"`
	Embeds []struct {
		Type        string `json:"type"`
		URL         string `json:"url"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Color       int    `json:"color"`
		Fields      []struct {
			Name   string `json:"name"`
			Value  string `json:"value"`
			Inline bool   `json:"inline"`
		} `json:"fields"`
		Author struct {
			Name         string `json:"name"`
			URL          string `json:"url"`
			IconURL      string `json:"icon_url"`
			ProxyIconURL string `json:"proxy_icon_url"`
		} `json:"author"`
		Thumbnail struct {
			URL      string `json:"url"`
			ProxyURL string `json:"proxy_url"`
			Width    int    `json:"width"`
			Height   int    `json:"height"`
		} `json:"thumbnail"`
	} `json:"embeds"`
}
