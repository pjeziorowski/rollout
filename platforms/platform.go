package platforms

// Platform represents a type for interacting with publishing platforms like Hashnode, Medium, Devto
type Platform interface {
	Publish(title string, markdown string, tags []string, canonicalURL string)
}
