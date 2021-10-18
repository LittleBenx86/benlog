package model

type Blog struct {
	Id             string `json:"id"`
	Title          string `json:"title"`
	Abstract       string `json:"abstract"`
	Author         string `json:"author"`
	Content        string `json:"content"`
	Version        string `json:"version"` // Make the blog could be managed by VCS
	CreateDateTime string `json:"create_date_time"`
	UpdateDateTime string `json:"update_date_time"`
}

type BlogCategory struct {
	Id             string             `json:"id"`
	Title          string             `json:"title"`
	Creator        string             `json:"creator"`
	CreateDateTime string             `json:"create_date_time"`
	UpdateDateTime string             `json:"update_date_time"`
	Status         BlogCategoryStatus `json:"status"`
}

type BlogCategoryStatus int

const (
	CATEGORY_UNKNOWN            BlogCategoryStatus = 0
	CATEGORY_DEFAULT_READONLY   BlogCategoryStatus = 1
	CATEGORY_CUSTOMIZED_IN_USE  BlogCategoryStatus = 2
	CATEGORY_CUSTOMIZED_DELETED BlogCategoryStatus = -1
)

type BlogMetadata struct {
	Clicks int `json:"clicks"` // Just open the blog
	Views  int `json:"views"`  // Completely read the entire blog
	Votes  int `json:"votes"`
}

type BlogComment struct {
}
