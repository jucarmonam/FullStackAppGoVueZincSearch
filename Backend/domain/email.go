package domain

type Email struct {
	Message_ID string `json:"Message_ID"`
	Date       string `json:"Date"`
	From       string `json:"From"`
	To         string `json:"To"`
	Subject    string `json:"Subject"`
	Body       string `json:"Body"`
}
