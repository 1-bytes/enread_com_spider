package parser

type JsonData struct {
	ID string `json:"id"`
	//ArticleID string    `json:"article_id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	ReleaseDate string    `json:"release_date"`
	Category    string    `json:"category"`
	SourceURL   string    `json:"source_url"`
	Paragraph   paragraph `json:"paragraph"`
}

type DictArticleModel struct {
	ID                  int    `gorm:"type:int; primaryKey; autoIncrement; unsigned; not null" json:"id"`
	Title               string `gorm:"type:varchar(255); not null" json:"title"`
	Author              string `gorm:"type:varchar(128); not null" json:"author"`
	ReleaseDate         string `gorm:"type:varchar(128); not null" json:"release_date"`
	MostRecentlyUpdated string `gorm:"type:varchar(128); not null" json:"most_recently_updated"`
}
