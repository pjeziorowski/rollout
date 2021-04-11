package platforms

import (
	"github.com/go-resty/resty/v2"
	"log"
)

// Medium is used to interact with Medium platform
type Medium struct {
	apiToken      string
	publicationID string
}

// NewMedium creates a new instance of Medium
func NewMedium(publicationID string, apiToken string) *Medium {
	return &Medium{apiToken: apiToken, publicationID: publicationID}
}

// Publish publishes an article on Medium
func (m *Medium) Publish(title string, markdown string, tags []string, canonicalURL string) {
	client := resty.New()

	r, err := client.R().
		SetBody(newCreateMediumPostPayload(title, markdown, tags, canonicalURL)).
		SetHeader("Authorization", "Bearer "+m.apiToken).
		Post("https://api.medium.com/v1/users/" + m.publicationID + "/posts")

	if err != nil {
		log.Fatal(err.Error())
	}

	if r.IsError() {
		log.Fatal(string(r.Body()))
	}
}

type createMediumPostPayload struct {
	Title         string   `json:"title"`
	ContentFormat string   `json:"contentFormat"`
	Content       string   `json:"content"`
	CanonicalURL  string   `json:"canonical_url"`
	Tags          []string `json:"tags"`
}

func newCreateMediumPostPayload(title string, content string, tags []string, canonicalURL string) *createMediumPostPayload {
	return &createMediumPostPayload{Title: title, ContentFormat: "markdown", Content: content, Tags: tags, CanonicalURL: canonicalURL}
}
