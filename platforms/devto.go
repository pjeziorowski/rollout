package platforms

import (
	"github.com/go-resty/resty/v2"
	"log"
)

// Devto is used to interact with dev.to platform
type Devto struct {
	apiToken string
}

// NewDevto creates a new instance of Devto
func NewDevto(apiToken string) *Devto {
	return &Devto{apiToken: apiToken}
}

// Publish publishes an article on dev.to
func (d *Devto) Publish(title string, markdown string, tags []string, canonicalURL string) {
	client := resty.New()

	r, err := client.R().
		SetBody(newCreateDevtoPostPayload(title, markdown, tags, canonicalURL)).
		SetHeader("api-key", d.apiToken).
		Post("https://dev.to/api/articles")

	if err != nil {
		log.Fatal(err.Error())
	}

	if r.IsError() {
		log.Fatal(string(r.Body()))
	}
}

type createDevtoPostPayload struct {
	Article devtoArticle `json:"article"`
}

func newCreateDevtoPostPayload(title string, markdown string, tags []string, canonicalURL string) *createDevtoPostPayload {
	return &createDevtoPostPayload{Article: devtoArticle{
		Title:        title,
		Markdown:     markdown,
		Published:    true,
		CanonicalURL: canonicalURL,
		Tags:         tags,
	}}
}

type devtoArticle struct {
	Title        string   `json:"title"`
	Markdown     string   `json:"body_markdown"`
	Published    bool     `json:"published"`
	CanonicalURL string   `json:"canonical_url"`
	Tags         []string `json:"tags"`
}
