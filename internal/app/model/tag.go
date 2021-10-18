package model

type Tag struct {
	Name           string    `json:"name"`
	Status         TagStatus `json:"status"`
	CreateDateTime string    `json:"create_date_time"`
}

type TagStatus int

const (
	TAG_UNKNOWN TagStatus = 0
	TAG_IN_USE  TagStatus = 1
	TAG_DELETED TagStatus = -1
)
