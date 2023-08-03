package controller

import (
	"Backend/domain"
	"Backend/repository"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
)

type EmailController struct {
	Ecr repository.EmailRepository
}

func (rs EmailController) Routes() chi.Router {
	r := chi.NewRouter()

	//r.Get("/", rs.List) // GET /posts - Read a list of posts.

	r.Route("/{id}", func(r chi.Router) {
		r.Use(EmailCtx)
		r.Get("/", rs.Get) // GET /emails/{key} - Read a single post by :id.
	})

	return r
}

func EmailCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "id", chi.URLParam(r, "id"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Request Handler - GET /emails/{id}.
func (rs EmailController) Get(w http.ResponseWriter, r *http.Request) {
	key := r.Context().Value("id").(string)

	body, err, resp := rs.Ecr.CallZincSearch(key)

	if err != nil {
		log.Fatal(err)
		return
	}

	defer resp.Body.Close()

	log.Println(resp.StatusCode)

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var zincSearchResponse domain.ZincSearchResponse
	err = json.Unmarshal(body, &zincSearchResponse)
	if err != nil {
		log.Fatal(err)
		return
	}

	var emails []domain.Email
	for _, emailHit := range zincSearchResponse.Hits.EmailHits {
		emails = append(emails, emailHit.Source)
	}

	response, err := json.Marshal(emails)
	if err != nil {
		// handle error
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	// Set the response content type
	w.Header().Set("Content-Type", "application/json")

	w.Write(response)
}
