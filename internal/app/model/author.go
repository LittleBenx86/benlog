package model

type Author struct {
	Id        string    `json:"id" gorm:"column:id;type int(15)"`
	Name      string    `json:"name"`
	Authority Authority `json:"authority"`
}
