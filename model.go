package main

type Page struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"` // 👈 Add this line

	Body struct {
		Storage struct {
			Value string `json:"value"`
		} `json:"storage"`
	} `json:"body"`

	Ancestors []Ancestor `json:"ancestors"`
}

type Ancestor struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type PageResponse struct {
	Results []Page `json:"results"`
	Links   struct {
		Next string `json:"next"`
	} `json:"_links"`
}

type PageNode struct {
	Page     Page
	Children []*PageNode
}
