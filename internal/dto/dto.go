package dto

type LinkState struct {
	URL    string `json:"url"`
	Status string `json:"status"`
}

type LinkCheckResult struct {
	ID      int         `json:"id"`
	Results []LinkState `json:"results"`
}

type NewLinkRequest struct {
	Links []string `json:"links"`
}

type NewLinksNumRequest struct {
	LinksList []int `json:"links_list"`
}
