package models

type Ping struct {
	Msg  string `json:"msg" gorm:"column:msg"`
	Time string `json:"time" gorm:"column:time"`
}
