package platforms

type Platform interface {
	Publish(title string, markdown string, tags []string, canonicalUrl string)
}
