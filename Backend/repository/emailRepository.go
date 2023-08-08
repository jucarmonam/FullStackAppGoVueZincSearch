package repository

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type EmailRepository struct{}

func (er EmailRepository) CallZincSearch(key string) ([]byte, error, *http.Response) {
	params := map[string]interface{}{
		"search_type": "match",
		"query": map[string]string{
			"term":  key,
			"field": "_all",
		},
		"sort_fields": []string{"-@timestamp"},
		"from":        0,
		"max_results": 20,
		"_source":     []string{"Message_ID", "From", "To", "Subject", "Date", "Body"},
	}

	body, _ := json.Marshal(params)
	dbUrl := os.Getenv("DATABASE_URL")
	req, err := http.NewRequest("POST", dbUrl, bytes.NewBuffer(body))

	if err != nil {
		log.Fatal(err)
		return nil, nil, nil
	}

	req.SetBasicAuth("admin", "Complexpass#123")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return body, err, resp
}
