package platforms

import (
	"github.com/go-resty/resty/v2"
	"log"
)

type Medium struct {
	apiToken      string
	publicationId string
}

func NewMedium(publicationId string, apiToken string) *Medium {
	return &Medium{apiToken: apiToken, publicationId: publicationId}
}

func (m *Medium) Publish(title string, markdown string, tags []string, canonicalUrl string) {
	client := resty.New()

	_, err := client.R().
		SetBody(NewCreateMediumPostPayload(title, markdown, tags, canonicalUrl)).
		SetHeader("Authorization", "Bearer "+m.apiToken).
		Post("https://api.medium.com/v1/users/" + m.publicationId + "/posts")

	if err != nil {
		log.Fatal(err.Error())
	}
}

type CreateMediumPostPayload struct {
	Title         string   `json:"title"`
	ContentFormat string   `json:"contentFormat"`
	Content       string   `json:"content"`
	CanonicalUrl  string   `json:"canonical_url"`
	Tags          []string `json:"tags"`
}

func NewCreateMediumPostPayload(title string, content string, tags []string, canonicalUrl string) *CreateMediumPostPayload {
	return &CreateMediumPostPayload{Title: title, ContentFormat: "markdown", Content: content, Tags: tags, CanonicalUrl: canonicalUrl}
}
