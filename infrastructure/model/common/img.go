package common

type Img struct {
	Src string `json:"src" gorm:"column:src; type:varchar(255)"`
	Tp  string `json:"type" gorm:"column:type; type:varchar(255)"`
}
