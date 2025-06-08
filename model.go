package main

type Page struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	Type       string     `json:"type"`
	Status     string     `json:"status"`
	Body       PageBody   `json:"body"`
	Ancestors  []Ancestor `json:"ancestors"`
	Space      Space      `json:"space"`
	Extensions struct {
		ContentRepresentation struct {
			View string `json:"view"`
		} `json:"content-representation"`
	} `json:"extensions"`
}

type PageBody struct {
	Storage struct {
		Value string `json:"value"`
	} `json:"storage"`
}

type Ancestor struct {
	ID string `json:"id"`
}

type Space struct {
	Key string `json:"key"`
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
