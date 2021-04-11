package platforms

import (
	"context"
	"github.com/machinebox/graphql"
	"log"
	"strings"
)

// Hashnode is used to interact with Hashnode platform
type Hashnode struct {
	apiToken      string
	publicationID string
}

// NewHashnode creates a new instance of Hashnode
func NewHashnode(apiToken string, publicationID string) *Hashnode {
	return &Hashnode{apiToken: apiToken, publicationID: publicationID}
}

// Publish publishes an article on Hashnode
func (h *Hashnode) Publish(title string, markdown string, tags []string, canonicalURL string) {
	hashnodeClient := graphql.NewClient("https://api.hashnode.com/")

	req := graphql.NewRequest(`
		mutation($publication: String!, $title: String!, $markdown: String!, $tags: [TagsInput]!, $canonicalURL: String!) {
		  createPublicationStory(
			hideFromHashnodeFeed: true
			publicationId: $publication
			input: {
			  isRepublished: {
			  originalArticleURL: $canonicalURL
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
	req.Var("publication", h.publicationID)
	req.Var("title", title)
	req.Var("markdown", markdown)
	hashnodeTags := h.hashnodeTags(tags)
	req.Var("tags", hashnodeTags)
	req.Var("canonicalURL", canonicalURL)

	if err := hashnodeClient.Run(context.Background(), req, nil); err != nil {
		log.Fatal(err)
	}
}

func (h *Hashnode) hashnodeTags(tags []string) []hashnodeTag {
	hashnodeClient := graphql.NewClient("https://api.hashnode.com/")
	req := graphql.NewRequest(`
		query {
  			tagCategories {
    			_id
				name
  			}
		}`) // TODO think of a cleaner way to manage GraphQL queries
	req.Header.Set("Authorization", h.apiToken)
	req.Var("publication", h.publicationID)

	var hashnodeTagsData hashnodeTags
	if err := hashnodeClient.Run(context.Background(), req, &hashnodeTagsData); err != nil {
		log.Fatal(err)
	}

	var desiredHashnodeTags = []hashnodeTag{}
	for _, hashnodeTag := range hashnodeTagsData.TagCategories {
		for _, desiredTag := range tags {
			if strings.ToLower(hashnodeTag.Name) == strings.ToLower(desiredTag) {
				desiredHashnodeTags = append(desiredHashnodeTags, hashnodeTag)
			}
		}
	}

	return desiredHashnodeTags
}

type hashnodeTagsData struct {
	Data hashnodeTags `json:"data"`
}

type hashnodeTags struct {
	TagCategories []hashnodeTag `json:"tagCategories"`
}

type hashnodeTag struct {
	Name string `json:"name"`
	ID   string `json:"_id"`
}
