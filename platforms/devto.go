package platforms

import (
	"github.com/go-resty/resty/v2"
	"log"
)

type Devto struct {
	apiToken string
}

func NewDevto(apiToken string) *Devto {
	return &Devto{apiToken: apiToken}
}

func (d *Devto) Publish(title string, markdown string, tags []string, canonicalUrl string) {
	client := resty.New()

	_, err := client.R().
		SetBody(NewCreateDevtoPostPayload(title, markdown, tags, canonicalUrl)).
		SetHeader("api-key", d.apiToken).
		Post("https://dev.to/api/articles")

	if err != nil {
		log.Fatal(err.Error())
	}
}

type CreateDevtoPostPayload struct {
	Article DevtoArticle `json:"article"`
}

func NewCreateDevtoPostPayload(title string, markdown string, tags []string, canonicalUrl string) *CreateDevtoPostPayload {
	return &CreateDevtoPostPayload{Article: DevtoArticle{
		Title:        title,
		Markdown:     markdown,
		Published:    true,
		CanonicalUrl: canonicalUrl,
		Tags:         tags,
	}}
}

type DevtoArticle struct {
	Title        string   `json:"title"`
	Markdown     string   `json:"body_markdown"`
	Published    bool     `json:"published"`
	CanonicalUrl string   `json:"canonical_url"`
	Tags         []string `json:"tags"`
}
