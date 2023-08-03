package domain

type ZincSearchResponse struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		EmailHits []struct {
			Source Email `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
