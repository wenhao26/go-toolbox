package model

type Article struct {
	Id         int    `json:"id"`
	TagId      int    `json:"tag_id"`
	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Content    string `json:"content"`
	CreatedAt  int    `json:"created_at"`
	ModifiedAt int    `json:"modified_at"`
}

func (Article) TableName() string {
	return "blog_article"
}
