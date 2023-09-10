package schemaorg

type Article struct {
	Context       string `json:"@context"`
	Type          string `json:"@type"`
	Headline      string `json:"headline"`
	Image         string `json:"image"`
	URL           string `json:"url"`
	DatePublished string `json:"datePublished"`
}
