package parser

type JsonData struct {
	ID        string    `json:"id"`
	ArticleID string    `json:"article_id"`
	Category  string    `json:"category"`
	SourceURL string    `json:"source_url"`
	Paragraph paragraph `json:"paragraph"`
}