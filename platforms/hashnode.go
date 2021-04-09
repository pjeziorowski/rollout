package platforms

import (
	"context"
	"github.com/machinebox/graphql"
	"log"
	"strings"
)

type Hashnode struct {
	apiToken      string
	publicationId string
}

func NewHashnode(apiToken string, publicationId string) *Hashnode {
	return &Hashnode{apiToken: apiToken, publicationId: publicationId}
}

func (h *Hashnode) Publish(title string, markdown string, tags []string, canonicalUrl string) {
	hashnodeClient := graphql.NewClient("https://api.hashnode.com/")

	req := graphql.NewRequest(`
		mutation($publication: String!, $title: String!, $markdown: String!, $tags: [TagsInput]!, $canonicalUrl: String!) {
		  createPublicationStory(
			hideFromHashnodeFeed: true
			publicationId: $publication
			input: {
			  isRepublished: {
			  originalArticleURL: $canonicalUrl
			  }
			  title: $title
			  contentMarkdown: $markdown
			  tags: $tags
			}
		  ) {
			success
		  }
		}`) // TODO think of a cleaner way to manage GraphQL queries
	req.Header.Set("Authorization", h.apiToken)
	req.Var("publication", h.publicationId)
	req.Var("title", title)
	req.Var("markdown", markdown)
	hashnodeTags := h.hashnodeTags(tags)
	req.Var("tags", hashnodeTags)
	req.Var("canonicalUrl", canonicalUrl)

	if err := hashnodeClient.Run(context.Background(), req, nil); err != nil {
		log.Fatal(err)
	}
}

func (h *Hashnode) hashnodeTags(tags []string) []HashnodeTag {
	hashnodeClient := graphql.NewClient("https://api.hashnode.com/")
	req := graphql.NewRequest(`
		query {
  			tagCategories {
    			_id
				name
  			}
		}`) // TODO think of a cleaner way to manage GraphQL queries
	req.Header.Set("Authorization", h.apiToken)
	req.Var("publication", h.publicationId)

	var hashnodeTagsData HashnodeTags
	if err := hashnodeClient.Run(context.Background(), req, &hashnodeTagsData); err != nil {
		log.Fatal(err)
	}

	var desiredHashnodeTags []HashnodeTag
	for _, hashnodeTag := range hashnodeTagsData.TagCategories {
		for _, desiredTag := range tags {
			if strings.ToLower(hashnodeTag.Name) == strings.ToLower(desiredTag) {
				desiredHashnodeTags = append(desiredHashnodeTags, hashnodeTag)
			}
		}
	}

	return desiredHashnodeTags
}

type HashnodeTagsData struct {
	Data HashnodeTags `json:"data"`
}

type HashnodeTags struct {
	TagCategories []HashnodeTag `json:"tagCategories"`
}

type HashnodeTag struct {
	Name string `json:"name"`
	Id   string `json:"_id"`
}
